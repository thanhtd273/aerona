package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"aerona.thanhtd.com/notification-service/internal/api/dto"
	"aerona.thanhtd.com/notification-service/internal/api/models"
	"aerona.thanhtd.com/notification-service/internal/constants"
	mongodb "aerona.thanhtd.com/notification-service/internal/db/mongo"
	"aerona.thanhtd.com/notification-service/internal/utils"
	"github.com/cenkalti/backoff"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TICKET_EVENT              = "ticket.events"
	NOTIFICATION_EVENT        = "notification.events"
	TICKET_CREATED_QUEUE      = "ticket.created"
	TICKET_FAILED_QUEUE       = "ticket.failed"
	NOTIFICATION_SENT_QUEUE   = "notification.sent"
	NOTIFICATION_FAILED_QUEUE = "notification.failed"
)

type TicketEvent struct {
	Type      string     `json:"type"`
	BookingId int64      `json:"booking_id"`
	Data      dto.Ticket `json:"ticket"`
}

type NotificationEvent struct {
	Type      string `json:"type"`
	BookingId int64  `json:"booking_id"`
}

type NotificationService struct {
	repo         *mongodb.NotificationRepository
	emailService *EmailService
	redisService *RedisService
	producer     *kafka.Producer
	consumer     *kafka.Consumer
	DLQTopic     string
	logger       *zap.Logger
}

type ProcessedStore struct {
	sync.Mutex
}

func NewNotificationService(repo *mongodb.NotificationRepository, emailService *EmailService, redisService *RedisService,
	producer *kafka.Producer, consumer *kafka.Consumer, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		repo:         repo,
		emailService: emailService,
		redisService: redisService,
		producer:     producer,
		consumer:     consumer,
		logger:       logger,
		DLQTopic:     constants.NOTIFICATION_DLQ_TOPIC,
	}
}

func (s *NotificationService) StartConsumer(ctx context.Context) {
	go func() {
		s.consumer.SubscribeTopics([]string{constants.TICKET_CREATED_TOPIC}, nil)
		s.deliveryReportHandler()
		s.processLoop(ctx)
	}()
}

func (s *NotificationService) SendNotification(ctx context.Context, booking dto.Booking) error {

	// Send email
	title := "Flight Booking Confirmation"
	content := fmt.Sprintf("Your e-ticket, bookingId=%s", booking.BookingId)
	s.logger.Debug("Booking info", zap.Any("booking", booking))
	err := s.emailService.SendEmail(booking.Contact.Email, title, content, false)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to send email: %v", err)
		booking.Status = constants.NOTIFICATION_FAILED
		return s.publishEvent(constants.BOOKING_STATUS_UPDATED_TOPIC, booking)
	}

	// Save on MongoDB
	notification := models.Notification{
		Title:     title,
		Content:   content,
		BookingId: booking.BookingId,
		Type:      constants.EMAIL,
	}
	err = s.SaveNotification(ctx, notification)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to save notification on MongoDB: %v", err)
		booking.Status = constants.NOTIFICATION_FAILED
		return s.publishEvent(constants.BOOKING_STATUS_UPDATED_TOPIC, booking)
	}

	s.logger.Sugar().Info("Notification sent successfully for booking %s", booking.BookingId)
	booking.Status = constants.NOTIFIED
	return s.publishEvent(constants.BOOKING_STATUS_UPDATED_TOPIC, booking)
}

func (s *NotificationService) SaveNotification(ctx context.Context, notification models.Notification) error {
	notification.Id = uuid.New().String()
	notification.Status = constants.NOTIFIED
	now := time.Now()
	notification.CreatedAt = &now
	_, err := s.repo.Create(ctx, notification)
	return err
}

func (s *NotificationService) FindById(ctx context.Context, notificationId string) (*models.Notification, error) {
	return s.repo.FindById(ctx, notificationId)
}

