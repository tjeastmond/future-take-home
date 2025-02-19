package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/store"
	"github.com/tjeastmond/future-take-home/utils"
)

type RouteHandler struct {
	store *store.Store
}

func NewRouteHandler(app *store.Store) *RouteHandler {
	return &RouteHandler{store: app}
}

func (h *RouteHandler) IDExtractor(c *gin.Context) {
	paramID := c.Param("trainer_id")
	trainerID, err := strconv.Atoi(paramID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Trainer ID"})
		c.Abort()
		return
	}

	c.Set("trainerID", trainerID)
	c.Next()
}

func (h *RouteHandler) trainerAppointments(c *gin.Context) {
	trainerIDStr := c.Query("trainer_id")
	trainerID, err := strconv.Atoi(trainerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainer_id"})
		return
	}

	results, _ := h.store.Appointments.GetAppointments(trainerID)
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"appointments": []store.Appointment{},
			"message":      "No appointments found for this trainer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments":           results,
		"trainer_id":             trainerID,
		"number_of_appointments": len(results),
	})
}

func (h *RouteHandler) createAppointment(c *gin.Context) {
	var newAppointment store.Appointment
	if err := c.ShouldBindJSON(&newAppointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAppointment.ToPst()

	if valid := utils.IsValidSlot(newAppointment.StartsAt, newAppointment.EndsAt); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time params"})
		return
	}

	if !h.store.Appointments.IsAvailable(&newAppointment) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "That time is already booked. Please select another time."})
		return
	}

	err := h.store.Appointments.CreateAppointment(&newAppointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment created successfully",
		"status":  "OK",
	})
}

func (h *RouteHandler) trainerAvailability(c *gin.Context) {
	trainerID, _ := c.Get("trainerID")
	startStr := c.Query("starts_at")
	endStr := c.Query("ends_at")

	start, end, valid := utils.ValidateAndParse(startStr, endStr)
	if valid == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time params"})
		return
	}

	availableSlots, _ := h.store.Appointments.GetAvailableSlots(trainerID.(int), start, end)

	c.JSON(http.StatusOK, gin.H{
		"trainer_id":   trainerID,
		"availability": availableSlots,
		"date_range": gin.H{
			"start": start,
			"end":   end,
		},
	})
}

func InitializeRoutes(router *gin.Engine, store *store.Store) {
	handler := NewRouteHandler(store)
	trainer := router.Group("/trainers")
	{
		trainer.GET("/:trainer_id", handler.IDExtractor, handler.trainerAppointments)
		trainer.POST("/:trainer_id", handler.IDExtractor, handler.createAppointment)
		trainer.GET("/:trainer_id/availability", handler.IDExtractor, handler.trainerAvailability)
	}
}
