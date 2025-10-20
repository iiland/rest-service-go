package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"rest-service/internal/models"
)

func TestPostgresSubscriptionRepository_GetSum(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()
	userID := uuid.New()
	start, end := "01-2023", "12-2023"
	serviceName := "Netflix"

	rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(100)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE")).
		WithArgs(end, start, userID.String(), serviceName).
		WillReturnRows(rows)

	sum, err := repo.GetSum(ctx, start, end, userID, serviceName)
	assert.NoError(t, err)
	assert.Equal(t, 100, sum)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSubscriptionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()
	sub := &models.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   "10-2025",
		EndDate:     nil,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO subscriptions`)).
		WithArgs(sub.ServiceName, sub.Price, sub.UserID.String(), sub.StartDate, sub.EndDate).
		WillReturnRows(rows)

	id, err := repo.Create(ctx, sub)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSubscriptionRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(1, "Netflix", 500, userID.String(), "10-2025", nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions")).
		WillReturnRows(rows)

	subs, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, "Netflix", subs[0].ServiceName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSubscriptionRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(1, "Netflix", 500, userID.String(), "10-2025", nil)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	sub, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, 1, sub.ID)
	assert.Equal(t, "Netflix", sub.ServiceName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSubscriptionRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()

	sub := &models.Subscription{
		ServiceName: "Netflix Updated",
		Price:       600,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     nil,
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscriptions SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 WHERE id=$6")).
		WithArgs(sub.ServiceName, sub.Price, sub.UserID.String(), sub.StartDate, sub.EndDate, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, 1, sub)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSubscriptionRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := &PostgresSubscriptionRepository{db: db}
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM subscriptions WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
