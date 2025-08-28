package geth_test

import (
	geth "geth"
	"testing"
)

func TestTransaction(t *testing.T) {
	geth.Transfer()
}

func TestBalance(t *testing.T) {
	geth.Balance("0x17B9E04320FAA5329f5e9Ad4f1A43DC3d836fcD2")
	geth.Balance("0xea7b53EEaFE29d122C043C56F88cCABbe3e2D6A7")
}

// go test -v -run=TestTransaction -count=1
// go test -v -run=TestBalance -count=1
