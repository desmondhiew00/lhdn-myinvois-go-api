// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/document-details": {
            "get": {
                "description": "Get the details of a document by documentId",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MyInvois"
                ],
                "summary": "Get Document Details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "documentId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document details",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/document-raw": {
            "get": {
                "description": "Get the raw document by documentId",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MyInvois"
                ],
                "summary": "Get Document Raw",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "documentId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document raw",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/get-invoice-qr-code": {
            "get": {
                "description": "Get the QR code of an invoice",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MyInvois"
                ],
                "summary": "Get Invoice QR Code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "documentId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Long ID",
                        "name": "longId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Invoice QR code",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/search-taxpayer-tin": {
            "get": {
                "description": "Search for a taxpayer TIN",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MyInvois"
                ],
                "summary": "Search Taxpayer TIN",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID Type",
                        "name": "idType",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID Value",
                        "name": "idValue",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Taxpayer Name",
                        "name": "taxpayerName",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Taxpayer TIN search result",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/submit-invoice": {
            "post": {
                "description": "Submit an invoice to the e-Invoicing system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "E-Invoicing"
                ],
                "summary": "Submit Invoice",
                "parameters": [
                    {
                        "description": "Invoice input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ubl.InvoiceUBLInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Invoice submitted successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/validate-taxpayer-tin": {
            "get": {
                "description": "Validate the Taxpayer TIN",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MyInvois"
                ],
                "summary": "Validate Taxpayer TIN",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID Type",
                        "name": "idType",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID Value",
                        "name": "idValue",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Taxpayer TIN validation result",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ubl.AddressUBLInput": {
            "type": "object",
            "properties": {
                "addressLine0": {
                    "type": "string"
                },
                "addressLine1": {
                    "type": "string"
                },
                "addressLine2": {
                    "type": "string"
                },
                "cityName": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "postalZone": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                }
            }
        },
        "ubl.BuyerUBLInput": {
            "type": "object",
            "properties": {
                "address": {
                    "$ref": "#/definitions/ubl.AddressUBLInput"
                },
                "contactNumber": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "idType": {
                    "type": "string"
                },
                "idValue": {
                    "description": "NRIC, BRN, PASSPORT, ARMY",
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "sstRegistrationNumber": {
                    "type": "string"
                },
                "tin": {
                    "type": "string"
                }
            }
        },
        "ubl.InvoiceLineItemInput": {
            "type": "object",
            "properties": {
                "classificationCode": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "number"
                },
                "subtotal": {
                    "type": "number"
                },
                "taxAmount": {
                    "type": "number"
                },
                "taxRate": {
                    "type": "number"
                },
                "unitCode": {
                    "type": "string"
                },
                "unitPrice": {
                    "type": "number"
                }
            }
        },
        "ubl.InvoiceUBLInput": {
            "type": "object",
            "properties": {
                "buyer": {
                    "$ref": "#/definitions/ubl.BuyerUBLInput"
                },
                "invoiceNumber": {
                    "type": "string"
                },
                "lineItems": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ubl.InvoiceLineItemInput"
                    }
                },
                "supplier": {
                    "$ref": "#/definitions/ubl.SupplierUBLInput"
                },
                "totalExcludingTax": {
                    "type": "number"
                },
                "totalIncludingTax": {
                    "type": "number"
                },
                "totalPayableAmount": {
                    "type": "number"
                },
                "totalTaxAmount": {
                    "type": "number"
                }
            }
        },
        "ubl.SupplierUBLInput": {
            "type": "object",
            "properties": {
                "address": {
                    "$ref": "#/definitions/ubl.AddressUBLInput"
                },
                "contactNumber": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "idType": {
                    "type": "string"
                },
                "idValue": {
                    "description": "NRIC, BRN, PASSPORT, ARMY",
                    "type": "string"
                },
                "msicCode": {
                    "type": "string"
                },
                "msicDescription": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "sst": {
                    "type": "string"
                },
                "tin": {
                    "type": "string"
                },
                "ttx": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
