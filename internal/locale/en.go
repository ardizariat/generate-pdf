package locale

// english is read-only after initialization.
var english = Labels{
	TypeTitles: map[string]string{
		"invoice":       "Invoice",
		"letter":        "Letter",
		"report":        "Report",
		"receipt":       "Receipt",
		"delivery-note": "Delivery Note",
		"certificate":   "Certificate",
	},
	TypeDescriptions: map[string]string{
		"invoice":       "Sales invoice with an item table and VAT",
		"letter":        "Formal letter with letterhead and signature",
		"report":        "Multi-page flowing report table",
		"receipt":       "Payment receipt with the amount in words",
		"delivery-note": "Delivery note with two signatures",
		"certificate":   "Landscape certificate with a verification QR code",
	},

	InvoiceNo: "Invoice No.",
	Date:      "Date",
	Seller:    "Seller",
	Buyer:     "Buyer",
	Item:      "Item",
	Qty:       "Qty",
	Price:     "Price",
	Subtotal:  "Subtotal",
	Tax:       "Tax (11%)",
	Total:     "Total",

	Ref:       "Ref",
	Subject:   "Subject",
	Attn:      "Attn.",
	Dear:      "Dear,",
	Sincerely: "Sincerely,",

	PagePattern: "Page {current} of {total}",

	No:    "No",
	Name:  "Name",
	City:  "City",
	Email: "Email",
	Value: "Value",

	Receipt:      "RECEIPT",
	ReceiptNo:    "No",
	ReceivedFrom: "Received from",
	Amount:       "Amount",
	InWords:      "In words",
	ForPayment:   "For payment",
	ReceivedBy:   "Received by",
	CurrencyName: "dollars",

	DeliveryNote: "DELIVERY NOTE",
	DeliveryNo:   "Delivery Note No.",
	To:           "To:",
	Unit:         "Unit",
	Note:         "Note",
	Sender:       "Sender",
	Receiver:     "Receiver",

	Certificate:         "CERTIFICATE",
	CertificateFallback: "This document was generated automatically and all its data is fictitious.",

	PromptFileCount:  "Number of files",
	PromptTypes:      "Document types (select one or more, used in order)",
	PromptPages:      "Pages per file",
	PromptInteger:    "Enter a whole number",
	PromptMinOneFile: "At least 1 file",
	PromptMinOnePage: "At least 1 page",
	PromptMinOneType: "Select at least one type",
	ProgressRows:     "[%d/%d] %s · target %d rows",
	ProgressPages:    "[%d/%d] %s · target %d pages",
	BatchFailed:      "failed  %s: %v\n",
	BatchOK:          "ok     %s  %d pages  %d bytes\n",
	BatchSummary:     "summary: %d succeeded, %d failed\n",
}
