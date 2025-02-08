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

	NODE_ENV := os.Getenv("NODE_ENV")
	if NODE_ENV == "" {
		os.Setenv("NODE_ENV", "development")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		os.Setenv("PORT", "8080")
	}
}
