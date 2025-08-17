package main

import (
	"encoding/json"
	"errors"
	"go/im/im_based_rabbitmq/common"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func InitLogger() {
	logFile, _ := os.OpenFile("log/im_server.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(logFile)
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 1 * time.Second,
	ReadBufferSize:   100,
	WriteBufferSize:  100,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 用户注册，不需要传参数，返回用户id
func RegistUser(w http.ResponseWriter, r *http.Request) {
	uid := time.Now().UnixMicro() //生成全局唯一的群ID（正规来讲，应该借助于mysql的自增id）
	err := GetRabbitMQ().RegistUser(uid, common.TypeUser)
	if err != nil {
		log.Printf("创建用户%d失败:%s", uid, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("创建用户失败"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatInt(uid, 10)))
		return
	}
}

// 把用户拉进群，群若不存在则先创建群
func JoinGroup(w http.ResponseWriter, r *http.Request) {
	var gid, uid int64
	var err error
	if gid, err = strconv.ParseInt(r.PathValue("gid"), 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("群id非法"))
		return
	}
	if uid, err = strconv.ParseInt(r.PathValue("uid"), 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("用户id非法"))
		return
	}
	if err = GetRabbitMQ().AddUser2Group(gid, uid); err != nil {
		log.Printf("用户%d入群%d失败：%s", uid, gid, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("用户入群失败"))
		return
	}
}

// 对非法内容进行拦截。比如机器人消息（发言频率过快）；包含欺诈、涉政等违规内容；涉嫌私下联系/交易等。
// 如果不需要拦截就返回nil。
func intercept(msg *common.Message) error {
	if len(msg.Content) == 0 {
		return errors.New("内容不能为空")
	}
	return nil
}

var (
	pongWait   = 5 * time.Second //等待pong的超时时间
	pingPeriod = 3 * time.Second //发送ping的周期，必须短于pongWait
)

func heartBeat(conn *websocket.Conn) {
	conn.SetPongHandler(func(appData string) error {
		// log.Println("receive pong")
		deadline := time.Now().Add(pongWait)
		conn.SetReadDeadline(deadline)
		// log.Printf("must read before %s", deadline.Format("2006-01-02 15:04:05"))
		return nil
	})

	err := conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		log.Printf("write ping error:%v\n", err)
		conn.WriteMessage(websocket.CloseMessage, nil)
	}

	ticker := time.NewTicker(pingPeriod)
LOOP:
	for {
		select { //通过select确保每次Ping的间隔是准确的2秒
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("write ping error:%v\n", err)
				conn.WriteMessage(websocket.CloseMessage, nil)
				break LOOP
			}
			// log.Println("write ping")
		}
	}
}

// 发言，通过websocket连接把发言内容传给web server，web server再传给MQ
func Speak(w http.ResponseWriter, r *http.Request) {
	//针对本次请求，创建一个单独的websocket连接
	conn, err := upgrader.Upgrade(w, r, nil) //http升级为websocket协议
	if err != nil {
		log.Printf("upgrade error: %v\n", err)
		return
	}
	//结束时关闭websocket连接
	defer func() {
		log.Println("close websocket connection")
		conn.Close()
	}()

	go heartBeat(conn) //心跳保持

	//接收用户的发言，写入RabbitMQ
	for {

		_, body, err := conn.ReadMessage() //如果对方主动断开连接或超时，该行会报错，for循环会退出
		if err != nil {
			log.Printf("read error:%v\n", err)
			break
		} else {
			var msg common.Message
			if err = json.Unmarshal(body, &msg); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("json参数解析失败"))
				return
			} else {
				msg.Content = strings.ReplaceAll(msg.Content, "\n", "  ") //换行符用空格代替。将来要用换行符来分隔每条Message，所以一条Message内部不能出现换行符
				msg.Id = time.Now().UnixMicro()                           //生成全局唯一的消息ID（正规来讲，应该使用SnowFlake等算法）
				msg.Time = time.Now().UnixMicro()                         //消息的生成时间以服务端接收到消息的时刻为准，防止客户端伪造时间
				if intercept(&msg) != nil {
					continue // 消息被拦截，不能转发给RabbitMQ
				}
				if err = GetRabbitMQ().Send(&msg, msg.To); err != nil {
					log.Printf("用户%s向%s发言失败：%s", msg.From, msg.To, err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("用户发言失败"))
					return
				}

				if strings.HasPrefix(msg.To, common.TypeUser) {
					// 在发言人自己这一端展示时不需要过拦截（intercept）检测
					if err = GetRabbitMQ().Send(&msg, msg.From); err != nil {
						log.Panicf("用户%s向%s发言失败：%s", msg.From, msg.To, err)
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("用户发言失败"))
						return
					}
				}
			}
		}
	}
}

func main() {
	InitLogger()
	mux := http.NewServeMux()
	// restful风格参数
	mux.HandleFunc("GET /regist_user", RegistUser)           //用户注册，不需要传参数，返回用户id
	mux.HandleFunc("GET /join_group/{gid}/{uid}", JoinGroup) //把用户拉进群，群若不存在则先创建群
	mux.HandleFunc("GET /speak", Speak)                      // 发言，通过websocket连接把发言内容传给web server，web server再传给MQ
	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
}

// go run ./im/im_based_rabbitmq/server
