package encryption_test

import (
	"bytes"
	"fmt"
	"go/base/encryption"
	"testing"
)

func TestMyEncryption1(t *testing.T) {
	key := [8]byte{34, 65, 12, 125, 65, 70, 54, 27}

	algo := encryption.NewMyEncryption1(key)
	plain := []byte("明月多情应笑我")
	cypher := algo.Encrypt(plain)
	fmt.Println(cypher)
	plain2, err := algo.Decrypt(cypher)
	fmt.Println(string(plain2))
	if err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(plain, plain2) { // 比较两个byte切片里的元素是否完全相等
			fmt.Println(len(plain2), string(plain2))
			t.Fail()
		}
	}
}

func TestMyEncryption2(t *testing.T) {
	key := [8]byte{34, 65, 12, 125, 65, 70, 54, 27}

	algo := encryption.NewMyEncryption2(key, encryption.NONE)
	plain := []byte("明月多情应笑我")
	cypher := algo.Encrypt(plain)
	fmt.Println(cypher)
	plain2, err := algo.Decrypt(cypher)
	fmt.Println(string(plain2))
	if err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(plain, plain2) { // 比较两个byte切片里的元素是否完全相等
			fmt.Println(len(plain2), string(plain2))
			t.Fail()
		}
	}
}

// go test -v ./encryption -run=^TestMyEncryption$ -count=1
