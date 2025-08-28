// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Variable {
    /**
    函数之外的变量是“状态变量”，变量值永久保存在合约存储空间中。状态变量默认是internal
    在函数执行过程中产生的变量为“局部变量”，函数退出后变量无效
    */
    int256 it; //int等价于int256，有int8, int16, int24 ... int256，以8位为步长递增。默认值为0
    uint256 uit = 7653; //uint等价于uint256
    bool gender; //默认值false

    /**
    基本原则：
    1. 需要永久存储的变量比临时变量消耗更多gas。所以状态变量比局部变量消耗更多gas
    2. 需要修改变量的操作比不修改变量的操作消耗更多gas。所以常量消耗的gas更少，view函数消耗的gas更少
    */
    string constant name = unicode"昝高伟"; //常量的值不可改变，且在声明常量时必须赋值。字符串如果包含汉字需要在前面加unicode
    uint256 immutable age = 30; //immutable不可变的，只能在声明或构造函数里赋值，即部署合约的时候它的值必须确定下来，不能再变

    // 构造函数在部署合约时会被调用。构造函数默认是public，所以不用显式地写public
    constructor(uint256 x) {
        age = x; //可以在构造里给immutable变量赋值
    }

    address c = 0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db; //20字节的以太坊地址。长度必须是20个字节，且不能随意构造
    address adr = address(this); //当前合约的地址
    enum Direction {
        East,
        West,
        South,
        North
    } //枚举，内部编码为0、1、2、3

    //结构体
    struct User {
        string Name;
        int256 Age;
        string City;
    }

    User user1 = User({Name: unicode"高伟", Age: 18, City: "BJ"}); //结构体。指定成员名称再赋值
    User user2 = User(unicode"高伟", 18, "BJ"); //结构体。直接赋值
}

/**
全局变量：
block.chainid--uint
block.coinbase--address  挖出当前区块的矿工地址
block.difficulty--uint  当前区块的难度
block.gaslimit--uint
block.number--uint  区块号
block.timestamp--uint  精确到秒。理论上它是创建区块的时间，它是块头散列值的输入项之一；实际中这个值矿工有很大的自由度，只要它位于前后两个区块的timestamp之间即可。
msg.sender--address  消息发送者（当前调用）
msg.value--uint  随消息发送的 wei 的数量
msg.data--bytes  完整的 calldata
msg.sig--bytes4  calldata的前4个字节
tx.gasprice--uint 
tx.origin--address  交易发起者（完全的调用链），只有一个合约时tx.origin就是msg.sender
*/
