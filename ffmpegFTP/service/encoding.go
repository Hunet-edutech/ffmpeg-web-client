package service

import (
	"bufio"
	"strings"

	os_exec "os/exec"
)

func startFFprobe(file *uploadRequestFile) error {
	metadata, err := os_exec.Command(
		`ffprobe`,
		`-v`, `error`,
		`-show_entries`,
		`format=duration`,
		`-sexagesimal`,
		"temp_"+file.CourseCd+"_"+file.UploadFileName).Output()
	if err != nil {
		return err
	}

	if err = setTotalTimeAndConvertTime(string(metadata), file); err != nil {
		return err
	}

	return err
}

func startFFmpeg(file *uploadRequestFile, encodingInfo uploadEncodingInfo) error {
	// 인코딩 (mp4, 700k, aac, 128k, 970x546)
	cmd := os_exec.Command(
		`ffmpeg`,
		`-progress`, `pipe:1`,
		`-y`,
		`-i`, "temp_"+file.CourseCd+"_"+file.UploadFileName,
		`-acodec`, encodingInfo.AudioCodec,
		`-b:a`, encodingInfo.AudioBitrate,
		`-vcodec`, encodingInfo.VideoCodec,
		`-b:v`, encodingInfo.VideoBitrate,
		`-s`, encodingInfo.Resolution,
		file.CourseCd+"_"+file.UploadFileName)

	file.RunningCmd = cmd
	runningFiles = append(runningFiles, *file)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	var text = ""
	for scanner.Scan() {
		m := scanner.Text()
		text += m
		if strings.Contains(m, "progress=") {
			sendProgressStatus(text, file)
			text = ""
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
