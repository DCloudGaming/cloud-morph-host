// main.js

// Modules to control application life and create native browser window
const { app, BrowserWindow } = require('electron')
const path = require('path')

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

//ipcMain.on will receive the “btnclick” info from renderprocess 
ipcMain.on("register", function (event, arg) {
  const { net, dialog, electron } = require('electron')

  const handleRegister = (path) => {
    event.sender.send("registerFinished", path);

    const registerURL = "";
    var postData = JSON.stringify([{ "app_name": "Neighbour", "app_path": path }]);
    console.log("Send HTTP register request to notepad ", postData)
    const request = net.request({
      method: 'POST',
      protocol: 'http:',
      hostname: 'localhost',
      port: 8082,
      path: '/registerApp',
      headers: {
        "Content-Type": "application/json",
      },
    })

    console.log(2);
    request.on('response', (response) => {
      console.log(`STATUS: ${response.statusCode}`)
      console.log(`HEADERS: ${JSON.stringify(response.headers)}`)
      response.on('data', (chunk) => {
        console.log(`BODY: ${chunk}`)
      })
      response.on('end', () => {
        console.log('No more data in response.')
      })
    })

    request.on('error', (error) => {
      console.error(error)
    })
    request.write(postData);
    request.end();
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