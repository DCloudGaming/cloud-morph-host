# cloud-morph-host

Windows client app for host on cloud-morph..

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

### 2. Run streamer

```
cd streamer
go run main.go
```

### 3. Run GUI, install Electron

```
cd gui
npm start
```

### 4. Run webapp

```
cd dcloud-webapp
npm install
npm start
```

### 5. Run full flow

- (First run) Setup Metamask extension
- Connect wallet
  - Click on Connect Wallet on React app
  - Connect and sign in Metamask pop-up
- Authorize Electron app
  - Get OTP on React app
  - Paste it to Electron app
  - Quick DB check
    - sqlite3
    - select \* from smart_otps;
    - We should see rows having both OTP and wallet addresses.
- (First run) Whitelist wallet address
  - sqlite3
  - INSERT INTO whitelisted_admins (id, wallet_address) VALUES (1, 'metamask-wallet-address');
- (First run) Register new app
  - Click Admin Update on React app navbar
  - Add some random app names
  - Save
- (First run) Register new app paths
  - Click Add app on Electron app
  - Pick a recently added app name
  - Choose the app
  - Save
- Start streaming
  - Go to http://localhost:3000/streams
  - Should see some cards
  - Click Start playing on one card
  - Should be directed to http://localhost:3000/play
  - Monitoring:
    - Console of React app: init/ice candidate/error logs
    - chrome://webrtc-internals/
- Notes
  - To restart the flow, sign out Metamask by
    - Click on Metamask extension icon
    - Click on Lock
