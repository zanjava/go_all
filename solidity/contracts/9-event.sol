// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

contract Event {
    // indexed参数最多只能有3个。通过使用indexed修饰符，可以为事件参数创建索引，以提高事件的可搜索性和可过滤性。这在开发DApp、进行日志分析和数据查询时非常有用。
    // topic是事件的签名--keccak256("Hello(address,string,bool,uint256)")，注意中间没有空格
    // 消耗gas: 状态变量>事件(indexed参数>非indexed参数)
    event Hello(address indexed adr, string b, bool indexed c, uint256 d);
    // 事件可以被继承

    function EmitEvent() external {
        // 通过发射事件把数据存储到交易的logs里。注意：事件的logs在合约里不能被访问，但是可以被订阅者读取
        emit Hello(msg.sender, "golang", true, uint256(4));
        // 在DAPP的应用中，如果监听了某事件，当事件发生时，会进行回调。
    }
}

/**
logs: [
	{
		"from": "0xe2899bddFD890e320e643044c6b95B9B0b84157A",        --合约地址
		"topic": "0x94fcec6ddb0f29943ad59ae3991562f81e8a7886413b8ded16dbf5d5e556b865",       --keccak256("Hello(address,string,bool,uint256)")
		"event": "Hello",    --事件名称
		"args": {            --事件参数
			"0": "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4",
			"1": "golang",
			"2": true,
			"3": "4",
			"adr": "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4",
			"b": "golang",
			"c": true,
			"d": "4"
		}
	}
]
 */
