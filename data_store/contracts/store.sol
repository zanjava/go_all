// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract DataStore {
    int num;
    event NumChange(int x);

    function setNum(int x) public {
        num = x;
        emit NumChange(num);
    }

    function getNum() public view returns (int) {
        return num;
    }
}
