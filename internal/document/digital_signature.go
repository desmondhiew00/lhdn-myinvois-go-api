package document

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/ubl"
	"golang.org/x/crypto/pkcs12"
)

type CertData struct {
	Cert     *x509.Certificate
	CertPEM  string
	KeyPEM   string
	CertHash string
}

var (
	certCache = make(map[string]*CertData)
	cacheMu   sync.Mutex
)

type DigitalSignature struct {
	SignUBL *ubl.SignatureUBL
}

func NewDigitalSignature() *DigitalSignature {
	return &DigitalSignature{
		SignUBL: ubl.NewSignatureUBL(),
	}
}

func (ds *DigitalSignature) GenerateSignature(certPath, certPass, documentStr string) error {
	certData, err := loadCert(certPath, certPass)
	if err != nil {
		return err
	}

	cert := certData.Cert
	certPEM := certData.CertPEM
	keyPEM := certData.KeyPEM

	issuer, err := buildIssuer(cert)
	if err != nil {
		return err
	}

	docHash := computeSHA256(documentStr)

	ds.SignUBL.SerialNumber = cert.SerialNumber.String()
	ds.SignUBL.IssuerName = issuer
	sigValue, err := signDocument(keyPEM, documentStr)
	if err != nil {
		return err
	}
	ds.SignUBL.SignatureValue = sigValue
	ds.SignUBL.CertRawContent = cleanPEMCertificate(certPEM)
	ds.SignUBL.CertDigest = certData.CertHash
	ds.SignUBL.DocDigest = docHash

	// Compute the properties digest.
	props := ds.SignUBL.QualifyingProperties()
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return err
	}
	ds.SignUBL.PropsDigest = computeSHA256(string(propsJSON))
	return nil
}

/* ------------------------ Internal Helper Functions ----------------------- */

// loadCert reads and decodes a certificate file (currently only supports PKCS#12),
// caching the result.
func loadCert(certPath, certPass string) (*CertData, error) {
	cacheKey := fmt.Sprintf("%s-%s", certPath, certPass)
	cacheMu.Lock()
	if data, exists := certCache[cacheKey]; exists {
		cacheMu.Unlock()
		return data, nil
	}
	cacheMu.Unlock()

	content, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(certPath))
	if ext != ".p12" && ext != ".pfx" {
		return nil, errors.New("unsupported certificate format")
	}

	key, cert, err := pkcs12.Decode(content, certPass)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PKCS#12 file: %w", err)
	}

	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	var keyPEM string
	switch k := key.(type) {
	case *rsa.PrivateKey:
		keyBytes := x509.MarshalPKCS1PrivateKey(k)
		keyPEM = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}))
	default:
		keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, err
		}
		keyPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}))
	}

	hash := sha256.Sum256(cert.Raw)
	certHash := base64.StdEncoding.EncodeToString(hash[:])

	data := &CertData{
		Cert:     cert,
		CertPEM:  certPEM,
		KeyPEM:   keyPEM,
		CertHash: certHash,
	}
	cacheMu.Lock()
	certCache[cacheKey] = data
	cacheMu.Unlock()
	return data, nil
}

// buildIssuer constructs an issuer string from the certificate attributes.
func buildIssuer(cert *x509.Certificate) (string, error) {
	var parts []string
	if cert.Issuer.CommonName != "" {
		parts = append(parts, "CN="+cert.Issuer.CommonName)
	}
	// Search for email using the OID for email.
	for _, name := range cert.Issuer.Names {
		if name.Type.String() == "1.2.840.113549.1.9.1" {
			if email, ok := name.Value.(string); ok && email != "" {
				parts = append(parts, "E="+email)
				break
			}
		}
	}
	if len(cert.Issuer.OrganizationalUnit) > 0 {
		parts = append(parts, "OU="+cert.Issuer.OrganizationalUnit[0])
	}
	if len(cert.Issuer.Organization) > 0 {
		parts = append(parts, "O="+cert.Issuer.Organization[0])
	}
	if len(cert.Issuer.Country) > 0 {
		parts = append(parts, "C="+cert.Issuer.Country[0])
	}
	return strings.Join(parts, ", "), nil
}

// computeSHA256 returns the base64-encoded SHA256 hash of the input string.
func computeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// cleanPEMCertificate removes PEM headers/footers and newlines.
func cleanPEMCertificate(pemCert string) string {
	cert := strings.ReplaceAll(pemCert, "-----BEGIN CERTIFICATE-----", "")
	cert = strings.ReplaceAll(cert, "-----END CERTIFICATE-----", "")
	return strings.TrimSpace(strings.ReplaceAll(cert, "\n", ""))
}

// signDocument signs the document string with the private key.
func signDocument(keyPEM, documentStr string) (string, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block for key")
	}
	var privKey interface{}
	var err error
	privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		privKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	digest := sha256.Sum256([]byte(documentStr))
	var signature []byte
	switch key := privKey.(type) {
	case *rsa.PrivateKey:
		signature, err = rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			return "", err
		}
	case *ecdsa.PrivateKey:
		var r, s *big.Int
		r, s, err = ecdsa.Sign(rand.Reader, key, digest[:])
		if err != nil {
			return "", err
		}
		signature, err = asn1.Marshal(struct{ R, S *big.Int }{r, s})
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("unsupported private key type")
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
