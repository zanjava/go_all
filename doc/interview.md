1.go语言底层原理

2.内存模型

- 垃圾回收机制  三色标记法,gc 针对于堆内存

- 内存逃逸,比如:数组超过128K就会分配到堆上(moved to heap);切片超过64K就会分配到堆上(escapes to heap)
  // go run -gcflags=-m
  // go tool trace ../data/trace.out

- 函数的入参、出参、局部变量一般在栈上;数组长度固定，通常在栈上;
  动态创建的结构体、切片、map在堆上。引用类型的全局变量内存分配在堆上，值类型的全局变量分配在栈上

- 向内存对齐:
  结构体的内存对齐原则是:最大成员变量的整倍数
  通过反射可以查看结构体里各变量的偏移量,和结构体所占总字节
  可以调整变量位置,优化内存管理

3.并发编程模型

  基于csp理论,通过goroutine 和channel 实现高效并发,支持M:N调度模型(m个协程映射到N个内核线程),实现秒级切换.

4.goroutine调度机制

  MPG:  M---machine  映射内核线程;P----processor ;G----goroutine

   运行队列  系统调用

5.channel实现原理

- 底层引用环形队列,可以指定容量，容量满时再写入数据会阻塞，容量空时再读取数据也会阻塞。
- 协程之间通信

6.sync包底层逻辑

- sync.Map:
  内部使用了两个主要的存储结构：`read`和`dirty`。`read`是一个只读的`map`，用于快速读取数据；`dirty`是一个可读写的`map`，用于存储新写入的数据。通过这种分离机制，`sync.Map`在保证并发安全的同时，还提高了读取性能

- sync.Pool
  一组可以单独保存和检索的临时对象重用池。它主要用于缓存已分配但未使用的对象，减轻垃圾收集器的压力
  内部使用一个双向链表来存储对象。当需要获取对象时，它会从链表中取出一个对象；如果没有可用对象，则会调用`New`方法创建一个新的对象

- sync.Once
  内部使用一个原子操作来确保`Do`方法只调用一次传入的函数.确保只执行一次

- sync.WaitGroup
  等待一组goroutine执行完成.通过`Add`、`Done`和`Wait`方法来管理goroutine的同步

- sync.Mutex
  利用互斥锁,保护共享资源  

7.框架 :

gorm 

gin/fiber

grpc,protobuf

kitex

8. Mysql Redis  数据库性能调优经验

9.mq

10.利用pprof性能瓶颈排查及使用对象池  (在线排查协程泄漏)

- -cpuprofile=data/cpu1  -memprofile=data/mem1  测试启动命令(基准测试)

- go tool pprof (-http=:5678) data/mem1,进入终端

- 常用交互命令：top n, peek func_name, list func_name
  
  --------------------------------------------------------------------

- impot   _ "net/http/pprof"

- func main() {  
  
      go http.ListenAndServe("127.0.0.1:5678", nil)  
      for {  
         Search()  
      }
  
  }

- pool   = sync.Pool{  
  
      New: func() any {  
         return make([]int, 10000)  
      },  
  
  }

- 浏览器访问127.0.0.1:5678/debug/pprof
  
  

11.solidity







# ***************go基础*************

1.Go中实现并发？

Go通过`goroutine`来实现并发。并发机制基于CSP模型,协程与协程之间通过`channel`进行通信。

2.Go中，如何处理error？

- 自定义error

- 创建error  `errors.New`

- 跟踪错误信息,fmt.Errorf(),对err进行包装

- errors.As背后其实就是类型断言errors.Is

3.在Go中，如何声明以及实现一个接口?

- 一组方法定义集合的接口,通过关键字type  XXX interface定义

- 通过某个类型实现接口中所有的方法

4.Go中，init方法的执行机制是什么？什么场景比较适合使用init方法？

场景:初始化全局变量,配置等

根据导入的包的顺序，依次从底层逐层往上执行`init`函数。main方法之前执行,全局变量之后

5.声明一个变量有哪些方式？

- var 

- :=

- 批量声明

6.Go里面有哪些基础类型？

* 数字类型：`byte`,`rune`,`int`, `float32`, `float64`, `complex64`, `complex128`

