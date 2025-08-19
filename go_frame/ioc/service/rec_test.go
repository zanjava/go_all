package service_test

import (
	"fmt"
	"go/frame/ioc/service"
	"strings"
	"testing"

	"github.com/golobby/container/v3"
)

func TestRec(t *testing.T) {
	var recommender service.INewsRecommender
	err := container.NamedResolve(&recommender, "sort_poster") //让IoC容器提供一个接口的具体实现
	if err != nil {
		t.Error(err)
		return
	} else {
		for _, news := range recommender.RecNews() {
			fmt.Printf("%s %d %s\n", news.ViewPostTime, news.UserId, news.Title)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	err = container.NamedResolve(&recommender, "sort_poster") //让IoC容器提供一个接口的具体实现
	if err != nil {
		t.Error(err)
		return
	} else {
		for _, news := range recommender.RecNews() {
			fmt.Printf("%s %d %s\n", news.ViewPostTime, news.UserId, news.Title)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	err = container.NamedResolve(&recommender, "sort_time")
	if err != nil {
		t.Error(err)
		return
	} else {
		for _, news := range recommender.RecNews() {
			fmt.Printf("%s %d %s\n", news.ViewPostTime, news.UserId, news.Title)
		}
	}

	err = container.NamedResolve(&recommender, "sort_time")
	if err != nil {
		t.Error(err)
		return
	} else {
		for _, news := range recommender.RecNews() {
			fmt.Printf("%s %d %s\n", news.ViewPostTime, news.UserId, news.Title)
		}
	}
}

// go test -v ./ioc/service -run=^TestRec$ -count=1
