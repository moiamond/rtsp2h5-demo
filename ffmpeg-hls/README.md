# 使用 HTML5 瀏覽 RTSP 串流 
> 利用 FFmpeg

## Demo

### Server side (Ubuntu with node)

#### 1. Clone the src code

```shell
$ git clone https://github.com/moiamond/rtsp2h5-demo
$ cd rtsp2h5-demo/ffmpeg-hls
```

#### 2. Launch a web server

```shell
$ python server.py
```
#### 3. Run 6 FFmpeg to generate HLS

```shell
$ ffmpeg -rtsp_transport tcp \
    -i rtsp://192.168.10.21/LV/ch7 \
    -c:v copy \
    -an \
    -f hls \
    -hls_time 2 \
    -hls_list_size 0  \
    -hls_flags delete_segments \
    ch7.m3u8
```

#### 4. Open browser

Navigate to `http://localhost:7000/app/viewer.html`