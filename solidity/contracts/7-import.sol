// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// import "./7-library.sol";   //会把这个文件里的所有library和contract全都引进来
import {Strings, A} from "./7-library.sol"; //显式指定引入的library或contract

// import的源文件也可以来自于网络，如import 'https://xxxxx/xxx.sol';

contract Test {
    function goo() public pure returns (string memory) {
        uint256 a = 123456;
        return Strings.toHexString(a, 20);
    }
}
