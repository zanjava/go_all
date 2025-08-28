// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// 空投合约。事先需要获得原始代币持有方的授权approve
contract MyAirdrop {
    function multiTransferToken(
        address tokenContract, //某个ERC20代币合约的地址
        address[] calldata addresses, //给这些账户发放代币
        uint256[] calldata amounts //每个账户对应的代币数量
    ) external {
        require(
            addresses.length == amounts.length,
            unicode"addresses和amounts的数组长度不相等"
        );
        ERC20 token = ERC20(tokenContract);
        uint256 amountSum = getSum(amounts); //计算需要发放出去的总量
        require(
            token.allowance(msg.sender, address(this)) >= amountSum, //原始代币的持有方给当前空投合约授予approve过多少代币。msg.sender必须是原始代币的持有方
            unicode"总发放量超过了授权量"
        );

        for (uint8 i = 0; i < addresses.length; i++) {
            token.transferFrom(msg.sender, addresses[i], amounts[i]); //msg.sender必须是原始代币的持有方。调用transferFrom()函数的是当前空投合约，不是msg.sender；如果想让msg.sender调transferFrom()需要使用delegatecall
        }
    }

    // 数组求和。calldata类似于memory，且函数内不能修改calldata参数
    function getSum(uint256[] calldata arr) public pure returns (uint256 sum) {
        for (uint256 i = 0; i < arr.length; i++) {
            sum = sum + arr[i];
        }
        // 已经声明了sum，不需要显式return了
    }
}
