package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/models"
	"github.com/tjeastmond/future-take-home/store"
)

// middleware to grab the ID from the URL and set it in the context
func IDExtractor(c *gin.Context) {
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

func getAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	appointment, exists := store.GetAppointmentByID(id.(int))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func getAppointments(c *gin.Context) {
	c.JSON(http.StatusOK, store.GetAllAppointments())
}

func createAppointment(c *gin.Context) {
	var newAppointment models.Appointment
	if err := c.ShouldBindJSON(&newAppointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidTime(newAppointment.StartedAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid appointment time. Must be :00 or :30 PST.",
		})
		return
	}

	newAppointment.EndedAt = newAppointment.StartedAt.Add(30 * time.Minute)
	appointmentID := store.AddAppointment(newAppointment)

	c.JSON(http.StatusCreated, gin.H{"id": appointmentID})
}

func updateAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	var updatedAppointment models.Appointment
	if err := c.ShouldBindJSON(&updatedAppointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if success := store.UpdateAppointment(id.(int), updatedAppointment); !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

func deleteAppointment(c *gin.Context) {
	id, _ := c.Get("id")
	if success := store.DeleteAppointment(id.(int)); !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func getTrainerAppointments(c *gin.Context) {}

func getTrainerAvailability(c *gin.Context) {
	trainerID := c.DefaultQuery("trainer_id", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	if trainerID == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trainer_id, start_date, and end_date are required",
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start_date format",
		})
		return
	}

	trainerIDInt, err := strconv.Atoi(trainerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid trainer_id",
		})
	}

	endTime, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end_date format",
		})
		return
	}

	appointments := store.GetAppointmentsForTrainer(trainerIDInt, startTime, endTime)
	availableSlots := getAvailableSlots(startTime, endTime, appointments)

	c.JSON(http.StatusOK, availableSlots)
}

func InitializeRoutes(router *gin.Engine) {
	appointments := router.Group("/appointments")
	{
		appointments.POST("/", createAppointment)
		appointments.GET("/", getAppointments)
		appointments.GET("/:id", IDExtractor, getAppointment)
		appointments.PUT("/:id", IDExtractor, updateAppointment)
		appointments.DELETE("/:id", IDExtractor, deleteAppointment)
		appointments.GET("/availability", getTrainerAvailability)
	}
}
