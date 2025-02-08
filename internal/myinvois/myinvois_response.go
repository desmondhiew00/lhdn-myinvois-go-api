package myinvois

/* --------------------------- ID Type -------------------------- */

type IdType string

const (
	NRIC     IdType = "NRIC"
	PASSPORT IdType = "PASSPORT"
	BRN      IdType = "BRN"
	ARMY     IdType = "ARMY"
)

/* --------------------------- Document Submission -------------------------- */

type DocumentSubmissionResponse struct {
	SubmissionUid     string              `json:"submissionUid"`
	AcceptedDocuments []AcceptedDocuments `json:"acceptedDocuments"`
	RejectedDocuments []RejectedDocuments `json:"rejectedDocuments"`
}

type AcceptedDocuments struct {
	UUID              string `json:"uuid"`
	InvoiceCodeNumber string `json:"invoiceCodeNumber"`
}

type RejectedDocuments struct {
	InvoiceCodeNumber string `json:"invoiceCodeNumber"`
	Error             string `json:"error"`
}

/* ------------------------------ Get Document ------------------------------ */

type DocumentSummary struct {
	UUID                  string  `json:"uuid"`
	SubmissionUid         string  `json:"submissionUid"`
	LongId                *string `json:"longId,omitempty"`
	InternalId            string  `json:"internalId"`
	TypeName              string  `json:"typeName"`
	TypeVersionName       string  `json:"typeVersionName"`
	IssuerTin             string  `json:"issuerTin"`
	IssuerName            string  `json:"issuerName"`
	ReceiverId            string  `json:"receiverId"`
	ReceiverName          string  `json:"receiverName"`
	DateTimeIssued        string  `json:"dateTimeIssued"`
	DateTimeReceived      string  `json:"dateTimeReceived"`
	DateTimeValidated     string  `json:"dateTimeValidated"`
	TotalExcludingTax     string  `json:"totalExcludingTax"`
	TotalDiscount         string  `json:"totalDiscount"`
	TotalNetAmount        string  `json:"totalNetAmount"`
	TotalPayableAmount    string  `json:"totalPayableAmount"`
	Status                string  `json:"status"`
	CancelDateTime        string  `json:"cancelDateTime"`
	RejectRequestDateTime string  `json:"rejectRequestDateTime"`
	DocumentStatusReason  string  `json:"documentStatusReason"`
	CreatedByUserId       string  `json:"createdByUserId"`
	Document              string  `json:"document"`
}

/* ----------------------------- Cancel Document ---------------------------- */

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type CancelDocumentResponse struct {
	UUID   *string `json:"uuid"`
	Status *string `json:"status"`
	Error  *Error  `json:"error"`
}

/* ------------------------------ Reject Document ------------------------------ */

type RejectDocumentResponse struct {
	UUID   *string `json:"uuid"`
	Status *string `json:"status"`
	Error  *Error  `json:"error"`
}

/* ------------------------------ Get Submission ------------------------------ */

type GetSubmissionResponse struct {
	SubmissionUid    string            `json:"submissionUid"`
	DocumentCount    int               `json:"documentCount"`
	DateTimeReceived string            `json:"dateTimeReceived"`
	OverallStatus    string            `json:"overallStatus"`
	DocumentSummary  []DocumentSummary `json:"documentSummary"`
}

/* ---------------------------- Recent Documents ---------------------------- */

type GetRecentDocumentsParams struct {
	PageNo             *int    `json:"pageNo,omitempty"`
	PageSize           *int    `json:"pageSize,omitempty"`
	SubmissionDateFrom *string `json:"submissionDateFrom,omitempty"`
	SubmissionDateTo   *string `json:"submissionDateTo,omitempty"`
	IssueDateFrom      *string `json:"issueDateFrom,omitempty"`
	IssueDateTo        *string `json:"issueDateTo,omitempty"`
	InvoiceDirection   *string `json:"InvoiceDirection,omitempty"`
	Status             *string `json:"status,omitempty"`
	DocumentType       *string `json:"documentType,omitempty"`
	ReceiverId         *string `json:"receiverId,omitempty"`
	ReceiverIdType     *string `json:"receiverIdType,omitempty"`
	IssuerIdType       *string `json:"issuerIdType,omitempty"`
	ReceiverTin        *string `json:"receiverTin,omitempty"`
	IssuerTin          *string `json:"issuerTin,omitempty"`
	IssuerId           *string `json:"issuerId,omitempty"`
}

type GetRecentDocumentsResponse struct {
	Result   []DocumentSummary `json:"result"`
	Metadata struct {
		TotalPages int `json:"totalPages"`
		TotalCount int `json:"totalCount"`
	} `json:"metadata"`
}
