package distributed_test

import (
	"context"
	distributed "go/frame/database/redis"
	"log/slog"
	"testing"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,  //redis默认会创建0-15号DB，这里使用默认的DB
		Username: "", //此处不需要用户名
		Password: "", //没有密码
	})
	//能ping成功才说明连接成功
	if err := client.Ping(context.Background()).Err(); err != nil {
		slog.Error("connect to redis failed", "error", err)
	} else {
		slog.Info("connect to redis")
	}
}

func TestStringValue(t *testing.T) {
	distributed.StringValue(context.Background(), client)
}

func TestStructValue(t *testing.T) {
	stu := &distributed.Student{Id: 1, Name: "zgw"}
	distributed.WriteStudent2Redis(client, stu)
	stu2 := distributed.GetStudentFromRedis(client, 1)
	if stu2.Id != stu.Id {
		t.Fail()
	}
	if stu2.Name != stu.Name {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	distributed.DeleteKey(context.Background(), client)
}

func TestScan(t *testing.T) {
	distributed.Scan(context.Background(), client)
}

func TestListValue(t *testing.T) {
	distributed.ListValue(context.Background(), client)
}

func TestSetgValue(t *testing.T) {
	distributed.SetValue(context.Background(), client)
}

func TestZSetValue(t *testing.T) {
	distributed.ZsetValue(context.Background(), client)
}

func TestHashTableValue(t *testing.T) {
	distributed.HashtableValue(context.Background(), client)
}

// go test -v ./database/redis -run=^TestStringValue$ -count=1
// go test -v ./database/redis -run=^TestStructValue$ -count=1
// go test -v ./database/redis -run=^TestDelete$ -count=1
// go test -v ./database/redis -run=^TestListValue$ -count=1
// go test -v ./database/redis -run=^TestSetgValue$ -count=1
// go test -v ./database/redis -run=^TestZSetValue$ -count=1
// go test -v ./database/redis -run=^TestHashTableValue$ -count=1
// go test -v ./database/redis -run=^TestScan$ -count=1
