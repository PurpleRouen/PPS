package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func checkRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("certificate")
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{
				"success": false,
				"message": "file not found",
			})
			return
		}

		fileReader, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "could not open file",
			})
			return
		}

		payload, err := extractPayload(fileReader)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"payload": gin.H{
				"firstName":  payload.FirstName,
				"lastName":   payload.LastName,
				"birthDate":  payload.BirthDate,
				"gender":     payload.Gender,
				"identifier": payload.Identifier,
				"expiresAt":  payload.ExpiryDate,
			},
		})
	}
}
