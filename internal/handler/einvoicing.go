package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/document"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/myinvois"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/ubl"
	"github.com/gin-gonic/gin"
)

type EInvoicingHandler struct {
	client *myinvois.EInvoicing
}

func NewEInvoicingHandler(clientId, clientSecret string) *EInvoicingHandler {
	return &EInvoicingHandler{
		client: myinvois.NewEInvoicing(clientId, clientSecret),
	}
}

// @Tags E-Invoicing
// @Summary Submit Invoice
// @Description Submit an invoice to the e-Invoicing system
// @Accept json
// @Produce json
// @Param input body ubl.InvoiceUBLInput true "Invoice input"
// @Success 200 {object} map[string]interface{} "Invoice submitted successfully"
// @Router /submit-invoice [post]
func (h *EInvoicingHandler) SubmitInvoice(c *gin.Context) {
	var input ubl.InvoiceUBLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	if input.InvoiceNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice number is required"})
		return
	}

	doc := document.NewDocumentInvoice()
	doc.SetCert(os.Getenv("CERT_PATH"), os.Getenv("CERT_PASS"))
	signedUBL, err := doc.Build(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error building invoice document: " + err.Error()})
		return
	}

	docJson, err := json.Marshal(signedUBL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling invoice document: " + err.Error()})
		return
	}

	result, err := h.client.SubmitDocument([]myinvois.DocumentSubmissionInput{
		{
			Document:   string(docJson),
			CodeNumber: input.InvoiceNumber,
		},
	})

	c.JSON(http.StatusOK, result)
}
