/**
 * Live video stream for HTML5 video. Uses FFMPEG to connect to H.264 camera stream,
 *
 * Ubuntu 14, node 0.11, ffmpeg 2.4.3
 * IP camera: spydroid App, H264 over RTSP via UDP
 * Note ffmpeg param frag_keyframe it will only work with chrome, firefox needs frag_keyframe+empty_moov
 *
 * sudo add-apt-repository ppa:mc3man/trusty-media
 * sudo apt-get update
 * sudo apt-get install ffmpeg
 *
 * start server:
 * node rtsp2fmp4.js 9000
 *
 * browse to:
 * http://localhost:9000
 */

var http = require('http');
var url = require('url');
var child_process = require('child_process');
var camera = {rtsp: null, liveStarted: false, mp4Headers: []};

// handle each client request by instantiating a new FFMPEG instance
http.createServer(function (request, response) {

    // We need to get the URL of our RSTP stream from request URL eg. http://localhost:8080/{streamIP}
    var streamURL = getStreamAddressFromUrl(request.url);

    // No stream address given - serve static HTML video player
    if (!streamURL) {
        console.log("GET / from: " + request.connection.remoteAddress);

        response.writeHead(200, {"Content-Type": "text/html"});
        response.end(
            '<video id="video" autoplay width="320" height="264" ></video><br />rtsp://<input id="rtsp">' +
            '<button onclick="document.getElementById(\'video\').src=\'http://'+ request.headers.host +'/\' + document.getElementById(\'rtsp\').value">Connect</button>');
        return true;
    }

    // Skip rubbish
    if (streamURL == 'favicon.ico') {
        response.writeHead(404);
        response.end();
        return true;
    }

    console.log("GET /"+ streamURL +" from: " + request.connection.remoteAddress + ':' + request.connection.remotePort);


    camera.rtsp = streamURL;

    if (typeof camera.liveffmpeg !== "undefined") {
        console.log("Changing live stream to: " + streamURL);
        camera.liveStarted = false;
        camera.liveffmpeg.kill();
    }

    // Range HTTP response needs a proper header
    response.writeHead(200, {
        // 'Transfer-Encoding': 'binary',
        "Connection": "keep-alive",        // This will not break TCP connection after sending chunks
        "Content-Type": "video/mp4",       // MIME Type
        //, 'Content-Length': chunksize,
        "Accept-Ranges": "bytes"           // Helps Chrome
    });

    if (camera.liveStarted == false) {

        // Reset buffered mp4 packets
        camera.mp4Headers = [];

        // Camera stream is remuxed to a MP4 stream for HTML5 video compatibility and segments are recorded for playback
        // For live streaming, create a fragmented MP4 file with empty moov (no seeking possible)
        camera.liveffmpeg = child_process.spawn("ffmpeg", [
            "-rtsp_transport",
                "tcp",
            "-i",
                "rtsp://" + camera.rtsp,
            "-vcodec",
                "copy",
            "-f",
                "mp4",
            "-movflags",
                "frag_keyframe", // <- Chrome, "frag_keyframe+empty_moov" <- firefox
             "-reset_timestamps",
                "1",
            "-vsync",
                "1",
            "-flags",
                "global_header",
            "-bsf:v",
                "dump_extra",
            "-y",
            "-" // output to stdout
        ], { detached: false });

        // Pipe FFMPEG's stdout directly to response
        camera.liveffmpeg.stdout.pipe(response);

        // Buffer initial x packets
        camera.liveffmpeg.stdout.on("data", function(data) {
            if (camera.mp4Headers.length < 3) {
                  camera.mp4Headers.push(data);
            }
        });

        camera.liveffmpeg.stderr.on("data", function (data) {
            console.log("FFMPEG -> " + data);
        });
        camera.liveffmpeg.on("exit", function (code) {
            console.log("Live FFMPEG terminated with code " + code);
        });
        camera.liveffmpeg.on("error", function (e) {
            console.log("Live FFMPEG system error: " + e);
        });
        camera.liveStarted = true;
    } else {

        // Re-send buffered mp4 packets
        for (var i = 0; i < camera.mp4Headers.length; i++) {
            response.write(camera.mp4Headers[i]);
        }
        console.log('Fresh conn, re-sent: ' + camera.mp4Headers.length);

        // Pipe FFMPEG's stdout directly to response
        camera.liveffmpeg.stdout.pipe(response);
    }

    return true;

}).listen(process.argv[2]);

var getStreamAddressFromUrl = function(requestUrl) {

    var reqUrl = url.parse(requestUrl, true);
    var streamIP = typeof reqUrl.pathname === "string" ? reqUrl.pathname.substring(1) : undefined;
    if (streamIP) {
        try {
            return decodeURIComponent(streamIP);
        } catch (exception) {
            // Can throw URI malformed exception.
            console.log("Live Camera Streamer bad request received - " + reqUrl);
        }
    }
};