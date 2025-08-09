package encryption

import (
	"errors"
	PKCS7 "go/base/basic"
)

// CBC分组加密算法
type MyEncryption2 struct {
	key       [8]byte
	BlockSize int
	BlockMode int
}

func NewMyEncryption2(key [8]byte, blockMode int) *MyEncryption2 {
	return &MyEncryption2{
		key:       key,
		BlockSize: 8,
		BlockMode: blockMode,
	}
}

func (en *MyEncryption2) Encrypt(plain []byte) []byte {
	plainPadding := PKCS7.Padding(plain, en.BlockSize) //明文末尾填充字节，长度成为BlockSize的整倍数
	cypher := make([]byte, len(plainPadding))
	prevCypher := make([]byte, en.BlockSize) //全0，任何数跟0异或还是自身
	for i := 0; i < len(plainPadding); i += en.BlockSize {
		begin := i
		end := i + en.BlockSize
		en.confuseBlock(plainPadding[begin:end], prevCypher)
		en.encryptBlock(plainPadding[begin:end], cypher[begin:end])
		copy(prevCypher, cypher[begin:end])
	}
	return cypher
}

// 分组混淆。plain和prevCypher的长度一定相等
func (en *MyEncryption2) confuseBlock(block, prevCypher []byte) {
	switch en.BlockMode {
	case CBC:
		for i := 0; i < len(block); i++ {
			block[i] ^= prevCypher[i]
		}
	default:
	}
}

// 加密一组数据。plain和cypher的长度一定都等于BlockSize
func (en *MyEncryption2) encryptBlock(plain, cypher []byte) {
	for i := 0; i < len(plain); i++ {
		cypher[i] = plain[i] ^ en.key[i%8]
	}
}

func (en *MyEncryption2) Decrypt(cypher []byte) ([]byte, error) {
	if len(cypher)%en.BlockSize != 0 {
		return nil, errors.New("密文长度不合法")
	}
	if len(cypher) == 0 {
		return []byte{}, nil
	}
	plainPadding := make([]byte, len(cypher))
	blockNum := len(cypher) / en.BlockSize //几组
	// 倒着解密，先解密最后一个分组
	for i := blockNum - 1; i >= 0; i-- {
		begin := i * en.BlockSize
		end := begin + en.BlockSize
		var prevCypher []byte //前一组密文
		if i == 0 {
			prevCypher = make([]byte, en.BlockSize) //全0
		} else {
			prevCypher = cypher[begin-en.BlockSize : end-en.BlockSize]
		}
		en.decryptBlock(plainPadding[begin:end], cypher[begin:end]) //解密得到明文，注意这还不是真正的明文
		en.deconfuseBlock(plainPadding[begin:end], prevCypher)      //跟上一组密文解混淆，得到的才是真正的明文
	}
	return PKCS7.Unpadding(plainPadding, en.BlockSize) //去除末尾填充的字节
}

func (en *MyEncryption2) deconfuseBlock(block, prevCypher []byte) {
	en.confuseBlock(block, prevCypher)
}

func (en *MyEncryption2) decryptBlock(plain, cypher []byte) {
	for i := 0; i < len(plain); i++ {
		plain[i] = cypher[i] ^ en.key[i%8]
	}
}
