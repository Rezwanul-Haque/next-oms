package models

type Order struct {
	ConsignmentID    string  `json:"order_consignment_id"`
	CreatedAt        string  `json:"order_created_at"`
	Description      string  `json:"order_description"`
	MerchantOrderID  string  `json:"merchant_order_id"`
	RecipientName    string  `json:"recipient_name"`
	RecipientAddress string  `json:"recipient_address"`
	RecipientPhone   string  `json:"recipient_phone"`
	Amount           float64 `json:"order_amount"`
	TotalFee         float64 `json:"total_fee"`
	Instruction      string  `json:"instruction"`
	OrderTypeID      int     `json:"order_type_id"`
	CodFee           float64 `json:"cod_fee"`
	PromoDiscount    float64 `json:"promo_discount"`
	Discount         float64 `json:"discount"`
	DeliveryFee      float64 `json:"delivery_fee"`
	Status           string  `json:"order_status"`
	OrderType        string  `json:"order_type"`
	ItemType         string  `json:"item_type"`
}
