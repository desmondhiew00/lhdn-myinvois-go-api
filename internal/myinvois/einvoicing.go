package myinvois

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type DocumentSubmissionInput struct {
	Document   string `json:"document"`
	CodeNumber string `json:"codeNumber"`
}

type DocumentSubmissionData struct {
	Format       string `json:"format"`
	CodeNumber   string `json:"codeNumber"`
	Document     string `json:"document"`
	DocumentHash string `json:"documentHash"`
}

type DocumentSubmissionsRequest struct {
	Documents []DocumentSubmissionData `json:"documents"`
}

type EInvoicing struct {
	*MyInvoisClient
}

func NewEInvoicing(clientId, clientSecret string) *EInvoicing {
	return &EInvoicing{
		MyInvoisClient: NewMyInvoisClient(clientId, clientSecret),
	}
}

func (e *EInvoicing) ValidateTaxpayerTin(idType, idValue string) (bool, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/taxpayer/validate", e.BaseURL)
	auth, err := e.BearerToken()
	if err != nil {
		return false, err
	}

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return false, err
	}
	q := reqURL.Query()
	q.Set("idType", idType)
	q.Set("idValue", idValue)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Return true if status 200.
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (e *EInvoicing) SubmitDocument(docs []DocumentSubmissionInput) (*DocumentSubmissionResponse, error) {
	auth, err := e.BearerToken()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/api/v1.0/documentsubmissions", e.BaseURL)

	var submissions []DocumentSubmissionData
	for _, d := range docs {
		submissions = append(submissions, DocumentSubmissionData{
			Format:       "JSON",
			CodeNumber:   d.CodeNumber,
			Document:     e.stringToBase64(d.Document),
			DocumentHash: e.sha256Hash(d.Document),
		})
	}
	payload := DocumentSubmissionsRequest{Documents: submissions}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return nil, e.parseHTTPError(resp)
	}
	var result DocumentSubmissionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}
	return &result, nil
}

func (e *EInvoicing) GetDocument(documentId string) (*DocumentSummary, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documents/%s/raw", e.BaseURL, documentId)
	resBody, err := e.handleGetRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result DocumentSummary
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *EInvoicing) GetDocumentDetails(documentId string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documents/%s/details", e.BaseURL, documentId)
	resBody, err := e.handleGetRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}

	return e.parseJSON(resBody)
}

func (e *EInvoicing) CancelDocument(documentId, reason string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documents/state/%s/state", e.BaseURL, documentId)

	payload := map[string]string{
		"status": "cancelled",
		"reason": reason,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	res, err := e.handlePutRequest(endpoint, body)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *EInvoicing) RejectDocument(documentId, reason string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documents/state/%s/state", e.BaseURL, documentId)
	payload := map[string]string{
		"status": "rejected",
		"reason": reason,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	res, err := e.handlePutRequest(endpoint, body)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return e.parseJSON(resBody)
}

func (e *EInvoicing) GetSubmission(submissionId string, pageNo, pageSize int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documentsubmissions/%s", e.BaseURL, submissionId)
	resBody, err := e.handleGetRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}

	return e.parseJSON(resBody)
}

func (e *EInvoicing) GetRecentDocuments(params *GetRecentDocumentsParams) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1.0/documents/recent", e.BaseURL)
	queryParams := url.Values{}
	if params != nil {
		if params.PageNo != nil {
			queryParams.Set("pageNo", fmt.Sprintf("%d", *params.PageNo))
		}
		if params.PageSize != nil {
			queryParams.Set("pageSize", fmt.Sprintf("%d", *params.PageSize))
		}
	}
	resBody, err := e.handleGetRequest(endpoint, &queryParams)
	if err != nil {
		return nil, err
	}

	return e.parseJSON(resBody)
}

func (e *EInvoicing) SearchTaxpayerTin(idType, idValue, taxpayerName string) (string, error) {
	auth, err := e.BearerToken()
	if err != nil {
		return "", err
	}
	endpoint := fmt.Sprintf("%s/api/v1.0/taxpayer/search/tin", e.BaseURL)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	q := reqURL.Query()
	q.Set("idType", idType)
	q.Set("idValue", idValue)
	if taxpayerName != "" {
		q.Set("taxpayerName", taxpayerName)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return empty string if not OK.
		return "", nil
	}

	var result struct {
		Tin string `json:"tin"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Tin, nil
}

func (e *EInvoicing) GetInvoiceQrCode(documentId, longId string) string {
	if os.Getenv("NODE_ENV") == "production" {
		return fmt.Sprintf("https://myinvois.hasil.gov.my/%s/share/%s", documentId, longId)
	}
	return fmt.Sprintf("https://preprod.myinvois.hasil.gov.my/%s/share/%s", documentId, longId)
}

func (e *EInvoicing) sha256Hash(document string) string {
	sum := sha256.Sum256([]byte(document))
	return hex.EncodeToString(sum[:])
}

func (e *EInvoicing) stringToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
