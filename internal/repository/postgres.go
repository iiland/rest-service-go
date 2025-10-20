package repository

import (
	"context"
	"database/sql"
	"fmt"
	"rest-service/internal/models"

	"github.com/google/uuid"
)

// PostgresSubscriptionRepository реализует интерфейс работы с подписками через PostgreSQL
type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

// GetSum подсчитывает сумму стоимости подписок за период с фильтрами
func (r *PostgresSubscriptionRepository) GetSum(ctx context.Context, start, end string, userID uuid.UUID, serviceName string) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`
	args := []interface{}{end, start}

	// Управляем параметрами запроса динамически
	if userID != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, userID.String())
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", len(args)+1)
		args = append(args, serviceName)
	}

	var sum int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	return sum, err
}

// Create добавляет новую подписку и возвращает сгенерированный ID
func (r *PostgresSubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID.String(), sub.StartDate, sub.EndDate).Scan(&sub.ID)
	return sub.ID, err
}

// GetAll возвращает все подписки
func (r *PostgresSubscriptionRepository) GetAll(ctx context.Context) ([]models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var userID string
		err = rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &userID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}
		sub.UserID, err = uuid.Parse(userID)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// GetByID возвращает подписку по ID
func (r *PostgresSubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var sub models.Subscription
	var userID string
	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &userID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	sub.UserID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Update изменяет данные подписки по ID
func (r *PostgresSubscriptionRepository) Update(ctx context.Context, id int, sub *models.Subscription) error {
	query := `UPDATE subscriptions SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 WHERE id=$6`
	result, err := r.db.ExecContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID.String(), sub.StartDate, sub.EndDate, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete удаляет подписку по ID
func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
