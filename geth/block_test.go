package geth_test

import (
	"geth"
	"testing"
)

func TestReadBlock(t *testing.T) {
	geth.ReadBlock()
}

// go test -v -run=TestReadBlock -count=1
