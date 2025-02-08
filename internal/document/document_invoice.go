package document

import (
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/constant"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/ubl"
)

type DocumentInvoice struct {
	DocumentBase
}

func NewDocumentInvoice() *DocumentInvoice {
	return &DocumentInvoice{}
}

func (d *DocumentInvoice) Build(input ubl.InvoiceUBLInput) (map[string]interface{}, error) {
	invoice := ubl.InvoiceUBL(input, constant.CurrencyCodeMYR)
	ubl.AddSupplierUBL(invoice, d.Supplier())
	ubl.AddBuyerUBL(invoice, input.Buyer)

	rawUBL := make(map[string]interface{})
	rawUBL["Invoice"] = []interface{}{invoice}
	d.addDocumentHeader(rawUBL)

	sigData, err := d.DigitalSignature(rawUBL)
	if err != nil {
		return nil, err
	}

	signedUBL := rawUBL

	if invList, ok := signedUBL["Invoice"].([]interface{}); ok && len(invList) > 0 {
		if invoiceMap, ok := invList[0].(map[string]interface{}); ok {
			invoiceMap["UBLExtensions"] = sigData["UBLExtensions"]
			invoiceMap["Signature"] = sigData["Signature"]
		}
	}

	return signedUBL, nil
}

func (d *DocumentInvoice) addDocumentHeader(rawUbl map[string]interface{}) {
	header := map[string]string{
		"_D": "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		"_A": "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		"_B": "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
	}

	for k, v := range header {
		rawUbl[k] = v
	}
}
