// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";

// NFT, Non-fungible tokens, 非同质化代币
contract MyNFT is ERC721 {
    constructor(
        string memory name,
        string memory symbol,
        uint256 tokenId
    ) ERC721(name, symbol) {
        _mint(msg.sender, tokenId); //铸造一枚NFT，归msg.sender所有
    }
}

// 如果safeTransferFrom的接收方是合约地址，则要求该合约实现IERC721Receiver接口
contract R is IERC721Receiver {
    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4) {}
}
