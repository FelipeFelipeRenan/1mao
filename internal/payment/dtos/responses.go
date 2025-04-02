package dtos

type PaymentResponse struct {
    ID        string `json:"id" example:"a1b2c3d4-e5f6-7890"`
    Status    string `json:"status" example:"pending"`
    GatewayID string `json:"gateway_id" example:"pi_123456789"`
}