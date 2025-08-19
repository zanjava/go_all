package database_test

import (
	"fmt"
	"go/frame/ioc/database"
	"testing"

	"github.com/golobby/container/v3"
)

func TestGetNews(t *testing.T) {
	var nd database.GormNews
	err := container.Fill(&nd) //让IoC容器提供一个bean(容器只负责给带container Tag的Field赋值)
	if err != nil {
		t.Error(err)
		return
	}
	news := nd.GetNews()
	if len(news) == 0 {
		t.Fail()
		return
	}
	for _, ele := range news {
		fmt.Println(ele.ViewPostTime, ele.UserId, ele.Title)
	}
}

// go test -v ./ioc/database -run=^TestGetNews$ -count=1
