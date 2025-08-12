package distributed

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	redis "github.com/redis/go-redis/v9"
)

// value是简单的string
func StringValue(ctx context.Context, client *redis.Client) {
	key := "name"
	value := "大乔乔"
	defer client.Del(ctx, key) //函数结束时删除redis上的key，不影响下次运行演示

	err := client.Set(ctx, key, value, 1*time.Second).Err() //1秒后失效。0表示永不失效
	checkError(err)

	client.Expire(ctx, key, 3*time.Second) //通过Expire设置3秒后失效。该方法对任意类型的redis value都适用
	time.Sleep(2 * time.Second)

	v2, err := client.Get(ctx, key).Result()
	checkError(err)
	fmt.Println(v2)

	err = client.Set(ctx, "age", 18, 1*time.Second).Err() //int写入redis后会转成string
	checkError(err)
	v3, err := client.Get(ctx, "age").Int()
	checkError(err)
	fmt.Printf("age=%d\n", v3)
}

type Student struct {
	Id   int
	Name string
}

func WriteStudent2Redis(client *redis.Client, stu *Student) error {
	if stu == nil {
		return nil
	}
	key := "STU_" + strconv.Itoa(stu.Id) //避免各种id冲突，用前缀区分
	v, err := sonic.Marshal(stu)
	if err != nil {
		return err
	}
	err = client.Set(context.Background(), key, string(v), 5*time.Minute).Err()
	return err
}

func GetStudentFromRedis(client *redis.Client, sid int) *Student {
	key := "STU_" + strconv.Itoa(sid) //避免各种id冲突，用前缀区分
	v, err := client.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil { //如果key不存在，会返回redis.Nil
			log.Println(err)
		}
		return nil
	}
	var stu Student
	err = sonic.Unmarshal([]byte(v), &stu)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &stu
}

func DeleteKey(ctx context.Context, client *redis.Client) {
	n, err := client.Del(ctx, "not_exissts").Result()
	if err == nil {
		fmt.Printf("删除%d个key\n", n)
	}
}

// value是List
func ListValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []interface{}{1, "中", 3, 4, 3, 1}    //各种数据类型混合
	err := client.RPush(ctx, key, values...).Err() //RPush向List右侧插入，LPush向List左侧插入。如果List不存在会先创建
	checkError(err)

	v2, err := client.LRange(ctx, key, 0, -1).Result() //截取，双闭区间。LRange表示List Range，即遍历List。0表示第一个，-1表示倒数第一个。v2是个[]string，即1,3,4存到redis里实际上是string
	checkError(err)
	fmt.Println(v2)
}

// value是Set
func SetValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []interface{}{1, "中", 3, 4, 3, 1}   //1,3,4存到redis里实际上是string
	err := client.SAdd(ctx, key, values...).Err() //SAdd向Set中添加元素,set里不允许出现重复元素
	checkError(err)

	//判断Set中是否包含指定元素
	var value any
	value = 1 //数字1会转成string再去redis里查找
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}
	value = "1"
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}
	value = 2
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}

	//遍历Set
	for _, ele := range client.SMembers(ctx, key).Val() {
		fmt.Println(ele)
	}

	key2 := "ids2"
	defer client.Del(ctx, key2)
	values = []interface{}{1, "中", "大", "乔"}
	err = client.SAdd(ctx, key2, values...).Err() //SAdd向Set中添加元素
	checkError(err)

	//差集
	fmt.Println("key - key2 差集")
	for _, ele := range client.SDiff(ctx, key, key2).Val() {
		fmt.Println(ele)
	}
	fmt.Println("key2 - key 差集")
	for _, ele := range client.SDiff(ctx, key2, key).Val() {
		fmt.Println(ele)
	}

	//交集
	fmt.Println("key & key2 交集")
	for _, ele := range client.SInter(ctx, key, key2).Val() {
		fmt.Println(ele)
	}
}

