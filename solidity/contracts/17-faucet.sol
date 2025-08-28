// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// @openzeppelin表示一个npm库
import "../node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol"; //实现了IERC20接口

// import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";  //可销毁(燃烧)的代币
// import "@openzeppelin/contracts/access/Ownable.sol";//可增发的代币
// import "@openzeppelin/contracts/access/AccessControl.sol";  //通过角色控制权限
// import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Capped.sol";//有总量上限
// import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Pausable.sol";//可暂停
// import "@openzeppelin/contracts/finance/VestingWallet.sol";//可锁仓
// 直接从github上引用
// import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/IERC20.sol";
// import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20.sol";

// abstract contract MyToken is IERC20 {}

contract MyFixedSupplyToken is ERC20 {
    constructor(
        string memory name,
        string memory symbol,
        uint256 totalSupply
    ) ERC20(name, symbol) {
        _mint(msg.sender, totalSupply); //铸币
    }
}

error HighFreqError(address account, uint256 lastTime, uint256 currentTime);

// ERC20代币的水龙头合约。限制代币的领取频次和单次领取额度
contract MyFaucet {
    // /ˈfɔːsɪt/
    address public tokenContract; // token合约地址
    uint256 public amountAllowed = 10; // 每次领取代币数，是个固定值
    mapping(address => uint256[]) public historyRequested; // 记录每个账号历次领取的时间，1天之内最多只能领取一次

    // SendToken事件
    event SendToken(address indexed Receiver, uint256 indexed Amount);

    // 部署时设定ERC20合约地址（比如上面那个MyFixedSupplyToken的合约地址）
    constructor(address _tokenContract) {
        tokenContract = _tokenContract;
    }

    // 用户领取代币函数
    function requestToken() external {
        uint256[] storage historyTimes = historyRequested[msg.sender]; //在查询mapping时会创建任何可能的key，并将value赋为零值。
        // 如果当前时间跟上一次的领取时间间隔少于一天，则回退。注意block.timestamp是当前块的打包时间，获取不到准确的当前时间
        // if (historyTimes.length>0 && block.timestamp - historyTimes[historyTimes.length-1] <= 1 days){
        if (
            historyTimes.length > 0 &&
            block.timestamp - historyTimes[historyTimes.length - 1] <= 3 //为方便测试，把领取间隔阈值调小为3秒
        ) {
            revert HighFreqError(
                msg.sender,
                historyTimes[historyTimes.length - 1],
                block.timestamp
            );
        }

        ERC20 token = ERC20(tokenContract); // 创建合约对象。把子合约的地址转成父合约的对象。不用写死为MyFixedSupplyToken，这样MyFaucet也可以作为其他ERC20合约的水龙头
        require(
            token.balanceOf(address(this)) >= amountAllowed, //水龙头合约的余额少于申领的数量
            unicode"水龙头余额不足！"
        );

        token.transfer(msg.sender, amountAllowed); // 发送token。注意：代币是从谁的账户上扣除的？谁调transfer()函数从谁的账户上扣除。那是msg.sender在调transfer()函数吗？不是。从水龙头合约里面调ERC20合约里的transfer()函数，则函数的调用方是水龙头合约，即代币是从水龙头合约的账户上扣除的。如果希望transfer()函数的调用是msg.sender则需要使用delegatecall
        historyTimes.push(block.timestamp); // 记录领取时间
        // historyRequested[msg.sender] = historyTimes; //不需要把historyTimes重新放回map里，因为当初取出来的value是storage，是指针，而且就算key之前不存在，根据key取value的时候已经自动把key加到map里了

        emit SendToken(msg.sender, amountAllowed); // 记录SendTOken事件
    }
}

// 部署完这两个合约后，先通过ERC20代币合约transfer代币到水龙头合约，然后其他用户就可以从水龙头合约去领取代币
