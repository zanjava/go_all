package geth

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Transfer 转账
func Transfer() {
	client := CreateClient()
	fromPrivate := "0x4b78d2d2b849cf7f26506174d8efec83f592f96dab2ce943839a554ada03018d"
	toAddress := "0xea7b53EEaFE29d122C043C56F88cCABbe3e2D6A7"
	privateKey, err := crypto.HexToECDSA(fromPrivate[2:])
	if err != nil {
		log.Fatalf("privateKey %s", err)
	}
	address, err := GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("address %s", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("address %s", err)
	}
	log.Printf("chainID=%d\n", chainID.Int64())

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatalf("nonce %s", err)
	}
	log.Printf("nonce=%d\n", nonce)

	var gasLimit uint64 = 21000
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("gasPrice %s", err)
	}
	log.Printf("gasPrice=%d\n", gasPrice)

	// 转账1个以太币
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), big.NewInt(1e18), gasLimit, gasPrice, nil) //仅仅是转账，不调用合约函数，所以data为nil

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

// 查看账户余额
func Balance(address string) {
	client := CreateClient()
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("账户余额", balance, "wei")
}
