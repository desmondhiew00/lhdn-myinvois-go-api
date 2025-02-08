package ubl

import (
	"fmt"
	"time"

	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/constant"
)

type InvoiceUBLInput struct {
	InvoiceNumber      string                 `json:"invoiceNumber"`
	TotalPayableAmount float64                `json:"totalPayableAmount"`
	TotalExcludingTax  float64                `json:"totalExcludingTax"`
	TotalIncludingTax  float64                `json:"totalIncludingTax"`
	TotalTaxAmount     float64                `json:"totalTaxAmount"`
	Supplier           SupplierUBLInput       `json:"supplier"`
	Buyer              BuyerUBLInput          `json:"buyer"`
	LineItems          []InvoiceLineItemInput `json:"lineItems"`
}

type InvoiceLineItemInput struct {
	ID                 int     `json:"id"`
	Quantity           float64 `json:"quantity"`
	UnitPrice          float64 `json:"unitPrice"`
	UnitCode           string  `json:"unitCode"`
	Subtotal           float64 `json:"subtotal"`
	TaxAmount          float64 `json:"taxAmount"`
	TaxRate            float64 `json:"taxRate"`
	ClassificationCode string  `json:"classificationCode"`
	Description        string  `json:"description"`
}

func InvoiceUBL(input InvoiceUBLInput, currencyCode string) map[string]interface{} {
	now := time.Now().UTC()
	invoice := map[string]interface{}{
		"ID":                   []interface{}{map[string]interface{}{"_": input.InvoiceNumber}},
		"IssueDate":            []interface{}{map[string]interface{}{"_": now.Format("2006-01-02")}},
		"IssueTime":            []interface{}{map[string]interface{}{"_": now.Format("15:04:05") + "Z"}},
		"InvoiceTypeCode":      []interface{}{map[string]interface{}{"_": constant.InvoiceTypeCode, "listVersionID": "1.1"}},
		"DocumentCurrencyCode": []interface{}{map[string]interface{}{"_": currencyCode}},
		"TaxCurrencyCode":      []interface{}{map[string]interface{}{"_": currencyCode}},
		"LegalMonetaryTotal": []interface{}{
			map[string]interface{}{
				"PayableAmount":      []interface{}{map[string]interface{}{"_": input.TotalPayableAmount, "currencyID": currencyCode}},
				"TaxExclusiveAmount": []interface{}{map[string]interface{}{"_": input.TotalExcludingTax, "currencyID": currencyCode}},
				"TaxInclusiveAmount": []interface{}{map[string]interface{}{"_": input.TotalIncludingTax, "currencyID": currencyCode}},
			},
		},
		"TaxTotal": []interface{}{
			map[string]interface{}{
				"TaxAmount": []interface{}{map[string]interface{}{"_": input.TotalTaxAmount, "currencyID": currencyCode}},
				"TaxSubtotal": []interface{}{
					map[string]interface{}{
						"TaxableAmount": []interface{}{map[string]interface{}{"_": input.TotalTaxAmount, "currencyID": currencyCode}},
						"TaxAmount":     []interface{}{map[string]interface{}{"_": input.TotalTaxAmount, "currencyID": currencyCode}},
						"TaxCategory": []interface{}{
							map[string]interface{}{
								"ID": []interface{}{map[string]interface{}{"_": "01"}},
								"TaxScheme": []interface{}{
									map[string]interface{}{
										"ID": []interface{}{map[string]interface{}{"_": "OTH", "schemeID": "UN/ECE 5153", "schemeAgencyID": "6"}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	addInvoiceItemsUBL(invoice, input.LineItems, currencyCode)
	return invoice
}

func addInvoiceItemsUBL(ubl map[string]interface{}, items []InvoiceLineItemInput, currencyCode string) {
	var lines []interface{}

	for _, item := range items {
		unitCode := item.UnitCode
		if unitCode == "" {
			unitCode = string(constant.UnitCodeOne)
		}

		line := map[string]interface{}{
			"ID":               []interface{}{map[string]interface{}{"_": fmt.Sprintf("%d", item.ID)}},
			"InvoicedQuantity": []interface{}{map[string]interface{}{"_": item.Quantity, "unitCode": unitCode}},
			"Price": []interface{}{
				map[string]interface{}{
					"PriceAmount": []interface{}{map[string]interface{}{"_": item.UnitPrice, "currencyID": currencyCode}},
				},
			},
			"ItemPriceExtension": []interface{}{
				map[string]interface{}{
					"Amount": []interface{}{map[string]interface{}{"_": item.Subtotal, "currencyID": currencyCode}},
				},
			},
			"LineExtensionAmount": []interface{}{map[string]interface{}{"_": item.Subtotal, "currencyID": currencyCode}},
			"TaxTotal": []interface{}{
				map[string]interface{}{
					"TaxAmount": []interface{}{map[string]interface{}{"_": item.TaxAmount, "currencyID": currencyCode}},
					"TaxSubtotal": []interface{}{
						map[string]interface{}{
							"TaxCategory": []interface{}{
								map[string]interface{}{
									"ID": []interface{}{map[string]interface{}{"_": "01"}},
									"TaxScheme": []interface{}{
										map[string]interface{}{
											"ID": []interface{}{map[string]interface{}{"_": "OTH", "schemeID": "UN/ECE 5153", "schemeAgencyID": "6"}},
										},
									},
								},
							},
							"TaxableAmount": []interface{}{map[string]interface{}{"_": item.Subtotal, "currencyID": currencyCode}},
							"TaxAmount":     []interface{}{map[string]interface{}{"_": item.TaxAmount, "currencyID": currencyCode}},
							"Percent":       []interface{}{map[string]interface{}{"_": item.TaxRate}},
						},
					},
				},
			},
			"Item": []interface{}{
				map[string]interface{}{
					"CommodityClassification": []interface{}{
						map[string]interface{}{
							"ItemClassificationCode": []interface{}{map[string]interface{}{"_": "9800.00.0010", "listID": "PTC"}},
						},
						map[string]interface{}{
							"ItemClassificationCode": []interface{}{map[string]interface{}{"_": item.ClassificationCode, "listID": "CLASS"}},
						},
					},
					"Description": []interface{}{map[string]interface{}{"_": item.Description}},
					"OriginCountry": []interface{}{
						map[string]interface{}{
							"IdentificationCode": []interface{}{map[string]interface{}{"_": "MYS"}},
						},
					},
				},
			},
		}

		lines = append(lines, line)
	}

	ubl["InvoiceLine"] = lines
}
