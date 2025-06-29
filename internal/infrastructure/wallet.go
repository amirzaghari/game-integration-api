package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WalletClient struct {
	BaseURL string
	APIKey  string
}

type WalletBalanceResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
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
	url := fmt.Sprintf("%s/api/v1/balance/%d", w.BaseURL, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", w.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("wallet service error: %s", resp.Status)
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
	url := fmt.Sprintf("%s/api/v1/withdraw", w.BaseURL)
	return w.doOperation(url, req)
}

func (w *WalletClient) Deposit(req WalletDepositRequest) (*WalletOperationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/deposit", w.BaseURL)
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
	httpReq.Header.Set("x-api-key", w.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		var errResp WalletErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("wallet service error: %s", errResp.Msg)
	}
	var result WalletOperationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
