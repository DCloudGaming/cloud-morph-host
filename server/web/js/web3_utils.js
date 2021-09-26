// TODO: Move to env file
const APP_BACKEND_URL = "http://localhost:8080/api";
const web3 = new Web3("https://cloudflare-eth.com");

const handleWalletClick = async () => {
    console.log("Get web3 Start1");
    if (window.ethereum) {
        try {
            const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
            // const accounts = await web3.eth.getAccounts();
            // setAccounts(accounts);
            console.log("Get web3 Start2");
            console.log(accounts);
            if (accounts.length > 0) {
                window.currentEthAccount = accounts[0];
                $('#walletAddress').html(accounts[0]);
                var user = await handleClick(window.currentEthAccount);
                // var signature = await handleSignMessage(user);

            }
        } catch (error) {
            if (error.code === 4001) {
                // User rejected request
            }
            // setError(error);
        }
    }
};

const handleClick = async (walletAddress) => {
    // --snip--
    var userExistResponse = await fetch(`${APP_BACKEND_URL}/users?` + new URLSearchParams({
            wallet_address: walletAddress
        }), {
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*'
        },
        method: 'GET'
    });

    if (userExistResponse) {
        return userExistResponse.json()
    } else {
        var user = await handleSignup(walletAddress)
        return user.json()
    }
};

const handleSignup = async (walletAddress) => {
    var res = await fetch(`${APP_BACKEND_URL}/users/signup`, {
        body: JSON.stringify({"wallet_address": walletAddress}),
        headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*'
        },
        method: 'POST'
    });
    return res;
};

const handleSignMessage = async (user) => {
    var walletAddress = user.WalletAddress;
    var nonce = user.nonce;
    var sig = await web3.eth.sign(web3.utils.sha3(`I am signing my one-time nonce: ${nonce}`), walletAddress);
    return sig;
};

const handleAuthenticate = ({ walletAddress, signature }) => {
    fetch(`${APP_BACKEND_URL}/users/auth`, {
        body: JSON.stringify({ "walletAddress": walletAddress, "signature": signature }),
        headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*'
        },
        method: 'POST'
    }).then(response => response.json());
};
