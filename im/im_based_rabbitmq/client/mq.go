package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/im/im_based_rabbitmq/common"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	mqConn    *amqp.Connection
	mqChannel *amqp.Channel //支持并发使用
}

type Set[T comparable] map[T]struct{}

var (
	mq     *RabbitMQ
	mqOnce sync.Once

	windowChan = make(chan []byte, 100) //从MQ拉下来的数据存到这个channel里

	groupBelong map[string]Set[string] //key是user, value是group set，即记录每个用户在哪些群里
)

// 填充groupBelong
func InitGroupBelong() {
	groupBelong = make(map[string]Set[string], 10)
	//filepath.Walk 会递归地遍历一个目录
	filepath.Walk(common.GroupMemberPath, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Mode().IsRegular() {
			gid := info.Name()
			file, _ := os.Open(subPath)
			defer file.Close()
			reader := bufio.NewReader(file)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					if len(line) > 0 {
						uid := string(bytes.Trim(line, "\n"))
						if len(uid) > 0 {
							if mp, exists := groupBelong[uid]; exists {
								mp[gid] = struct{}{}
							} else {
								groupBelong[uid] = make(Set[string], 10)
								groupBelong[uid][gid] = struct{}{}
							}
						}
					}
					break
				} else {
					if len(line) > 0 {
						uid := string(bytes.Trim(line, "\n"))
						if len(uid) > 0 {
							if mp, exists := groupBelong[uid]; exists {
								mp[gid] = struct{}{}
							} else {
								groupBelong[uid] = make(map[string]struct{}, 10)
								groupBelong[uid][gid] = struct{}{}
							}
						}
					}
				}
			}
		}
		return nil
	})
}

func GetRabbitMQ() *RabbitMQ {
	mqOnce.Do(func() {
		//连接RabbitMQ
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", common.MqUser, common.MqPass))
		if err != nil {
			log.Panicf("connect to RabbitMQ failed: %s", err)
		}
		//基于连接创建Channel
		ch, err := conn.Channel()
		if err != nil {
			log.Panicf("open channel failed: %s", err)
		}
		mq = &RabbitMQ{
			mqConn:    conn,
			mqChannel: ch,
		}
	})
	return mq
}

// 把MQ的指定队列里拉取消息，存入本地文件和channel(仅属于当前会话窗口的存入channel)
// map是引用类型，但sync.Map不是，sync.Map就是个普通的结构体，要在函数里面修改这个结构体需要传指针
func (mq *RabbitMQ) readQueue(queue, qtype string, path string, fileMap *sync.Map) {
	deliveryCh, err := mq.mqChannel.Consume(
		queue+"_computer", //queue, PC端聊天
		"",                //consumer
		false,             //auto-ack
		false,             //exclusive
		false,             //no-local
		false,             //no-wait
		nil,               //args
	)
	if err != nil {
		log.Printf("regist consumer failed: %s", err)
	} else {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			log.Printf("开始消费队列%s", queue+"_computer")
			for delivery := range deliveryCh { //读取MQ里的消息
				var msg common.Message
				err := json.Unmarshal(delivery.Body, &msg)
				if err != nil {
					log.Printf("Unmarshal message failed: %s", err)
					continue
				}
				// 把消息写入本地文件
				var file *os.File
				session := msg.From
				if qtype == common.TypeGroup {
					session = msg.To
				}
				v, exists := fileMap.Load(session) // 一个from对应一个文件。TODO 为防止单个文件无阻增大，需要按照日期分割文件
				if !exists {
					fout, err := os.OpenFile(path+"/"+session, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
					if err == nil {
						file = fout
						log.Printf("准备写入文件%s", path+"/"+session)
						fileMap.Store(session, fout)
					}
				} else {
					file = v.(*os.File)
				}
				file.Write(delivery.Body) //把消息写入文件
				file.WriteString("\n")
				delivery.Ack(false) // 成功写入本地文件后，再ACK，标记消息消费成功

				// 把仅属于当前会话窗口的消息写入channel
				log.Printf("from %s to %s windowSession %s", msg.From, msg.To, windowTo)
				if queue == msg.From || // 用户自己发出去的消息
					(strings.HasPrefix(windowTo, common.TypeUser) && windowTo == msg.From) || // 对方发给当前用户的
					(strings.HasPrefix(windowTo, common.TypeGroup) && windowTo == msg.To) { // 发到这个群里面的
					log.Printf("消息id %d 写入channel", msg.Id)
					windowChan <- delivery.Body
				}
			}
		}()
	}
}

// 从MQ里接收发送给uid的消息，写入本地文件和channel(仅属于当前会话窗口的存入channel)
func (mq *RabbitMQ) Receive(uid int64) {
	fileMap := new(sync.Map)
	queue := common.TypeUser + strconv.FormatInt(uid, 10) //一对一发给uid的消息全在这个队列里
	var userPath = common.ReceiveUserPath + queue
	os.MkdirAll(userPath, os.ModePerm)
	mq.readQueue(queue, common.TypeUser, userPath, fileMap)

	// uid在这些群里面，群里有消息时，也需要写入本地文件
	if gids, exists := groupBelong[strconv.FormatInt(uid, 10)]; exists {
		log.Printf("用户%d在这些群里%v", uid, gids)
		for gid := range gids {
			queue := common.TypeGroup + gid + "_" + common.TypeUser + strconv.FormatInt(uid, 10) //发到群gid的消息全在这个队列里，当然这个队列有很多副本，每个群成员都拥有一个副本
			mq.readQueue(queue, common.TypeGroup, userPath, fileMap)
		}
	}
}

// 释放MQ连接
func (mq *RabbitMQ) Release() {
	if mq.mqChannel != nil {
		mq.mqChannel.Close()
	}
	if mq.mqConn != nil {
		mq.mqConn.Close()
	}
}
