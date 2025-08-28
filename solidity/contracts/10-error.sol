// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// 定义异常。可以放在合约外面，也可以放在合约里面
error DivError(uint256 a, uint256 b);

contract Math {
    uint256 public s = 0;

    function Div(uint256 a, uint256 b) public returns (uint256) {
        s++;
        if (b == 0) {
            // 异常必须搭配revert（交易回退）来使用
            revert DivError(a, b);
            // revert();  //revert()可以单独使用
        }
        return a / b;
    }

    // require携带异常描述，这个string越长消耗的gas就越多
    function Div2(uint256 a, uint256 b) public returns (uint256) {
        s++;
        require(b > 0, "b must more than 0"); //require如果条件不满足就会抛出异常，异常会导致交易回退revert
        return a / b;
    }

    // assert不携带异常描述
    function Div3(uint256 a, uint256 b) public pure returns (uint256) {
        assert(b > 0); //assert如果条件不满足就会抛出异常，异常会导致交易回退revert
        return a / b;
    }
}

// require 应该被用于函数中检查条件，assert 用于预防不应该发生的情况，即不应该使条件错误。
