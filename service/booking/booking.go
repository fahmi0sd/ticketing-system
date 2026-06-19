package booking

import (
	"time"

	routepkg "github.com/fahmi0sd/ticketing-system/service/route"
	userpkg "github.com/fahmi0sd/ticketing-system/service/user"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
)

type Booking struct {
	ID         int             `gorm:"primaryKey" json:"id"`
	UserID     int             `json:"user_id"`
	RouteID    int             `json:"route_id"`
	Quantity   int             `json:"quantity"`
	TotalPrice float64         `json:"total_price"`
	Status     Status          `json:"status"`
	PaymentURL string          `json:"payment_url,omitempty"`
	ExternalID string          `json:"external_id,omitempty"`
	ExpiredAt  *time.Time      `json:"expired_at,omitempty"`
	PaidAt     *time.Time      `json:"paid_at,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	User       *userpkg.User   `json:"user,omitempty"  gorm:"foreignKey:UserID"`
	Route      *routepkg.Route `json:"route,omitempty" gorm:"foreignKey:RouteID"`
}

type CreateRequest struct {
	RouteID  int `json:"route_id"  validate:"required"`
	Quantity int `json:"quantity"  validate:"required,min=1,max=10"`
}

type WebhookPayload struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
}