* 字符串类型：`string`

* 布尔类型：`bool`

* 派生类型：`array`, `slice`, `map`, `struct`, `pointer`, `function`, `channel`, `interface`

其中`byte`,`rune`,`string`这三个类型可以相互转换。

7.两个nil是否相等？

两个`nil`值的比较结果取决于它们的类型。例如两个`interface{}`类型的`nil`值可以比较，相等。

8.怎么声明一个常量？

使用`const`关键字声明常量。

9.怎么声明一个方法？

方法是和类型关联的函数

10.怎么创建一个切片？切片和数组有什么区别？

使用`[]`声明切片，但不指定长度

切片底层是一个结构体，它持有一个数组的指针，并且额外维护了len和cap字段。

11.map

- 使用`make`或直接初始化。

- 使用`for`range`遍历，直接使用key访问value，用`delete`删除元素。

- `slice`、`map`、`function`不能作为map的key。因为它们是引用类型，不能比较是否相等

12.channel

- 使用`make`创建channel，第二个参数是指定channel长度。

- 无缓冲channel在发送和接收时必须同步

- 缓冲channel则允许异步发送。

- 缓冲区满前，可以随时写入，不会阻塞，直到缓冲区满后，写入需要等缓冲区有空余位置。
  如果缓冲区没有数据，那么读取会阻塞，只要有数据为止。

13.怎么关闭一个channel？读取一个关闭的channel会发生什么？怎么判断一个channel已经被关闭了？

使用`close`关闭channel，读取关闭的channel会返回零值，判断是否关闭可通过第二个返回值

14.什么是goroutine？怎么使用goroutine？怎么停止一个goroutine？

- 并发执行单元。使用`go`关键字启动

- 可以通过`context`或`channel`控制`goroutine`的停止。

15.怎么编译应用为一个可执行二进制文件？

 go build 生成exe文件

16.Go中nil的切片(Slice)和空的切片有何不同？

- `nil`的切片没有分配任何底层数组，长度和容量为0；

- 空切片已经分配了底层数组，长度为0，但容量可能大于0。

17.`append`用于向切片添加元素

18.Go中nil的map和空map有何不同？

- `nil`的map没有分配内存，不能直接存储键值对；只是声明map

- 空map可以存储键值对，但初始为空。通过make创建

19.函数返回局部变量的指针是否安全？

是安全的。Go的内存管理机制会自动扩展局部变量的生命周期。

20.数组、slice、struct允许并发修改（可能会脏写），并发修改map可能会（不是一定会）发生fatal error，recover()只能捕获panic，但不能捕获fatal error。如果需要并发修改map请使用sync.Map

21.多协程如何协调结束工作  以及如何避免重复关闭channel?

- close channel  进行广播  

- 关闭在单独协程里进行通知,同时发送方和接收方在多路监听里(避免阻塞),往关闭的channel里发送关闭信息
  
  

# *****************中级****************

1.如何从panic中恢复？或者说如何拦截panic？

- 使用recover()从panic中恢复,需配合defer使用.

- 只能在panic中使用,在本协程里生效

2.Golang中make和new的区别？

- `new`：分配内存，但不初始化。它返回的是类型的指针，通常用于值类型（如数组、结构体）。

- `make`：只用于创建`slice`、`map`和`channel`，并且会初始化这些数据结构，返回的不是指针而是数据类型的引用。

3.Golang中函数和方法的接受者可以是值和指针，它们有什么区别？

* 值接收者：函数操作接收者的副本，不能修改原始值。

* 指针接收者：函数操作接收者的引用，可以修改原始值。

4.Golang中非接口的任意类型T()都能够调用*T的方法吗？反过来呢？

- 可以。如果是值类型T，可以调用指针接收者的方法，Go会自动将值转为指针。

- 指针类型也可以调用值接收者的方法。

- 但是只有接受者是指针类型时，才可以修改原始值。

5.调用函数传入结构体时，应该传值还是指针？

传指针

- 传值会拷贝整个结构体，性能上可能较慢，尤其是大结构体

- 传指针可以避免拷贝并允许修改原值

6.uintptr和unsafe.Pointer的区别？

- `unsafe.Pointer`：是一种通用指针类型，用于不同类型指针之间的转换，不参与指针运算。

- `uintptr`：是一个整数类型，用于存储指针值，可以进行指针运算，但与垃圾回收无关

7.描述一下defer关键字的作用？什么时候会用到defer？

- 延迟执行函数,当前的函数返回后执行

- 常用于资源释放操作（如关闭文件、解锁互斥锁）

8.同一个方法中声明多个defer，defer的方法怎么执行？

多个`defer`按照后进先出的顺序执行（LIFO）

9.循环内部执行defer会发生什么？

在循环中使用`defer`，每次循环都会推入一个`defer`操作，这些操作会在方法结束后按逆序执行

10.  defer和return的执行先后顺序如何？

`defer`的执行顺序是在`return`语句之后，函数退出之前

11.defer是否可以修改返回值？

可以。如果函数声明了命名返回值，`defer`中可以修改它

12.如何创建一个自定义类型？它的作用是什么？

- 使用`type`关键字定义自定义类型

- 让现有类型更具语义或行为

13.什么是闭包(closure)？可以在哪些地方使用？

- 在函数内声明的匿名函数

- 可以捕获和引用外部作用域变量

- 常用于回调函数或工厂函数

14.select关键字的作用？描述一下使用场景

- 监听多个`channel`的操作(多路复用)

- 通常配合无限循环的for使用

15.类型如何强转？如何判断变量是否为指定类型？如果变量可能是多种类型，怎么判断？

- type `Conversion`

- `type assertion`     type.()

- 使用`switch`来处理多种类型的变量

16.并发情况下时候可以使用基础数据结构map和slice以及array？会发生什么？

- map不行,发生panic

- 建议使用`sync.Mutex`或者`sync.Map`来保证并发安全

17.atomic包下提供的类型怎么使用？

- 对基础类型的原子操作

- 用于无锁并发编程

18.怎么使用context包在一个请求作用域内携带上下文相关的数据？如何读取？

使用`context.WithValue`将数据添加到`context`中，使用`context.Value`读取

19.strings包提供了哪些处理字符串的方法？

* `strings.Split`：分割字符串

* `strings.Join`：连接字符串

* `strings.Contains`：检查子字符串是否存在

* `strings.Replace`：替换子字符串

20.如果处理一个操作需要花费一些时间处理，应该怎么添加超时检查？如果还没到超时时间，需要取消这个正在执行的操作，应该怎么做？

- 使用`context.WithTimeout`设置超时 ;context.WithCancel 调用cancel，触发Done,去监听    

- `context.CancelFunc`取消正在执行的操作

21.`json.Unmarshal`解析JSON，使用`json.Marshal`生成JSON。

22.在Go中怎么创建一个单元测试用例？怎么在命令行下单独执行某个指定的测试用例

- 使用`testing`包编写测试函数，函数名以`Test`开头。

- 使用命令go test -run 执行指定测试用例

23.Golang的内存模型，为什么小对象多了会造成gc压力？怎么解决

- 增加gc的频率,导致性能下降

- 使用对象池sync.Pool,避免频繁创建对象

24.协程、线程、进程的区别？为什么Goroutine比线程轻量？

- Goroutine：轻量级的线程，Go运行时管理，调度更高效

- 线程：进程的执行单元，共享内存。

- 进程：独立运行的程序，具有自己的内存空间。
  Goroutine比线程轻量是因为它们由Go运行时管理，不需要操作系统的线程调度开销。

25.为什么channel可以做到线程安全？

-     内部实现使用了同步机制

- 保证了发送和接收操作的原子性

26.Go中切片如何扩容？

- 当切片容量不足时，Go会自动扩容，通常是将容量翻倍

- 使用`append`函数时，如果超出当前容量，Go会分配一个新的底层数组并复制元素

27.Go中map如何扩容？

`map`会自动扩容，重新分配哈希桶来存储更多键值对



# *****************高级****************

1.Go中如何排查内存泄漏和性能问题？

- **pprof**：`net/http/pprof`包可以用于生成程序的性能剖析数据（CPU、内存、goroutine等）

- **go tool pprof**：通过`go tool pprof`可以分析内存和 CPU 使用情况

- **trace**：使用`runtime/trace`生成追踪文件，分析程序的调度、内存分配等
  
        go run -gcflags=-m
        go tool trace ../data/trace.out

2.Go中channel的底层数据结构与实现？

Go 中的 `channel` 是通过队列实现的，有缓冲和无缓冲之分。

* **无缓冲 channel**：底层使用一个`goroutine`队列，读写必须同步，即发送方等待接收方，反之亦然。

* **有缓冲 channel**：底层维护一个环形队列，写操作将数据存入队列中，读操作从队列中取数据。如果缓冲区满了，写操作阻塞；如果缓冲区为空，读操作阻塞。 底层使用`lock-free`数据结构与`sync/atomic`原子操作来保证并发安全。

3.什么是Data Race？go中如何定位和避免？

**Data Race** 是指两个或多个 `goroutine` 并发访问共享变量，并且至少一个是写操作，且这些操作没有适当的同步机制。

定位：

* 使用 `-race` 选项检测数据竞争。
    go run -race main.go

避免方法：

* 使用互斥锁 `sync.Mutex` 来保护共享资源。

* 使用 `channel` 进行消息传递，避免直接共享数据。

4.如何交叉编译？

只需设置目标系统的 `GOOS` 和 `GOARCH` 环境变量

5.条件编译?

编译标签（build tags）实现条件编译

6.Golang的GC过程？

Golang 使用的垃圾回收器是三色标记-清除算法：

* **标记阶段**：从根对象开始，标记所有可达对象。

* **清除阶段**：清除所有未标记的对象，释放其内存。

* **并发GC**：标记操作与程序的执行并行，减少停顿时间。

* **增量GC**：GC 过程以小步增量执行，避免大停顿。

`GOGC` 环境变量可以控制垃圾回收的频率。

7.Golang中Goroutine是如何调度的？

Go 运行时使用**GMP 模型**（Goroutine、Machine、Processor）来调度 `goroutine`。

* **G**（Goroutine）：代表每个任务或执行单元。

* **M**（Machine）：表示底层的操作系统线程。

* **P**（Processor）：代表逻辑处理器，负责将 `goroutine` 分配给 `M`。 调度器通过`work stealing`、抢占式调度、避免锁竞争来提高效率。

8.描述一下Golang并发机制？以及它所使用的CSP并发模型？

Go 并发是基于 **CSP（Communicating Sequential Processes）** 模型的，核心是通过 `goroutine` 和 `channel` 进行并发编程：

* **goroutine**：轻量级线程，由 Go 运行时管理。

* **channel**：用于 `goroutine` 之间的通信，避免了共享内存的直接操作。

在 CSP 模型中，任务通过消息传递进行通信，确保线程安全。

9.怎么查看goroutine的数量？怎么限制goroutine的数量？

查看 Goroutine 数量：

* 可以使用 `runtime.NumGoroutine()` 查看当前的 `goroutine` 数量。

限制 Goroutine 数量：

* 使用 `sync.WaitGroup` 来控制并发数

10.Go中如何实现一个协程池?

- 使用带缓冲的 `channel`

- sync.Pool

11.Go中select的实现原理？

- `select` 通过监听多个 `channel` 上的操作来实现多路复用

- 底层实现为遍历所有 `case`，然后根据 `channel` 的状态（阻塞或可用）来决定执行哪个分支

- 内部使用了锁和同步机制来避免并发访问冲突

12.Golang的栈空间管理？

- 自动管理的，初始栈很小（一般是 2KB 左右），当 Goroutine 需要更多栈空间时，Go 会动态扩展栈空间（通过分配新的栈并拷贝旧栈数据）

- 栈是连续的，它的动态增长使得 Goroutine 非常轻量。

13.描述下Go的对象在内存中的分配过程？

Go 中对象的内存分配分为两种情况：

* **栈上分配**：局部变量和小对象通常分配在栈上。

* **堆上分配**：当对象超出作用域或者无法确定生命周期时，会在堆上分配。

Go 通过逃逸分析决定变量是否逃逸到堆。运行时通过 `mallocgc` 函数分配堆内存，并由垃圾回收器管理堆内存的释放。



# *****************框架****************

### gorm:

1.Gorm中，查询一条数据，如果没有查到会怎么样？

user := User{City: "HongKong", Id: 3} //Id会自动放到Where条件里，其他非0字段不会

    tx := db.

        Select("uid,city,gender,keywords").....

   <mark> if tx.Error != nil {</mark>

        if !errors.Is(tx.Error, <mark>gorm.ErrRecordNotFound</mark>) {

            slog.Error("读DB失败", "error", tx.Error)

        } else {

            slog.Info("查无结果")

        }

    } else {

        if tx.<mark>RowsAffected </mark>> 0 {

            fmt.Printf("read结果：%+v\n", user)

        } else {

            slog.Info("查无结果", "user", user)

        }

    }

2.Gorm更新数据时数据位0值被忽略了如何解决？

使用 `Select()` 方法明确指定要更新的字段，或者<mark>使用 `Updates` 方法可以指定要更新的列</mark>。
    db.Model(&user).Updates(User{Name: "Tom", Age: 0})

或者明确选择字段：
    db.Model(&user).Select("Age").Updates(User{Age: 0})

3.Gorm与原生方式相比有什么优势？

**自动迁移**：通过 `AutoMigrate` 自动创建和更新数据库表结构。

**ORM 特性**：通过<mark>结构体映射数据库表，避免手写 SQL</mark>。

**事务管理**：<mark>内置的事务支持，更方便的控制</mark>。

**查询构造器**：<mark>提供链式调用的方式构造查询条件，减少 SQL 拼接的风险</mark>。

4.Gorm如何实现1对多和多对多的关系映射？

**一对多**：使用 `hasMany` 关系，定义结构体中的关联字段。

**多对多**：使用 `many2many` 关系，定义中间表。

5.Gorm查询数据库数据时，如何避免N + 1查询问题？

在使用一对多或多对多关系查询时，<mark>GORM 的 `Preload` 方法可以帮助你避免 N + 1 查询问题</mark>。`Preload` 会提前加载关联的数据，避免每次访问关联字段时单独发出 SQL 查询。
    var users []Userdb.Preload("Posts").Find(&users)

这会在一个 SQL 查询中加载 `users` 和他们的 `posts`，避免多次查询。

6.如何使用Gorm进行事务管理？

<mark>GORM 提供了显式的事务管理功能。你可以通过 `Begin` 开启事务，并使用 `Commit` 或 `Rollback` 结束事务。</mark>

7.Gorm中Preload方法和Joins方法的区别是什么？二者各自适合的场景？

* **`Preload`**：通过单独的 SQL 查询来加载关联数据，适用于你不需要复杂条件的关联查询。它会先查询主表，再查询关联表的数据。
    db.Preload("Posts").Find(&users)

* **`Joins`**：通过 SQL JOIN 语句将数据表连接起来一次性查询所有数据，适合复杂的联表查询，性能可能更高。
    db.Joins("JOIN posts ON posts.user_id = users.id").Find(&users)

**选择场景**：

* 当需要<mark>一次性加载</mark>关联表数据并且<mark>没有复杂的条件时</mark>，使用 `Preload`。

* 当需要在查询中使用复杂条件、筛选或排序时，使用 `Joins`。

8.如果结构体和数据库中的表名称对应不上时如何解决？

可以通过 GORM 的 `Table` 方法或<mark> `gorm:"tablename"` 标签显式指定结构体对应的表名称</mark>。

* 通过标签指定表名：
    type User struct {    ID   uint    Name string}func (User) TableName() string {    return "custom_user_table"}

* 在查询时动态指定表名：
    db.Table("custom_user_table").Find(&users)

这样即使结构体名与数据库表名不一致，也可以正确映射表。



### gin:

1.如何编一个Docker镜像？

创建 `Dockerfile`，定义应用的构建步骤。

使用 `docker build` 构建镜像。

2.Gin框架如何文件上传？

- <mark>Gin 支持文件上传，使用 `<mark><mark>ctx.FormFile("file")</mark></mark>` 获取上传的文件并保存到指定位置</mark>。

- <mark>ctx.MultipartForm()  files := form.File["files"] </mark> SaveUploadedFile()

3.Gin如何解决跨域问题？如何配置？

Gin 可以通过<mark>中间件来解决跨域问题</mark>，常用的方式是使用 `gin-contrib/cors` 包。

安装：
    go get github.com/gin-contrib/cors

配置 CORS 中间件：
    import "github.com/gin-contrib/cors"func main() {    r := gin.Default()    r.Use(cors.Default()) // 默认允许所有源    r.Run(":8080")}

也可以自定义 CORS 配置：
    r.Use(cors.New(cors.Config{    AllowOrigins: []string{"http://example.com"},    AllowMethods: []string{"GET", "POST"},    AllowHeaders: []string{"Origin", "Content-Type"},}))

4.Gin支持哪些HTTP请求方式？

Gin 支持以下 HTTP 请求方式：

* `GET`

* `POST`

* `PUT`

* `DELETE`

* `PATCH`

* `OPTIONS`

* `HEAD`

5.如何在Gin中处理GET和POST请求参数？

**GET 请求参数**：使用 `c.Query("param")` 获取 URL 中的查询参数。

**POST 请求参数**：使用 `c.PostForm("param")` 获取表单数据，或 `<mark>c.BindJSON(&obj)` </mark>获取 JSON 数据。
    // 处理 GET 请求参数func handleGet(c *gin.Context) {    name := c.Query("name") // ?name=foo    c.String(http.StatusOK, "Hello %s", name)}// 处理 POST 请求参数func handlePost(c *gin.Context) {    name := c.PostForm("name")    c.String(http.StatusOK, "Received %s", name)}

6.Gin框架中如何实现路由？

Gin 使用 `<mark>router := gin.Default()` 创建路由实例</mark>，通过 `router.GET`、`router.POST` 等方法来定义路由，并将处理逻辑绑定到相应的 URL 路径和 HTTP 方法。
   <mark> r := gin.Default()</mark>r.GET("/ping", func(c *gin.Context) {    c.JSON(200, gin.H{        "message": "pong",    })})r.Run(":8080")

7.Gin框架的错误处理方式是怎样的？

Gin 使用 `Context.AbortWithError()` 或 <mark>`Context.Abort()</mark>` 来处理错误。可以在中间件中或控制器中间接调用。可以通过自定义中间件统一处理错误。
    func errorMiddleware(c *gin.Context) {    if err := someOperation(); err != nil {        c.AbortWithError(http.StatusInternalServerError, err)    }}

8.Gin框架如何处理并发请求？

Gin 天然支持并发处理请求，每个请求在一个独立的 `goroutine` 中运行。Gin 使用 `Context` 进行请求上下文的管理，不会引起并发冲突。

9.Gin框架中Context的作用是什么？

`Context` <mark>在 Gin 中用于管理请求的上下文信息</mark>。它提供了：

* 读写<mark>请求参数</mark>（GET/POST 参数、路径参数等）

* <mark>处理响应</mark>

* 存储和传递<mark>中间件之间共享的数据</mark>

* <mark>管理请求生命周期（如超时、错误处理）</mark>
  
  
  
  

### RabbitMQ:

###### 一.默认交换机,模式为direct:

producer

1.连接mq

2.创建channel (可以创建多个channel ,往同一个对列写数据)

3.声明队列(队列需持久化),指定队列名

4.发送消息(可以开多个协程,并行处理,消息也需要持久化)  , exg "",为默认的Exchange（direct类型），这种Exchange会把消息传递给routing key指定的Queue名



consumer

1.连接mq

2.1.创建channel (可以创建多个channel ,对同一队列中的数据进行消费)

(2.2.ch.Qos  可以设置prefetch count。一个消费方最多能有多少条未ack的消息)

3.声明队列(队列需持久化),指定队列名-----------producer consumer 看他们谁先创建

4.创建消费者(指定channel和队列名,auto-ack(no-ack)设置伟false)

5.消费消息(需手动确认消息是否被消费,也可以开多个协程,并行处理)



###### 二.交换机为direct

producer

1.连接mq

2.创建channel 

3.声明交换机(类型direct,exg需持久化)

4.发送消息(指定交换机和route key,消息持久化)



consumer

1.连接mq

2.创建channel

3.声明队列(队列需持久化,队列名为空,Server指定一个随机（且唯一）的队列名)

4.队列和Exchange建立绑定关系(队列名,由server指定),同时指定routing key

5.创建消费者(指定channel和队列名,auto-ack(no-ack)设置伟false)

6.消费消息(需手动确认消息是否被消费,也可以开多个协程,并行处理)



###### 三.交换机为faout

producer

1.连接mq

2.创建channel

3.声明交换机(类型faout,exg需持久化)

4.发送消息(指定交换机和route key为空,消息持久化)



consumer

1.连接mq

2.创建channel

3.声明队列(队列需持久化,队列名为空,Server指定一个随机（且唯一）的队列名)

4.队列和Exchange建立绑定关系(队列名,由server指定),routing key为空(fout模式下会忽略routing key)

5.创建消费者(指定channel和队列名,auto-ack(no-ack)设置伟false)

6.消费消息(需手动确认消息是否被消费,也可以开多个协程,并行处理)



###### 四.交换机为topic

producer

1.连接mq

2.创建channel

3.声明交换机(类型topic,exg需持久化)

4.发送消息(指定交换机和route key(正则匹配),消息持久化)

consumer

1.连接mq

2.创建channel

3.声明队列(队列需持久化,队列名为空,Server指定一个随机（且唯一）的队列名)

4.队列和Exchange建立绑定关系(队列名,由server指定),同时指定routing key

5.创建消费者(指定channel和队列名,auto-ack(no-ack)设置伟false)

6.消费消息(需手动确认消息是否被消费,也可以开多个协程,并行处理)





### Kafka:

- kafka是一个集群,由多个broker组成,每个broker可以部署单个机器上(端口号不同),实际中一台机器上部署一个broker,它由多个topic组成,一种业务是一个topic,    而同一个topic又是由多个partition组成(是为了负载均衡),分为leader和多个follower,提高可靠性

- 一个group对应一个使用数据的业务方,即每个group消费一份完整的topic数据.一个topic可以由多个group消费

- 对于同一个group而言,一个partition只能由一个consumer来消费,一个consumer可以消费多个partition。consumer可以指定某个partiton消费,但不能同时指定group ID.partition与consumer的配对关系自动调整,由hashring算法实现

- partition内部的消息是有序的,越新的消息offset越大.consumer每消费partition里的消息都会上报commit.也可以间隔性上报.consumer重启时,kafka会根据group上次提交的最大offset开始消费

- 对于生产者而言,要指定topic,写数据时决定选择哪个partition进行写入:1.可以显示指定partition  2.没有指定,是根据key的哈希算法选择一个partition写入(同一个key只能写入同一个partition).3没指定partition,也没有key,是按时间片轮转选择

- 

- 

- producer询问kafka cluster,得到特定的topic的特 

- producer把数据发给leader,leader将数据写入本地磁盘

- follower从leader上拉取数据,把数据写入本地磁盘

- follower向leader返回ack,确认成功

- leader向producer返回ack,确认成功

- 

- 

- 消息顺序消费,不要依赖kafka的顺序性...消息要写入同一个partition,又对应同一个consumer.即要保证partition数目不发生变化,要保证consumer不重启,且没有新的consumer加入.
  
  

    

### RocketMQ:

producer  cluster                                                                nameServer   cluster
                    proxy cluster

consumer     cluster                                                            broker   cluster(主从,负载均衡)



1. 先启动nameServer ->broker->proxy 

2. producer:一条消息只能属于一个topic ,只能打一个tag

3. consumer:有消费组group,consumer属于哪个消费组,一个group下的consumer平分所有消息

4. 订阅关系:consumer  通过topic  tag 去订阅消费消息 

5. 消息类型:
   Normal
   DELAY:设置消息延迟投递时间
   FIFO:要设置消息组,保证生产顺序性
   TRANSACTION:手动producer.BeginTransaction()   commit  rollback
   


