/**
 * PRIVATE_KEY={private_key} ts-node scripts/timelock_deploy.ts "{network}"
 */
 import { ethers, providers } from 'ethers'

 import {
   DcloToken__factory,
   TimeLockRegistry__factory,
   AirdropRegistry__factory
 } from 'contracts'
 
 async function deployTimeLockAndAirdropRegistry () {
   const txnArgs = { gasLimit: 5_000_000, gasPrice: 16_000_000_000 }
   const provider = new providers.InfuraProvider(process.argv[2], 'e33335b99d78415b82f8b9bc5fdc44c0')
   const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider)
 
   const dcloTokenImpl = await (await new DcloToken__factory(wallet).deploy(txnArgs)).deployed()
   console.log(`DcloToken Impl at: ${dcloTokenImpl.address}`)
   const timeLockRegistry = await (await new TimeLockRegistry__factory(wallet).deploy(txnArgs)).deployed()
   const airdropRegistry = await (await new AirdropRegistry__factory(wallet).deploy(txnArgs)).deployed()
   await(await timeLockRegistry.initialize(dcloTokenImpl.address, txnArgs)).wait()
   await(await airdropRegistry.initialize(dcloTokenImpl.address, txnArgs)).wait()
 }
 
 deployTimeLockAndAirdropRegistry().catch(console.error)
