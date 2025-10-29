1. evm 存储结构

stack====>>>32bit   256位

memeroy

storage===>>>存储再链上  永久====>>状态变量=》stateDB

基本数据类型====>>>>长度固定  实际数据就在stack上

整数

bool 

地址 

固定长度字节

枚举


引用类型===>>>长度未知或者超过32bit  实际数据就在memory或storage上，
              stack只有keccack256类型的哈希

数组

字符串

结构体  

映射


constant/immutable===>>>>> 存储再字节码/代码中








