// include the ipc module to communicate with main process.
const ipcRenderer = require("electron").ipcRenderer;
const axios = require("axios");
axios.defaults.withCredentials = true;

const apppathText = document.getElementById("apppathText");

const registerButton = document.getElementById("registerButton");
registerButton.addEventListener("click", function () {
  var arg = "secondparam";

  //send the info to main process . we can pass any arguments as second param.
  // ipcRender.send will pass the information to main process. Here is event to open file dialog
  ipcRenderer.send("register", {});
});

const connectWalletButton = document.getElementById("connectWalletButton");
connectWalletButton.addEventListener("click", async function () {
  var otp = document.getElementById('walletOTP').value;
  var response = ipcRenderer.sendSync("connectWallet", otp);
  console.log(response);
  document.getElementById("walletAddressValue").innerText = "Wallet Address: " + response.WalletAddress

  // TODO : Save in more secure place. For now can't save in cookie due to bug
  localStorage.setItem("WalletAddress", response.WalletAddress);
  localStorage.setItem("Token", response.Token);
});

//ipcRenderer.on will receive the “btnclick-task-finished'” info from main process
ipcRenderer.on("registerFinished", function (event, param) {
  console.log(param);
});

const setupModal = (modalId, btnId, spanId) => {
  var modal = document.getElementById(modalId);
  var btn = document.getElementById(btnId);
  var span = document.getElementById(spanId);
  console.log(modal);

  btn.onclick = function () {
    modal.style.display = "block";
  };

  span.onclick = function () {
    modal.style.display = "none";
  };

  window.onclick = function (event) {
    if (event.target == modal) {
      modal.style.display = "none";
    }
  };
};

setupModal("addAppModal", "addAppButton", "addAppClose");
setupModal("connectWalletModal", "connectWalletButtonParent", "connectWalletClose");
