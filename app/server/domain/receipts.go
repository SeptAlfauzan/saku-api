package domain

// FORMATTED RESPONSE
type Receipt struct {
	MerchantName    *string       `json:"merchant_name"`
	TransactionDate *string       `json:"transaction_date"`
	TransactionTime *string       `json:"transaction_time"`
	Currency        *string       `json:"currency"`
	Subtotal        *float64      `json:"subtotal"`
	Tax             *float64      `json:"tax"`
	Discount        *float64      `json:"discount"`
	Total           *float64      `json:"total"`
	PaymentMethod   *string       `json:"payment_method"`
	Items           []ReceiptItem `json:"items"`
}

type ReceiptItem struct {
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}
