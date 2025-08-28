package geth

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

func ReadBlock() {
	client := CreateClient()

	//取得最后一个区块的header信息
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("最后一个区块的编号", header.Number.String())
	fmt.Println("最后一个区块的Hash", header.Hash().String())
	fmt.Println()

	// header, err = client.HeaderByNumber(context.Background(), header.Number) //根据区块编号获得区块信息
	block, err := client.BlockByNumber(context.Background(), header.Number) //根据区块编号获得区块信息
	// block, err = client.BlockByHash(context.Background(), header.Hash())    //根据区块哈希值获得区块信息
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块编号", block.Number().Uint64()) //block基本把header的信息都包含了
	fmt.Println("打包区块的时间戳(精确到秒)", block.Time())  // 矿工有自由发挥的空间，时间戳不一定是准确的，只要链上各区块的时间戳是递增的即可
	fmt.Println("难度", block.Difficulty().Uint64())
	fmt.Println("哈希值", block.Hash().Hex()) //跟header.Hash()是一样的
	fmt.Println("交易数量", len(block.Transactions()))
	fmt.Println()

	//遍历区块里的所有交易
	for _, tx := range block.Transactions() {
		fmt.Println("链ID", tx.ChainId().Uint64())        //链ID是为了防止重放攻击，不同链的ID不一样
		fmt.Println("交易哈希值", tx.Hash().Hex())            //每条交易都有个哈希值，按照Merkle Tree的规则计算出整个区块的哈希值
		fmt.Println("交易金额(单位:wei)", tx.Value().String()) //1 ether = 10^18 wei
		fmt.Println("Gas的单价", tx.GasPrice().Uint64())
		fmt.Println("总支出", tx.Cost().Int64()) // 燃气费用 + value
		fmt.Println("交易Nonce值", tx.Nonce())
		fmt.Println("交易的输入", tx.Data()) //前4个字节是函数签名（也是函数ID），剩下的是函数参数
		if tx.To() != nil {
			fmt.Println("钱转给谁了", tx.To().Hex())
		}
		fmt.Println("ChainId", tx.ChainId().String())

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()) // receipt/rɪˈsiːt/:交易执行结果(发票，收据)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("交易是否执行成功", receipt.Status)
		fmt.Println("交易日志:")
		for i, log := range receipt.Logs {
			fmt.Println(i, log.Address.Hex(), common.Bytes2Hex(log.Data)) //indexed参数会被记录在log.Topics里，非indexed参数会被记录在log.Data里
		}
		break //只打印第一条交易
	}
	fmt.Println()

	//直接根据区块哈希值来获得区块里的交易数目
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块里的交易数目", count)

	//遍历区块里的所有交易
	for i := uint(0); i < count; i++ { //count:区块里的交易数
		//直接根据交易Hash查询交易信息
		// tx2, err := client.TransactionInBlock(context.Background(), block.Hash(), i)
		// tx2, isPending, err := client.TransactionByHash(context.Background(), tx.Hash())
		tx2, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash("0x0409614f627621c3c08c205b3d53756a12b816a5ca7d4a985b63740b1f1ae1ce"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("交易哈希值", tx2.Hash().Hex())
		fmt.Println("交易是否被挂起", isPending)
		break
	}
}
