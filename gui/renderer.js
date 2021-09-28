// include the ipc module to communicate with main process.
const ipcRenderer = require("electron").ipcRenderer;

const registerButton = document.getElementById("registerButton");
const apppathText = document.getElementById("apppathText");

registerButton.addEventListener("click", function () {
  var arg = "secondparam";

  //send the info to main process . we can pass any arguments as second param.
  // ipcRender.send will pass the information to main process. Here is event to open file dialog
  ipcRenderer.send("register", arg);
});

const connectWalletButton = document.getElementById("connectWalletButton");
connectWalletButton.addEventListener("click", function () {
  var arg = "secondparam";

  ipcRenderer.send("connectWallet", arg);
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
setupModal("connectWalletModal", "connectWalletButton", "connectWalletClose");
