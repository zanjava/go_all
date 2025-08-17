package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/im/im_based_rabbitmq/common"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"encoding/json"

	"github.com/desertbit/grumble"
	"github.com/nsf/termbox-go" // 在Linux shell里面表现比较好，在windows的power shell里对按键感知不灵敏
)

var (
	app *grumble.App
)

var (
	windowFrom string                   //当前登录者是谁
	windowTo   string                   //当前打开的会话窗口是在跟谁聊
	waterfall  = make(chan struct{}, 1) //是否要终止消息瀑布流的展示（发言之前需要终止消息瀑布流）
	ugReg      *regexp.Regexp           //用正则表达式判断msg.To是否合法
)

func init() {
	ugReg = regexp.MustCompile(`^[u|g]\d+$`)
}

func listenForKeyPress() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlA: //按下Ctrl+A键时
				waterfall <- struct{}{} //终止消息瀑布流
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func InitApp() {
	app = grumble.New(&grumble.Config{
		Name:        "dqq-im",
		Description: "基于RabbitMQ的聊天系统",
	})
	app.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		err := termbox.Init()
		if err != nil {
			fmt.Println("Failed to initialize termbox:", err)
			os.Exit(1)
		}
		go listenForKeyPress()
		return nil
	})
	app.AddCommand(&grumble.Command{
		Name: "regist",
		Help: "注册新用户",
		Run: func(c *grumble.Context) error {
			uid := RegistUser()
			if uid <= 0 {
				c.App.Println("注册失败")
			} else {
				c.App.Println("注册成功，用户id是", uid)
			}
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "join",
		Help: "加入群组",
		Flags: func(f *grumble.Flags) { // join -u 3 -g 1744695500829846
			f.Int64("u", "uid", 0, "用户id")
			f.Int64("g", "gid", 0, "群组id")
		},
		Run: func(c *grumble.Context) error {
			success := JoinGroup(c.Flags.Int64("gid"), c.Flags.Int64("uid"))
			if !success {
				c.App.Println("加入群组失败")
			} else {
				c.App.Println("加入群组成功")
			}
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "login",
		Help: "登录",
		Args: func(a *grumble.Args) { // login 1
			a.Int64("uid", "用户id")
		},
		Run: func(c *grumble.Context) error {
			uid := c.Args.Int64("uid")
			windowFrom = common.TypeUser + strconv.FormatInt(uid, 10) // 更新windowFrom
			GetRabbitMQ().Receive(uid)
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "send",
		Help: "发送消息",
		Flags: func(f *grumble.Flags) {
			f.String("t", "to", "", "消息发送给谁，发给个人请用u+uid，发给群请使用g+群id")
			f.String("m", "message", "", "消息内容")
		},
		Run: func(c *grumble.Context) error {
			// 流程检查
			if len(windowFrom) == 0 {
				c.App.Println("请先登录")
				return nil
			}
			// 参数合法性检查
			if len(c.Flags.String("to")) == 0 {
				c.App.Println("请指定to")
				c.App.Println()
				return nil
			}
			if !ugReg.MatchString(c.Flags.String("to")) {
				c.App.Println("请输入合理的用户或群组标识")
				return nil
			}
			if len(c.Flags.String("message")) == 0 {
				c.App.Println("请指定message")
				c.App.Println()
				return nil
			}

			//聊天对象发生了变化（即切换了聊天窗口）
			if windowTo != c.Flags.String("to") {
				log.Println("关闭老channel，打开新channel")
				close(windowChan)
				windowChan = make(chan []byte, 100)

				windowTo = c.Flags.String("to") // 更新windowTo
			}

			// 把消息发给web server
			msg := &common.Message{
				From:    windowFrom,
				To:      windowTo,
				Content: c.Flags.String("message"),
			}
			err := Send(msg)
			if err != nil {
				c.App.Println("消息发送失败")
				c.App.Println()
			}

			// 展示消息瀑布流
			dup := make(Set[int64], 1000) //按消息ID排重
			// 1. 先从文件里读取历史消息，展示给用户
			var userPath = common.ReceiveUserPath + windowFrom
			file, err := os.OpenFile(userPath+"/"+windowTo, os.O_CREATE|os.O_APPEND|os.O_RDONLY, os.ModePerm)
			if err != nil {
				log.Printf("open file %s failed: %s", userPath+"/"+windowTo, err)
			} else {
				msgs := make([]*common.Message, 0, 1000)
				reader := bufio.NewReader(file)
				for {
					line, _, err := reader.ReadLine()
					if err != nil {
						if err != io.EOF {
							log.Printf("读文件%s失败:%s", userPath+"/"+windowTo, err)
						}
						putLine2Msgs(line, &msgs, dup)
						break
					} else {
						putLine2Msgs(line, &msgs, dup)
					}
				}
				file.Close()
				sort.Slice(msgs, func(i, j int) bool {
					return msgs[i].Time < msgs[j].Time //按时间顺序排序
				})
				if len(msgs) > 10 { //只截取最后的10条进行展示
					msgs = msgs[len(msgs)-10:]
				}
				// 把文件里的内容展示给用户看
				for _, msg := range msgs {
					printMsg(c, msg)
				}
			}
			// 2. 再把channel里的内容展示给用户看，直到用户按下Ctrl+A键时退出
			close(waterfall)                   // 废弃老的waterfall
			waterfall = make(chan struct{}, 1) // 启用一个新的waterfall
		LOOP:
			for {
				select {
				case ele := <-windowChan: // 这部分消息不能保证是按时间排序的
					var msg common.Message
					err = json.Unmarshal(ele, &msg)
					if err != nil {
						log.Printf("Unmarshal line failed: %s", err)
					} else {
						if _, exists := dup[msg.Id]; !exists { //按消息ID排重
							dup[msg.Id] = struct{}{}
							printMsg(c, &msg)
						}
					}
				case <-waterfall:
					break LOOP
				}
			}
			return nil // 退出瀑布流
		},
	})
}

func putLine2Msgs(line []byte, msgs *[]*common.Message, dup Set[int64]) {
	if len(line) > 0 {
		var msg common.Message
		err := json.Unmarshal(bytes.Trim(line, "\n"), &msg)
		if err != nil {
			log.Printf("Unmarshal line failed: %s", err)
		} else {
			if _, exists := dup[msg.Id]; !exists { //按消息ID排重
				dup[msg.Id] = struct{}{}
				*msgs = append(*msgs, &msg)
			}
		}
	}
}

func printMsg(c *grumble.Context, msg *common.Message) {
	show := fmt.Sprintf("%s  %s: %s", time.Unix(msg.Time/1e6, 1e3*msg.Time%1e6).Format("2006-01-02 15:04:05.000000"), msg.From, msg.Content)
	c.App.Println(show)
}

func CloseShell() {
	if app != nil {
		app.Close()
	}
}
