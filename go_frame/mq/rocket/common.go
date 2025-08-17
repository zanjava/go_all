package common

const (
	NORMAL_TOPIC      = "user_click" // 主题名称，需要先创建好主题 ./mqadmin.cmd topicList -n localhost:9876
	DELAY_TOPIC       = "user_click_delay"
	FIFO_TOPIC        = "user_click_fifo"
	TRANSACTION_TOPIC = "user_click_tx"

	Endpoint = "localhost:8081" // proxy的ip和端口号
)
