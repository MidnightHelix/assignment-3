package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Response struct {
	Water struct {
		Value  int    `json:"value"`
		Status string `json:"status"`
	} `json:"water"`
	Wind struct {
		Value  int    `json:"value"`
		Status string `json:"status"`
	} `json:"wind"`
}

func getRandomNumber(min, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return random.Intn(max-min+1) + min
}

func updateJSONFile() {
	data := Status{
		Water: getRandomNumber(1, 100),
		Wind:  getRandomNumber(1, 100),
	}
	file, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile("status.json", file, 0644)
}

// Function to handle status request
func statusHandler(c *gin.Context) {
	// Read status from JSON file
	file, err := os.Open("status.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	var status Status
	err = json.NewDecoder(file).Decode(&status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Determine status for water
	var waterStatus string
	if status.Water < 5 {
		waterStatus = "Aman"
	} else if status.Water >= 6 && status.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	// Determine status for wind
	var windStatus string
	if status.Wind < 6 {
		windStatus = "Aman"
	} else if status.Wind >= 7 && status.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	response := Response{
		Water: struct {
			Value  int    `json:"value"`
			Status string `json:"status"`
		}{Value: status.Water, Status: waterStatus},
		Wind: struct {
			Value  int    `json:"value"`
			Status string `json:"status"`
		}{Value: status.Wind, Status: windStatus},
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/public", "./public")

	// Define status endpoint
	r.GET("/status", statusHandler)

	// Update JSON file every 15 seconds
	go func() {
		for {
			updateJSONFile()
			time.Sleep(15 * time.Second)
		}
	}()

	// Start server
	fmt.Println("Server is running on port 3000")
	r.Run(":3000")
}
