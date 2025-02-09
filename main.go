package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type PPSPayload struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	BirthDate  string `json:"birthdate"`
	Gender     string `json:"gender"`
	Identifier string `json:"pps_identifier"`
	ExpiryDate string `json:"expiry_date"`
}

func (first *PPSPayload) IsSame(second *PPSPayload) bool {
	return first.FirstName == second.FirstName &&
		first.LastName == second.LastName &&
		first.BirthDate == second.BirthDate &&
		first.Gender == second.Gender &&
		first.Identifier == second.Identifier &&
		first.ExpiryDate == second.ExpiryDate
}

func main() {
	r := gin.Default()
	r.POST("/check-pdf", checkRoute())

	err := r.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
