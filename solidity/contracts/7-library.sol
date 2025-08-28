// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

library Strings {
    bytes16 private constant _HEX_SYMBOLS = "0123456789abcdef";

    // 把数字转成16进制的字符串
    function toHexString(
        uint256 value,
        uint256 length
    ) public pure returns (string memory) {
        bytes memory buffer = new bytes(2 * length + 2);
        buffer[0] = "0";
        buffer[1] = "x";
        for (uint256 i = 2 * length + 1; i > 1; --i) {
            buffer[i] = _HEX_SYMBOLS[value & 0xf];
            value >>= 4;
        }
        require(value == 0, "Strings: hex length insufficient");
        return string(buffer);
    }
}

contract A {
    function moo(uint256 x) private pure {}

    function goo() private pure returns (string memory) {
        uint32 a = 123456;
        //虽然moo()接收uint256,但是可以传uint32
        moo(a); // 由于goo()调用了moo(),如果moo()不是pure，则goo()也不能是pure。
        // 使用library的第一种方式：直接加前缀library.调用库里的函数、状态变量、结构体等
        return Strings.toHexString(a, 20);
    }

    // 使用library的第二种方式：using Strings for dataType，相当于把dataType作为库函数的第一个参数，类似于python类方法的第一个参数self
    using Strings for uint256;

    function foo() private pure returns (string memory) {
        uint256 a = 123456;
        return a.toHexString(20);
    }
}
