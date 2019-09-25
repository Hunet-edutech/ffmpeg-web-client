package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// "strings"
	"path/filepath"

	"os"
	os_exec "os/exec"

	"github.com/Hunet-edutech/ffmpeg-web-client/ffmpegFTP/logio"

	"github.com/gorilla/websocket"
)

var currentWebsocket *websocket.Conn

var runningFiles []uploadRequestFile
var basePath string

type msg struct {
	FileIndex      int    `json:"fileIndex"`
	FileName       string `json:"fileName"`
	ProgressTime   string `json:"progressTime"`
	ProgressStatus string `json:"progressStatus"`
	TotalTime      string `json:"totalTime"`
	TimeRemaining  string `json:"timeRemaining"`
	ProgressRate   string `json:"progressRate"`
}

type uploadRequestFile struct {
	FileIndex      int    `json:"fileIndex"`
	CourseCd       string `json:"courseCd"`
	BasePath       string
	UploadFileName string `json:"uploadFileName"`
	RunningCmd     *os_exec.Cmd
	TotalTime      string
	ConvertTime    int
}

type uploadEncodingInfo struct {
	VideoCodec   string `json:"videoCodec"`
	VideoBitrate string `json:"videoBitrate"`
	AudioCodec   string `json:"audioCodec"`
	AudioBitrate string `json:"audioBitrate"`
	Resolution   string `json:"resolution"`
}

//cross domain 문제 해결
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Max-Age", "3600")
}

func echo(w http.ResponseWriter, r *http.Request) {
	// 여러개 브라우저 동시 접속시
	initProgramStatus()

	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	currentWebsocket = c

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logio.Info.Print("read:", err)
			break
		}
		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)
	}

	// 웹소켓 연결 해제시 초기화 작업
	initProgramStatus()
}

func fileExist(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	files := uploadRequestFile{}
	files.CourseCd = r.URL.Query().Get("courseCd")
	files.UploadFileName = r.URL.Query().Get("uploadFileName")

	check := fileExistCheck(files)
	w.Write([]byte(check))
}

func serviceCheck(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.Write([]byte("1"))
}

func transcoding(w http.ResponseWriter, r *http.Request) {
	//cors 허용
	setupResponse(&w, r)
	switch r.Method {
	case http.MethodPost:

		encodingInfo := uploadEncodingInfo{}
		requestFile := uploadRequestFile{}
		requestFile.BasePath = basePath

		r.ParseMultipartForm(32 << 20)

		tempFileInfo := r.FormValue("uploadFileInfo")
		if err := json.Unmarshal([]byte(tempFileInfo), &requestFile); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		tempEncodingInfo := r.FormValue("uploadEncodingInfo")
		if err := json.Unmarshal([]byte(tempEncodingInfo), &encodingInfo); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		requestfile, _, err := r.FormFile("uploadFileData")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if err = requestFileIO(requestfile, &requestFile); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// FFprobe Start
		if err = startFFprobe(&requestFile); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FFmpeg Start
		if err = startFFmpeg(&requestFile, encodingInfo); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// FTP upload
		if err = uploadFTP(requestFile); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if err = removeRunningFile(requestFile); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

	case http.MethodDelete:
		decoder := json.NewDecoder(r.Body)
		requestFile := uploadRequestFile{}
		if err := decoder.Decode(&requestFile); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		check := stopRunningEncoding(requestFile)

		w.Write([]byte(strconv.FormatBool(check)))

	default:
		w.Write([]byte("wrong access\n"))
	}
}

func StartServer() {

	basePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	os.Chdir(basePath)
	logio.Info.Print("Executable dir path : " + basePath)

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/fileExist", fileExist)
	http.HandleFunc("/serviceCheck", serviceCheck)
	http.HandleFunc("/transcoding", transcoding)

	http.ListenAndServe("localhost:5050", nil)

}
