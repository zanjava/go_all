// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0;

contract Ballot {
    // 选民
    struct Voter {
        uint256 weight; // 该选民手上的票权重是多少
        bool voted; // 是否已投票
        uint256 vote; // 投给了第几个候选人
    }

    // 提案
    struct Proposal {
        uint256 voteCount;
    }

    address chairPerson; //主持人
    mapping(address => Voter) voters; // 选民
    Proposal[] proposals; //所有候选人的得票情况。数组类型

    enum Phase {
        Init,
        Regs,
        Vote,
        Done
    } //枚举，内部编码为0、1、2、3。Regs登记
    Phase public state = Phase.Init; //投票活动现处于哪个阶段

    // 构造函数(部署合约时会自动调用构造函数)。初始化三个状态变量
    constructor(uint256 numProposals) {
        require(
            numProposals >= 3,
            "candidate count should not less than three"
        );
        chairPerson = msg.sender; // 创建合约的账户即为主持人
        voters[chairPerson].weight = 2; //主持人也是选民，并且他的权重为2
        for (uint256 prop = 0; prop < numProposals; prop++) {
            proposals.push(Proposal(0)); //push往不定长数组里添加元素。Proposal(0)初始化了一个结构体
        }
        state = Phase.Regs; // 直接进入登录选民的阶段
    }

    function changeState(Phase x) public onlyChair {
        // 只有chairPerson可以改变状态
        // require(msg.sender == chairPerson, "only chairPerson can change state");
        if (x < state) revert(); //状态只能前进，不能后退
        state = x;
    }

    modifier validPhase(Phase x) {
        require(state == x);
        _;
    }

    modifier onlyChair() {
        require(msg.sender == chairPerson, "only chairPerson can change state");
        _;
    }

    // 登记选民
    function register(address voter) public validPhase(Phase.Regs) {
        // chairPerson才有登记权。如果已投过票了则直接返回（否则会把他标记为未投票状态）
        if (msg.sender != chairPerson || voters[voter].voted) {
            revert(); //回退
        }
        voters[voter].weight = 1; // 普通选民权重为1
        voters[voter].voted = false; // 初始为未投票状态
    }

    // 给候选人投票
    function vote(uint256 candidate) public validPhase(Phase.Vote) {
        Voter storage sender = voters[msg.sender];
        if (sender.voted || candidate >= proposals.length) {
            revert();
        }
        sender.vote = candidate;
        sender.voted = true;
        proposals[candidate].voteCount += sender.weight;
    }

    // 获得最终赢得选举的人
    function getWinner() public view validPhase(Phase.Done) returns (uint256) {
        uint256 maxCount = 0; //谁票数最多，谁赢得选举
        uint256 winner = proposals.length;
        for (uint256 i = 0; i < proposals.length; i++) {
            if (proposals[i].voteCount > maxCount) {
                maxCount = proposals[i].voteCount;
                winner = i;
            }
        }
        // 一般用require来验证数据、计算和参数值，而assert来处理异常，assert失败要比require回滚浪费更多的区块链执行成本
        assert(maxCount >= proposals.length / 2); //必须获得超过半数的投票，否则回退交易
        return winner;
    }
}
