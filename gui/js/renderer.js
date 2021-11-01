// include the ipc module to communicate with main process.
const ipcRenderer = require("electron").ipcRenderer;

const apppathText = document.getElementById("apppathText");

const connectWalletButton = document.getElementById("connectWalletButton");
connectWalletButton.addEventListener("click", async function () {
  let otp = document.getElementById('walletOTP').value;
  let response = ipcRenderer.sendSync("connectWallet", otp);
  document.getElementById("walletAddressValue").innerText = "Wallet Address: " + response.WalletAddress

  // TODO : Save in more secure place. For now can't save in cookie due to bug
  localStorage.setItem("WalletAddress", response.WalletAddress);
  localStorage.setItem("Token", response.Token);
});

const saveButton = document.getElementById("saveButton");
saveButton.addEventListener("click", async function () {
  const appRows = Array.from(document.getElementsByClassName("app-row"));
  const appPaths = [];
  const appNames = [];
  appRows.forEach((appRow) => {
    const {selectedElement, nameElement, pathElement} = getAppRowElements(
        appRow.id
    );
    // TODO: Store checked/unchecked info in backend too.
    if (selectedElement.checked) {
      appPaths.push(pathElement.value);
      appNames.push(nameElement.options[nameElement.selectedIndex].innerText)
    }
  });
  ipcRenderer.sendSync("registerApps", {
    walletAddress: localStorage.getItem("WalletAddress"),
    token: localStorage.getItem("Token"),
    appPaths: appPaths, appNames: appNames
  })
});

//ipcRenderer.on will receive the “btnclick-task-finished'” info from main process
ipcRenderer.on("registerFinished", function (event, param) {
  let appPathText = document.getElementById(param.appPathText);
  appPathText.value = param.Path;
});

const addNewAppButton = document.getElementById("addnew");
addNewAppButton.addEventListener("click", async function () {
  let response = ipcRenderer.sendSync("getAllowedApps");
  addAppRow(0, response.AllowedApps);
})

const setupModal = (modalId, btnId, spanId) => {
  let modal = document.getElementById(modalId);
  let btn = document.getElementById(btnId);
  let span = document.getElementById(spanId);
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
