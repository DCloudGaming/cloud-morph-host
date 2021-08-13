# cloud-morph-host
Windows client app for host on cloud-morph.

## Design
![screenshot](docs/img/dclouddiagram.png)  
[Edit in draw.io](https://drive.google.com/file/d/1MuF32rcGpRHmpQrA0_MX2IgTkY6Evv7J/view?usp=sharing)

## Getting started
### 1. Start server 
```bigquery
cd server
go run main.go
```
### 2. Start player client
Open `http://localhost:8080/`
### 3. Start streamer host
```bigquery
cd streamer
go run main.go
```