package db

import (
	"next-oms/app/domain"
	"next-oms/app/serializers"
	"next-oms/infra/conn/db/models"
	"next-oms/infra/errors"
)

func (dc DatabaseClient) SaveOrder(order *domain.Order) (*domain.Order, *errors.RestErr) {
	mOrder := &models.Order{
		ConsignmentID:    order.ConsignmentID,
		Description:      order.Description,
		MerchantOrderID:  order.MerchantOrderID,
		RecipientName:    order.RecipientName,
		RecipientAddress: order.RecipientAddress,
		RecipientPhone:   order.RecipientPhone,
		Amount:           order.Amount,
		TotalFee:         order.TotalFee,
		Instruction:      order.Instruction,
		OrderTypeID:      order.OrderTypeID,
		CodFee:           order.CodFee,
		PromoDiscount:    order.PromoDiscount,
		Discount:         order.Discount,
		DeliveryFee:      order.DeliveryFee,
		Status:           order.Status,
		OrderType:        order.OrderType,
		ItemType:         order.ItemType,
	}

	res := dc.DB.Model(&models.Order{}).Create(&mOrder)

	if res.Error != nil {
		dc.lc.Error("error occurred when create order", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return order, nil
}

func (dc DatabaseClient) GetOrders(filters *serializers.ListFilters) (domain.Orders, *errors.RestErr) {
	var resp domain.Orders

	var totalRows int64 = 0
	tableName := "orders"
	stmt := applyFilters(dc.DB, tableName, filters, false)
	countStmt := applyFilters(dc.DB, tableName, filters, true)

	res := stmt.Find(&resp)
	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no orders found")
	}

	if res.Error != nil {
		dc.lc.Error("error occurred when getting orders", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	filters.Results = resp

	// count all data
	errCount := countStmt.Model(&models.Order{}).Count(&totalRows).Error
	if errCount != nil {
		dc.lc.Error("error occurred when getting total orders count", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	filters.TotalRows = totalRows
	filters.CalculateTotalPageAndRows(totalRows)
	filters.GeneratePagesPath()

	return resp, nil
}
