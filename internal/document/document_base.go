package document

import (
	"encoding/json"
	"errors"
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
