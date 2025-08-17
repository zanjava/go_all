# RocketMQ运维命令
```shell
cd cd rocketmq-all-5.3.1-bin-release/bin
./mqadmin.cmd updateTopic -n localhost:9876 -c DefaultCluster -t user_click # 创建(更新)Topic，类型默认为Normal
./mqadmin.cmd topicList -n localhost:9876  # 查看所有Topic
./mqadmin.cmd deleteTopic -n localhost:9876 -c DefaultCluster -t user_click  # 删除Topic
./mqadmin.cmd sendMessage  -n localhost:9876 -t user_click -p hello -k k1 -c t1  # 发送消息
./mqadmin.cmd printMsg -n localhost:9876 -t user_click  # 打印消息
./mqadmin.cmd consumeMessage -n localhost:9876 -t user_click  # 消费消息
./mqadmin.cmd clusterList -n localhost:9876   # 查看集群
./mqadmin.cmd updateSubGroup -n localhost:9876 -c DefaultCluster -g recommend_biz # 创建(更新)订阅关系
./mqadmin.cmd consumerProgress -n localhost:9876 -c DefaultCluster -t user_click  # 查看消费进度. Diff就是堆积的消息量
./mqadmin.cmd resetOffsetByTime -n localhost:9876 -g recommend_biz -t user_click -s -1   # 重置消费点位offset。s=-1表示CONSUME_FROM_LAST_OFFSET，s=-2表示CONSUME_FROM_FIRST_OFFSET，s=-3表示CONSUME_FROM_TIMESTAMP：第一次启动从指定时间点（即-s后面指定一个时间戳）位置消费，后续再启动接着上次消费的进度开始消费。不做指定的话默认半小时之前堆积的消息开始消费（启动一个全新的消费者，默认是3）
# but！！！go-client不支持重置offset
```  

创建不同类型的Topic
```shell
./mqadmin.cmd updateTopic -n localhost:9876 -c DefaultCluster -t user_click_delay -a +message.type=DELAY
./mqadmin.cmd updateTopic -n localhost:9876 -c DefaultCluster -t user_click_fifo -a +message.type=FIFO
./mqadmin.cmd updateTopic -n localhost:9876 -c DefaultCluster -t user_click_tx -a +message.type=TRANSACTION
```