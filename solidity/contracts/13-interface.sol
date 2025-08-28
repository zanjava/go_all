// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

interface Animal {
    // 接口没有状态变量，没有构造函数

    // 接口中的函数必须为external，不能有函数体
    function Eat() external returns (int256);

    function Drink() external returns (int256);
}

// 接口可以继承另一个接口
interface Human is Animal {
    function Think() external returns (int256);
}

// 合约“实现”接口也是用is
// 由于没有把接口里的函数全部实现，所以合约必须用abstract修饰。接口和抽象合约不能部署
abstract contract Female is Human {
    function Eat() public pure returns (int256) {
        return 1;
    }
}

// 这里的is既是“实现”也是“继承”，注意先后顺序：接口(Human)在前，合约(Female)在后，Human->Female->Male
contract Male is Human, Female {
    function Drink() public pure returns (int256) {
        // 在合约里面，函数的可见性没必须还是external，可以是external或public，不能是private或internal，因为要确保外部可以调用
        return 1;
    }

    function Think() external pure returns (int256) {
        return 1;
    }
}
