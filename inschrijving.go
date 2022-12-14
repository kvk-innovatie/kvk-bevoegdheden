package bevoegdheden

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
	"github.com/kvk-innovatie/kvk-bevoegdheden/soap"
)

var ErrInschrijvingNotFound = errors.New("inschrijving niet gevonden op basis van het KVK nummer")

func getFilePath(kvkNummer string) string {
	// make it possible to annotate cached files in the filename as postfix
	// cached files start with a kvkNummer
	cachePath := "./cache-inschrijvingen/"
	files, err := ioutil.ReadDir(cachePath)
	if err != nil {
		return ""
	}
	var filePath string
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasPrefix(file.Name(), kvkNummer) {
				filePath = cachePath + file.Name()
				break
			}
		}
	}
	return filePath
}

func GetInschrijving(kvkNummer, cert, key string, useCache bool, env string) (*models.OphalenInschrijvingResponse, error) {
	cachePath := "cache-inschrijvingen"
	ophalenInschrijvingResponse := models.OphalenInschrijvingResponse{}

	if useCache {
		filePath := getFilePath(kvkNummer)
		respBody, err := os.ReadFile(filePath)
		if err == nil {
			fmt.Println("using cache")
			envelope := soap.NewEnvelope(&ophalenInschrijvingResponse)

			if err := xml.Unmarshal(respBody, envelope); err != nil {
				panic(err)
			}
			r := envelope.Body.Content.(*models.OphalenInschrijvingResponse)

			r.InschrijvingXML = string(respBody)
			return r, nil
		}
	}

	if cert == "" || key == "" {
		return nil, errors.New("no certificate or private key, so no connection possible with HRDS")
	}

	fmt.Println("not using cache")

	wsseInfo, authErr := soap.NewWSSEAuthInfo(cert, key)
	if authErr != nil {
		fmt.Printf("Auth error: %s\n", authErr.Error())
		return nil, authErr
	}

	ophalenInschrijvingRequest := models.OphalenInschrijvingRequest{
		KvkNummer: kvkNummer,
	}

	url := "https://webservices.preprod.kvk.nl/postbus2"
	toAddress := "http://es.kvk.nl/KVK-DataservicePP/2015/02"

	if env == "prd" {
		url = "https://webservices.kvk.nl/postbus2"
		toAddress = "http://es.kvk.nl/KVK-Dataservice/2015/02"
	}

	soapReq := soap.NewRequest("http://es.kvk.nl/ophalenInschrijving", url, ophalenInschrijvingRequest, &ophalenInschrijvingResponse, nil)

	soapReq.AddHeader(soap.ActionHeader{
		ID:    "_2",
		Value: "http://es.kvk.nl/ophalenInschrijving",
	})
	soapReq.AddHeader(soap.MessageIDHeader{
		ID:    "_3",
		Value: "uuid:" + uuid.New().String(),
	})
	soapReq.AddHeader(soap.ToHeader{
		ID:      "_4",
		Address: toAddress,
	})

	soapReq.SignWith(wsseInfo)

	certificate, _ := tls.X509KeyPair([]byte(cert), []byte(key))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
			},
		},
	}

	soapClient := soap.NewClient(client)

	soapResp, err := soapClient.Do(context.Background(), soapReq)
	if err != nil {
		fmt.Printf("Unable to validate: %s\n", err.Error())
		return nil, err
	} else if soapResp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to validate (status code invalid): %d\n", soapResp.StatusCode)
		return nil, err
	} else if ophalenInschrijvingResponse.Meldingen.Fout != nil {
		fmt.Printf("SOAP fault experienced during call: %s\n", ophalenInschrijvingResponse.Meldingen.Fout.Omschrijving)
		if ophalenInschrijvingResponse.Meldingen.Fout.Code == "IPD0004" {
			return nil, ErrInschrijvingNotFound
		}
		return nil, errors.New(ophalenInschrijvingResponse.Meldingen.Fout.Omschrijving)
	}

	if useCache {
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			os.MkdirAll(cachePath, 0700)
		}
		_ = os.WriteFile(cachePath+"/"+kvkNummer+".xml", soapResp.RespBody, 0644)
	}

	ophalenInschrijvingResponse.InschrijvingXML = string(soapResp.RespBody)

	// ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	// jsonMA, _ := json.MarshalIndent(ma, "", "  ")
	// _ = os.WriteFile("inschrijving.json", jsonMA, 0644)

	// data, _ := os.ReadFile("all.json")
	// bvgn := []models.Bevoegdheid{}
	// _ = json.Unmarshal(data, &bvgn)

	return &ophalenInschrijvingResponse, nil
}
