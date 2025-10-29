1.区块链 


保留记录和执行合同的分布式数据库技术
加密使数据不被篡改
共享账本

分布式账本的一致性和信任：共识机制

数据是如何分布的？ ===>>>分布式，各节点共享账本的副本

区块链中的数据表示状态，区块链使用事务将数据的状态从一个值更改为另一个值。
（在整个区块链网络中，区块链都会发送事务。 每个节点都会获取事务的副本）

块是区块链中存储事务信息的数据群集

如何确保账本不可变？=====>>>>加密哈希,在生成下一个块的哈希时包含上一个块的哈希值,这些块会通过哈希链接在一起

2.智能合约：

智能合约是区块链中存储的一种程序

智能合约的地址====>>>>20 个字节  ,合约的唯一标识

智能合约的主要属性和优点:透明  不可变性  分发(该合约的输出由该网络的节点验证。 合约状态可以公开显示。 在某些情况下，甚至可以看到“私有”变量)

3.使用OpenZeppelin创建代币
代币；可交换  不可互换（NFT）

合约标准：
ERC20 ERC777

NFT： ERC721

ERC1155:管理多种代币类型

4.以太坊

公共区块链网络  分布式计算平台（全球性的、开放）  
EVM  solidity
ETH 支付交易费用和网络服务的

POW(复杂数学问题)===>>>  pos (抵押 ETH) 
能耗 安全性 效率  

5.钱包  
非托管型钱包 ：用户完全控制钱包密钥（助记词或私钥）的钱包  ----去中心化钱包
托管型钱包  中心化钱包

冷钱包：不触网
热钱包；联网使用  安全性低

基础代币：原生代币 如 ETH  比特币的 BTC，支付交易费用
合约代币: 由智能合约创建 


钱包地址：资金存取的公开标识，由公钥通过加密算法生成
公钥：私钥对应的非保密配对
私钥：秘密数字，是区块链身份和资产控制的关键


区块链网络的数据完整性：加密  时间戳  区块链的不可变性特性 共识机制

区块链项目的升级和数据迁移：软分叉  硬分叉  使用智能合约可升级的设计，只更新合约逻辑部分



========================================================================

比特币闪电网络=====>>>>L2(目的是扩容)  计算在链下，把最终结果(Rollup)提交L1

加密算法 

哈希 --->>>默克尔树

非对称算法 ====server生成密钥对  将公钥发给client   client对msg(公钥)加密   server（私钥）对msg解密

签名----------->>>>>私钥加密  公钥解密

oracle预言机chainlink ------>>>>>>去中心化预言机网络 通过调chainlink合约拿到线下的数据，再传到线上的合约。比如实时的价格  

钱包---->>>>非对称加密   私钥---钥匙   加解密算法----锁  
            公钥---银行卡号，公开的
            助记词----私钥另一种展现形式
            keystore----加密的私钥
            密码
钱包是区块链的入口，进行签名 管理资产和交易

智能合约------自动执行  不可篡改  透明可验证  安全性

变量的作用域---------局部变量  状态变量  全局变量 
变量的存储位置 -------storage  memory  calldata  

error----自定义异常节省gas 

data location----------storage   memory   calldata 
修改链上数据---storage  相当于指针
只读取-----memory   拷贝 


继承的顺序性--->>>>  最上层最优先

⽗合约构造函数的初始化顺序---->>>>>由继承的顺序决定，⽽⾮调⽤顺序

constant----声明 在编译时就已确定且永不改变值
immutable----部署时确定值，并在部署后不可改变

call other contract -----------合约地址  vs 合约类型

library-------代码重用 数据类型增强  using for指令 

keccak256-----签名  生成唯一Id
abi.encode   
abi.encodePacked--添加⼀个整数参数,改变顺序，避免哈希冲突

A----->MultiCall--------(static) call ----->B  msg.sender=MultiCall合约地址
A----->MultiDelegatecall--------delegatecall ----->B  msg.sender=A

multi delegate call 重复调用可能引发问题

gas优化技巧：
使⽤calldata代替memory：通过改变变量存储位置来减少燃⽓消耗。
循环内部变量优化：在循环开始前将状态变量加载到内存，循环结束后再更新状态变量。
表达式短路（ShortCircuiting）：优化条件判断逻辑，避免不必要的计算。
循环增量简化：使⽤ ++i 代替 i + 1 来减少操作。
缓存数组⻓度：将数组⻓度存储在局部变量中，减少每次循环的计算量。
数组元素加载到内存：将频繁访问的数组元素预先加载到变量中。

质押收益---->>>>单币质押  流动性质押  借贷质押 
两种代币（质押代币  奖励代币）==>>>奖励机制设置===>>>用户质押代币

receive（）是msg.data为空且存在receive方法时触发
fallback（）是不匹配任何方法名，或者msg.data为空且合约没实现receive时触发

