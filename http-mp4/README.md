# 使用 HTML5 瀏覽 RTSP 串流 
> 利用 FFmpeg

## Demo

### Server side (Ubuntu with node)

#### 1. Clone the src code

```shell
$ git clone https://github.com/moiamond/rtsp2h5-demo
$ cd rtsp2h5-demo/http-mp4
```

#### 2. Run HTTP server

```shell
$ python server.py
```

#### 3. Run 6 node.js APPs to server RTSP to fragmented mp4

```shell
$ node rtsp2fmp4.js 9001
$ node rtsp2fmp4.js 9002
$ node rtsp2fmp4.js 9003
$ node rtsp2fmp4.js 9004
$ node rtsp2fmp4.js 9005
$ node rtsp2fmp4.js 9006
```

#### 4. View in the browser

使用瀏覽器開啟 `hppt://{SERVER IP}:9000/app/viewer.html`


低延遲

```
      ~0s delay
RTSP -----------> HTTP-fMP4
```

![lantency](../pics/low-latency.png)
