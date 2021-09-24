// TODO: Move to env file
const APP_BACKEND_URL = "http://127.0.0.1:8080/api";

const handleWalletClick = async () => {
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
                handleClick(window.currentEthAccount);
            }
        } catch (error) {
            if (error.code === 4001) {
                // User rejected request
            }
            // setError(error);
        }
    }
};

const handleClick = (walletAddress) => {
    // --snip--
    fetch(`${APP_BACKEND_URL}/users?` + new URLSearchParams({
            wallet_address: walletAddress
        }), {
        headers: {
            'Content-Type': 'application/json'
        },
        method: 'GET'
    }).then(response => response.json())
        // If yes, retrieve it. If no, create it.
        .then(
            users => (users.length ? users[0] : handleSignup(walletAddress))
        )
        // Popup MetaMask confirmation modal to sign message
        .then(handleSignMessage)
        // Send signature to back end on the /auth route
        .then(handleAuthenticate)
    // --snip--
};

const handleSignup = walletAddress => {
    fetch(`${APP_BACKEND_URL}/users/signup`, {
        body: JSON.stringify({"wallet_address": walletAddress}),
        headers: {
            'Content-Type': 'application/json'
        },
        method: 'POST'
    }).then(response => response.json());
};

const handleSignMessage = ({ walletAddress, nonce }) => {
    return new Promise((resolve, reject) =>
        web3.personal.sign(
            web3.fromUtf8(`I am signing my one-time nonce: ${nonce}`),
            walletAddress,
            (err, signature) => {
                if (err) return reject(err);
                return resolve({ walletAddress, signature });
            }
        )
    );
};

const handleAuthenticate = ({ walletAddress, signature }) => {
    fetch(`${APP_BACKEND_URL}/users/auth`, {
        body: JSON.stringify({ "walletAddress": walletAddress, "signature": signature }),
        headers: {
            'Content-Type': 'application/json'
        },
        method: 'POST'
    }).then(response => response.json());
};
