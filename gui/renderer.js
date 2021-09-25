// include the ipc module to communicate with main process.
const ipcRenderer = require('electron').ipcRenderer;

const registerButton = document.getElementById('registerButton');
const apppathText = document.getElementById('apppathText');

registerButton.addEventListener('click', function () {
    var arg = "secondparam";

    //send the info to main process . we can pass any arguments as second param.
    // ipcRender.send will pass the information to main process. Here is event to open file dialog
    ipcRenderer.send("register", arg);
});

//ipcRenderer.on will receive the “btnclick-task-finished'” info from main process 
ipcRenderer.on('registerFinished', function (event, param) {
    console.log(param);
    apppathText.value = param;
});