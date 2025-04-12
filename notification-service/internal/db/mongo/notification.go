package mongodb

import (
	"context"
	"fmt"

	"aerona.thanhtd.com/notification-service/internal/api/models"
	"aerona.thanhtd.com/notification-service/internal/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(client *mongo.Client) *NotificationRepository {
	collection := configs.GetCollection(client, "aerona_notification_database", "notifications")
	return &NotificationRepository{
		collection: collection,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	_, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to save notification to mongodb: %v", err)
	}
	return &notification, nil
}

func (r *NotificationRepository) FindById(ctx context.Context, notificationId string) (*models.Notification, error) {
	var notification models.Notification
	err := r.collection.FindOne(ctx, bson.D{{Key: "id", Value: notificationId}}).Decode(&notification)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("not found notification for id=%s: %v", notificationId, err)
	}
	return &notification, nil
}
