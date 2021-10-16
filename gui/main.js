// main.js

// Modules to control application life and create native browser window
const { app, BrowserWindow } = require('electron')
const path = require('path')
const axios = require("axios");
axios.defaults.withCredentials = true;

function createWindow() {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: true,
      contextIsolation: false,
      devTools: true
    }
  })

  // and load the index.html of the app.
  mainWindow.loadFile('index.html')

  // Open the DevTools.
  mainWindow.webContents.openDevTools()
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  createWindow()

  // Run Streamer app
  // var child = require('child_process').execFile;
  // var executablePath = "../streamer/main.exe";
  // child(executablePath, function (err, data) {
  //   if (err) {
  //     console.error(err);
  //     return;
  //   }
  //   console.log("Loaded streamer");
  //   console.log(data.toString());
  // });

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') app.quit()
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
const { ipcMain } = require('electron'); // include the ipc module to communicate with render process ie to receive the message from render process
const { hostname } = require('os')

//ipcMain.on will receive the “btnclick” info from renderprocess 
ipcMain.on("btnclick", function (event, arg) {
  //create new window
  var newWindow = new BrowserWindow({
    width: 450, height: 300, show:
      false, webPreferences: {
        webSecurity: false, plugins:
          true, nodeIntegration: false
      }
  });  // create a new window

  var facebookURL = "https://www.facebook.com"; // loading an external url. We can load our own another html file , like how we load index.html earlier

  newWindow.loadURL(facebookURL);
  newWindow.show();

  // inform the render process that the assigned task finished. Show a message in html
  // event.sender.send in ipcMain will return the reply to renderprocess
  event.sender.send("btnclick-task-finished", "yes");
});

ipcMain.on("connectWallet", async function(event, arg) {
  var response = await axios({
    method: "POST",
    url: "http://localhost:8080/api/users/verifyOTP",
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      otp: arg
    },
    withCredentials: true
  })

  await axios({
    method: "POST",
    url: "http://localhost:8082/updateToken",
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      token: response.data.token
    },
    withCredentials: true
  })

  console.log(response.data);
  event.returnValue = {
    WalletAddress: response.data.wallet_address,
    Token: response.data.token
  }

});

ipcMain.on("registerApps", async function(event, arg) {
  console.log(arg);
  var response = await axios({
    method: "POST",
    url: "http://localhost:8080/api/apps/registerApp",
    headers: {
      "Content-Type": "application/json"
    },
    data: {
      wallet_address: arg.walletAddress,
      token: arg.token,
      app_paths: arg.appPaths,
      app_names: arg.appNames
    },
    withCredentials: true
  })

});

ipcMain.on("getAllowedApps", async function(event, arg) {
  var response = await axios({
    method: "GET",
    url: "http://localhost:8080/api/users/getAdminSettings",
    headers: {
      "Content-Type": "application/json",
    }
  })

  event.returnValue = {
    AllowedApps: response.data.allowed_apps
  }
});

ipcMain.on("getRegisteredApps", async function(event, arg) {
  var response = await axios({
    method: "GET",
    url: `http://localhost:8080/api/apps?wallet_address=${arg}`,
    headers: {
      "Content-Type": "application/json",
    }
  })

  event.returnValue = {
    AppMetas: response.data
  }
})

//ipcMain.on will receive the “btnclick” info from renderprocess 
ipcMain.on("register", function (event, arg) {
  const { net, dialog, electron } = require('electron')

  const handleRegister = (path) => {
    event.sender.send("registerFinished", {
      Path: path,
      appPathText: arg.appPathText,
      id: arg.id
    });
  }

  dialog.showOpenDialog({ properties: ['openFile', 'multiSelections'] }).then(result => {
    console.log(result.filePaths[0])
    // inform the render process that the assigned task finished. Show a message in html
    // event.sender.send in ipcMain will return the reply to renderprocess
    handleRegister(result.filePaths[0])
  }).catch(err => {
    console.log(err)
  })

});
