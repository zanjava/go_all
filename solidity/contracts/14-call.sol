// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract A {
    // payable可以接收value
    function add(uint256 a) external payable returns (uint256) {
        return a;
    }

    // 函数重载。函数签名是keccak256("add(uint256,uint256)")
    function add(uint256 a, uint256 b) external payable returns (uint256) {
        return a + b;
    }

    event Log(string msg);

    // 如果外界调用的函数不存在，会执行fallback
    fallback() external {
        emit Log("fallback method mot found!!!");
    }
}

// 调用其他合约
contract B {
    event Sig(bytes msg);
    event Sel(bytes msg);

    // 第一种调用方式
    function callAdd1(
        address contractAddress,
        uint256 a
    ) external payable returns (uint256) {
        // 把地址转成合约，再调合约里的函数
        return A(contractAddress).add(a);
        // 调用函数的时候还可以指定value(向对方转账)和gas
        // return A(contractAddress).add{value: msg.value, gas: 5000}(a);  //这里指定了value，则函数必须是payable
    }

    // 第二种调用方式。测试时contractName直接赋contractAddress即可
    function callAdd2(
        A contractName,
        uint256 a
    ) external payable returns (uint256) {
        return contractName.add{value: msg.value, gas: 5000}(a);
    }

    // 第三种调用方式，low level方式--call
    function callAdd3(
        address contractAddress,
        uint256 a,
        uint256 b
    ) external payable returns (uint256) {
        bytes memory signature = abi.encodeWithSignature( //通过签名选择函数
                "add(uint256,uint256)",
                a,
                b
            );
        emit Sig(signature);
        bytes memory selector = abi.encodeWithSelector( //或者通过函数选择器(selector)，取哈希之后的前4个字节
                bytes4(keccak256("add(uint256,uint256)")),
                a,
                b
            );
        emit Sel(selector);

        (bool success, bytes memory data) = contractAddress.call{
            value: msg.value
        }(
            //signature和selector实际上是相等的
            signature //如果函数签名不存在，会命中对方的fallback
            // selector
        );
        if (success) {
            return abi.decode(data, (uint256)); // 对结果进行反解
        } else {
            return 0;
        }
    }
}
