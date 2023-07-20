// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFTCollection is ERC721URIStorage, Ownable {
    uint256 public tokenIdCounter = 0;

    event CollectionCreated(address indexed collection, string name, string symbol);
    event TokenMinted(address indexed collection, address indexed recipient, uint256 indexed tokenId, string tokenUri);

    constructor(string memory name, string memory symbol) ERC721(name, symbol) {
        emit CollectionCreated(address(this), name, symbol);
    }

    function mintToken(address recipient, string memory tokenUri) public onlyOwner {
        uint256 newTokenId = tokenIdCounter;
        _mint(recipient, newTokenId);
        _setTokenURI(newTokenId, tokenUri);
        tokenIdCounter++;

        emit TokenMinted(address(this), recipient, newTokenId, tokenUri);
    }
}

