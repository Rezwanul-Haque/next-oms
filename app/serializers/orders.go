package serializers

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"next-oms/infra/errors"
)

type OrderReq struct {
	StoreID            int     `json:"store_id"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name"`
	RecipientPhone     string  `json:"recipient_phone"`
	RecipientAddress   string  `json:"recipient_address"`
	RecipientCity      int     `json:"recipient_city"`
	RecipientZone      int     `json:"recipient_zone"`
	RecipientArea      int     `json:"recipient_area"`
	DeliveryType       int     `json:"delivery_type"`
	ItemType           int     `json:"item_type"`
	SpecialInstruction string  `json:"special_instruction"`
	ItemQuantity       int     `json:"item_quantity"`
	ItemWeight         float64 `json:"item_weight"`
	AmountToCollect    float64 `json:"amount_to_collect"`
	ItemDescription    string  `json:"item_description"`
}

type OrderResp struct {
	ConsignmentID   string  `json:"consignment_id"`
	MerchantOrderID string  `json:"merchant_order_id"`
	OrderStatus     string  `json:"order_status"`
	DeliveryFee     float64 `json:"delivery_fee"`
}

func (o OrderReq) Validate() error {
	return v.ValidateStruct(&o,
		// Fields with hardcoded values
		v.Field(&o.StoreID, v.Required, v.In(131172)),  // Must be 131172
		v.Field(&o.RecipientCity, v.Required, v.In(1)), // Must be 1
		v.Field(&o.RecipientZone, v.Required, v.In(1)), // Must be 1
		v.Field(&o.RecipientArea, v.Required, v.In(1)), // Must be 1
		v.Field(&o.DeliveryType, v.Required, v.In(48)), // Must be 48
		v.Field(&o.ItemType, v.Required, v.In(2)),      // Must be 2
		v.Field(&o.ItemQuantity, v.Required, v.In(1)),  // Must be 1
		v.Field(&o.ItemWeight, v.Required, v.In(0.5)),  // Must be 0.5

		// Fields with user input (required)
		v.Field(&o.RecipientName, v.Required),                             // Required
		v.Field(&o.RecipientPhone, v.Required, v.By(validatePhoneNumber)), // Required and custom phone validation
		v.Field(&o.RecipientAddress, v.Required),                          // Required
		v.Field(&o.AmountToCollect, v.Required),                           // Required
	)
}

// Custom phone number validation function
func validatePhoneNumber(value interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return errors.NewError("invalid")
	}

	if len(phone) < 11 {
		return errors.NewError("phone number must be at least 11 digits")
	}

	return nil
}
