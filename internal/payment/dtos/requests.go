package dtos

type CreatePaymentRequest struct {
    BookingID string `json:"booking_id" binding:"required,uuid" example:"a1b2c3d4-e5f6-7890"`
    Amount    int64  `json:"amount" binding:"required,min=1" example:"10000"` // Em centavos (R$100.00 = 10000)
    Method    string `json:"method" binding:"required,oneof=card pix" example:"card"`
}