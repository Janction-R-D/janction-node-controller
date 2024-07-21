package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"math/big"
)

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

//func VerifySignature(address, msg, sign string) bool {
//	if !strings.HasPrefix(sign, "0x") {
//		sign = "0x" + sign
//	}
//
//	signAddress := common.HexToAddress(address)
//
//	message := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
//	data := []byte(message)
//	hash := crypto.Keccak256Hash(data)
//	signature := hexutil.MustDecode(sign)
//	if signature[64] != 27 && signature[64] != 28 {
//		return false
//	}
//	signature[64] -= 27
//
//	sigPublicKey, err := crypto.SigToPub(hash.Bytes(), signature)
//	if err != nil {
//		return false
//	}
//	sigPublicKeyAddr := crypto.PubkeyToAddress(*sigPublicKey)
//	return signAddress == sigPublicKeyAddr
//}

func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case int:
		amount = decimal.NewFromFloat(float64(v))
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// sign and verify
func signMessage(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	messageHash := accounts.TextHash([]byte(message))
	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		return "", err
	}
	signature[crypto.RecoveryIDOffset] += 27
	return hexutil.Encode(signature), nil
}

func VerifySignature(fromAddress, message, signatureHex string) bool {
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false
	}
	signature[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	messageHash := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return false
	}

	if common.HexToAddress(fromAddress) != crypto.PubkeyToAddress(*pubKey) {
		return false
	}
	return true
}
