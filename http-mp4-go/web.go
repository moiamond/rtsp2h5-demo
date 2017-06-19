package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func runFFmpeg(w http.ResponseWriter, r *http.Request, channel string) {
	fmt.Println("runFFmpeg:")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")

	cmd := exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		"rtsp://192.168.10.21/LV/"+channel,
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
		"-")
	printCommand(cmd)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = w

	// Start command asynchronously
	err := cmd.Start()
	printError(err)

	if _, err := w.Write(randomBytes.Bytes()); err != nil {
		log.Println("unable to write output.")
	}

	cmd.Wait()
	log.Println("leave ...")
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func ch1(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch1")
}

func ch2(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch2")
}

func ch3(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch3")
}

func ch4(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch4")
}

func ch5(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch5")
}

func ch7(w http.ResponseWriter, r *http.Request) {
	runFFmpeg(w, r, "ch7")
}

func main() {
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))

	http.HandleFunc("/192.168.10.21/LV/ch1", ch1)
	http.HandleFunc("/192.168.10.21/LV/ch2", ch2)
	http.HandleFunc("/192.168.10.21/LV/ch3", ch3)
	http.HandleFunc("/192.168.10.21/LV/ch4", ch4)
	http.HandleFunc("/192.168.10.21/LV/ch5", ch5)
	http.HandleFunc("/192.168.10.21/LV/ch7", ch7)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
