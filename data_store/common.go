package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 根据私钥计算得到账户地址
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("address=%s\n", address)
	return address, nil
}

// value:本次操作的转账金额(in wei)，当调用合约里的payable函数时，value需要大于0
func GenTransactOpt(client *ethclient.Client, private string, value int64) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(private[2:])
	if err != nil {
		return nil, fmt.Errorf("privateKey %w", err)
	}
	address, err := GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("address %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("address %w", err)
	}
	log.Printf("chainID=%d\n", chainID.Int64())

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("nonce %w", err)
	}
	log.Printf("nonce=%d\n", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("gasPrice %w", err)
	}
	log.Printf("gasPrice=%d\n", gasPrice)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("auth %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(value) // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}
