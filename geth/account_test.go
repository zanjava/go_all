package geth_test

import (
	"fmt"
	geth "geth"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestCreateAccount(t *testing.T) {
	geth.CreateAccount()
}

func TestSignature(t *testing.T) {
	plainText := "zgw"
	signature := geth.SignData(plainText)
	if !geth.VerifySign(plainText, signature) {
		t.Fail()
	}
}

func TestKeccak256WithPadding(t *testing.T) {
	var value int64 = 82
	var secret int64 = 2026
	digest := geth.Keccak256WithPadding(value, secret)
	fmt.Println("盲拍出价", hexutil.Encode(digest))
	fmt.Println()

	value = 57
	secret = 2025
	digest = geth.Keccak256WithPadding(value, secret)
	fmt.Println("盲拍出价", hexutil.Encode(digest))
	fmt.Println()
}

// go test -v -run=TestKeccak256 -count=1
// go test -v -run=TestCreateAccount -count=1
// go test -v -run=TestSignature -count=1
// go test -v -run=TestKeccak256WithPadding -count=1
