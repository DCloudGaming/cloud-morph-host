ClaimableContract.sol: Authentication baseline, check owner
ERC20.sol: Token baseline, maintain cap, current supply, balance map
TimeLockedToken.sol: Logic for locked tokens and over-time linear unlock
TimeLockedRegistry.sol: Holds treasury/ vault with locked strategy and vested/ released overtime
AirdropRegistry.sol: Holds treasury for airdrop/ release airdrop rewards.
DcloToken.sol: main entry of contract

timelock_and_airdrop_deploy.ts : js utils to deploy contracts to test net
