package consts

import "next-oms/app/serializers"

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

const OrderPending = "Pending"

var ItemTypeMap = map[int]string{
	1: "Electronics",
	2: "Clothing",
	3: "Groceries",
	4: "Books",
	5: "Furniture",
}

func GetItemTypeDescription(itemTypeID int) string {
	if description, exists := ItemTypeMap[itemTypeID]; exists {
		return description
	}
	return "Unknown Item Type"
}

var OrderTypeMap = map[int]string{
	1: "Standard Delivery",
	2: "Express Delivery",
	3: "Same-Day Delivery",
	4: "Scheduled Delivery",
	5: "Return Order",
}

func GetOrderTypeDescription(orderTypeID int) string {
	if description, exists := OrderTypeMap[orderTypeID]; exists {
		return description
	}
	return "Unknown Order Type"
}

// OrderTypeDeliveryMap defines a mapping of order type IDs to discount fees (as percentages).
var OrderTypeDeliveryMap = map[int]float64{
	1: 0.0,  // Standard Delivery - No Discount
	2: 5.0,  // Express Delivery - 5% Discount
	3: 10.0, // Same-Day Delivery - 10% Discount
	4: 15.0, // Scheduled Delivery - 15% Discount
	5: 20.0, // Return Order - 20% Discount
}

func CalculateDeliveryFee(orderTypeID int, baseFee float64) float64 {
	delivery, exists := OrderTypeDeliveryMap[orderTypeID]
	if !exists {
		delivery = 0.0
	}
	return baseFee - (baseFee * delivery / 100)
}

var itemTypeFees = map[int]float64{
	1: 5.0,  // Regular item
	2: 10.0, // Fragile item
	3: 15.0, // Oversize item
}

const codFee = 5.0        // Flat COD fee
const weightFeeRate = 2.0 // Fee per kg
const zoneSurcharge = 3.0 // Flat zone surcharge for specific zones

func CalculateTotalFee(order serializers.OrderReq) float64 {
	// Base delivery fee
	deliveryFee, exists := OrderTypeDeliveryMap[order.DeliveryType]
	if !exists {
		deliveryFee = 0 // Default to 0 if DeliveryType not recognized
	}

	// Item type fee
	itemFee, exists := itemTypeFees[order.ItemType]
	if !exists {
		itemFee = 0 // Default to 0 if ItemType not recognized
	}

	// Weight-based fee
	weightFee := order.ItemWeight * weightFeeRate

	// COD Fee (if applicable)
	codFeeApplied := 0.0
	if order.AmountToCollect > 0 {
		codFeeApplied = codFee
	}

	// Zone Surcharge (if applicable)
	zoneFee := 0.0
	if order.RecipientZone == 1 { // Example: zone 1 requires surcharge
		zoneFee = zoneSurcharge
	}

	// Total fee calculation
	totalFee := deliveryFee + itemFee + weightFee + codFeeApplied + zoneFee
	return totalFee
}
