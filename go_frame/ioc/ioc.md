## 依赖注入(Dependency Injection)
非自己主动初始化依赖，而通过外部来传入依赖。  
## 控制反转(Inversion of Control)
对象A依赖于对象B，不中由A主动去创建对象B，而是由IoC容器去创建，控制权颠倒过来了，所以称之为“控制反转”。  
## 依赖注入和控制反转的关系
- 控制反转是一种思想，依赖注入是一种设计模式。
- IoC框架使用依赖注入作为实现控制反转的方式，但是控制反转还有其他的实现方式，例如说ServiceLocator。
- 在wire框架中，何时去创建依赖是由业务代码自己控制的（在调用InitHandler()时创建了中间每一层的依赖），wire框架帮我们建立了依赖链路；IoC框架也帮我们建立了依赖链路，但创建依赖的时机不是由业务代码自己控制的（在调用Resolve()时Singleton结构体可能在之前某个时间已经创建好了），依赖是由IoC容器创建的。不论是wire还是ioc框架，我们只需要把每个环节的依赖关系告诉框架，比如B->F,A->G,F->A，框架会自动梳理出完整的依赖链条B->F->A->G，当我们去向框架索要B对象时，框架会依次创建好G、A、F、B。
## golobby IOC框架
- 告诉框架如何创建对象：`container.Singleton`、`container.Transient`、`container.NamedSingleton`、`container.NamedTransient`。这些代码一般放在`init()`函数中。
- 向框架索要对象：`container.Resolve`、`container.Fill`