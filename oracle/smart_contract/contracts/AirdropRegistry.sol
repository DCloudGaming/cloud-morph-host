pragma solidity 0.6.10;

import {ClaimableContract} from "./common/ClaimableContract.sol";

import {TimeLockedToken} from "./TimeLockedToken.sol";

contract AirdropRegistry is ClaimableContract {
    TimeLockedToken public token;
    uint                   public maxAirdropBalance = 0;
    uint                   public remainAirdropBalance = 0;

    mapping(address => uint256) public airdropDistribution;

    event DropToken(address receiver, uint256 distribution);

    function initialize(TimeLockedToken _token) external {
        require(msg.sender == _token.owner, "Airdrop must be created by token contract owner");
        token = _token;
        owner = _token.owner;
        initialized = true;
    }    
    
    // Contract addresses only holds balances deposited by players. For airdrops, we must extract tokens from owner address
    function upgradeAirdropBalance(uint amount) external onlyOwner {
        require(token.mint(owner, amount));
        remainAirdropBalance += amount;
        maxAirdropBalance += amount;
    }

    // Contract addresses only holds balances deposited by players. For airdrops, we must extract tokens from owner address
    function airDrop(address[] recipients, uint[] amounts) external onlyOwner {
        for (uint i = 0; i < recipients.length; i++) {
            assert(remainAirdropBalance >= amounts[i]);
            assert(token.transfer(recipients[i],amounts[i]));
            remainAirdropBalance -= amounts[i];
        }
    }

}
