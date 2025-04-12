package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"

	"aerona.thanhtd.com/ticket-service/internal/api/constants"
	"aerona.thanhtd.com/ticket-service/internal/api/dto"
	"aerona.thanhtd.com/ticket-service/internal/api/models"
	"aerona.thanhtd.com/ticket-service/internal/db/mongodb"
	"aerona.thanhtd.com/ticket-service/internal/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type TicketService struct {
	repo         *mongodb.TicketRepository
	s3Service    *S3Service
	pdfService   *PDFService
	redisService *RedisService
	producer     *kafka.Producer
	consumer     *kafka.Consumer
	processed    *ProcessedStore
	DLQTopic     string
	logger       *zap.Logger
}

type ProcessedStore struct {
	sync.Mutex
	processed map[string]bool
}

func NewTicketService(repo *mongodb.TicketRepository, s3Service *S3Service, pdfService *PDFService, redisService *RedisService,
	producer *kafka.Producer, consumer *kafka.Consumer, logger *zap.Logger) *TicketService {
	return &TicketService{
		repo:         repo,
		s3Service:    s3Service,
		pdfService:   pdfService,
		redisService: redisService,
		producer:     producer,
		consumer:     consumer,
		processed:    NewProcessedStore(),
		DLQTopic:     constants.TICKET_DLQ_TOPIC,
		logger:       logger,
	}
}

func (s *TicketService) StartConsumer(ctx context.Context) {
	go func() {
		s.consumer.SubscribeTopics([]string{constants.PAYMENT_COMPLETED_TOPIC}, nil)
		s.logger.Sugar().Infof("Listening on %s topic", constants.PAYMENT_COMPLETED_TOPIC)
		s.deliveryReportHandler()
		s.processLoop(ctx)
	}()
}

func (s *TicketService) processTicketWithRetry(booking dto.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operation := func() error {
		err := s.CreateTicket(ctx, booking)
		return err
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Second

	return backoff.Retry(operation, backoff.WithContext(b, ctx))
}

func (s *TicketService) PublishEvent(topic string, booking dto.Booking) error {
	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(string(booking.BookingId)),
		Value:          data,
	}, nil)
}

func (s *TicketService) CreateTicket(ctx context.Context, booking dto.Booking) error {
	ticket := models.Ticket{
		PNR:          booking.PNR,
		BookingId:    booking.BookingId,
		FlightId:     booking.Flight.FlightId,
		Passengers:   booking.Passengers,
		TicketNumber: "0",
	}
	// Export pdf
	departure := booking.Flight.Departure
	arrival := booking.Flight.Arrival

	duration := arrival.Scheduled.Sub(*departure.Scheduled)
	totalMinutes := int(duration.Minutes())
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	passengers := make([]Passenger, 0, len(booking.Passengers))
	for _, passenger := range booking.Passengers {
		passengers = append(passengers, Passenger{
			Title:          "Mr",
			FullName:       fmt.Sprintf("%s %s", strings.ToUpper(passenger.FirstName), strings.ToUpper(passenger.LastName)),
			Type:           "Adult",
			Route:          fmt.Sprintf("%s - %s", departure.AirportCode, arrival.AirportCode),
			HandLuggage:    "7 KG carry-on baggage",
			CheckedLuggage: "No checked baggage accepted",
			TotalWeight:    "23 KG Luggage",
		})
	}
	ticketData := TicketData{
		DateTime:         departure.Scheduled.Format("Monday, 2 January 2006"),
		DepartureTime:    fmt.Sprintf("%d:%d", departure.Scheduled.Hour(), departure.Scheduled.Minute()),
		DepartureCity:    departure.City,
		DepartureCode:    departure.AirportCode,
		DepartureAirport: fmt.Sprintf("%s Airport", departure.AirportName),
		Duration:         fmt.Sprintf("%02d:%02d", hours, minutes),
		ArrivalDate:      arrival.Scheduled.Format("02 January"),
		ArrivalTime:      fmt.Sprintf("%d:%d", arrival.Scheduled.Hour(), arrival.Scheduled.Minute()),
		ArrivalCity:      arrival.City,
		ArrivalCode:      arrival.AirportCode,
		ArrivalAirport:   fmt.Sprintf("%s Airport", arrival.AirportName),
		BookingID:        booking.BookingId,
		PNR:              booking.PNR,
		Refundable:       "No refund",
		Passengers:       passengers,
	}
	pdfBytes, err := s.pdfService.GenerateTicketPDF(ticketData)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to export e-ticket: %v", err)
		booking.Status = constants.TICKET_FAILED
		return s.PublishEvent(constants.TICKET_FAILED_TOPIC, booking)
	}

	// Upload to Amazon S3
	pdfUrl, err := s.s3Service.UploadBytes(ctx, pdfBytes, fmt.Sprintf("e-ticket/%s.pdf", ticketData.BookingID))
	if err != nil {
		s.logger.Sugar().Errorf("failed to upload e-ticket on Amazon S3: %v", err)
		booking.Status = constants.TICKET_FAILED
		return s.PublishEvent(constants.TICKET_FAILED_TOPIC, booking)
	}
	ticket.PdfUrl = pdfUrl
	// Create QR info
	qrCode := ""
	ticket.QrCode = qrCode

	_, err = s.repo.Create(ctx, ticket)
	if err != nil {
		s.logger.Sugar().Errorf("failed to upload e-ticket on Amazon S3: %v", err)
		booking.Status = constants.TICKET_FAILED
		return s.PublishEvent(constants.BOOKING_STATUS_UPDATED, booking)
	}
	booking.Status = constants.TICKET_ISSUED
	err = s.PublishEvent(constants.TICKET_CREATED_TOPIC, booking)
	if err != nil {
		return err
	}
	return s.PublishEvent(constants.BOOKING_STATUS_UPDATED, booking)
}

