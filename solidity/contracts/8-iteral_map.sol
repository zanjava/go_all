// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// 可迭代的mapping。以key是string,value是uint16为例
library IterableMapping {
    struct KeyFlag {
        string key;
        bool deleted; // 默认情况下为false
    }

    struct IndexValue {
        uint256 keyIndex; //下标从1开始
        uint16 value;
    }

    struct itmap {
        /**
        原生的mapping是 
        key-->value
        "100"-->18
        "200"-->45

        data是
                keyIndex,  value
        data["100"]={1,      18}
        data["200"]={2,      45}

        keys是
        keys=[{"100", false}, {"200", false}]
        */
        mapping(string => IndexValue) data; // 在mapping中并不存储原始的key，而是存储key的keccak256哈希值，所以不存在key的集合的概念，也不能遍历mapping。在查询mapping时会创建任何可能的key，并将value赋为零值，所以mapping没有长度，不支持检查某个key是否存在，不支持删除某个key。
        KeyFlag[] keys; // 这个数组存在的意义是：将来可以遍历map
        uint256 size;
    }

    function set(
        itmap storage self, //虽然函数是public，但库函数的第一个参数只能是storage
        string memory key,
        uint16 value
    ) public returns (bool replaced) {
        uint256 keyIndex = self.data[key].keyIndex; // mapping的key不存在时value默认为0值
        self.data[key].value = value;
        if (keyIndex > 0) {
            // 不为0值，说明key之前就是存在的
            return true; //覆盖了之前的值
        } else {
            // key之前不存在
            self.keys.push(KeyFlag({key: key, deleted: false}));
            self.data[key].keyIndex = self.keys.length;
            self.size++;
            return false;
        }
    }

    function get(
        itmap storage self,
        string memory key
    ) public view returns (uint16 value) {
        value = self.data[key].value;
    }

    function remove(
        itmap storage self,
        string memory key
    ) public returns (bool) {
        uint256 keyIndex = self.data[key].keyIndex;
        if (keyIndex == 0) {
            return false; // key不存在，没有删除任何元素
        } else {
            delete self.data[key]; //从mapping里删除（把相应的value置0）
            self.keys[keyIndex - 1].deleted = true; //在keys数组里标记为已删除（逻辑删除）
            self.size--;
            return true;
        }
    }

    /**
    遍历map需要使用以下3个函数
    */
    // prevIndex是数组keys的下标，从keyIndex往后找第一个deleted=false的元素的下标
    function iterate_next(
        itmap storage self,
        int256 keyIndex
    ) public view returns (uint256) {
        keyIndex++; //++后可能会越界，所以iterate_next()返回的可能是一个非法的index
        assert(keyIndex >= 0);
        while (
            uint256(keyIndex) < self.keys.length &&
            self.keys[uint256(keyIndex)].deleted
        ) {
            keyIndex++;
        }
        return uint256(keyIndex);
    }

    // 从数组里keys里找第一个deleted=false的元素的下标
    function iterate_start(itmap storage self) public view returns (uint256) {
        return iterate_next(self, int256(-1));
    }

    // 循环退出条件
    function iterate_valid(
        itmap storage self,
        uint256 keyIndex
    ) public view returns (bool) {
        return keyIndex < self.keys.length;
    }
}

contract B {
    IterableMapping.itmap public data;
    using IterableMapping for IterableMapping.itmap;
    event Log(string key, uint16 value);

    function set(string memory k, uint16 v) public returns (uint256) {
        data.set(k, v);
        return data.size;
    }

    function remove(string memory k) public {
        data.remove(k);
    }

    function get(string memory k) public view returns (uint16) {
        return data.get(k);
    }

    function iterate() public {
        uint256 keyIndex = data.iterate_start();
        while (data.iterate_valid(keyIndex)) {
            string memory key = data.keys[keyIndex].key;
            uint16 value = get(key);
            emit Log(key, value);
            keyIndex = data.iterate_next(int256(keyIndex));
        }
    }
}
