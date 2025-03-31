package domain

import _ "gorm.io/gorm"

type Status string

const (
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusFailed   Status = "failed"
	StatusRefunded Status = "refunded"
)

//	 Transaction representa um transação de serviço
//		@Description	Modelo completo de transação
//		@name			Transaction
//		@model			Transaction
type Transaction struct {
	ID            string `json:"id" gorm:"primaryKey"`
	BookingID     string `json:"booking_id" gorm:"unique;not null"`
	Amount        int64  `json:"amount" gorm:"not null"`
	Currency      string `json:"currency" gorm:"not null"`
	Status        Status `json:"status" gorm:"not null"`
	PaymentMethod string `json:"payment_method" gorm:"not null"`
	GatewayID     string `json:"gateway_id" gorm:"not null"`
}
