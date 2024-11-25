package bevoegdheden

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

func SearchCompanies(searchTerm, clientID, clientSecret, authServerURL string, useCache bool) ([]byte, error) {
	searchTermRegexp := regexp.MustCompile(`^[\w\-\s\.']+$`)
	if searchTerm == "" || !searchTermRegexp.MatchString(searchTerm) {
		fmt.Println("invalid searchterm")
		return nil, errors.New("invalid searchterm")
	}

	cachePath := "cache-search"

	if useCache {
		searchTerm = filepath.Clean(searchTerm)
		cachedBody, err := os.ReadFile(cachePath + "/" + searchTerm + ".json")
		if err == nil {
			fmt.Println("using cache")
			return cachedBody, nil
		}
	}

	termEscaped := url.PathEscape(searchTerm)

	match, _ := regexp.MatchString("^[0-9]{8}$", termEscaped)

	category := "name"
	if match {
		category = "id"
	}

	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authServerURL + "/token",
	}

	tok, err := conf.Token(context.Background())
	if err != nil {
		return nil, err
	}
	if !tok.Valid() {
		return nil, errors.New("invalid token")
	}

	url := fmt.Sprintf("https://api.signicat.com/info/lookup/organizations/search?countries=NL&source=kvk&%s=%s&limit=12", category, termEscaped)

	const clientConnectTimeout = time.Second * 10
	client := &http.Client{
		Transport: SafeTransport(clientConnectTimeout),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// if useCache {
	// 	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
	// 		os.MkdirAll(cachePath, 0700)
	// 	}
	// 	_ = os.WriteFile(cachePath+"/"+searchTerm+".json", body, 0644)
	// }

	return body, nil
}
