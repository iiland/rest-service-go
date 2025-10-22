package handlers

import (
	"log"
	"net/http"
	"rest-service/internal/models"
	"rest-service/internal/repository"
	"strconv"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

// Create godoc
// @Summary Create a new subscription
// @Description Create a new subscription with JSON body
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription data"
// @Success 201 {object} map[string]int "id of created subscription"
// @Failure 400 {object} map[string]string "invalid input"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.repo.Create(c.Request.Context(), &sub)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetAll godoc
// @Summary Get all subscriptions
// @Description Retrieve list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	subs, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching all subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get subscription details by its ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "invalid id"
// @Failure 404 {object} map[string]string "subscription not found"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	sub, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		} else {
			log.Printf("Error getting subscription by ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update godoc
// @Summary Update subscription
// @Description Update subscription by ID with JSON body
// @Tags subscriptions
// @Accept json
// @Param id path int true "Subscription ID"
// @Param subscription body models.Subscription true "Subscription data"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "invalid id or input"
// @Failure 404 {object} map[string]string "subscription not found"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.repo.Update(c.Request.Context(), id, &sub)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		} else {
			log.Printf("Error updating subscription: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "invalid id"
// @Failure 404 {object} map[string]string "subscription not found"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		} else {
			log.Printf("Error deleting subscription: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// GetSum godoc
// @Summary Get total cost sum for subscriptions
// @Description Get sum of subscription prices filtered by start/end dates, user ID and service name
// @Tags subscriptions
// @Produce json
// @Param start query string true "Start date in MM-YYYY"
// @Param end query string true "End date in MM-YYYY"
// @Param user_id query string true "User UUID"
// @Param service_name query string true "Service Name"
// @Success 200 {object} map[string]int "sum total cost"
// @Failure 400 {object} map[string]string "missing or invalid parameters"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /subscriptions/sum [get]
func (h *SubscriptionHandler) GetSum(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	// Проверяем обязательные параметры
	if start == "" || end == "" || userIDStr == "" || serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start, end, user_id, and service_name params are required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	sum, err := h.repo.GetSum(c.Request.Context(), start, end, userID, serviceName)
	if err != nil {
		log.Printf("Error getting sum: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sum": sum})
}
