pragma solidity 0.6.10;

import {SafeMath} from "@openzeppelin/contracts/math/SafeMath.sol";
import {TimeLockedToken} from "./TimeLockedToken.sol";

/**
 * @dev The DcloToken contract is a claimable contract where the
 * owner can only mint or transfer ownership. 
 * Tolerates dilution to slash stake and accept rewards.
 */
contract DcloToken is TimeLockedToken {
    using SafeMath for uint256;

    uint256 constant MAX_SUPPLY = 145000000000000000;
  // for users who don't want refund to their address but keep the balance in sharedVault
    mapping(address => uint256) public sharedVault;
    mapping(uint256 => (address, address, uint) public gameSessionInfo;
 
    function _transfer(
        address _from,
        address _to,
        uint256 _amount
    ) internal override {
        // check if recipient is not the contract itself
        require(_to != address(this), "Can't transfer to the contract itself");
        super._transfer(_from, _to, _amount);
    }

    /**
     * This is necessary to set ownership for proxy
     */
    function initialize() public {
        require(!initalized, "already initialized");
        owner_ = msg.sender;
        initalized = true;
    }

    function depositVault(uint256 _amount) external {
        assert(transfer(address(this), _amount));
        sharedVault[msg.sender] += _amount;
    }

    function withdrawVault(uint256 _amount) external {
        require(sharedVault[msg.sender] >= _amount, "Insufficient balance");
        _transfer(address(this), msg.sender, _amount);
        sharedVault[msg.sender] -= _amount;
    }

    // Contract addresses only holds balances deposited by players. For airdrops, we must extract tokens from owner address
    function depositGameSession(uint256 _amount, address _streamer, uint256 gameSessionId) external {
        if (!gameSessionId[gameSessionId].isValue) {
            if (sharedVault[msg.sender] < _amount) {
                assert(transfer(address(this), _amount - sharedVault[msg.sender]));
                sharedVault[msg.sender] = _amount;
            }
            gameSessionInfo[gameSessionId] = (_streamer, msg.sender, _amount);
        }
    }

    function releaseGameSession(uint256 _amount, uint256 gameSessionId, bool fromVault) external onlyOwner {
        if (gameSessionId[gameSessionId].isValue) {
            streamer = gameSessionId[gameSessionId][0];
            player = gameSessionId[gameSessionId][1];
            maxAmount = gameSessionId[gameSessionId][2];
            require(_amount <= maxAmount, "Cannot claim more than initial deposit");
            assert(_transfer(address(this), streamer, _amount));
            // TODO: prevent potential issue
            sharedVault[player] -=  _amount;
        }
    }

    /**
     * Can never mint more than MAX_SUPPLY = 1.45 billion
     */
    function mint(address _to, uint256 _amount) external onlyOwner {
        if (totalSupply.add(_amount) <= MAX_SUPPLY) {
            _mint(_to, _amount);
        } else {
            revert("Max supply exceeded");
        }
    }

    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }

    function decimals() public override pure returns (uint8) {
        return 8;
    }

    function rounding() public pure returns (uint8) {
        return 8;
    }

    function name() public override pure returns (string memory) {
        return "DCloud Gaming";
    }

    function symbol() public override pure returns (string memory) {
        return "DCLO";
    }
}