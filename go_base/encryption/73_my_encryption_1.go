package encryption

import (
	"errors"
	PKCS7 "go/base/basic"
)

// MyEncryption1 is a simple XOR encryption algorithm
// It uses a fixed key of 8 bytes and operates on blocks of 8 bytes.
// 对称加密
type MyEncryption1 struct {
	key [8]byte
}

func NewMyEncryption1(key [8]byte) *MyEncryption1 {
	return &MyEncryption1{
		key: key,
	}
}

func (en *MyEncryption1) BlockSize() int {
	return 8
}

func (en *MyEncryption1) Encrypt(plain []byte) []byte {
	plainPadding := PKCS7.Padding(plain, en.BlockSize()) //明文末尾填充字节，长度成为BlockSize的整倍数
	cypher := make([]byte, len(plainPadding))
	for i := 0; i < len(plainPadding); i++ {
		// 异或加密，要求plainPadding和key的长度相同
		cypher[i] = plainPadding[i] ^ en.key[i%8]
	}
	return cypher
}

func (en *MyEncryption1) Decrypt(cipher []byte) ([]byte, error) {
	if len(cipher)%en.BlockSize() != 0 {
		return nil, errors.New("invalid ciphertext length")
	}
	if len(cipher) == 0 {
		return nil, errors.New("empty ciphertext")
	}

	plain := make([]byte, len(cipher))
	for i := 0; i < len(cipher); i++ {
		plain[i] = cipher[i] ^ en.key[i%8]
	}
	unpadded, err := PKCS7.Unpadding(plain, en.BlockSize())
	return unpadded, err
}
