package mongodb

import (
	"context"
	"time"

	"aerona.thanhtd.com/ticket-service/internal/api/models"
	"aerona.thanhtd.com/ticket-service/internal/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(client *mongo.Client) *TicketRepository {
	collection := configs.GetCollection(client, "aerona_ticket_database", "tickets")
	return &TicketRepository{collection: collection}
}

func (r *TicketRepository) Create(ctx context.Context, ticket models.Ticket) (*models.Ticket, error) {

	now := time.Now()
	ticket.CreatedAt = &now
	r.collection.InsertOne(ctx, ticket)
	return &ticket, nil
}

func (r *TicketRepository) GetAllTickets(ctx context.Context) ([]models.Ticket, error) {
	var tickets []models.Ticket
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &tickets)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
