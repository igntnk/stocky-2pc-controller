package models

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	Uuid     string         `json:"uuid"`
	Comment  string         `json:"comment"`
	UserId   string         `json:"userId"`
	StaffId  string         `json:"staffId"`
	Price    float64        `json:"price"`
	Products []OrderProduct `json:"products"`
}

type OrderProduct struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type OrderFilter struct {
	Limit  int         `json:"limit" validate:"min=1,max=100"`
	Offset int         `json:"offset" validate:"min=0"`
	Status OrderStatus `json:"status,omitempty" validate:"omitempty,oneof=new processing completed cancelled"`
}

type OrderCreateRequest struct {
	UserID   string              `json:"user_id" validate:"required,uuid"`
	StaffID  string              `json:"staff_id" validate:"required,uuid"`
	Comment  string              `json:"comment" validate:"max=500"`
	Products []OrderProductInput `json:"products" validate:"required,min=1,dive"`
}

type OrderProductInput struct {
	ProductID uuid.UUID `json:"product_id" validate:"required,uuid"`
	Amount    int       `json:"amount" validate:"required,min=1,max=100"`
}

type OrderUpdateRequest struct {
	Comment *string      `json:"comment,omitempty" validate:"omitempty,max=500"`
	Status  *OrderStatus `json:"status,omitempty" validate:"omitempty,oneof=new processing completed cancelled"`
}

type OrderResponse struct {
	ID           string          `json:"id"`
	Comment      string          `json:"comment"`
	UserID       string          `json:"user_id"`
	StaffID      string          `json:"staff_id"`
	OrderCost    float64         `json:"order_cost"`
	Status       OrderStatus     `json:"status"`
	CreationDate string          `json:"creation_date"`
	FinishDate   *string         `json:"finish_date,omitempty"`
	Products     []ProductDetail `json:"products"`
}

type ProductDetail struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ProductCode string  `json:"product_code"`
	Amount      int     `json:"amount"`
	TotalPrice  float64 `json:"total_price"`
}
