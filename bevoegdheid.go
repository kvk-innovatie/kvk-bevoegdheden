package bevoegdheden

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
)

func isMinderjarig(gb string) bool {
	geboorteDatum, err := time.Parse(time.RFC3339, gb[6:]+"-"+gb[3:5]+"-"+gb[0:2]+"T00:00:00+00:00")
	if err != nil {
		fmt.Println(err)
		return true
	}
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return geboorteDatum.After(eighteenYearsAgo)
}

func validateInput(kvkNummer string, identityNP models.IdentityNP) error {
	kvkRegexp := regexp.MustCompile(`^\d{8}$`)
	if kvkNummer == "" || !kvkRegexp.MatchString(kvkNummer) {
		fmt.Println("kvk nummer klopt niet")
		return ErrInvalidInput
	}
	if identityNP.Voornamen == "" || identityNP.Geslachtsnaam == "" || identityNP.Geboortedatum == "" || len(identityNP.Geboortedatum) != 10 {
		fmt.Println("persoonsgegevens kloppen niet")
		return ErrInvalidInput
	}
	gb := identityNP.Geboortedatum
	_, err := time.Parse(time.RFC3339, gb[6:]+"-"+gb[3:5]+"-"+gb[0:2]+"T00:00:00+00:00")
	if err != nil {
		fmt.Println(err)
		return ErrInvalidInput
	}
	return nil
}

func GetBevoegdheid(kvkNummer string, identityNP models.IdentityNP, clientID, clientSecret, authServerURL string, useCache bool, env string) (*models.BevoegdheidResponse, error, string) {
	err := validateInput(kvkNummer, identityNP)
	if err != nil {
		return nil, err, ""
	}

	ophalenInschrijvingResponse, err, rawResponse := GetInschrijving(kvkNummer, clientID, clientSecret, authServerURL, true, env)
	if err != nil {
		return nil, err, rawResponse
	}

	bevoegdheidUittreksel := &models.BevoegdheidUittreksel{}
	paths := &models.Paths{}
	bevoegdheidResponse := &models.BevoegdheidResponse{
		BevoegdheidUittreksel: bevoegdheidUittreksel,
		Paths:                 paths,
	}

	getBevoegdheidUittreksel(bevoegdheidUittreksel, paths, ophalenInschrijvingResponse, identityNP)

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	bevoegdheidResponse.Inschrijving = ma

	return bevoegdheidResponse, nil, rawResponse
}

func GetLPID(kvkNummer string, identityNP models.IdentityNP, clientID, clientSecret, authServerURL string, useCache bool, env string) (*models.BevoegdheidResponse, error, string) {
	ophalenInschrijvingResponse, err, rawResponse := GetInschrijving(kvkNummer, clientID, clientSecret, authServerURL, true, env)
	if err != nil {
		return nil, err, rawResponse
	}

	bevoegdheidUittreksel := &models.BevoegdheidUittreksel{}
	paths := &models.Paths{}
	bevoegdheidResponse := &models.BevoegdheidResponse{
		BevoegdheidUittreksel: bevoegdheidUittreksel,
		Paths:                 paths,
	}

	getLPID(bevoegdheidUittreksel, ophalenInschrijvingResponse)

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	bevoegdheidResponse.Inschrijving = ma

	return bevoegdheidResponse, nil, rawResponse
}

func GetCompanyCertificate(kvkNummer string, identityNP models.IdentityNP, clientID, clientSecret, authServerURL string, useCache bool, env string) (*models.BevoegdheidResponse, error, string) {
	ophalenInschrijvingResponse, err, rawResponse := GetInschrijving(kvkNummer, clientID, clientSecret, authServerURL, true, env)
	if err != nil {
		return nil, err, rawResponse
	}

	bevoegdheidUittreksel := &models.BevoegdheidUittreksel{}
	paths := &models.Paths{}
	bevoegdheidResponse := &models.BevoegdheidResponse{
		BevoegdheidUittreksel: bevoegdheidUittreksel,
		Paths:                 paths,
	}

	getCompanyCertificate(bevoegdheidUittreksel, paths, ophalenInschrijvingResponse, identityNP)

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	bevoegdheidResponse.Inschrijving = ma

	return bevoegdheidResponse, nil, rawResponse
}
