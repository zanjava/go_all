package rabbitmq

const (
	//guest账号只能用于连接localhost
	// user = "guest"
	// pass = "guest"
	//可以通过网页后台创建一个admin账号
	User = "zgw"
	Pass = "123456"
	Host = "localhost"
	Port = "5672"
)

var (
	ExchangeName1 = "excg1"
	ExchangeName2 = "excg2"
	ExchangeName3 = "excg3"
)