func (s *TicketService) GetAllTickets(ctx context.Context) ([]models.Ticket, error) {
	return s.repo.GetAllTickets(ctx)
}

func (s *TicketService) deliveryReportHandler() {
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

func (s *TicketService) processLoop(ctx context.Context) {
	defer s.consumer.Close()

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
			s.logger.Sugar().Debugf("Saved on Redis: %v", booking.BookingId)
			if err != nil {
				s.logger.Sugar().Errorf("failed to check processed status in Redis: %v", err)
				s.moveToDLQ(msg)
				s.consumer.CommitMessage(msg)
				continue
			}

			if processed {
				s.logger.Sugar().Infof("Booking %s already processed (idempotency)", booking.BookingId)
				// s.moveToDLQ(msg)
				s.consumer.CommitMessage(msg)
				continue
			}

			err = s.processTicketWithRetry(booking)
			if err != nil {
				s.logger.Sugar().Errorf("Failed to process booking %d: %v", booking.BookingId, err)
				s.moveToDLQ(msg)
			} else {
				if err := s.redisService.MarkProcessed(ctx, booking.BookingId); err != nil {
					s.logger.Sugar().Errorf("Failed to mark booking %s as processed in Redis: %v", booking.BookingId, err)
				}
			}

			s.consumer.CommitMessage(msg)
		}
	}
}

func NewProcessedStore() *ProcessedStore {
	return &ProcessedStore{
		processed: make(map[string]bool),
	}
}

func (s *TicketService) moveToDLQ(msg *kafka.Message) {
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

func (s *TicketService) validateBookingInfo(booking *dto.Booking) bool {
	return booking != nil && booking.BookingId != "" && booking.PNR != "" && booking.NumOfPassengers >= 0 &&
		booking.TotalPrice >= 0 && booking.Currency != "" && s.validateFlight(&booking.Flight) &&
		s.validateContact(&booking.Contact) && s.validatePassengers(booking.Passengers)
}

func (s *TicketService) validateContact(contact *dto.Contact) bool {

	if contact == nil {
		return false
	}
	return contact.FirstName != "" && contact.LastName != "" &&
		utils.ValidatePhoneNumber(contact.Phone) && utils.ValidateEmail(contact.Email)
}

func (s *TicketService) validatePassengers(passengers []dto.Passenger) bool {
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

func (s *TicketService) validateFlight(flight *dto.Flight) bool {
	// TODO
	return flight != nil
}
