// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// 在线拍卖和去中心化市场是区块链的理想应用场景
// 盲拍--投标期间发送的是出价的哈希值
contract BlindAuction {
    event Log(string name, bytes32 value);
    struct Bid {
        bytes32 blindBid; //出价的哈希值。keccak256哈希的结果是32B
        uint256 deposit; //押金>=出价
    }

    enum Phase {
        Init,
        Bidding, //出价阶段
        Reveal, //公布真实出价
        Done
    } //对应的内部编码是0、1、2、3
    Phase public currentPhase = Phase.Init;
    event BiddingStarted();
    event RevealStarted();
    event AuctionEnded(address winner, uint256 highestBid);

    address payable beneficiary; //受益人，即他要拍卖东西
    bool paid;
    mapping(address => Bid) bids; //记录每个竞拍者的出价情况

    uint256 public highestBid = 0; //最高出价为多少
    address public highestBidder; //出价最高者是谁

    mapping(address => uint256) depositReturns;

    //修饰符
    modifier validPhase(Phase phase) {
        require(currentPhase == phase);
        _; // 被修饰的函数内联到这里
    }
    modifier onlyBeneficiary() {
        require(msg.sender == beneficiary);
        _;
    }

    constructor() {
        beneficiary = payable(msg.sender); // 把address转成address payable
        currentPhase = Phase.Bidding;
    }

    // 竞拍阶段只能一步一步往前推进
    function advancePhase() public onlyBeneficiary {
        if (currentPhase == Phase.Done) {
            currentPhase = Phase.Init; //如果已经结束了，则可以从头再来（方便测试）
        } else {
            currentPhase = Phase(1 + uint(currentPhase)); //枚举和uint互相转换
        }
        if (currentPhase == Phase.Bidding) emit BiddingStarted();
        else if (currentPhase == Phase.Reveal) emit RevealStarted();
        else if (currentPhase == Phase.Done)
            emit AuctionEnded(highestBidder, highestBid);
    }

    // 盲拍，出价是哈希之后的。调用该函数时，外部账户需要给当前的合约账户转让押金(payable)
    function bid(bytes32 blindBid) public payable validPhase(Phase.Bidding) {
        bids[msg.sender] = Bid({
            blindBid: blindBid, //哈希之后的出价
            deposit: msg.value //转给合约的钱就是押金
        });
    }

    // 竞拍者公布自己的真实出价
    function reveal(
        uint256 value,
        bytes32 secret
    ) public validPhase(Phase.Reveal) {
        bytes32 digest = keccak256(abi.encodePacked(value, secret)); // 不是直接对value进行哈希，而是会再引入另外一个256位的变量secret。如果只是简单的对value施行keccak256哈希，则很容易被暴力破解，因为拍卖的价格范围是很小的。
        if (digest == bids[msg.sender].blindBid) {
            // 验证通过
            assert(bids[msg.sender].deposit >= value); // 押金>=出价
            depositReturns[msg.sender] = bids[msg.sender].deposit;
            if (value > highestBid) {
                // 记录谁是出价最高者
                highestBid = value;
                highestBidder = msg.sender;
            }
        } else {
            // 验证不通过就不退还押金
            emit Log("digest", digest);
            emit Log("expected", bids[msg.sender].blindBid);
        }
    }

    // 每个投标人都需要调一下withDraw()，拿走自己的押金
    function withDraw(address payable bidder) public validPhase(Phase.Done) {
        uint256 value = depositReturns[bidder];
        if (value > 0) {
            if (bidder != highestBidder) {
                bidder.transfer(value);
            } else {
                bidder.transfer(value - highestBid); // 中标者出价的部分不能拿走
            }
            depositReturns[bidder] = 0; //押金清0，防止重复调用withDraw()
        }
    }

    // 把钱付给受益人
    function payForBeneficiary() public validPhase(Phase.Done) onlyBeneficiary {
        if (!paid) {
            beneficiary.transfer(highestBid);
            paid = true; //防止重复调用payForBeneficiary()
        }
    }
}

/**
测试流程
step 0. 以第一个ACCOUNT部署合约，它是beneficiary
step 1. ACCOUNT选中第二个，调bid()函数，参数为 0x132b650dab39844eee3aa1523da44a8f3b1e4cce52ba190c003f066a2bdba009 ，VALUE置为100
step 2. ACCOUNT选中第三个，调bid()函数，参数为 0xc1c98903cc546809ece3df6a04ff468df4e9428cbc17e9bdd89b42a4a563a4c3 ，VALUE置为100
step 3. ACCOUNT选中第一个，调changeState()函数，参数为2
step 4. ACCOUNT选中第二个，调reveal()函数，参数为value = 82, secret = 0x00000000000000000000000000000000000000000000000000000000000007ea （可以把value或secret故意写错，看看有没有emit Log）
step 5. ACCOUNT选中第三个，调reveal()函数，参数为value = 57, secret = 0x00000000000000000000000000000000000000000000000000000000000007e9
step 6. ACCOUNT选中第一个，调changeState()函数，参数为3
step 7. 调highestBid()函数，确认为82
step 8. highestBidder()函数，确认为第二个ACCOUNT
step 9. 调payForBeneficiary()函数
step 10. 调withDraw()函数，参数为第二个ACCOUNT
step 11. 调withDraw()函数，参数为第三个ACCOUNT
*/
