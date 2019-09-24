package service

import (
	"github.com/tj3828/ffmpeg-web-client/ffmpegFTP/logio"
	"fmt"
	"strconv"
	"strings"
)

func sendProgressStatus(progress string, file *uploadRequestFile) error {

	files := msg{}
	files.FileIndex = file.FileIndex
	files.TotalTime = file.TotalTime

	key := strings.Split(progress, "=")
	if key[0] == "out_time" {
		files.ProgressTime = key[1]
		files.ProgressStatus = "continue"

		timeArr := strings.Split(key[1], ":")
		hour, err := strconv.Atoi(timeArr[0])
		if err != nil {
			return err
		}
		minute, err := strconv.Atoi(timeArr[1])
		if err != nil {
			return err
		}
		tempSec := timeArr[2]
		second, err := strconv.Atoi(tempSec[:2])
		if err != nil {
			return err
		}

		currentTime := hour*60*60 + minute*60 + second
		logio.Info.Print(strconv.Itoa(hour) + " / " + strconv.Itoa(minute) + " / " + strconv.Itoa(second) + " / " + strconv.Itoa(file.ConvertTime))
		files.ProgressRate = fmt.Sprintf("%.2f", float64(currentTime)/float64(file.ConvertTime)*100)

		if err = currentWebsocket.WriteJSON(files); err != nil {
			return err
		}
	} else if key[0] == "progress" && key[1] == "end" {
		files.ProgressStatus = key[1]
		files.ProgressRate = "100"

		if err := currentWebsocket.WriteJSON(files); err != nil {
			return err
		}
	}
	return nil
}
