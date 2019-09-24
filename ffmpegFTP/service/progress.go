package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tj3828/ffmpeg-web-client/ffmpegFTP/logio"
)

func sendProgressStatus(progress string, file *uploadRequestFile) error {

	files := msg{}
	files.FileIndex = file.FileIndex
	files.TotalTime = file.TotalTime

	files.ProgressStatus = progress[strings.Index(progress, "progress=")+9:]
	if files.ProgressStatus == "end" {
		files.ProgressRate = "100"
		if err := currentWebsocket.WriteJSON(files); err != nil {
			return err
		}
		return nil
	}

	strSpeed := progress[strings.Index(progress, "speed=")+6 : strings.Index(progress, "speed=")+10]
	floSpeed, err := strconv.ParseFloat(strSpeed, 32)
	if err != nil {
		return err
	}

	files.ProgressTime = progress[strings.Index(progress, "time=")+5 : strings.Index(progress, "time=")+16]

	timeArr := strings.Split(files.ProgressTime, ":")
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

	remainingTimeSecond := int((float64(file.ConvertTime) - float64(currentTime)) / floSpeed)
	logio.Info.Print("Time Remaining Second : " + fmt.Sprintf("%d", remainingTimeSecond))

	secondsToTime(remainingTimeSecond, &files)
	if err = currentWebsocket.WriteJSON(files); err != nil {
		return err
	}

	return nil
}

func secondsToTime(input int, file *msg) {
	days := math.Floor(float64(input) / 60 / 60 / 24)
	seconds := input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if days > 0 {
		file.TimeRemaining = strconv.Itoa(int(days)) + "일"
	} else if hours > 0 {
		file.TimeRemaining = strconv.Itoa(int(hours)) + "시간 " + strconv.Itoa(int(minutes)) + "분"
	} else if minutes > 0 {
		file.TimeRemaining = strconv.Itoa(int(minutes)) + "분 " + strconv.Itoa(int(seconds)) + "초"
	} else {
		file.TimeRemaining = strconv.Itoa(int(seconds)) + "초"
	}
}
