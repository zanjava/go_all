// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Modifier {
    address owner;
    bool locked;
    int256 public a;

    constructor() {
        owner = msg.sender; //部署合约的账户即为owner
        locked = false;
    }

    // modifier用处一：权限控制
    modifier onlyOwner() {
        require(msg.sender == owner, unicode"只有owner才能权限调用此函数");
        _; //修饰的函数代码内联到这里。_;不一定非得放到modifier末尾
    }

    // modifier用处二：参数校验
    modifier checkInt(int256 i) {
        require(i > 0, "int arg must more than 0");
        _;
    }

    // 一个函数可以有多个modifier
    function changeA(int256 x) public onlyOwner checkInt(x) {
        a = x;
    }

    // modifier用处三：防止重入
    modifier noReentrancy() {
        require(!locked, unicode"函数禁止重入");
        locked = true; //获得锁
        _; //执行函数
        locked = false; //释放锁
    }

    // modifier和visibility的先后顺序无所谓
    function sub(int256 m) public noReentrancy {
        a -= m; //由于revert导致-m没有成功
        if (m > 1) {
            sub(m - 1);
        }
    }
    /**
    revert
        The transaction has been reverted to the initial state.
    Reason provided by the contract: "函数禁止重入".
    */
}