func (s *NotificationService) deliveryReportHandler() {
	go func() {
		for e := range s.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					s.logger.Sugar().Errorf("Delivery failed s%s: %v", *ev.TopicPartition.Topic, ev.TopicPartition.Error)
				} else {
					s.logger.Sugar().Infof("Delivered to %s [%d]", *ev.TopicPartition.Topic, ev.TopicPartition.Partition)
				}
			}
		}
	}()
}

func (s *NotificationService) processLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Shutting down consumer")
			return
		default:
			msg, err := s.consumer.ReadMessage(5 * time.Second)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				s.logger.Sugar().Errorf("Consumer error: %v", err)
				continue
			}

			var booking dto.Booking
			if err := json.Unmarshal(msg.Value, &booking); err != nil {
				s.logger.Sugar().Errorf("Failed to unmarshal booking: %v", err)
				s.moveToDLQ(msg)
				s.consumer.CommitMessage(msg)
				continue
			}
			if !s.validateBookingInfo(&booking) {
				s.logger.Sugar().Errorf("Booking with ID=%s is invalid", booking.BookingId)
				// s.moveToDLQ(msg)
				s.consumer.CommitMessage(msg)
				continue
			}
			processed, err := s.redisService.IsProcessed(ctx, booking.BookingId)

			if err != nil {
				s.logger.Sugar().Errorf("failed to check processed status in Redis: %v", err)
				s.moveToDLQ(msg)
				s.consumer.CommitMessage(msg)
				continue
			}

			if processed {
				s.logger.Sugar().Infof("Booking %s already processed (idempotency)", booking.BookingId)
				s.consumer.CommitMessage(msg)
				continue
			}

			err = s.processNotificationWithRetry(booking)
			if err != nil {
				s.logger.Sugar().Errorf("Failed to process notification %d: %v", booking.BookingId, err)
				s.moveToDLQ(msg)
			} else {
				if err := s.redisService.MarkProcessed(ctx, booking.BookingId); err != nil {
					s.logger.Sugar().Errorf("Failed to mark notification %s as processed in Redis: %v", booking.BookingId, err)
				}
			}

			s.consumer.CommitMessage(msg)
		}
	}
}

func (s *NotificationService) publishEvent(topic string, booking dto.Booking) error {
	data, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(string(booking.BookingId)),
		Value:          data,
	}, nil)
}

func (s *NotificationService) moveToDLQ(msg *kafka.Message) {
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.DLQTopic, Partition: kafka.PartitionAny},
		Key:            msg.Key,
		Value:          msg.Value,
	}, nil)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to send to DLQ: %v", err)
	} else {
		s.logger.Sugar().Infof("Moved message to DLQ: %s", s.DLQTopic)
	}
}

func (s *NotificationService) processNotificationWithRetry(booking dto.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operation := func() error {
		return s.SendNotification(ctx, booking)
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Second

	return backoff.Retry(operation, backoff.WithContext(b, ctx))
}

func (s *NotificationService) validateBookingInfo(booking *dto.Booking) bool {
	return booking != nil && booking.BookingId != "" && booking.PNR != "" && booking.NumOfPassengers >= 0 &&
		booking.TotalPrice >= 0 && booking.Currency != "" && s.validateFlight(&booking.Flight) &&
		s.validateContact(&booking.Contact) && s.validatePassengers(booking.Passengers)
}

func (s *NotificationService) validateContact(contact *dto.Contact) bool {

	if contact == nil {
		return false
	}
	return contact.FirstName != "" && contact.LastName != "" &&
		utils.ValidatePhoneNumber(contact.Phone) && utils.ValidateEmail(contact.Email)
}

func (s *NotificationService) validatePassengers(passengers []dto.Passenger) bool {
	if passengers == nil {
		return false
	}
	for _, passenger := range passengers {
		if passenger.FirstName == "" || passenger.LastName == "" ||
			passenger.Nationality == "" || passenger.PassportNumber == "" || passenger.DayOfBirth == nil {
			return false
		}
	}
	return true
}

func (s *NotificationService) validateFlight(flight *dto.Flight) bool {
	// TODO
	return flight != nil
}