math  unchecked 禁⽤溢出和下溢检查  节省gas  

透明可升级代理：
delegatecall 的本质===>>>>delegatecall 会执行 目标合约的代码（CounterV1），但使用的是当前合约（BuggyProxy）的上下文和存储。

1. 对⻬存储布局-----确保所有实现合约与代理合约的存储布局⼀致(状态变量)
2. 改进回退函数---------修改回退函数以便能够返回数据

hardhat:
nodejs: [https://nodejs.org/en](https://nodejs.org/en)
yarn: `npm install -g yarn`

yarn init y
yarn add -D hardhat
npx hardhat --init
编写合约 --->>> npx hardaht compile
测试脚本----->>> npx hardhat test
部署脚本------>>>>npx hardhat run scripts/deploy.ts --network localhost  (npx hardhat node---启动本地网络)


连接小狐狸钱包-----webpack打包
yarn add -D webpack webpack-cli ts-loader html-webpack-plugin dotenv webpack-dev-server
npx webpack serve   
     
合约----json文件--->>> deploy---bytecode---->>>区块链（网络）
Dapp(浏览器)---->>>>metamask--->>区块链（网络）

代理合约部署任意合约------>>>>>   代理合约通过合约字节码拿到目标合约地址->加载目标合约（msg.sender是代理合约的地址）->代理合约要执行（目标合约地址，调用方法的字节码）->目标合约d的地址就是当前部署合约的地址了。


foundry:

export  PRIVATE_KEY=
forge create Counter --rpc-url sepolia --private-key $PRIVATE_KEY --broadcast
forge script ./script/Counter.s.sol:CounterScript --rpc-url sepolia --private-key $PRIVATE_KEY --broadcast


智能合约安全：

账户模型---外部账户（EOA）  合约账户（CA）

账户状态---nonce（该账户发出的交易数量）  banace（以太币余额）  storageRoot（存储的是默克尔树根节点的哈希值，用于存储和验证状态）
codehash(evm code的哈希) (外部账户的codehash为空)

节点---区块链---区块---交易

交易 ---->>>>状态转移  
交易的核心结构：
nonce
gaspraice 
gaslimit
to---接收地址
value---以太金额
data---消息调用的输入数据
signature

EVM----stack 256bit *1024  
       memory 读写8bit  256bit 
       storage  256bit-256bit 键值对

evm执行交易的步骤：
   合约编译--->>abi和bytecode---初始化一个evm环境，生成一笔交易tx
   1.加载合约字节码
   2.字节码转化操作码  逐条执行
   3.pc移动，gas计费  消耗b
   4.stac进行操作码运算
   5.memory临时存储中间数据
   6.交易过程和结果在storage永久存储 ---contract address 和runtime code

调用交易---生成tx(to:合约地址  data:函数选择器) 

ABI编码----abi.encode()   abi.encodePacked()

abi.decode()

abi.encodeWithSignature("withdraw(address,uint256)", to, amount);
abi.encodeWithSelector(Vault.withdraw.selector, to, amount);
abi.encodeCall(Vault.withdraw, (to,amount));

低级调用----call  delegateCall  staticCall

函数相关: 目标地址.call{value:X,gas:y}(calldata)
          目标地址.delegatecall(calldata)
          目标地址.staticall(calldata)
eth相关：
       目标地址.transfer(amount)
       目标地址.send(amount)
       目标地址.balance


call-----通用调用方式  可以调用其他合约函函数或发送以太
        是在目标合约的上下文中运行
        可指定calldata外，还可以指定gas value参数
delegate call----委托调用   存储布局要和目标合约一致
                 在当前合约的上下中运行
                 目标合约的存储（状态变量）不会改变
                 只可以指定calldata
static call-----调用合约不允许修改区块链的任何状态
                调用合约只能是pure 或view

转账类型：
ETH转账---非合约交易
Token转账（ERC20）---合约交易
token等资产转账对应的是合约中的状态变量
safeTransferFrom调用时，执行onERC721Received回调


      
send----gas限制，2300  不回滚   不推荐使用
transfer----  gas限制，2300 失败回滚  简单场景使用
call--无gas限制  不回滚  复杂场景使用


合约的生命周期；
创建阶段----部署交易，运行构造函数，生成runtime code 上链
运行阶段----调用函数，改变状态，执行外部调用函数
终止阶段----销毁交易，清除相关的状态，转移剩余的ETH

合约自毁----selfdesruct  强制转移ETH到指定地址

tx.origin -----EOA  原始交易的发起者
msg.sender----EOA or CA   当前执行上下文的调用者

合约攻击模式
npx hardhat vars set INFURA_API_KEY  XXX
const INFURA_API_KEY = vars.get("INFURA_API_KEY");


solidity 插槽机制