package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/desmondhiew00/lhdn-myinvois-go-api/docs"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SignRequest struct {
	Document string `json:"document"`
}

func main() {
	InitEnv()
	port := os.Getenv("PORT")

	r := gin.Default()

	myInvoisHandler := handler.NewMyInvoisHandler(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	eInvoicingHandler := handler.NewEInvoicingHandler(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	// Document API
	r.POST("/document/invoice", handler.InvoiceDocument)

	// MyInvois API
	r.GET("/document-raw", myInvoisHandler.GetDocumentRaw)
	r.GET("/document-details", myInvoisHandler.GetDocumentDetails)
	r.GET("/get-invoice-qr-code", myInvoisHandler.GetInvoiceQrCode)
	r.GET("/search-taxpayer-tin", myInvoisHandler.SearchTaxpayerTin)
	r.GET("/validate-taxpayer-tin", myInvoisHandler.ValidateTaxpayerTin)
	r.POST("/submit-invoice", eInvoicingHandler.SubmitInvoice)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("\n--------------------------------")
	fmt.Println("API URL: http://localhost:" + port)
	fmt.Println("Swagger URL: http://localhost:" + port + "/swagger/index.html")
	fmt.Println("--------------------------------\n ")

	r.Run(":" + port)

	// http.HandleFunc("/invoice-document", eInvoiceHandler)
	// http.HandleFunc("/document-raw", getDocumentHandler)
	// log.Printf("Server listening on port %s", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))

}
