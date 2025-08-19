package database_test

import (
	"fmt"
	database "go/frame/dependency_injection/database/gorm"
	"testing"
)

func TestGetNews(t *testing.T) {
	conn, _ := database.GetMysqlDB()
	nd := database.NewGormNews(conn)
	news := nd.GetNews()
	if len(news) == 0 {
		t.Fail()
	}
	for _, ele := range news {
		fmt.Println(ele.ViewPostTime, ele.UserId, ele.Title)
	}
}

// go test -v ./dependency_injection/database/gorm -run=^TestGetNews$ -count=1
