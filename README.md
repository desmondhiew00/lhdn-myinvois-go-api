# LDHN MyInvois Go Web Service

A web server built using the [Gin](https://github.com/gin-gonic/gin) framework in Golang. 
To integrate with [LHDN MyInvois APIs](https://sdk.myinvois.hasil.gov.my/einvoicingapi) to manage tasks such as creating invoices, documents, and handling digital signatures.

## Features

#### Document Submission

- [x] **Digital Signatures**: Digital signature generation for documents
- [x] **Submit Invoice**: Submit invoice to LHDN MyInvois
- [ ] **Credit Note**:
- [ ] **Debit Note**:
- [ ] **Refund Note**:
- [ ] **Self-Billed Invoice**:
- [ ] **Self-Billed Credit Note**:
- [ ] **Self-Billed Debit Note**:
- [ ] **Self-Billed Refund Note**:

#### e-Invoice APIs

- [x] **Submit Documents**
- [ ] **Cancel Document**
- [ ] **Reject Document**
- [x] **Get Recent Documents**
- [ ] **Get Submission**
- [x] **Get Document**
- [x] **Get Document Details**
- [x] **Search Documents**
- [x] **Search Taxpayer's TIN**

## Environment Variables

Create a `.env` file to configure API keys and other environment variables:

```env
PORT=8080
NODE_ENV=production

CERT_PATH=certificate-path # e.g. ./certificate.p12
CERT_PASS=certificate-password

CLIENT_ID=myinvois-client-id
CLIENT_SECRET=myinvois-client-secret

# Supplier Information (for Document UBL)
SUPPLIER_NAME=''
SUPPLIER_TIN=''
SUPPLIER_ID_TYPE=''
SUPPLIER_ID_VALUE=''
SUPPLIER_MSIC_CODE=''
SUPPLIER_MSIC_DESCRIPTION=''
SUPPLIER_SST=''
SUPPLIER_TTX=''
SUPPLIER_BUSINESS_ACTIVITY_DESCRIPTION=''
SUPPLIER_CONTACT_NO=''
SUPPLIER_EMAIL=''
SUPPLIER_ADDRESS=''
SUPPLIER_CITY=''
SUPPLIER_STATE=''
SUPPLIER_POSTAL_CODE=''
SUPPLIER_COUNTRY=''

```

## API Endpoints

- Refer: [MyInvois e-Invoice APIs](https://sdk.myinvois.hasil.gov.my/einvoicingapi)
- Swagger: http://localhost:{port}/swagger/index.html

| Method | Endpoint               | Description                  |
| ------ | ---------------------- | ---------------------------- |
| POST   | /document/invoice      | Get Invoice Document         |
| POST   | /submit-invoice        | Submit Invoice               |
| GET    | /get-invoice-qr-code   | Get Invoice QR Code          |
| -      |                        |                              |
| GET    | /document-raw          | Get Invoice Document Raw     |
| GET    | /document-details      | Get Invoice Document Details |
| GET    | /get-recent-documents  | Get Recent Documents         |
| GET    | /search-taxpayer-tin   | Search Taxpayer TIN          |
| GET    | /validate-taxpayer-tin | Validate Taxpayer TIN        |

## License

This project is licensed under the MIT License.
