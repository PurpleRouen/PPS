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
	if dataUrl[len(dataUrl)-1] != '=' {
		dataUrl += "="
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

	return &payload, nil
}
