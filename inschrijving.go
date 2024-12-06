package bevoegdheden

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
	"golang.org/x/oauth2/clientcredentials"
)

var ErrInschrijvingNotFound = errors.New("inschrijving niet gevonden op basis van het KVK nummer")
var ErrParseInschrijving = errors.New("inschrijving parse error")

type KvkDataServiceResponse struct {
	MetaData struct {
		RawJson string `json:"rawJson,omitempty"`
	} `json:"metadata,omitempty"`
}

// IsDisallowedIP parses the ip to determine if we should allow the HTTP client to continue
func IsDisallowedIP(hostIP string) bool {
	ip := net.ParseIP(hostIP)
	return ip.IsMulticast() || ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate()
}

// SafeTransport uses the net.Dial to connect, then if successful check if the resolved
// ip address is disallowed. We do this due to hosts such as localhost.lol being resolvable to
// potentially malicious URLs. We allow connection only for resolution purposes.
func SafeTransport(timeout time.Duration) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				return nil, err
			}
			ip, _, _ := net.SplitHostPort(c.RemoteAddr().String())
			if IsDisallowedIP(ip) {
				return nil, errors.New("ip address is not allowed")
			}
			return c, err
		},
		DialTLS: func(network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{Timeout: timeout}
			c, err := tls.DialWithDialer(dialer, network, addr, &tls.Config{})
			if err != nil {
				return nil, err
			}

			ip, _, _ := net.SplitHostPort(c.RemoteAddr().String())
			if IsDisallowedIP(ip) {
				return nil, errors.New("ip address is not allowed")
			}

			err = c.Handshake()
			if err != nil {
				return c, err
			}

			return c, c.Handshake()
		},
		TLSHandshakeTimeout: timeout,
	}
}

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

func GetInschrijving(kvkNummer, clientID, clientSecret, authServerURL string, useCache bool, env string) (*models.OphalenInschrijvingResponse, error, string) {
	var hrResponse models.HrResponse
	if useCache {
		filePath := getFilePath(kvkNummer)
		respBody, err := os.ReadFile(filePath)
		if err == nil {
			fmt.Println("using cache")

			if err := json.Unmarshal(respBody, &hrResponse); err != nil {
				return nil, ErrParseInschrijving, string(respBody)
			}

			return &hrResponse.Envelope.Body.OphalenInschrijvingResponse, nil, string(respBody)
		}
	}
	authServerURL = strings.TrimRight(authServerURL, "/")

	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authServerURL + "/token",
	}
	fmt.Printf("ClientID URL: %#v\n", conf.ClientID)

	fmt.Printf("Token URL: %#v\n", conf.TokenURL)
	tok, err := conf.Token(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if !tok.Valid() {
		fmt.Printf("token invalid. got: %#v", tok)
	}

	url := fmt.Sprintf("https://api.signicat.com/info/lookup/countries/nl/organizations/%s?source=kvk-dataservice&rawJSON=true", kvkNummer)

	const clientConnectTimeout = time.Second * 10
	client := &http.Client{
		Transport: SafeTransport(clientConnectTimeout),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp, err := client.Do(req)
	if resp.StatusCode == 404 {
		return nil, ErrInschrijvingNotFound, ""
	} else if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var kvkDataServiceResponse KvkDataServiceResponse
	err = json.Unmarshal(body, &kvkDataServiceResponse)
	if err != nil {
		panic(err)
	}

	// if useCache {
	// 	cachePath := "cache-inschrijvingen"
	// 	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
	// 		os.MkdirAll(cachePath, 0700)
	// 	}
	// 	_ = os.WriteFile(cachePath+"/"+kvkNummer+".json", []byte(kvkDataServiceResponse.MetaData.RawJson), 0644)
	// }
	err = json.Unmarshal([]byte(kvkDataServiceResponse.MetaData.RawJson), &hrResponse)
	if err != nil {
		return nil, ErrParseInschrijving, kvkDataServiceResponse.MetaData.RawJson
	}

	return &hrResponse.Envelope.Body.OphalenInschrijvingResponse, nil, kvkDataServiceResponse.MetaData.RawJson
}
