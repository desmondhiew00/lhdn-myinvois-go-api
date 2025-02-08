package ubl

import "time"

type SignatureUBL struct {
	SigningTime    time.Time
	SerialNumber   string
	IssuerName     string
	CertRawContent string
	SignatureValue string
	CertDigest     string
	PropsDigest    string
	DocDigest      string
}

func NewSignatureUBL() *SignatureUBL {
	return &SignatureUBL{SigningTime: time.Now().UTC()}
}

func (s *SignatureUBL) SignedInfo() map[string]interface{} {
	return map[string]interface{}{
		"SignatureMethod": []interface{}{
			map[string]interface{}{
				"_":         "",
				"Algorithm": "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			},
		},
		"Reference": []interface{}{
			map[string]interface{}{
				"Type": "http://uri.etsi.org/01903/v1.3.2#SignedProperties",
				"URI":  "#id-xades-signed-props",
				"DigestMethod": []interface{}{
					map[string]interface{}{
						"_":         "",
						"Algorithm": "http://www.w3.org/2001/04/xmlenc#sha256",
					},
				},
				"DigestValue": []interface{}{map[string]interface{}{"_": s.PropsDigest}},
			},
			map[string]interface{}{
				"Type": "",
				"URI":  "",
				"DigestMethod": []interface{}{
					map[string]interface{}{
						"_":         "",
						"Algorithm": "http://www.w3.org/2001/04/xmlenc#sha256",
					},
				},
				"DigestValue": []interface{}{map[string]interface{}{"_": s.DocDigest}},
			},
		},
	}
}

// func (s *SignatureUBL) SignatureValue() map[string]interface{} {
// 	return map[string]interface{}{"_": s.SignatureValue}
// }

func (s *SignatureUBL) KeyInfo() map[string]interface{} {
	return map[string]interface{}{
		"X509Data": []interface{}{
			map[string]interface{}{
				"X509Certificate": []interface{}{map[string]interface{}{"_": s.CertRawContent}},
				"X509SubjectName": []interface{}{map[string]interface{}{"_": s.IssuerName}},
				"X509IssuerSerial": []interface{}{
					map[string]interface{}{
						"X509IssuerName":   []interface{}{map[string]interface{}{"_": s.IssuerName}},
						"X509SerialNumber": []interface{}{map[string]interface{}{"_": s.SerialNumber}},
					},
				},
			},
		},
	}
}

func (s *SignatureUBL) QualifyingProperties() map[string]interface{} {
	return map[string]interface{}{
		"Target": "signature",
		"SignedProperties": []interface{}{
			map[string]interface{}{
				"Id": "id-xades-signed-props",
				"SignedSignatureProperties": []interface{}{
					map[string]interface{}{
						"SigningTime": []interface{}{map[string]interface{}{"_": s.SigningTime.Format(time.RFC3339)}},
						"SigningCertificate": []interface{}{
							map[string]interface{}{
								"Cert": []interface{}{
									map[string]interface{}{
										"CertDigest": []interface{}{
											map[string]interface{}{
												"DigestMethod": []interface{}{
													map[string]interface{}{
														"_":         "",
														"Algorithm": "http://www.w3.org/2001/04/xmlenc#sha256",
													},
												},
												"DigestValue": []interface{}{map[string]interface{}{"_": s.CertDigest}},
											},
										},
										"IssuerSerial": []interface{}{
											map[string]interface{}{
												"X509IssuerName":   []interface{}{map[string]interface{}{"_": s.IssuerName}},
												"X509SerialNumber": []interface{}{map[string]interface{}{"_": s.SerialNumber}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *SignatureUBL) SignatureInformation() map[string]interface{} {
	return map[string]interface{}{
		"ID":                    []interface{}{map[string]interface{}{"_": "urn:oasis:names:specification:ubl:signature:1"}},
		"ReferencedSignatureID": []interface{}{map[string]interface{}{"_": "urn:oasis:names:specification:ubl:signature:Invoice"}},
		"Signature": []interface{}{
			map[string]interface{}{
				"Id": "signature",
				"Object": []interface{}{
					map[string]interface{}{"QualifyingProperties": []interface{}{s.QualifyingProperties()}},
				},
				"KeyInfo":        []interface{}{s.KeyInfo()},
				"SignatureValue": []interface{}{map[string]interface{}{"_": s.SignatureValue}},
				"SignedInfo":     []interface{}{s.SignedInfo()},
			},
		},
	}
}

func (s *SignatureUBL) UBLExtension() []map[string]interface{} {
	return []map[string]interface{}{
		{"UBLExtension": []interface{}{
			map[string]interface{}{
				"ExtensionURI": []interface{}{
					map[string]interface{}{"_": "urn:oasis:names:specification:ubl:dsig:enveloped:xades"},
				},
				"ExtensionContent": []interface{}{
					map[string]interface{}{
						"UBLDocumentSignatures": []interface{}{
							map[string]interface{}{
								"SignatureInformation": []interface{}{s.SignatureInformation()},
							},
						},
					},
				},
			},
		}},
	}
}

func (s *SignatureUBL) Signature() []map[string]interface{} {
	return []map[string]interface{}{
		{"ID": []interface{}{map[string]interface{}{"_": "urn:oasis:names:specification:ubl:signature:Invoice"}},
			"SignatureMethod": []interface{}{map[string]interface{}{"_": "urn:oasis:names:specification:ubl:dsig:enveloped:xades"}},
		},
	}
}
