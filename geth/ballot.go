package geth

import (
	"context"
	"fmt"
	"geth/bin/ballot"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 选举阶段的枚举，跟合约里的保持一致
const (
	Init = iota //0
	Regs        //1
	Vote        //2
	Done        //3
)

var (
	// 这些key需要跟你本地的Ganache里的账户保持一致
	ContractAddress    = "0x2eB4fe08912A1EeD7c2e65D47750E72AFe3438B0" // 合约部署成功后，会返回一个地址，日志里会显示是哪个账户部署了该合约，那个账户就是ChairPerson
	ChairPersonPrivate = "0x4b78d2d2b849cf7f26506174d8efec83f592f96dab2ce943839a554ada03018d"
	Voter1Private      = "0xe902a1560aabd57aa71b6dc838c7f875df78124cb6f9bb60f950cda2178e8f94"
	Voter2Private      = "0x6d5fae2b4cc6e50e17e18e6c81ffaa847f2b6b126c46a0d6b72ed3131238cae4"
	Voter3Private      = "0x7ddfcdb57308fecd61a9cf2e6ca3c48503ab139d5efac487ee5868192a40cfae"
)

// 只有chairPerson可以调用ChangeState()
func ChangeState(client *ethclient.Client, phase int) error {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := ballot.NewBallot(contractAddress, client) //根据合约的地址加载合约
	if err != nil {
		return fmt.Errorf("NewContracts %w", err)
	}

	opt, err := GenTransactOpt(client, ChairPersonPrivate, 0)
	if err != nil {
		return err
	}

	_, err = contract.ChangeState(opt, uint8(phase))
	if err != nil {
		return fmt.Errorf("ChangeState %w", err)
	} else {
		return nil
	}
}

// 任何人都可以调用GetState()
func GetState(client *ethclient.Client) (int, error) {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := ballot.NewBallot(contractAddress, client) //根据合约的地址加载合约
	if err != nil {
		return -1, err
	}

	state, err := contract.State(new(bind.CallOpts))
	return int(state), err
}

// 登记选民。只有chairPerson有权限执行
func Regist(client *ethclient.Client, privateKey string) error {
	private, err := crypto.HexToECDSA(privateKey[2:])
	if err != nil {
		return fmt.Errorf("privateKey %w", err)
	}
	address, err := GetAddressFromPrivateKey(private)
	if err != nil {
		return fmt.Errorf("address %w", err)
	}

	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := ballot.NewBallot(contractAddress, client) //根据合约的地址加载合约
	if err != nil {
		return fmt.Errorf("NewContract %w", err)
	}

	opt, err := GenTransactOpt(client, ChairPersonPrivate, 0)
	if err != nil {
		return err
	}

	_, err = contract.Register(opt, address)
	if err != nil {
		return fmt.Errorf("register %w", err)
	} else {
		log.Printf("选民%s注册成功\n", address)
		return nil
	}
}

// 选民voter给candidate投票
func VoteFor(client *ethclient.Client, voterPrivateKey string, candidate int) error {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := ballot.NewBallot(contractAddress, client) //根据合约的地址加载合约
	if err != nil {
		return fmt.Errorf("NewContract %w", err)
	}

	opt, err := GenTransactOpt(client, voterPrivateKey, 0)
	if err != nil {
		return err
	}

	_, err = contract.Vote(opt, big.NewInt(int64(candidate)))
	if err != nil {
		return fmt.Errorf("vote %w", err)
	} else {
		log.Printf("给候选人%d投票成功\n", candidate)
		return nil
	}
}

// 查看最终哪个候选人获胜了
func GetWinner(client *ethclient.Client) (int, error) {
	contractAddress := common.HexToAddress(ContractAddress)
	contract, err := ballot.NewBallot(contractAddress, client) //根据合约的地址加载合约
	if err != nil {
		return -1, fmt.Errorf("NewContract %w", err)
	}

	winner, err := contract.GetWinner(new(bind.CallOpts))
	if err != nil {
		return -1, fmt.Errorf("winner %w", err)
	} else {
		return int(winner.Int64()), nil
	}
}

type Proposal struct {
	Voter     string
	Candidate int
}

// 从区块链上取出每个voter的投票记录
func Audit(client *ethclient.Client) []*Proposal {
	contractAbi, err := abi.JSON(strings.NewReader(string(ballot.BallotABI)))
	if err != nil {
		log.Println(err)
		return nil
	}

	//取得最后一个区块的header信息
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println("最后一个区块的编号", header.Number.String())

	proposals := make([]*Proposal, 0, 100)
	// 倒查10个区块
	for i := 0; i < 10; i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(0).Sub(header.Number, big.NewInt(int64(i)))) //根据区块编号获得区块信息
		if err != nil {
			log.Println(err)
			continue
		}

		//遍历区块里的所有交易
		for _, tx := range block.Transactions() {
			if tx.To().Hex() != ContractAddress { // 只关心本合约上的Transact
				continue
			}
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()) // 交易执行结果
			if err != nil {
				log.Fatal(err)
			}
			if receipt.Status != 1 { // 只关心执行成功的Transact
				continue
			}

			inputBytes := tx.Data()
			methodSignature, data := inputBytes[:4], inputBytes[4:] // 前4个字节是函数签名，也是函数ID
			method, err := contractAbi.MethodById(methodSignature)
			if err != nil {
				log.Println(err)
				continue
			}
			if method.Name != "vote" { //只关心对vote函数的调用
				continue
			}

			var args = make(map[string]interface{}) //key是参数名，value是参数值
			err = method.Inputs.UnpackIntoMap(args, data)
			if err != nil {
				log.Println(err)
				continue
			}
			for key, value := range args {
				if key == "candidate" { // 只关心candidate这个参数
					candidate, ok := value.(*big.Int)
					if ok {
						sender, err := GetFrom(tx)
						if err == nil {
							proposals = append(proposals, &Proposal{Candidate: int(candidate.Int64()), Voter: sender})
						}
					}
				}
				fmt.Printf("%s: %+v\n", key, value)
			}
		}
	}
	return proposals
}
