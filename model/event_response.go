package model

import "time"

type Event struct {
	Object            string    `json:"object"`
	ID                string    `json:"id"`
	Livemode          bool      `json:"livemode"`
	Location          string    `json:"location"`
	WebhookDeliveries []string  `json:"webhook_deliveries"`
	Data              EventData `json:"data"`
	Key               string    `json:"key"`
	CreatedAt         time.Time `json:"created_at"`
	TeamUID           string    `json:"team_uid"`
	UserUID           string    `json:"user_uid"`
}

type EventData struct {
	Object                   string            `json:"object"`
	ID                       string            `json:"id"`
	Location                 string            `json:"location"`
	Amount                   int               `json:"amount"`
	AuthorizationType        *string           `json:"authorization_type"`
	AuthorizedAmount         int               `json:"authorized_amount"`
	CapturedAmount           int               `json:"captured_amount"`
	AcquirerReferenceNum     *string           `json:"acquirer_reference_number"`
	Net                      int               `json:"net"`
	Fee                      int               `json:"fee"`
	FeeVAT                   int               `json:"fee_vat"`
	Interest                 int               `json:"interest"`
	InterestVAT              int               `json:"interest_vat"`
	FundingAmount            int               `json:"funding_amount"`
	RefundedAmount           int               `json:"refunded_amount"`
	TransactionFees          TransactionFees   `json:"transaction_fees"`
	PlatformFee              PlatformFee       `json:"platform_fee"`
	Currency                 string            `json:"currency"`
	FundingCurrency          string            `json:"funding_currency"`
	IP                       *string           `json:"ip"`
	Refunds                  Refunds           `json:"refunds"`
	Link                     *string           `json:"link"`
	Description              *string           `json:"description"`
	Metadata                 map[string]string `json:"metadata"`
	Card                     *string           `json:"card"`
	Source                   Source            `json:"source"`
	Schedule                 *string           `json:"schedule"`
	LinkedAccount            *string           `json:"linked_account"`
	Customer                 *string           `json:"customer"`
	Dispute                  *string           `json:"dispute"`
	Transaction              string            `json:"transaction"`
	FailureCode              *string           `json:"failure_code"`
	FailureMessage           *string           `json:"failure_message"`
	Status                   string            `json:"status"`
	AuthorizeURI             string            `json:"authorize_uri"`
	ReturnURI                string            `json:"return_uri"`
	CreatedAt                time.Time         `json:"created_at"`
	PaidAt                   time.Time         `json:"paid_at"`
	AuthorizedAt             time.Time         `json:"authorized_at"`
	ExpiresAt                time.Time         `json:"expires_at"`
	ExpiredAt                *time.Time        `json:"expired_at"`
	ReversedAt               *time.Time        `json:"reversed_at"`
	ZeroInterestInstallments bool              `json:"zero_interest_installments"`
	Branch                   *string           `json:"branch"`
	Terminal                 *string           `json:"terminal"`
	Device                   *string           `json:"device"`
	Authorized               bool              `json:"authorized"`
	Capturable               bool              `json:"capturable"`
	Capture                  bool              `json:"capture"`
	Disputable               bool              `json:"disputable"`
	Livemode                 bool              `json:"livemode"`
	Refundable               bool              `json:"refundable"`
	PartiallyRefundable      bool              `json:"partially_refundable"`
	Reversed                 bool              `json:"reversed"`
	Reversible               bool              `json:"reversible"`
	Voided                   bool              `json:"voided"`
	Paid                     bool              `json:"paid"`
	Expired                  bool              `json:"expired"`
	CanPerformVoid           bool              `json:"can_perform_void"`
	ApprovalCode             *string           `json:"approval_code"`
}

type TransactionFees struct {
	FeeFlat string `json:"fee_flat"`
	FeeRate string `json:"fee_rate"`
	VATRate string `json:"vat_rate"`
}

type PlatformFee struct {
	Fixed      *string `json:"fixed"`
	Amount     *string `json:"amount"`
	Percentage *string `json:"percentage"`
}

type Refunds struct {
	Object   string    `json:"object"`
	Data     []string  `json:"data"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Total    int       `json:"total"`
	Location string    `json:"location"`
	Order    string    `json:"order"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
}

type Source struct {
	Object                   string        `json:"object"`
	ID                       string        `json:"id"`
	Livemode                 bool          `json:"livemode"`
	Location                 string        `json:"location"`
	Amount                   int           `json:"amount"`
	Barcode                  *string       `json:"barcode"`
	Bank                     *string       `json:"bank"`
	CreatedAt                time.Time     `json:"created_at"`
	Currency                 string        `json:"currency"`
	Email                    *string       `json:"email"`
	Flow                     string        `json:"flow"`
	InstallmentTerm          *string       `json:"installment_term"`
	IP                       string        `json:"ip"`
	AbsorptionType           *string       `json:"absorption_type"`
	Name                     *string       `json:"name"`
	MobileNumber             *string       `json:"mobile_number"`
	PhoneNumber              *string       `json:"phone_number"`
	PlatformType             *string       `json:"platform_type"`
	ScannableCode            ScannableCode `json:"scannable_code"`
	Billing                  *string       `json:"billing"`
	Shipping                 *string       `json:"shipping"`
	Items                    []string      `json:"items"`
	References               *string       `json:"references"`
	ProviderReferences       ProviderRefs  `json:"provider_references"`
	StoreID                  *string       `json:"store_id"`
	StoreName                *string       `json:"store_name"`
	TerminalID               *string       `json:"terminal_id"`
	Type                     string        `json:"type"`
	ZeroInterestInstallments *bool         `json:"zero_interest_installments"`
	ChargeStatus             string        `json:"charge_status"`
	ReceiptAmount            *string       `json:"receipt_amount"`
	Discounts                []string      `json:"discounts"`
	PromotionCode            *string       `json:"promotion_code"`
}

type ScannableCode struct {
	Object string   `json:"object"`
	Type   string   `json:"type"`
	Image  Document `json:"image"`
}

type Document struct {
	Object      string    `json:"object"`
	Livemode    bool      `json:"livemode"`
	ID          string    `json:"id"`
	Deleted     bool      `json:"deleted"`
	Filename    string    `json:"filename"`
	Location    string    `json:"location"`
	Kind        string    `json:"kind"`
	DownloadURI string    `json:"download_uri"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProviderRefs struct {
	ReferenceNumber1 string  `json:"reference_number_1"`
	ReferenceNumber2 *string `json:"reference_number_2"`
}
