package handler

import (
	"net/http"

	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/myinvois"
	"github.com/gin-gonic/gin"
)

type MyInvoisHandler struct {
	client *myinvois.EInvoicing
}

func NewMyInvoisHandler(clientId, clientSecret string) *MyInvoisHandler {
	return &MyInvoisHandler{
		client: myinvois.NewEInvoicing(clientId, clientSecret),
	}
}

// @Tags MyInvois
// @Summary Get Document Raw
// @Description Get the raw document by documentId
// @Accept json
// @Produce json
// @Param documentId query string true "Document ID"
// @Success 200 {object} map[string]interface{} "Document raw"
// @Router /document-raw [get]
func (h *MyInvoisHandler) GetDocumentRaw(c *gin.Context) {
	documentId := c.Query("documentId")
	if documentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "documentId is required"})
		return
	}

	res, err := h.client.GetDocument(documentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags MyInvois
// @Summary Get Document Details
// @Description Get the details of a document by documentId
// @Accept json
// @Produce json
// @Param documentId query string true "Document ID"
// @Success 200 {object} map[string]interface{} "Document details"
// @Router /document-details [get]
func (h *MyInvoisHandler) GetDocumentDetails(c *gin.Context) {
	documentId := c.Query("documentId")
	if documentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "documentId is required"})
		return
	}

	res, err := h.client.GetDocumentDetails(documentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting document details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags MyInvois
// @Summary Validate Taxpayer TIN
// @Description Validate the Taxpayer TIN
// @Accept json
// @Produce json
// @Param idType query string true "ID Type"
// @Param idValue query string true "ID Value"
// @Success 200 {object} map[string]interface{} "Taxpayer TIN validation result"
// @Router /validate-taxpayer-tin [get]
func (h *MyInvoisHandler) ValidateTaxpayerTin(c *gin.Context) {
	idType := c.Query("idType")
	idValue := c.Query("idValue")

	res, err := h.client.ValidateTaxpayerTin(idType, idValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating taxpayer TIN: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags MyInvois
// @Summary Search Taxpayer TIN
// @Description Search for a taxpayer TIN
// @Accept json
// @Produce json
// @Param idType query string true "ID Type"
// @Param idValue query string true "ID Value"
// @Param taxpayerName query string false "Taxpayer Name"
// @Success 200 {object} map[string]interface{} "Taxpayer TIN search result"
// @Router /search-taxpayer-tin [get]
func (h *MyInvoisHandler) SearchTaxpayerTin(c *gin.Context) {
	idType := c.Query("idType")
	idValue := c.Query("idValue")
	taxpayerName := c.Query("taxpayerName")

	res, err := h.client.SearchTaxpayerTin(idType, idValue, taxpayerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching taxpayer TIN: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags MyInvois
// @Summary Get Invoice QR Code
// @Description Get the QR code of an invoice
// @Accept json
// @Produce json
// @Param documentId query string true "Document ID"
// @Param longId query string true "Long ID"
// @Success 200 {object} map[string]interface{} "Invoice QR code"
// @Router /get-invoice-qr-code [get]
func (h *MyInvoisHandler) GetInvoiceQrCode(c *gin.Context) {
	documentId := c.Query("documentId")
	longId := c.Query("longId")

	if documentId == "" || longId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "documentId and longId are required"})
		return
	}

	res := h.client.GetInvoiceQrCode(documentId, longId)
	c.JSON(http.StatusOK, res)
}

// @Tags MyInvois
// @Summary Get Recent Documents
// @Description Get the recent documents
// @Accept json
// @Produce json
// @Param pageNo query string false "Page Number"
// @Param pageSize query string false "Page Size"
// @Success 200 {object} map[string]interface{} "Recent documents"
// @Router /get-recent-documents [get]
func (h *MyInvoisHandler) GetRecentDocuments(c *gin.Context) {
	params := map[string]string{}
	if c.Query("pageNo") != "" {
		params["pageNo"] = c.Query("pageNo")
	}
	if c.Query("pageSize") != "" {
		params["pageSize"] = c.Query("pageSize")
	}

	res, err := h.client.GetRecentDocuments(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting recent documents: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// TODO: ---
// RejectDocument
// GetSubmission
