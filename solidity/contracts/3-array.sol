// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Array {
    int256[3] arr; //固定长度

    //可变长度
    int256[] brr;
    int256[] crr = [int256(1), 2, 3]; //数组里元素的类型是以第一个元素为准的，默认是uint8
    int256[] drr = new int256[](3); //目前长度为3，里面有3个0

    function arrayOP() private returns (int256) {
        int256 sum = 0;
        // 只有可变长度的数组才能调用push和pop
        brr.push(1); // 向数组尾部追加元素
        brr.push(); // 默认追加0
        brr.push(2);
        brr.push(3);
        crr.push(3);
        drr.push(3);
        brr.pop(); // 移除尾部元素

        uint256 L = brr.length; // 定长和变量数组都可以调length
        for (uint256 i = 0; i < L; i++) {
            sum += brr[i];
        }

        return sum;
    }

    // 注意bytesN和bytesN不能进行加减运算！！！bytesN的N最多只到32
    bytes1 n = 0xA5; // 1个字节
    bytes2 m = 0xA5A5; // 2个字节
    bytes32 big =
        0xA5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5; //整好32个字节，不能多，也不能少！！！

    //变长字节数组（数组里每个元素是一个字节）
    bytes bs1 = "1234567890123456789012345678901234567890"; //变长字节数组，且数组里的元素可修改
    bytes bs2 = unicode"890😃语言"; //含有ASCII之外的编码则需要加unicode
    bytes bs3 = new bytes(3); //变长字节数组，里面有3个0，且数组里的元素可修改

    //变长字节数组（数组里每个元素是一个bytes，bytes[]相当于是二维的bytes）
    bytes[] bs4 = [bytes("abc"), "efg", "ge"]; //bytes等价于bytes1。注意不存在byte[]这种类型
    bytes32[] bs5 = [bytes32(unicode"华人"), unicode"中国人"]; //数组的每个元素是32字节，32字节足够表示一个汉字

    //string不能通过下标访问某个元素，也不能调.length获得长度(得先把string强制转成bytes才能获取其长度)
    string public ss1 = "1234567890123456789012345678901234567890";

    //定长字节数组
    bytes5[10] fbs;

    function bytesOP() external {
        m = 0x1234;
        fbs[0] = 0xA5A5A5A5A5; //必须是5个字节，不能多，也不能少！！！

        bs1.push(0xB5); //只能追加1个字节，注意不能追加一个汉字
        bs2.push(0xA5);
        bs3.push(0xA5);
        bs4.push(bytes("x"));
        bs5.push(bytes32("x"));

        bs1[0] = 0xC5;
        bs2[0] = 0xC5;
        bs3[0] = 0xC5;
        bs4[0] = bytes("x");
        bs5[0] = bytes32("x");

        //通过把string转为bytes storage，可修改string里的内容及长度
        //bytes和string都不是值类型。引用类型作为函数参数或局部变量，最好都显式指定storage或memory
        bytes storage bs7 = bytes(ss1);
        bs7[0] = "a";
    }
}
