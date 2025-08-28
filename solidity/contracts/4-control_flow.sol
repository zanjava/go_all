// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract ControlFlow {
    function IF() public pure returns (int) {
        int a = 9;
        if (a > 10) return 1;
        // 只有一行时可以不加{}
        else return 0;
    }

    function FOR() public pure returns (int) {
        int sum = 0;
        int[3] memory arr = [int(1), 4, 7]; //=右边的数组是在函数内部创建的，是局部变量，所以=左边必须是memory
        for (uint i = 0; i < arr.length; i++) sum += arr[i];
        return sum;
    }

    function WHILE() public pure returns (int) {
        int sum = 0;
        int[3] memory arr = [int(1), 4, 7];
        uint i = 0;
        while (i < arr.length) {
            //先判断再执行
            sum += arr[i];
            i++;
        }
        return sum;
    }

    function DO_WHILE() public pure returns (int) {
        int sum = 0;
        int[3] memory arr = [int(1), 4, 7];
        uint i = 0;
        do {
            //先执行再判断
            sum += arr[i];
            i++;
        } while (i < arr.length);
        return sum;
    }

    // 三元运算符。当条件满足时返回前者(:前面的)，否则返回后者
    function ternary() public pure returns (uint) {
        return 9 > 10 ? 0 : 1;
    }
}
