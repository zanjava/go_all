package geth_test

import (
	"fmt"
	geth "geth"
	"log"
	"testing"
)

var (
	client = geth.CreateClient()
)

func TestRegist(t *testing.T) {
	for _, address := range []string{geth.Voter1Private, geth.Voter2Private, geth.Voter3Private} {
		err := geth.Regist(client, address) //先确保现在是登记阶段
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestState(t *testing.T) {
	x := geth.Regs
	err := geth.ChangeState(client, x)
	if err != nil {
		t.Fatal(err)
	}
	state, err := geth.GetState(client)
	if err != nil {
		t.Fatal(err)
	}
	if state != x {
		t.Fatalf("expected %d but %d\n", x, state)
	}

	// err = geth.ChangeState(client, geth.Regs) // 设一个更小的state，肯定会失败
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	t.Fatal()
	// }
}

func TestVote(t *testing.T) {
	// 先确保进入到投票阶段
	err := geth.ChangeState(client, geth.Vote)
	if err != nil {
		t.Fatal(err)
	}

	err = geth.VoteFor(client, geth.ChairPersonPrivate, 1) //ChairPerson的一票顶两票，这是合约里规定的
	if err != nil {
		t.Fatal(err)
	}
	err = geth.VoteFor(client, geth.Voter1Private, 1)
	if err != nil {
		t.Fatal(err)
	}
	err = geth.VoteFor(client, geth.Voter2Private, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = geth.VoteFor(client, geth.Voter3Private, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = geth.VoteFor(client, geth.Voter3Private, 3)
	if err == nil {
		log.Print("Voter3第二次投票居然成功了！")
		t.Fail()
	}
}

func TestWinner(t *testing.T) {
	// 先确保投票已经结束
	err := geth.ChangeState(client, geth.Done)
	if err != nil {
		t.Fatal(err)
	}

	winner, err := geth.GetWinner(client)
	if err != nil {
		t.Fatal(err)
	} else {
		log.Printf("获胜人 %d\n", winner)
	}
}

func TestAudit(t *testing.T) {
	proposals := geth.Audit(client)
	for _, proposal := range proposals {
		fmt.Printf("voter %s proposal %d\n", proposal.Voter, proposal.Candidate)
	}
}

// go test -v -run=TestRegist -count=1
// go test -v -run=TestState -count=1
// go test -v -run=TestVote -count=1
// go test -v -run=TestWinner -count=1
// go test -v -run=TestAudit -count=1
