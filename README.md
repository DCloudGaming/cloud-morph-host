# cloud-morph-host

Windows client app for host on cloud-morph.

## Design

![screenshot](docs/img/dclouddiagram.png)  
[Edit in draw.io](https://drive.google.com/file/d/1MuF32rcGpRHmpQrA0_MX2IgTkY6Evv7J/view?usp=sharing)

## Getting started

### 0. Setup

```
git clone --recurse-submodules https://github.com/DCloudGaming/cloud-morph-host.git
```

To update:

```
git pull --recurse-submodules
```

### 1. Start server

```
cd server
go run main.go
```

### 2. Run GUI, install Electron

```
cd gui
npm start
```

### 3. Run webapp

```
cd dcloud-webapp
npm install
npm start
```

### 4. Open App on browser

1. Open `http://localhost:8080/play`
2. Click `Register` Button in GUI
3. Click Notepad entry

## FAQ

If cannot run, there maybe some duplicate process running. We can fix but in the mean time, we can restart everything
