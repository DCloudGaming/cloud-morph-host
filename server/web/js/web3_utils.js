const getWeb3Sync = async () => {
    console.log("Get web3 Start1");
    if (window.ethereum) {
        try {
            const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
            // setAccounts(accounts);
            console.log("Get web3 Start2");
            console.log(accounts);
            if (accounts.length > 0) {
                window.currentEthAccount = accounts[0];
                $('#walletAddress').html(accounts[0]);
            }
        } catch (error) {
            if (error.code === 4001) {
                // User rejected request
            }
            // setError(error);
        }
    }
}
