package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/models"
	"github.com/tjeastmond/future-take-home/store"
)

type RouteHandler struct {
	store *store.MemoryStore
}

func NewRouteHandler(ms *store.MemoryStore) *RouteHandler {
	return &RouteHandler{store: ms}
}

func (h *RouteHandler) IDExtractor(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		c.Abort()
		return
	}
	c.Set("id", id)
	c.Next()
}

func (h *RouteHandler) getAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	appointment, exists := h.store.GetAppointmentByID(id.(int))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func (h *RouteHandler) getAppointments(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetAllAppointments())
}

func (h *RouteHandler) createAppointment(c *gin.Context) {
	var newAppointment models.Appointment
	errorMessage := ""

	if err := c.ShouldBindJSON(&newAppointment); err != nil {
		errorMessage = err.Error()
	}

	if !isValidTime(newAppointment.StartedAt) {
		errorMessage = "appointment time. Must be :00 or :30 PST."
	}

	if !h.store.IsAvailable(newAppointment) {
		errorMessage = "That time is already booked. Please select another time."
	}

	if errorMessage != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	newAppointment.EndedAt = newAppointment.StartedAt.Add(30 * time.Minute)
	appointmentID, err := h.store.AddAppointment(newAppointment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": appointmentID})
}

func (h *RouteHandler) updateAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	var updatedAppointment models.Appointment
	if err := c.ShouldBindJSON(&updatedAppointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if success := h.store.UpdateAppointment(id.(int), updatedAppointment); !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

func (h *RouteHandler) deleteAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	if success := h.store.DeleteAppointment(id.(int)); !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *RouteHandler) getTrainerAvailability(c *gin.Context) {
	err, trainerID, startTime, endTime := getTrainerIdAndDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request params"})
		return
	}

	availableSlots := h.store.GetTrainerAvailability(trainerID, startTime, endTime)

	c.JSON(http.StatusOK, gin.H{
		"trainer_id":   trainerID,
		"availability": availableSlots,
		"dates": gin.H{
			"start_date": startTime,
			"end_date":   endTime,
		},
	})
}

func (h *RouteHandler) getTrainerAppointments(c *gin.Context) {
	err, trainerID, startTime, endTime := getTrainerIdAndDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request params"})
		return
	}

	appointments := h.store.GetAppointmentsForTrainer(trainerID, startTime, endTime)

	c.JSON(http.StatusOK, gin.H{
		"trainer_id":   trainerID,
		"appointments": appointments,
		"dates": gin.H{
			"start_date": startTime,
			"end_date":   endTime,
		},
	})
}

func InitializeRoutes(router *gin.Engine, ms *store.MemoryStore) {
	handler := NewRouteHandler(ms)
	appointments := router.Group("/appointments")
	{
		appointments.GET("/", handler.getAppointments)
		appointments.POST("/", handler.createAppointment)
		appointments.GET("/:id", handler.IDExtractor, handler.getAppointment)
		appointments.PUT("/:id", handler.IDExtractor, handler.updateAppointment)
		appointments.DELETE("/:id", handler.IDExtractor, handler.deleteAppointment)
	}

	trainers := router.Group("/trainers")
	{
		trainers.GET("/:id/availability", handler.IDExtractor, handler.getTrainerAvailability)
		trainers.GET("/:id/appointments", handler.IDExtractor, handler.getTrainerAppointments)
	}
}

// Utility Functions

func getTrainerIdAndDates(c *gin.Context) (error, int, time.Time, time.Time) {
	_, start, end := getTodaysStartAndEnd()
	trainerID, _ := c.Get("id")

	// default to today's date
	startDate := c.DefaultQuery("start_date", start.Format(time.RFC3339))
	endDate := c.DefaultQuery("end_date", end.Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return err, 0, time.Time{}, time.Time{}
	}

	endTime, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return err, 0, time.Time{}, time.Time{}
	}

	return nil, trainerID.(int), startTime, endTime
}

func isValidTime(t time.Time) bool {
	_, offset := t.Zone()
	if offset != -8*3600 {
		return false
	}

	hours := t.Hour()
	minutes := t.Minute()
	weekday := t.Weekday()

	return (minutes == 0 || minutes == 30) &&
		hours >= 8 && hours < 17 &&
		weekday != time.Saturday &&
		weekday != time.Sunday
}

func getTodaysStartAndEnd() (error, time.Time, time.Time) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return err, time.Time{}, time.Time{}
	}

	now := time.Now().In(location)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, location)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, location)

	return nil, startOfDay, endOfDay
}
