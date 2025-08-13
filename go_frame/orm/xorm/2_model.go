package xorm

import (
	"time"
)

// 默认情况下，xorm的命名规则使用SnakeMapper，即不论对于结构体还是成员变量，驼峰转为蛇形就是对应的表名和列名。
// 建表、建索引、导数据等操作不建议后端开发人员通过go代码完成，这种重要操作应该由DBA完成。即Migrate功能不建议使用
type User struct {
	Id         int `xorm:"pk autoincr"` //显式标记为主键。名为Id的int64成员默认就是主键，且是自增的。如果没有xorm tag，则Insert时会给User的Id赋值；如果有xorm tag,则tag里必须包含autoincr, Insert时才会给Id赋值
	UserId     int `xorm:"uid"`         //显式指定列名。如果列名与其他关键字冲突，用单引号括起来
	Degree     string
	Keywords   []string  `xorm:"json"`    //对于array, slice, map和内嵌结构体，xorm会自动执行json序列化/反序列化（不需要标注`xorm:"json"`），可对应到DB里的char、varchar、text或Blob。[]byte（type byte = uint8）除外，[]byte会对应到DB里的Blob类型。标注上`xorm:"json"`的好处在于如果DB里该列内容为空，则查询时json反序列化不会报错，否则slice类型要求在DB里至少得有一个[]，map和struct要求在DB里至少得有一个{}
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	UpdateTime time.Time `xorm:"updated"` //这个Field将在Insert或Update时自动赋值为当前时间
	DeleteTime time.Time `xorm:"deleted"` //软删除
	Gender     string
	City       string
	Version    int    `xorm:"version"`
	Province   string `xorm:"-"` //该字段不进行映射
}

// 显式指定表名
func (User) TableName() string {
	return "xorm_user"
}

// 显式指定表名或列名  的优先级要高于默认的Mapper名称映射方式
