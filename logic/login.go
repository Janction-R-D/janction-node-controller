package logic

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"janction/dao/postgres"
	"janction/model"
	"janction/setting"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetTokenOrLogin() (*string, error) {
	jwt, err := postgres.GetJWT()
	if err == nil && time.Now().Before(jwt.ExpiredAt) {
		return &jwt.Token, nil
	}

	nonce, err := getNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	token, err := login(*nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	return token, nil
}

func getNonce() (*string, error) {
	target := setting.Config.UrlConfig.JanctionBackend + "/auth/nonce"
	return performNonce(target)
}

func login(nonce string) (*string, error) {
	privateKeyStr := os.Getenv("PRIVATE_KEY")

	target := setting.Config.UrlConfig.JanctionBackend + "/auth/login"

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	message, err := siwe.InitMessage(
		"janction.io",
		address,
		"https://janction.io",
		nonce,
		map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	prepare := message.String()
	signature, err := signMessage(prepare, privateKey)
	if err != nil {
		return nil, err
	}

	token, err := performLogin(target, &model.ReqLogin{
		Message:   prepare,
		Signature: hexutil.Encode(signature),
		IsNode:    false,
	})
	if err != nil {
		return nil, err
	}
	err = postgres.SaveOrReplaceJWT(*token, time.Now().Add(8760*time.Hour))
	if err != nil {
		return token, err
	}
	return token, nil
}

func performNonce(url string) (*string, error) {
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Code int                `json:"code"`
		Msg  string             `json:"msg"`
		Data model.RespGetNonce `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Code != 1000 {
		return nil, fmt.Errorf("failed to get nonce: %s", response.Msg)
	}

	return &response.Data.Nonce, nil
}

func performLogin(url string, req *model.ReqLogin) (*string, error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data model.RespLogin `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Code != 1000 {
		return nil, fmt.Errorf("request failed: %s", response.Msg)
	}

	return &response.Data.Token, nil
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sign := signHash([]byte(message))
	signature, err := crypto.Sign(sign.Bytes(), privateKey)

	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return signature, nil
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}
