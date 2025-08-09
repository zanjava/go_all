package client_test

import (
	"go/base/socket/client"
	"testing"
)

func TestTcpClient(t *testing.T) {
	client.TcpClient()
}

func TestUdpClient(t *testing.T) {
	client.UdpClient()
}

func TestTcpLongConnection(t *testing.T) {
	client.TcpLongConnection()
}

func TestUdpLongConnection(t *testing.T) {
	client.UdpLongConnection()
}

func TestTcpStick(t *testing.T) {
	client.TcpStick()
}

func TestUdpConnectionCurrent(t *testing.T) {
	client.UdpConnectionCurrent()
}

func TestUdpRpcClient(t *testing.T) {
	client.UdpRpcClient()
}

// go test -v ./socket/client -run=^TestTcpClient$ -count=1
// go test -v ./socket/client -run=^TestUdpClient$ -count=1
// go test -v ./socket/client -run=^TestTcpLongConnection$ -count=1
// go test -v ./socket/client -run=^TestUdpLongConnection$ -count=1
// go test -v ./socket/client -run=^TestTcpStick$ -count=1
// go test -v ./socket/client -run=^TestUdpConnectionCurrent$ -count=1
// go test -v ./socket/client -run=^TestUdpRpcClient$ -count=1
