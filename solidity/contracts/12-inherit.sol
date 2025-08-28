// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract A {
    // 父类中的Public和internal成员可以被子类继承。
    string public name;
    event Log(string msg); //事件可以被继承

    constructor(string memory n) {
        name = n;
    }

    // 父类中标记为virtual的函数可被子类重写
    function foo() public virtual {
        emit Log("A.foo called");
    }
}

contract C is
    A // is表示C继承了A
{
    string public title;

    // 先调用父类的构造函数，再执行自己的构造函数
    constructor(string memory t, string memory n) A(n) {
        title = t;
    }
}

contract D is
    A("gw") // is表示C继承了A。直接在合约声明的时候调用父类的构造函数，把参数写死
{
    string public title;

    constructor(string memory t) {
        title = t;
    }
}

contract B {
    string public city;
    event Log2(string msg);

    constructor(string memory c) {
        city = c;
    }

    function foo() public virtual {
        emit Log2("B.foo called");
    }
}

contract E is
    B,
    A //继承多个父类。构造函数的调用顺序 B->A->E
{
    string public title;

    // 在构造函数这儿调用父类的构造函数
    constructor(string memory t, string memory n, string memory c) A(n) B(c) {
        title = t;
    }

    // 在override()里指定要重写哪个（或哪几个）父类的方法。重写不能改变方法的可见性
    // 由于E同时继承了B和A，而B和A里都有foo()，则在E里必须重写foo()，否则在E里调用foo()时不知道该调谁的foo()
    function foo() public virtual override(A, B) {
        // 同时它自己又被声明为virtual，表示它自己还可以被子类重写(override)
        emit Log("E.foo called");
    }

    function think() public {
        foo(); //调用自己的foo()函数
        super.foo(); //调用最近的父合约函数，即A的
        B.foo(); // 显示指定调用哪个父类的函数
    }

    // modifier也可以被继承
    modifier checkArg(int256 arg) virtual {
        require(arg > 0, "arg must more than zero");
        _;
    }
}

contract F is E("CEO", "zcy", "ZZ") {
    function foo() public override(E) {
        // F的直接父类只有一个E，所以这里的(E)可以不写
        emit Log("F.foo called");
    }

    function goo() public {
        super.foo(); // 通过super指定调用父类的函数
        think();
    }

    // modifier也可以被override
    modifier checkArg(int256 arg) override {
        require(arg > 10, "arg must more than 10");
        _;
    }
}
