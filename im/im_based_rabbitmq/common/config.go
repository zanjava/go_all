package common

const (
	TypeUser  = "u" //用户
	TypeGroup = "g" //群
)

const (
	//连接RabbitMQ使用的账号
	MqUser = "zgw"
	MqPass = "123456"
)

const (
	GroupMemberPath = "data/im/server/group/" //群成员信息保存到此目录的文件里
	ReceiveUserPath = "data/im/client/user/"  //从MQ拉取的消息保存到本地文件中
)
