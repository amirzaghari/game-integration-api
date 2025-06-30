package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	WalletBalanceEndpoint  = "/api/v1/balance"
	WalletWithdrawEndpoint = "/api/v1/withdraw"
	WalletDepositEndpoint  = "/api/v1/deposit"
	WalletAPIKeyHeader     = "x-api-key"
	WalletContentType      = "application/json"
)

type WalletClient struct {
	BaseURL string
	APIKey  string
}

type WalletBalanceResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func (w *WalletBalanceResponse) UnmarshalJSON(data []byte) error {
	type Alias WalletBalanceResponse
	aux := &struct {
		Balance interface{} `json:"balance"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Balance.(type) {
	case string:
		balance, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid balance format: %v", err)
		}
		w.Balance = balance
	case float64:
		w.Balance = v
	case int:
		w.Balance = float64(v)
	case int64:
		w.Balance = float64(v)
	default:
		return fmt.Errorf("unsupported balance type: %T", v)
	}

	return nil
}

type WalletWithdrawRequest struct {
	Currency     string `json:"currency"`
	Transactions []struct {
		Amount    float64 `json:"amount"`
		BetID     int     `json:"betId"`
		Reference string  `json:"reference"`
	} `json:"transactions"`
	UserID int64 `json:"userId"`
}

type WalletDepositRequest struct {
	Currency     string `json:"currency"`
	Transactions []struct {
		Amount    float64 `json:"amount"`
		BetID     int     `json:"betId"`
		Reference string  `json:"reference"`
	} `json:"transactions"`
	UserID int64 `json:"userId"`
}

type WalletOperationResponse struct {
	Balance      float64 `json:"balance"`
	Transactions []struct {
		ID        int    `json:"id"`
		Reference string `json:"reference"`
	} `json:"transactions"`
}

type WalletErrorResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewWalletClient(baseURL, apiKey string) *WalletClient {
	return &WalletClient{BaseURL: baseURL, APIKey: apiKey}
}

func (w *WalletClient) GetBalance(userID int64) (*WalletBalanceResponse, error) {
	url := fmt.Sprintf("%s%s/%d", w.BaseURL, WalletBalanceEndpoint, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(WalletAPIKeyHeader, w.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errResp WalletErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Msg != "" {
			return nil, fmt.Errorf("wallet service error: %s", errResp.Msg)
		}
		return nil, fmt.Errorf("wallet service error: status %d", resp.StatusCode)
	}

	var result WalletBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (w *WalletClient) GetBalanceStr(walletID string) (*WalletBalanceResponse, error) {
	var id int64
	_, err := fmt.Sscan(walletID, &id)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID: %v", err)
	}
	return w.GetBalance(id)
}

func (w *WalletClient) Withdraw(req WalletWithdrawRequest) (*WalletOperationResponse, error) {
	url := fmt.Sprintf("%s%s", w.BaseURL, WalletWithdrawEndpoint)
	return w.doOperation(url, req)
}

func (w *WalletClient) Deposit(req WalletDepositRequest) (*WalletOperationResponse, error) {
	url := fmt.Sprintf("%s%s", w.BaseURL, WalletDepositEndpoint)
	return w.doOperation(url, req)
}

func (w *WalletClient) doOperation(url string, req interface{}) (*WalletOperationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(WalletAPIKeyHeader, w.APIKey)
	httpReq.Header.Set("Content-Type", WalletContentType)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errResp WalletErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Msg != "" {
			return nil, fmt.Errorf("wallet service error: %s", errResp.Msg)
		}
		return nil, fmt.Errorf("wallet service error: status %d", resp.StatusCode)
	}

	var result WalletOperationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
