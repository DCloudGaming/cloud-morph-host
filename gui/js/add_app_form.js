// TODO: Replace this mock data
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

function addAppRow(serverId) {
  appRowId++;
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
  nameElement.innerHTML = sampleAllowedApps.map(
    (app) => `<option value=${app.id}>${app.name}</option>`
  );
  pathElement.id = `apppathText-${id}`;
  registerElement.id = `registerButton-${id}`;
}

function removeAppRow(id) {
  d = document;
  var ele = d.getElementById(getAppRowId(id));
  console.log(id);
  console.log(ele);
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
function prefillAddAppForm(registeredApp) {
  const wrapper = document.getElementById("appRowWrapper");
  registeredApp.forEach((element) => {
    addAppRow(element.id);
    const { selectedElement, nameElement, pathElement, appRow } =
      getAppRowElements(getAppRowId(element.id));
    selectedElement.checked = element.selected;
    nameElement.value = element.id;
    pathElement.value = element.path;
  });
}

function updateApps() {
  const appRows = Array.from(document.getElementsByClassName("app-row"));
  const body = appRows.map((appRow) => {
    const { selectedElement, nameElement, pathElement } = getAppRowElements(
      appRow.id
    );
    return {
      id: appRow.id, // ID from server if it's an existing entry, random ID otherwise
      name: nameElement.options[nameElement.selectedIndex].innerText,
      selected: selectedElement.checked,
      path: pathElement.value,
    };
  });

  console.log("Sending to server: ", body);
}

// Main
prefillAddAppForm(sampleRegisteredApps);
