package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ContractAddress = "0xdb7b5cFbdF67776797CC36176F331a7730cBEBd8"
	Account1Private = "0x9e60fafd9657b796ffdac86472797721d44bed5dfccd3ae31d2e79547ee3329d"
)

func SetNumber(client *ethclient.Client, x int64) error {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := NewMain(contractAddress, client) //根据合约的地址加载合约。生成go代码的时候用的是--pkg=main，所以这里是NewMain
	if err != nil {
		return fmt.Errorf("NewMain %w", err)
	}

	auth, err := GenTransactOpt(client, Account1Private, 0)
	if err != nil {
		return err
	}

	_, err = contract.SetNum(auth, big.NewInt(x))
	if err != nil {
		return fmt.Errorf("SetNum %w", err)
	} else {
		return nil
	}
}

func GetNumber(client *ethclient.Client) (int64, error) {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := NewMain(contractAddress, client) //根据合约的地址加载合约。生成go代码的时候用的是--pkg=main，所以这里是NewMain
	if err != nil {
		return 0, fmt.Errorf("NewMain %w", err)
	}

	bi, err := contract.GetNum(new(bind.CallOpts))
	if err != nil {
		return 0, fmt.Errorf("SetNum %w", err)
	} else {
		return bi.Int64(), nil
	}
}

func ListenEvent() {
	// 根据ABI解析事件log
	contractAbi, err := abi.JSON(strings.NewReader(string(MainABI))) //加载合约的abi
	if err != nil {
		log.Fatal(err)
	}
	type Event struct {
		//事件参数有几个，这里就定义几个
		X *big.Int
	}

	client, err := ethclient.Dial("ws://127.0.0.1:7545") //websocket，server可主动向client推送数据
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xdb7b5cFbdF67776797CC36176F331a7730cBEBd8")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress}, //查询特定合约的区块链数据
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs) //订阅合约的事件，事件发生时，会将事件log放入logs通道
	if err != nil {
		log.Fatal(err)
	}

	//监听logs通道，获取事件log
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			bs, _ := vLog.MarshalJSON() //这行代码没用，仅仅是为了演示如何将vLog转换为json格式
			fmt.Println("完整log", string(bs))

			fmt.Println("log的Data", vLog.Data)
			var event Event
			// 根据ABI解析事件log
			err := contractAbi.UnpackIntoInterface(&event, "NumChange", vLog.Data)
			if err == nil {
				fmt.Println("X=", event.X.Uint64())
			} else {
				log.Println(err)
			}
		}
	}
}
