package distributed_test

import (
	"context"
	distributed "go/frame/database/redis"
	"testing"
)

func TestPubSub(t *testing.T) {
	distributed.PubSub(context.Background(), client)
}

// go test -v ./database/redis -run=^TestPubSub$ -count=1
