package service

import (
	"os"
	"strings"

	"github.com/tj3828/ffmpeg-web-client/ffmpegFTP/logio"

	"github.com/dutchcoders/goftp"
)

func connectFTP() (*goftp.FTP, error) {
	// Connect hsftp
	ftp, err := goftp.Connect("localhost:2121")
	if err != nil {
		logio.Info.Print("goftp Connect err : " + err.Error())
		return nil, err
	}

	// defer ftp.Close()

	// Username / password authentication
	if err = ftp.Login("admin", "123456"); err != nil {
		logio.Info.Print("goftp Login err : " + err.Error())
		return nil, err
	}
	return ftp, nil
}

func fileExistCheck(file uploadRequestFile) string {
	connectedFtp, err := connectFTP()
	if err != nil {
		return err.Error()
	}
	defer connectedFtp.Close()

	// Change current directory
	if err := connectedFtp.Cwd("/GoTest/" + file.CourseCd); err != nil {
		if err := connectedFtp.Mkd("/GoTest/" + file.CourseCd); err != nil {
			logio.Info.Print(err.Error)
			return err.Error()
		}
	}

	// Get directory listing
	var files []string
	files, err = connectedFtp.List("")
	if err != nil {
		logio.Info.Print(err.Error)
		return err.Error()
	}

	// Check upload a file exist
	for _, name := range files {
		temp := strings.Split(name, " ")
		tempFileName := temp[len(temp)-1]
		tempFileName = strings.TrimRight(tempFileName, "\r\n")
		if strings.EqualFold(tempFileName, file.UploadFileName) {
			return "false"
		}
	}

	return "true"
}

func uploadFTP(reqFiles uploadRequestFile) error {
	connectedFtp, err := connectFTP()
	if err != nil {
		return err
	}
	defer connectedFtp.Close()

	// Upload a file
	var file *os.File
	var url = reqFiles.BasePath + "\\" + reqFiles.CourseCd + "_" + reqFiles.UploadFileName

	file, err = os.Open(url)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = connectedFtp.Stor("/GoTest/"+reqFiles.CourseCd+"/"+reqFiles.UploadFileName, file); err != nil {
		return err
	}

	return nil
}
