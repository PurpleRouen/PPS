package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"image"
	"image/png"
	"mime/multipart"
	"net/url"
	"strings"
)

func extractPayload(fileReader multipart.File) (*PPSPayload, error) {
	pages := []string{"1"}
	rawImages, err := pdfcpu.ExtractImagesRaw(fileReader, pages, nil)
	if err != nil {
		return nil, err
	}

	var qrCodeImg image.Image
	var convertErr error
	for _, rawImage := range rawImages {
		for _, img := range rawImage {
			if img.ObjNr < 100 {
				continue
			}
			qrCodeImg, convertErr = png.Decode(img.Reader)
			if convertErr != nil {
				continue
			} else {
				break
			}
		}
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(qrCodeImg)
	if err != nil {
		return nil, err
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("no qr code found")
	}

	unsafeUrl, err := url.Parse(result.String())
	if err != nil {
		return nil, err
	}

	if unsafeUrl.Host != "pps.athle.fr" {
		return nil, errors.New("invalid host")
	}

	queryUrl := unsafeUrl.RawQuery
	dataUrl := queryUrl[5:]

	padding := len(dataUrl) % 4
	if padding > 0 {
		dataUrl += strings.Repeat("=", 4-padding)
	}

	decodedString, err := base64.URLEncoding.DecodeString(dataUrl)
	if err != nil {
		return nil, err
	}

	var payload PPSPayload
	err = json.Unmarshal(decodedString, &payload)
	if err != nil {
		return nil, err
	}

	payload.FirstName = strings.ToUpper(payload.FirstName)
	payload.LastName = strings.ToUpper(payload.LastName)
	payload.Gender = strings.ToUpper(payload.Gender)

	verifiedPayload, err := verifyAuthenticity(*unsafeUrl)
	if err != nil {
		return nil, err
	}

	if !payload.IsSame(verifiedPayload) {
		return nil, errors.New("payload does not match")
	}

	return verifiedPayload, nil
}
