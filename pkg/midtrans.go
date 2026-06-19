package midtrans

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	serverKey  string
	httpClient *http.Client
	snapURL    string
}

func NewClient(serverKey, snapURL string) *Client {
	return &Client{
		serverKey:  serverKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		snapURL:    snapURL,
	}
}

type snapRequest struct {
	TransactionDetails transactionDetail `json:"transaction_details"`
}

type transactionDetail struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type snapResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

func (c *Client) CreatePayment(externalID string, amount float64) (paymentURL string, err error) {
	payload := snapRequest{
		TransactionDetails: transactionDetail{
			OrderID:     externalID,
			GrossAmount: amount,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.snapURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(c.serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+encoded)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call Midtrans Snap API: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Midtrans error (status %d): %s", res.StatusCode, string(resBody))
	}

	var result snapResponse
	if err := json.Unmarshal(resBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	return result.RedirectURL, nil
}
