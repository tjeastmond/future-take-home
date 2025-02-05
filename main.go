package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Appointment struct {
	ID        int    `json:"id"`
	TrainerID int    `json:"trainer_id"`
	ClientId  int    `json:"user_id"`
	StartsAt  string `json:"started_at"`
	EndsAt    string `json:"ended_at"`
}

func parseAppointmentsFile(filename string) (error, []Appointment) {
	file, err := os.Open(filename)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	var appointments []Appointment
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appointments)
	return err, appointments
}

func listAppointments(c *gin.Context) {
	file, err := os.Open("docs/appointments.json")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	defer file.Close()

	err, appointments := parseAppointmentsFile("docs/appointments.json")

	c.JSON(http.StatusOK, gin.H{
		"data": appointments,
	})
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})

	router.GET("/list", listAppointments)

	router.Run(":8080")
}
