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

func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	subs, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching all subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

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
