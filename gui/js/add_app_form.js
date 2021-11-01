const sampleAllowedApps = [
  {
    id: 1, // ID in the allowed apps list, from server
    name: "Fred",
  },
  {
    id: 2,
    name: "Lorde",
  },
  {
    id: 3,
    name: "Doja",
  },
  {
    id: 4,
    name: "Troye",
  },
  {
    id: 5,
    name: "Billie",
  },
];

const sampleRegisteredApps = [
  {
    id: 1, // ID in the registered app list, from server
    name: "Fred",
    votes: 3,
    selected: true,
    path: "sample/path",
  },
  {
    id: 2,
    name: "Lorde",
    votes: 3,
    selected: false,
    path: "sample/path",
  },
  {
    id: 3,
    name: "Doja",
    votes: 3,
    selected: true,
    path: "sample/path",
  },
  {
    id: 5,
    name: "Billie",
    votes: 3,
    selected: true,
    path: "sample/path",
  },
];

// Utils
var appRowId = 1;

function getAppRowId(eleId) {
  return `appRow-${eleId}`;
}

function addAppRow(serverId, allowedApps, chosenName="") {
  const id = serverId != 0 ? serverId : appRowId;
  var newDiv = document.createElement("div");
  newDiv.id = getAppRowId(id);
  newDiv.className = "row space-between mb-2 app-row";

  const removeButton = `<button onclick=removeAppRow(${id})>&times;</button>`;
  newDiv.innerHTML = document.getElementById("appRow").innerHTML + removeButton;
  document.getElementById("appRowWrapper").appendChild(newDiv);

  const { nameElement, pathElement, registerElement } = getAppRowElements(
    getAppRowId(id)
  );
  nameElement.id = `app-${id}`;
  nameElement.innerHTML = allowedApps.map(
      (app_name) => {
        if (chosenName == "") {
          return `<option value=${id}>${app_name}</option>`;
        }
        if (app_name == chosenName) {
          return `<option value=${id} selected>${app_name}</option>`;
        }
      }
  );

  pathElement.id = `apppathText-${id}`;
  registerElement.id = `registerButton-${id}`;
  appRowId++;
  registerElement.addEventListener("click", function () {
    //send the info to main process . we can pass any arguments as second param.
    // ipcRender.send will pass the information to main process. Here is event to open file dialog
    ipcRenderer.send("register", {appPathText: pathElement.id, id: id});
  });
}

function removeAppRow(id) {
  d = document;
  var ele = d.getElementById(getAppRowId(id));
  var parentEle = d.getElementById("appRowWrapper");
  parentEle.removeChild(ele);
}

function getAppRowElements(id) {
  const appRow = document.getElementById(id);
  const selectedElement = appRow.children[0];
  const nameElement = appRow.children[1];
  const pathElement = appRow.children[2];
  const registerElement = appRow.children[3];
  const removeElement = appRow.children[4];

  return {
    appRow,
    selectedElement,
    nameElement,
    pathElement,
    registerElement,
    removeElement,
  };
}

// Handlers
function prefillAddAppForm() {
  let getAllowedAppsResponse = ipcRenderer.sendSync("getAllowedApps");
  let allowedApps = getAllowedAppsResponse.AllowedApps;

  let getRegisteredAppsResponse = ipcRenderer.sendSync("getRegisteredApps", localStorage.getItem("WalletAddress"));
  let registeredApps = getRegisteredAppsResponse.AppMetas;

  const wrapper = document.getElementById("appRowWrapper");
  registeredApps.forEach((element) => {
    addAppRow(appRowId, allowedApps, chosenName=element.app_name);
    const { selectedElement, nameElement, pathElement, appRow } =
      getAppRowElements(getAppRowId(appRowId-1));
    selectedElement.checked = true;
    nameElement.value = appRowId;
    pathElement.value = element.app_path;
  });
}

// Main
prefillAddAppForm();
