package database_test

import (
	"fmt"
	database "go/post/database/gorm"
	"testing"
)

var (
	nid = 4
)

func TestPostNews(t *testing.T) {
	var err error
	nid, err = database.PostNews(7, "查尔斯·汤姆森", "作为自由独立的国家，它们完全有权宣战、缔和、结盟、通商和独立国家有权去做的一切行动")
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("新闻发布成功， news id=%d\n", nid)
	}
}

func TestUpdateNews(t *testing.T) {
	err := database.UpdateNews(nid, "独立", "当任何形式的政府对这些目标具破坏作用时，人民便有权利改变或废除它，以建立一个新的政府；其赖以奠基的原则，其组织权力的方式，务使人民认为唯有这样才最可能获得他们的安全和幸福。")
	if err != nil {
		t.Error(err)
	}
}

func TestGetNewsById(t *testing.T) {
	news := database.GetNewsById(nid)
	if news == nil {
		t.Fatalf("could not get news by id %d", nid)
	} else {
		fmt.Printf("%+v", *news)
	}

	tmpNid := -1
	news = database.GetNewsById(tmpNid)
	if news != nil {
		t.Fatalf("get news by id %d", tmpNid)
	}
}

func TestGetNewsByUid(t *testing.T) {
	news := database.GetNewsByUid(uid)
	if len(news) == 0 {
		t.Fatalf("could not get news by uid %d", uid)
	} else {
		for _, ele := range news {
			fmt.Printf("%+v\n", *ele)
		}
	}

	tmpUid := -1
	news = database.GetNewsByUid(tmpUid)
	if len(news) > 0 {
		t.Fatalf("get news by uid %d", tmpUid)
	}
}

func TestGetNewsByPage(t *testing.T) {
	pageNo, pageSize := 0, 2
	fmt.Println("--------------", pageNo, "--------------")
	total, news := database.GetNewsByPage(pageNo, pageSize)
	if total == 0 {
		t.Fatal("got zero news")
	} else {
		if len(news) > 0 {
			for _, ele := range news {
				fmt.Printf("%+v\n", *ele)
			}
		} else {
			fmt.Printf("got zero news, pageNo %d pageSize %d", pageNo, pageSize)
			return
		}
	}

	pageNo, pageSize = 1, 2
	fmt.Println("--------------", pageNo, "--------------")
	total, news = database.GetNewsByPage(pageNo, pageSize)
	if total == 0 {
		t.Fatal("got zero news")
	} else {
		if len(news) > 0 {
			for _, ele := range news {
				fmt.Printf("%+v\n", *ele)
			}
		} else {
			fmt.Printf("got zero news, pageNo %d pageSize %d", pageNo, pageSize)
			return
		}
	}

	pageNo, pageSize = 2, 2
	fmt.Println("--------------", pageNo, "--------------")
	total, news = database.GetNewsByPage(pageNo, pageSize)
	if total == 0 {
		t.Fatal("got zero news")
	} else {
		if len(news) > 0 {
			for _, ele := range news {
				fmt.Printf("%+v\n", *ele)
			}
		} else {
			fmt.Printf("got zero news, pageNo %d pageSize %d", pageNo, pageSize)
			return
		}
	}

}

func TestDeleteNews(t *testing.T) {
	err := database.DeleteNews(nid)
	if err != nil {
		t.Fatal(err)
	}

	news := database.GetNewsById(nid)
	if news != nil {
		t.Fail()
		return
	}

	err = database.DeleteNews(nid)
	if err == nil {
		t.Fatalf("新闻%d第二次删除成功！", nid)
	} else {
		fmt.Printf("新闻%d第二次删除失败：%s", nid, err)
	}
}

// go test -v ./post/database/gorm -run=^TestPostNews$ -count=1
// go test -v ./post/database/gorm -run=^TestUpdateNews$ -count=1
// go test -v ./post/database/gorm -run=^TestGetNewsById$ -count=1
// go test -v ./post/database/gorm -run=^TestGetNewsByUid$ -count=1
// go test -v ./post/database/gorm -run=^TestGetNewsByPage$ -count=1
// go test -v ./post/database/gorm -run=^TestDeleteNews$ -count=1
