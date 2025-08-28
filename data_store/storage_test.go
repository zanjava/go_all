package main

import (
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client *ethclient.Client
)

func init() {
	var err error
	client, err = ethclient.Dial("http://127.0.0.1:7545") //本地Ganache测试链
	if err != nil {
		log.Fatal(err)
	}
}

func TestStore(t *testing.T) {
	const N int64 = 2025
	SetNumber(client, N)
	if x, err := GetNumber(client); err != nil {
		t.Error(err)
	} else if x != N {
		t.Errorf("%d!=%d", x, N)
	}
}

// go test -v -run=TestStore -count=1
