// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Function {
    bool gender;

    function(bool) f = setGender; //函数本身可以作为一种数据类型，函数类型默认是内部函数，因此不需要声明internal关键字。

    /**
    public修饰的变量和函数，任何用户或者合约都能调用和访问。public状态变量会自动生成一个getter()函数，用于查询数值。
    private修饰的变量和函数，只能在其所在的合约中调用和访问，即使是其子合约也没有权限访问。
    internal 和 private 类似，不过， 如果某个合约继承自其父合约，这个合约即可以访问父合约中定义的“内部”函数。
    external 与public 类似，只不过这些函数只能在合约之外调用，内部可以通过this.xxx()访问。外部函数由一个地址和一个函数签名组成。
    状态变量默认是internal，而函数的visibility需要显式指定。
     */
    function setGender(bool g) private {
        gender = g;
    }

    // 在函数执行过程中产生的变量为“局部变量”，函数退出后变量无效。
    function setMale(uint256 arg) public returns (address) {
        uint256 tmp = arg + 100; // tmp和arg都是局部变量
        delete gender; //delete变量会把变量赋为“零值”。bool的零值是false
        if (tmp > block.number) {
            // block、msg、tx是全局变量，在任何地方都可以访问到，用于获取区块链的相关信息
            gender = true;
        } else {
            gender = false;
        }
        return tx.origin;
    }

    // 函数的重载：可以存在同名的函数，只要入参不同就行，因为函数签名是keccak256("setMale(uint256,bool)")，函数返回值不会用于函数签名
    function setMale(uint256 arg, bool b) public {}
    // function setMale(uint256 arg) public {}   不合法的代码

    // returns里可以指定多个值
    function multiReturn() external pure returns (uint256, bool, int) {
        return (5, true, 9);
    }

    /**
    constant、view、pure三个函数修饰词的作用是告诉编译器，函数不改变/不读取状态变量，这样函数执行就可以不消耗gas了（是完全不消耗！）
    view的作用和constant一模一样，可以读取状态变量但是不能改；
    pure则更为严格，pure修饰的函数不能改也不能读状态变量，不能通过this调用当前合约的函数，否则编译通不过。
     */

    // 可以指定return的变量名也可以不指定
    function multiReturn(
        bool b
    ) internal view returns (uint256 arg, bool m, int) {
        if (b) {
            arg = 10;
        } else {
            // 访问external函数，得加this。同时当前函数不能是pure
            (arg, m, ) = this.multiReturn(); //函数返回多个值，此处只用到第一个和第二个值
        }
        return (arg, m, 9);
    }

    /**
    payable变量通过调transfer或send可接收转账，payable函数被调用时需要给函数所在的合约转账。
     */

    address c = 0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db; //20字节的以太坊地址。
    address payable pc = payable(c); // payable address多了transfer和send两个成员方法，用于接收转账

    // payable表示调函数的时候必须携带value，给合约账户转账
    uint public a;
    function deposit() public payable {
        a = msg.value; // 函数的调用方给合约转了这么多钱（单位Wei）
    }

    function withdraw(address payable x) public {
        address myAddress = address(this); //可以把合约转成地址
        if (myAddress.balance >= 10) {
            //balance表示一个地址的余额
            x.transfer(10); //从合约地址向x地址转账10Wei。x必须是payable地址才能调用transfer
        }
    }
}