// value是ZSet(有序的Set)
func ZsetValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []redis.Z{{Member: "张三", Score: 70.0}, {Member: "李四", Score: 100.0}, {Member: "王五", Score: 80.0}} //Score是用来排序的,比如把时间戳赋给score
	err := client.ZAdd(ctx, key, values...).Err()
	checkError(err)

	//遍历ZSet，按Score有序输出Member
	for _, ele := range client.ZRange(ctx, key, 0, -1).Val() {
		fmt.Println(ele)
	}
}

// value是哈希表(即map)
func HashtableValue(ctx context.Context, client *redis.Client) {
	student1 := map[string]interface{}{"Name": "张三", "Age": 18, "Height": 173.5}
	err := client.HMSet(ctx, "学生1", student1).Err() //前缀H表示HashTable。redis-server4.0之后的版本可以直接使用HSet
	checkError(err)
	student2 := map[string]interface{}{"Name": "李四", "Age": 20, "Height": 180.0}
	err = client.HMSet(ctx, "学生2", student2).Err()
	checkError(err)

	age, err := client.HGet(ctx, "学生2", "Age").Int() //指定redis的key以及map里的key
	checkError(err)
	fmt.Printf("age=%d\n", age)

	for field, value := range client.HGetAll(ctx, "学生1").Val() { //GetAll表示获取完整的map
		fmt.Printf("field:%s  value:%s\n", field, value)
	}

	client.Del(ctx, "学生1")
	client.Del(ctx, "学生2")
}

func checkError(err error) {
	if err != nil {
		if err == redis.Nil { //读redis发生error，大部分情况是因为key不存在
			fmt.Println("key不存在")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// 遍历Key。直接用keys命令是全库遍历,redis是单线程的，会阻塞很长时间。
//
// SCAN cursor [MATCH pattern] [COUNT count]。scan命令注意事项：
//
// 1、当游标返回0时，表示迭代结束。第一次 Scan 时指定游标为 0，表示开启新的一轮迭代。cursor是HashTable槽位里的值，并不是递增的。
//
// 2、count表示一次遍历多少个key，这些key可能全部不能匹配patten。count设为10000比较合适，count越大总耗时越短，但是单次查询阻塞的时间越长。
//
// 3、返回的结果可能会有重复，需要客户端去重复，这点非常重要。
//
// 4、遍历的过程中如果有数据修改，改动后的数据能不能遍历到是不确定的。
//
// 5、单次返回的结果是空的并不意味着遍历结束，而要看返回的游标值是否为0。
func Scan(ctx context.Context, client *redis.Client) {
	if client == nil {
		log.Printf("connect redis failed")
		os.Exit(1)
	}
	const (
		MID = "_dqq_"
	)
	for i := 0; i < 10; i++ {
		//构造10个key,都匹配模式*_dqq_*
		key := strconv.Itoa(i) + MID + strconv.Itoa(i)
		err := client.Set(ctx, key, "1", 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
	// 删除创建的那10个key
	// defer func() {
	// 	for i := 0; i < 10; i++ {
	// 		key := strconv.Itoa(i) + MID + strconv.Itoa(i)
	// 		client.Del(ctx, key)
	// 	}
	// }()
	const COUNT = 100 //遍历的批次大小，建议设为10000
	var cursor uint64 = 0
	dup := make(map[string]struct{}, 10) //对遍历出来的key排重
	for {
		// 取出所有match上patten的key。如果match参数设为空，则是遍历库里的所有key
		keys, c, err := client.Scan(ctx, cursor, "*"+MID+"*", COUNT).Result()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("cursor %d keys count %d\n", c, len(keys))
		for _, key := range keys {
			dup[key] = struct{}{}
		}
		if c == 0 {
			break
		}
		cursor = c //本次scan返回的cursor，作为下一次scan使用的cursor
	}
	fmt.Println("total", len(dup))
	for key := range dup {
		fmt.Println(key)
	}
}

// redis 应用场景
// 1、分布式锁
// 2、缓存
