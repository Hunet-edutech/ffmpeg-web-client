package service

import (
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

func initProgramStatus() {
	var tempRunningFiles []uploadRequestFile
	copier.Copy(&tempRunningFiles, &runningFiles)
	for _, f := range tempRunningFiles {
		stopRunningEncoding(f)
	}
	runningFiles = []uploadRequestFile{}
}

func setTotalTimeAndConvertTime(progress string, file *uploadRequestFile) error {

	key := strings.Split(progress, "=")

	temp := key[1]
	key[1] = temp[:len(temp)-13]

	file.TotalTime = key[1]

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

	totalTime := hour*60*60 + minute*60 + second

	file.ConvertTime = totalTime

	return nil
}

// Request 파일 저장
func requestFileIO(file multipart.File, files *uploadRequestFile) error {
	savedFileName := files.BasePath + "\\temp_" + files.CourseCd + "_" + files.UploadFileName
	f, err := os.OpenFile(savedFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, file)

	return err
}

// 인코딩 중지시 프로세스와 파일 삭제
func stopRunningEncoding(file uploadRequestFile) bool {
	for i, f := range runningFiles {
		if f.FileIndex == file.FileIndex {
			kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(f.RunningCmd.Process.Pid))
			if err := kill.Start(); err != nil {
				return false
			}

			if err := kill.Wait(); err != nil {
				return false
			}

			for {
				if err := os.Remove(f.BasePath + "\\temp_" + f.CourseCd + "_" + f.UploadFileName); err == nil {
					break
				}
			}
			for {
				if err := os.Remove(f.BasePath + "\\" + f.CourseCd + "_" + f.UploadFileName); err == nil {
					break
				}
			}
			runningFiles = append(runningFiles[:i], runningFiles[i+1:]...)
			return true

		}
	}
	return false
}

// ftp 업로드 완료 후 요청받은 파일과 인코딩한 파일 삭제
func removeRunningFile(file uploadRequestFile) error {
	for i, f := range runningFiles {
		if f.FileIndex == file.FileIndex {
			runningFiles = append(runningFiles[:i], runningFiles[i+1:]...)
			if err := os.Remove(file.BasePath + "\\temp_" + f.CourseCd + "_" + f.UploadFileName); err != nil {
				return err
			}
			if err := os.Remove(file.BasePath + "\\" + f.CourseCd + "_" + f.UploadFileName); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
