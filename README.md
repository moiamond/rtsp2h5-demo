# 使用 HTML5 瀏覽 RTSP 串流 


## [FFmpeg to genetate HLS](./ffmpeg-hls/README.md)

URL: `http://{SERVER IP}:7000/app/viewer.html`

Server side:
* Use python to host a web server
* Use ffmpeg to generate HLS

Client side:
* Use Video.js to playback HLS

## [NGINX RTMP Module + FFmpeg](./nginx-rtmp/README.md)

URL: `http://{SERVER IP}:8000/app/viewer.html`

Server side:
* Use NGINX to host a web server
* Use nginx-rtmp-module to host a RTMP server
* Use ffmpeg to restreaming RTSP to RTMP server

Client side:
* Use Video.js to playback HLS

## [FFmpeg to generate HTTP-MP4](./http-mp4/README.md)

URL: `http://{SERVER IP}:9000/app/viewer.html`

Server side:
* Use python to host a web server
* Use nodejs to response resuqsted RTSP stream to Fragmented mp4

Client side:
* Use native video tag to playback Fragmented mp4

