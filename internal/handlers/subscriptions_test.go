package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rest-service/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockSubscriptionRepository — mock реализации интерфейса
type MockSubscriptionRepository struct {
	CreateFunc  func(ctx context.Context, sub *models.Subscription) (int, error)
	GetAllFunc  func(ctx context.Context) ([]models.Subscription, error)
	GetByIDFunc func(ctx context.Context, id int) (*models.Subscription, error)
	UpdateFunc  func(ctx context.Context, id int, sub *models.Subscription) error
	DeleteFunc  func(ctx context.Context, id int) error
	GetSumFunc  func(ctx context.Context, start, end string, userID uuid.UUID, serviceName string) (int, error)
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	return m.CreateFunc(ctx, sub)
}
func (m *MockSubscriptionRepository) GetAll(ctx context.Context) ([]models.Subscription, error) {
	return m.GetAllFunc(ctx)
}
func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockSubscriptionRepository) Update(ctx context.Context, id int, sub *models.Subscription) error {
	return m.UpdateFunc(ctx, id, sub)
}
func (m *MockSubscriptionRepository) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockSubscriptionRepository) GetSum(ctx context.Context, start, end string, userID uuid.UUID, serviceName string) (int, error) {
	return m.GetSumFunc(ctx, start, end, userID, serviceName)
}

func TestSubscriptionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockSubscriptionRepository{
		CreateFunc: func(ctx context.Context, sub *models.Subscription) (int, error) {
			return 1, nil
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.POST("/subscriptions", handler.Create)

	sub := models.Subscription{UserID: uuid.New(), ServiceName: "Test", Price: 10}
	body, _ := json.Marshal(sub)
	req, _ := http.NewRequest("POST", "/subscriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["id"])
}

func TestSubscriptionHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	subs := []models.Subscription{{ID: 1, UserID: uuid.New(), ServiceName: "Test", Price: 10}}
	mockRepo := &MockSubscriptionRepository{
		GetAllFunc: func(ctx context.Context) ([]models.Subscription, error) {
			return subs, nil
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.GET("/subscriptions", handler.GetAll)

	req, _ := http.NewRequest("GET", "/subscriptions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Subscription
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Test", resp[0].ServiceName)
}

func TestSubscriptionHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sub := &models.Subscription{ID: 1, UserID: uuid.New(), ServiceName: "Test", Price: 10}
	mockRepo := &MockSubscriptionRepository{
		GetByIDFunc: func(ctx context.Context, id int) (*models.Subscription, error) {
			if id == 1 {
				return sub, nil
			}
			return nil, sql.ErrNoRows
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.GET("/subscriptions/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/subscriptions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Subscription
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, sub.ID, resp.ID)
	assert.Equal(t, sub.ServiceName, resp.ServiceName)

	// Тест на 404 при несуществующем ID
	req, _ = http.NewRequest("GET", "/subscriptions/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockSubscriptionRepository{
		UpdateFunc: func(ctx context.Context, id int, sub *models.Subscription) error {
			if id != 1 {
				return sql.ErrNoRows
			}
			return nil
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.PUT("/subscriptions/:id", handler.Update)

	updateSub := models.Subscription{ServiceName: "Updated", Price: 20, UserID: uuid.New()}
	body, _ := json.Marshal(updateSub)
	req, _ := http.NewRequest("PUT", "/subscriptions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Тест на 404 при попытке обновления несуществующего ID
	req, _ = http.NewRequest("PUT", "/subscriptions/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockSubscriptionRepository{
		DeleteFunc: func(ctx context.Context, id int) error {
			if id != 1 {
				return sql.ErrNoRows
			}
			return nil
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.DELETE("/subscriptions/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/subscriptions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Тест на 404 при удалении несуществующего ID
	req, _ = http.NewRequest("DELETE", "/subscriptions/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_GetSum(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixedUUID := uuid.New() // фиксируем UUID для запроса и мока

	mockRepo := &MockSubscriptionRepository{
		GetSumFunc: func(ctx context.Context, start, end string, userID uuid.UUID, serviceName string) (int, error) {
			// Проверяем, что приходит ожидаемый uuid
			if userID != fixedUUID {
				return 0, nil
			}
			return 150, nil
		},
	}
	handler := NewSubscriptionHandler(mockRepo)
	router := gin.New()
	router.GET("/subscriptions/total-cost", handler.GetSum)

	url := "/subscriptions/total-cost?start=01-2023&end=12-2023&user_id=" + fixedUUID.String() + "&service_name=Netflix"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 150, resp["sum"])
}
