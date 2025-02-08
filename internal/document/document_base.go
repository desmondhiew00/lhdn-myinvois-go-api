package document

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/ubl"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/util"
)

type DocumentBase struct {
	CertPath string
	CertPass string
}

func (d *DocumentBase) SetCert(certPath, certPass string) {
	d.CertPath = certPath
	d.CertPass = certPass
}

func (d *DocumentBase) DigitalSignature(ubl map[string]interface{}) (map[string]interface{}, error) {
	if d.CertPath == "" || d.CertPass == "" {
		return nil, errors.New("certificate path or password is not set")
	}

	ds := NewDigitalSignature()
	ublBytes, err := json.Marshal(ubl)
	if err != nil {
		return nil, err
	}
	err = ds.GenerateSignature(d.CertPath, d.CertPass, string(ublBytes))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UBLExtensions": ds.SignUBL.UBLExtension(),
		"Signature":     ds.SignUBL.Signature(),
	}, nil
}

func (d *DocumentBase) Supplier() ubl.SupplierUBLInput {
	supplier := ubl.SupplierUBLInput{
		Name:            os.Getenv("SUPPLIER_NAME"),
		TIN:             os.Getenv("SUPPLIER_TIN"),
		IdType:          os.Getenv("SUPPLIER_ID_TYPE"),
		IdValue:         os.Getenv("SUPPLIER_ID_VALUE"),
		SST:             os.Getenv("SUPPLIER_SST"),
		TTX:             os.Getenv("SUPPLIER_TTX"),
		ContactNumber:   os.Getenv("SUPPLIER_CONTACT_NO"),
		Email:           os.Getenv("SUPPLIER_EMAIL"),
		MSICCode:        os.Getenv("SUPPLIER_MSIC_CODE"),
		MSICDescription: os.Getenv("SUPPLIER_MSIC_DESCRIPTION"),
		Address: ubl.AddressUBLInput{
			AddressLine0: os.Getenv("SUPPLIER_ADDRESS"),
			CityName:     os.Getenv("SUPPLIER_CITY"),
			State:        os.Getenv("SUPPLIER_STATE"),
			PostalZone:   os.Getenv("SUPPLIER_POSTAL_CODE"),
			Country:      util.NonEmpty(os.Getenv("SUPPLIER_COUNTRY"), "Malaysia"),
		},
	}
	return supplier
}
