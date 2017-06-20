# 使用 HTML5 瀏覽 RTSP 串流 
> 利用 FFmpeg

## Demo

### Server side (Ubuntu with Go)

#### 1. Clone the src code

```shell
$ git clone https://github.com/moiamond/rtsp2h5-demo
$ cd rtsp2h5-demo/websocket-mse-go
```

#### 2. Run HTTP server

```shell
$ go run main.go
```

#### 3. View in the browser

使用瀏覽器開啟 `hppt://{SERVER IP}:9090/app/viewer.html`

