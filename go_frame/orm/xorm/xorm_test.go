package xorm_test

import (
	"go/frame/orm/xorm"
	"log/slog"
	"os"
	"sync"
	"testing"
)

var (
	engine = xorm.CreateEngine("localhost", "test", "tester", "123456", 3306)
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
				},
			),
		),
	)
}

func TestXormQuickStart(t *testing.T) {
	xorm.XormQuickStart()
}

func TestCreate(t *testing.T) {
	xorm.Create(engine)
}

func TestDelete(t *testing.T) {
	xorm.Delete(engine)
}

func TestUpdate(t *testing.T) {
	xorm.Update(engine)
}

func TestUpdateByVersion(t *testing.T) {
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			xorm.UpdateByVersion(engine)
		}()
	}
	wg.Wait()
}

func TestRead(t *testing.T) {
	xorm.Read(engine)
}

func TestReadWithStatistics(t *testing.T) {
	xorm.ReadWithStatistics(engine)
}

func TestTransaction(t *testing.T) {
	xorm.Transaction(engine)
}

func TestRawSql(t *testing.T) {
	xorm.RawSelect(engine)
	xorm.RawExec(engine)
}

func TestHandleError(t *testing.T) {
	xorm.HandleError(engine)
}

// go test -v ./orm/xorm -run=^TestXormQuickStart$ -count=1
// go test -v ./orm/xorm -run=^TestCreate$ -count=1
// go test -v ./orm/xorm -run=^TestDelete$ -count=1
// go test -v ./orm/xorm -run=^TestUpdate$ -count=1
// go test -v ./orm/xorm -run=^TestUpdateByVersion$ -count=1
// go test -v ./orm/xorm -run=^TestRead$ -count=1
// go test -v ./orm/xorm -run=^TestReadWithStatistics$ -count=1
// go test -v ./orm/xorm -run=^TestTransaction$ -count=1
// go test -v ./orm/xorm -run=^TestRawSql$ -count=1
// go test -v ./orm/xorm -run=^TestHandleError$ -count=1
