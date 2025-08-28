package geth

import (
	"bytes"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// ECC(私钥)->公钥，哈希(公钥)->账户地址
func CreateAccount() (string, string) {
	privateKey, err := crypto.GenerateKey() //私钥，256bit。每次生成的私钥都是不一样的
	if err != nil {
		log.Fatal(err)
	}

	// ECDSA: 椭圆曲线数字签名算法
	privateKeyBytes := crypto.FromECDSA(privateKey) //把私钥转换成字节，方便打印
	fmt.Println("私钥", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.PublicKey                 // 公钥，520bit。公钥的第一个字节都是0x04
	publicKeyBytes := crypto.FromECDSAPub(&publicKey) //把公钥转换成字节，方便打印
	fmt.Println("公钥", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(publicKey).Hex() //账户地址，160bit，20B
	fmt.Println("账户地址", address)

	// 第一种Hash方式
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(publicKeyBytes[1:]) //去掉前面的第一个字节0x04
	digest := hasher.Sum(nil)
	fmt.Println("账户地址", hexutil.Encode(digest[12:])) //丢弃前12个字节，以太坊地址是20个字节
	// 第二种Hash方式
	digest = crypto.Keccak256Hash(publicKeyBytes[1:]).Bytes() //注意：Keccak256Hash()可以接收多个[]byte
	fmt.Println("账户地址", hexutil.Encode(digest[12:]))

	return hexutil.Encode(privateKeyBytes), address
}

func Keccak256WithPadding(numbers ...int64) []byte {
	array := make([][]byte, 0, len(numbers))
	for _, i := range numbers {
		bs := common.LeftPadBytes(big.NewInt(i).Bytes(), 32) //左侧补0，补够32个字节
		fmt.Println(hexutil.Encode(bs))
		array = append(array, bs)
	}
	hasher := crypto.Keccak256Hash(array...)
	return hasher.Bytes()
}

// 对数据进行签名
func SignData(text string) []byte {
	private := "0x9b23fdae2f1837b9e106a7ac1a814736c3d56a374690812c6c6057328b84f20d" //私钥
	privateKey, err := crypto.HexToECDSA(private[2:])                               //去掉0x前缀
	if err != nil {
		log.Fatalf("privateKey %s", err)
	}

	hash := crypto.Keccak256Hash([]byte(text)) // 原始text可能很大，直接执行非对称加密计算开销很大，所以先哈希
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("Sign %s", err)
	}
	return signature
}

func VerifySign(text string, signature []byte) bool {
	public := "0x0463a7b636a72e8dbc0ee786b9651cd886bad1d1be365bfc88f111e5f39e0125919134af36517c334ef7465f2d3cd5c5e67bf64cf4d0a9da98cb8014bab6f5c5c3" //跟私钥对应的公钥
	publicKeyBytes := common.Hex2Bytes(public[2:])

	hash := crypto.Keccak256Hash([]byte(text))
	// 注意这里是根据明文和签名，反算出公钥，然后跟传入的公钥比较
	publicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatalf("Ecrecover %s", err)
	}

	return bytes.Equal(publicKey, publicKeyBytes)
}
