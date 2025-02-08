package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func InitEnv() {
	godotenv.Load()

	requiredEnvVars := []string{
		"CLIENT_ID",
		"CLIENT_SECRET",
		"CERT_PATH",
		"CERT_PASS",

		"SUPPLIER_NAME",
		"SUPPLIER_TIN",
		"SUPPLIER_ID_TYPE",
		"SUPPLIER_ID_VALUE",
		"SUPPLIER_MSIC_CODE",
		"SUPPLIER_BUSINESS_ACTIVITY_DESCRIPTION",
		"SUPPLIER_CONTACT_NO",
		"SUPPLIER_ADDRESS",
		"SUPPLIER_CITY",
		"SUPPLIER_STATE",
		"SUPPLIER_POSTAL_CODE",
		"SUPPLIER_COUNTRY",
	}

	missingEnvVars := []string{}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingEnvVars = append(missingEnvVars, envVar)
		}
	}

	if len(missingEnvVars) > 0 {
		log.Fatalf("Environment variables not set: [%s]", strings.Join(missingEnvVars, ", "))
	}

	// Check is cert file exists
	if _, err := os.Stat(os.Getenv("CERT_PATH")); os.IsNotExist(err) {
		log.Fatalf("Certificate file does not exist: %s", os.Getenv("CERT_PATH"))
	}

	NODE_ENV := os.Getenv("NODE_ENV")
	if NODE_ENV == "" {
		os.Setenv("NODE_ENV", "development")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		os.Setenv("PORT", "8080")
	}
}
