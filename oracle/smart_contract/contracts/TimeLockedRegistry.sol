pragma solidity 0.6.10;

import {ClaimableContract} from "./common/ClaimableContract.sol";

import {TimeLockedToken} from "./TimeLockedToken.sol";

/**
 * @dev This contract allows owner to register distributions for a TimeLockedToken
 *
 * To register a distribution, register method should be called by the owner.
 * claim() should then be called by account registered to recieve tokens under lockup period
 * If case of a mistake, owner can cancel registration
 */
 
contract TimeLockRegistry is ClaimableContract {
    // time locked token
    TimeLockedToken public token;

    mapping(address => uint256) public registeredDistributions;

    event Register(address receiver, uint256 distribution);
    event Cancel(address receiver, uint256 distribution);
    event Claim(address account, uint256 distribution);

    function initialize(TimeLockedToken _token) external {
        require(!initalized, "Already initialized");
        token = _token;
        owner_ = msg.sender;
        initalized = true;
    }

    function register(address receiver, uint256 distribution) external onlyOwner {
        require(receiver != address(0), "Zero address");
        require(distribution != 0, "Distribution = 0");
        require(registeredDistributions[receiver] == 0, "Distribution for this address is already registered");

        // register distribution in mapping
        registeredDistributions[receiver] = distribution;

        // transfer tokens from owner
        require(token.transferFrom(msg.sender, address(this), distribution), "Transfer failed");

        // emit register event
        emit Register(receiver, distribution);
    }

    function cancel(address receiver) external onlyOwner {
        require(registeredDistributions[receiver] != 0, "Not registered");

        uint256 amount = registeredDistributions[receiver];

        // set distribution mapping to 0
        delete registeredDistributions[receiver];

        require(token.transfer(msg.sender, amount), "Transfer failed");

        emit Cancel(receiver, amount);
    }

    function claim() external {
        require(registeredDistributions[msg.sender] != 0, "Not registered");

        uint256 amount = registeredDistributions[msg.sender];

        delete registeredDistributions[msg.sender];

        // register lockup in TimeLockedToken
        // this will transfer funds from this contract and lock them for sender
        token.registerLockup(msg.sender, amount);

        emit Claim(msg.sender, amount);
    }
}