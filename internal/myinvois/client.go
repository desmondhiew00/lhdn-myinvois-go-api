package myinvois

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type TaxpayerLoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type MyInvoisClient struct {
	BaseURL           string
	accessToken       string
	accessTokenExpiry time.Time
	clientId          string
	clientSecret      string
}

func NewMyInvoisClient(clientId, clientSecret string) *MyInvoisClient {
	baseURL := "https://preprod-api.myinvois.hasil.gov.my"
	if os.Getenv("NODE_ENV") == "production" {
		baseURL = "https://api.myinvois.hasil.gov.my"
	}

	return &MyInvoisClient{
		BaseURL:      baseURL,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (p *MyInvoisClient) loginAsTaxpayer() (string, error) {
	endpoint := p.BaseURL + "/connect/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", p.clientId)
	data.Set("client_secret", p.clientSecret)
	data.Set("scope", "InvoicingAPI")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %d - %s", resp.StatusCode, string(body))
	}

	var loginResp TaxpayerLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	p.accessToken = loginResp.AccessToken
	p.accessTokenExpiry = time.Now().Add(time.Duration(loginResp.ExpiresIn) * time.Second)
	return p.accessToken, nil
}

func (p *MyInvoisClient) getAccessToken() (string, error) {
	if p.accessToken != "" && !p.accessTokenExpiry.IsZero() {
		// Refresh token if expiring in less than 10 minutes.
		if time.Now().Add(10 * time.Minute).After(p.accessTokenExpiry) {
			return p.loginAsTaxpayer()
		}
		return p.accessToken, nil
	}
	return p.loginAsTaxpayer()
}

func (p *MyInvoisClient) BearerToken() (string, error) {
	token, err := p.getAccessToken()
	if err != nil {
		return "", err
	}
	return "Bearer " + token, nil
}

func (p *MyInvoisClient) parseHTTPError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("%d - %s", resp.StatusCode, string(body))
}

func (p *MyInvoisClient) parseJSON(body []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *MyInvoisClient) handleGetRequest(endpoint string, params *url.Values) ([]byte, error) {
	auth, err := p.BearerToken()
	if err != nil {
		return nil, err
	}

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if params != nil {
		reqURL.RawQuery = params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseHTTPError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *MyInvoisClient) handlePostRequest(endpoint string, body []byte) (*http.Response, error) {
	auth, err := p.BearerToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseHTTPError(resp)
	}

	return resp, nil
}

func (p *MyInvoisClient) handlePutRequest(endpoint string, body []byte) (*http.Response, error) {
	auth, err := p.BearerToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseHTTPError(resp)
	}

	return resp, nil
}

func (p *MyInvoisClient) handleDeleteRequest(endpoint string) (*http.Response, error) {
	auth, err := p.BearerToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.parseHTTPError(resp)
	}
	return resp, nil
}
