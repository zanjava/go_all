package geth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 从文件中读取infura的API key
func GetInfuraKey() string {
	fin, err := os.Open("D:/software/key/infura")
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()
	bs, err := io.ReadAll(fin)
	if err != nil {
		log.Fatal(err)
	}
	return string(bs)
}

func CreateClient() *ethclient.Client {
	// url := "https://sepolia.infura.io/v3/" + GetInfuraKey() //这里使用sepolia，mainnet经常连不上
	url := "HTTP://127.0.0.1:7545" //本地测试Ganache网
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return client
	}
}

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

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("auth %w", err)
	}

	opt.Nonce = big.NewInt(int64(nonce))
	opt.Value = big.NewInt(value) // in wei
	opt.GasLimit = uint64(300000) // in units
	opt.GasPrice = gasPrice

	return opt, nil
}

// 获取Transaction的Sender
func GetFrom(tx *types.Transaction) (string, error) {
	if from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx); err == nil {
		return from.Hex(), nil
	} else {
		return "", err
	}
}
