package handler

import (
	"net/http"
	"os"

	document "github.com/desmondhiew00/lhdn-myinvois-go-api/internal/document"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/ubl"
	"github.com/gin-gonic/gin"
)

// @Tags Document
// @Summary Generate Invoice Document
// @Description Generate an invoice document with digital signature
// @Accept json
// @Produce json
// @Param input body ubl.InvoiceUBLInput true "Invoice input"
// @Success 200 {object} map[string]interface{} "Signed invoice document"
// @Router /document/invoice [post]
func InvoiceDocument(c *gin.Context) {

	var input ubl.InvoiceUBLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	doc := document.NewDocumentInvoice()
	doc.SetCert(os.Getenv("CERT_PATH"), os.Getenv("CERT_PASS"))
	signedUBL, err := doc.Build(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error building invoice document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": signedUBL})
}
