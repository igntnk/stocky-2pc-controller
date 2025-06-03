package models

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
