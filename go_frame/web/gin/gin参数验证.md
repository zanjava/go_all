# GIN参数验证
详细文档参见 https://pkg.go.dev/github.com/go-playground/validator

```Go
type User struct {
    Name string `form:"name" binding:"required"` //required:必须上传name参数
    Score int    `form:"score" binding:"gt=0"`    //score必须为正数    
    Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`       //自定义验证before_today，日期格式东8区
    Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"` //毕业时间要晚于入学时间
}
```  
范围约束  
- 对于数值，约束其取值。min, max, eq, ne, gt, gte, lt, lte, oneof=6 8  
- 对于字符串、切片、数组和map，约束其长度。len=10, min=6, max=10, gt=10   

跨字段约束  
- 跨字段就在范围约束的基础上加field后缀，如 gtfield=Enrollment
- 如果还跨结构体(cross struct)就在跨字段的基础上在field前面加cs，如 
```Go
type Inner struct {
	StartDate time.Time
}

type Outer struct {
	InnerStructField *Inner
	CreatedAt time.Time      `binding:"ltecsfield=InnerStructField.StartDate"`
}
```

字符串约束  
- contains包含子串
- containsany包含任意unicode字符， containsany=abcd
- containsrune包含rune字符， containsrune= ☻
- excludes不包含子串
- excludesall不包含任意的unicode字符，excludesall=abcd
- excludesrune不包含rune字符，excludesrune=☻
- startswith以子串为前缀
- endswith以子串为后缀  

唯一性约束unique  
- 对于数组和切片，约束没有重复的元素
- 对于map，约束没的重复的value
- 对于元素类型为结构体的切片，unique约束结构体对象的某个字段不重复，通过unqiue=field指定这个字段名。如 ``` Friends []User `binding:"unique=Name"` ```


