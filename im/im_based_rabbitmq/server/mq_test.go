package main_test

import (
	"go/im/im_based_rabbitmq/common"
	main "go/im/im_based_rabbitmq/server"
	"testing"
	"time"
)

func TestRegistUser(t *testing.T) {
	defer main.GetRabbitMQ().Release()
	var uid int64 = 1
	if err := main.GetRabbitMQ().RegistUser(uid, common.TypeUser); err != nil {
		t.Error(err)
	}

	uid = 2
	if err := main.GetRabbitMQ().RegistUser(uid, common.TypeUser); err != nil {
		t.Error(err)
	}

	uid = 3
	if err := main.GetRabbitMQ().RegistUser(uid, common.TypeUser); err != nil {
		t.Error(err)
	}
}

func TestAddUser2Group(t *testing.T) {
	defer main.GetRabbitMQ().Release()
	if err := main.GetRabbitMQ().AddUser2Group(1, 1, 2, 3); err != nil {
		t.Error(err)
	}
}

func TestSend(t *testing.T) {
	defer main.GetRabbitMQ().Release()
	msg := common.Message{From: "u3", To: "u1", Time: time.Now().UnixMicro(), Content: "hello, 最近忙吗？"}
	if err := main.GetRabbitMQ().Send(&msg, msg.To); err != nil {
		t.Error(err)
	}

	msg = common.Message{From: "u2", To: "g1", Time: time.Now().UnixMicro(), Content: "大家好，我是2号用户，请多关照！"}
	if err := main.GetRabbitMQ().Send(&msg, msg.To); err != nil {
		t.Error(err)
	}
}

// go test -v ./im/im_based_rabbitmq/server -run=^TestRegistUser$ -count=1
// go test -v ./im/im_based_rabbitmq/server -run=^TestAddUser2Group$ -count=1
// go test -v ./im/im_based_rabbitmq/server -run=^TestSend$ -count=1
