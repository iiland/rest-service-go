package repository

import (
	"context"
	"rest-service/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) (int, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, id int, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
	GetSum(ctx context.Context, start, end string, userID uuid.UUID, serviceName string) (int, error)
}
