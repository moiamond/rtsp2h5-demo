package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var addr = flag.String("addr", ":9090", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func serveWs(w http.ResponseWriter, r *http.Request, channel string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	cmd := exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		"rtsp://192.168.10.21/LV/"+channel,
		"-vcodec",
		"copy",
		"-an",
		"-f",
		"mp4",
		"-movflags",
		"empty_moov+default_base_moof+frag_keyframe", // <- Chrome, "frag_keyframe+empty_moov" <- firefox
		"-reset_timestamps",
		"1",
		"-vsync",
		"1",
		"-flags",
		"global_header",
		"-bsf:v",
		"dump_extra",
		"-y",
		"-")
	printCommand(cmd)

	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	err = cmd.Start()
	checkError(err)

	for {
		err = ws.WriteMessage(websocket.BinaryMessage, randomBytes.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(time.Second)
	}
	cmd.Process.Kill()
}

func ch1(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch1")
}

func ch2(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch2")
}

func ch3(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch3")
}

func ch4(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch4")
}

func ch5(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch5")
}

func ch7(w http.ResponseWriter, r *http.Request) {
	serveWs(w, r, "ch7")
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/192.168.10.21/LV/ch1", ch1)
	http.HandleFunc("/192.168.10.21/LV/ch2", ch2)
	http.HandleFunc("/192.168.10.21/LV/ch3", ch3)
	http.HandleFunc("/192.168.10.21/LV/ch4", ch4)
	http.HandleFunc("/192.168.10.21/LV/ch5", ch5)
	http.HandleFunc("/192.168.10.21/LV/ch7", ch7)

	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
