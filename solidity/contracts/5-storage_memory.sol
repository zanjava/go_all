// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.1;

contract Storage {
    uint a; //状态变量不能用storage或memory修饰，它就是storage，它就是存储在区块链上的

    function foo() private pure returns (uint) {
        uint b = 4; //值类型的局部变量不能用用storage或memory修饰，它就是memory，它就是临时存在的变量
        return b;
    }

    //结构体
    struct User {
        string Name;
        int256 Age;
    }

    User user1 = User({Name: unicode"高伟", Age: 18}); //结构体。指定成员名称再赋值

    function changeUser() public {
        User storage u = user1; //如果改成memory则修改的是拷贝，并不影响状态变量user1
        u.Age = 28;
    }

    function getAge() public view returns (int256) {
        User memory u = user1;
        return u.Age;
    }

    function soo() private pure returns (int) {
        User memory u = User({Name: "dqq", Age: 18}); //=右边的结构体是在函数内部创建的，是局部变量，所以=左边必须是memory
        return u.Age;
    }

    // 参数或返回值为数组（包括bytes和string）、mapping或结构体时，必须指定是memory还是storage
    function noo() private pure returns (User memory) {
        User memory u = User({Name: "dqq", Age: 18});
        return u;
    }

    // 注意，这个函数不能用pure修饰，因为将来可能把状态变量传给本函数，同时又是storage，所以会修改状态变量，与pure的承诺矛盾
    // 参数或返回值为数组（包括bytes和string）、mapping或结构体时，必须指定是memory还是storage
    function goo(uint[] storage arr) private {
        uint[] storage brr = arr; //storage类似于指针，修改一个，也会影响另一个
        brr[0] = 888;
    }

    // 可以用pure修饰，因为是memory，传是的拷贝
    function moo(uint[] memory arr) private pure {
        arr[0] = 888;
    }

    // calldata类似于memory，同时calldata参数不能修改，只能读取
    function boo(uint[] calldata arr) private pure returns (uint) {
        uint i = arr[0];
        return i;
    }

    // public和external函数的入参和出参(返回值)不能有storage变量，因为从外部传“变量地址”很不安全
    function qoo(User storage arg) internal returns (User storage) {
        User storage u1 = arg;
        User memory u2 = User({Name: "gw", Age: 18}); //在函数内初始化的引用类型(数组（包括bytes和string）、mapping或结构体)的变量不能用storage修饰，它只能是memory
        u1.Name = u2.Name;
        return u1;
    }

    string ss1 = "1234567890123456789012345678901234567890";
    string ss2 = unicode"😃高HA";
    bytes bs1 = "1234567890123456789012345678901234567890";
    bytes bs2 = new bytes(3);

    function getString() public view returns (bytes memory) {
        bytes memory tmp = bytes(ss2);
        return tmp;
    }

    function bytesOP() private view returns (string memory, bytes memory) {
        string memory join1 = string.concat(ss1, ss2, "golang");
        bytes memory join2 = bytes.concat(bs1, bs2);
        return (join1, join2);
    }
}

/**
状态变量不能用storage或memory修饰；值类型的局部变量不能用用storage或memory修饰
storage：存储在链上，消耗的gas多，合约中的状态变量
memory：内存中，消耗的gas少，局部变量中的变长的数据类型必须用memory修饰，如string、bytes、array、struct、mapping。storage相当于指针，memory相当于拷贝
calldata：内存中，消耗的gas少，与memory不同的是其不能修改，常用做函数的参数
*/
