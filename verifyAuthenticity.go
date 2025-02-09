package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

func verifyAuthenticity(safeUrl url.URL) (*PPSPayload, error) {
	response, err := http.DefaultClient.Get(safeUrl.String())
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("could not fetch data")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	payload := &PPSPayload{}

	doc.Find("dl > div").Each(func(i int, s *goquery.Selection) {
		title := s.Find("dt").Text()
		value := s.Find("dd").Text()

		switch strings.TrimSpace(title) {
		case "Nom":
			split := strings.Split(value, " ")
			payload.FirstName = strings.TrimSpace(split[0])
			payload.LastName = strings.TrimSpace(split[1])
		case "Date de naissance":
			payload.BirthDate = strings.TrimSpace(value)
		case "Sexe":
			if strings.TrimSpace(value) == "Masculin" {
				payload.Gender = "MALE"
			} else {
				payload.Gender = "FEMALE"
			}
		case "Numéro d'attestation":
			payload.Identifier = strings.TrimSpace(value)
		case "Valable jusqu'au":
			payload.ExpiryDate = strings.TrimSpace(value)
		}
	})

	return payload, nil
}
