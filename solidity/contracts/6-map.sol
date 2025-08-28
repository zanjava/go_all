// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Mapping {
    struct User {
        string Name;
    }
    // 在mapping中并不存储原始的key，而是存储key的keccak256哈希值，所以不存在key的集合的概念，也不能遍历mapping。在查询mapping时会创建任何可能的key，并将value赋为零值，所以mapping没有长度，不支持检查某个key是否存在，不支持删除某个key。
    mapping(address => User) userInfo; //mapping的key必须是内置类型，value可以是任意类型

    function set() public {
        address x = address(0x5B38Da6a701c568545dCfcB03FcB875f56beddC4);
        // address x = msg.sender;
        User memory user = User("gw");
        userInfo[x] = user;
    }

    // public函数的入参和出参(返回值)不能有storage变量，因为跨合约传“变量地址”很不安全
    function get() public view returns (User memory) {
        address x = address(0x5B38Da6a701c568545dCfcB03FcB875f56beddC4);
        // address x = msg.sender;
        return userInfo[x];
    }

    function remove() public {
        address x = address(0x5B38Da6a701c568545dCfcB03FcB875f56beddC4);
        // address x = msg.sender;
        delete userInfo[x]; //把相应的value置0
    }
}
