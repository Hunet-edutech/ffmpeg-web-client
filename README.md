# FFmpeg-Web-Client  

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://golang.org/)

A program registered in the Windows service is executed as a server, which transcodes (encodes) an video using the resources of a web client and transmits the video to a ftp server. 

## Browsers support

| [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_48x48.png" alt="IE / Edge" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>IE / Edge | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/chrome/chrome_48x48.png" alt="Chrome" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>Chrome |
| --------- | --------- |
| IE10, IE11| last version

## Install

    go get github.com/tj3828/ffmpeg-web-client

## Run

1. To run the program, you need to download libraries.

    #### for service

    * [golang/sys](https://godoc.org/golang.org/x/sys) : windows service
    * [gorilla/websocket](https://github.com/gorilla/websocket) : websocket 
    * [jinzhu/copier](https://github.com/jinzhu/copier) : deep copy
    * [dutchcoders/goftp](https://github.com/dutchcoders/goftp) : ftp client
    * [ffmpeg.exe / ffprobe.exe](https://ffmpeg.zeranoe.com/builds/) : ffmpeg and ffprobe .exe file 

    #### for test

    * [goftp/server](https://github.com/goftp/server) : ftp server for test
    

2. Relocate ffmpeg.exe / ffprobe.exe to directory of server.go

3. Run your ftp server for testing

        exampleftpd -root /tmp

4. Check your ftp server's connection

        ftp > open localhost 2121

5. Check your websocket connection of 'sample.html' 

## API Reference

* ftp server 

    | Domain      | Port           | 
    | ----------  | ---------------|
    | /           | 2121           | 

* default : localhost:5050

* To connect websocket 

    * Request 

            Get /echo

    * Progress Json Data

        | Field            | Description                                | Optional   |
        | ---------------- | -------------------------------------------| ---------- |
        | `fileIndex`      | file index                                 | no         |
        | `fileName`       | file name                                  | yes        |
        | `progressTime`   | current encoding time                      | yes        |
        | `progressStatus` | current encoding status (continue / end)   | no         |
        | `totalTime`      | file total time                            | yes        |
        | `progressRate`   | encoding progress rate ( 00.00)            | no         |

* To check a File existence of ftp server

    * Request

            Get /fileExist

        Query parameters :

        | Field            | Description                            | Optional   |
        | ---------------- | ---------------------------------------| ---------- |
        | `uploadFileName` | file name to check(encoded file name)   | no         |
        | `courseCd`       | course code                            | no         |

    * Response

            Content-Type : text/plain;
            exmaple : true(not exist) / false(exist) / error

* To requst encoding and uploading a file

    * Request

            Post /transcoding
            Content-Type : multipart/form-data; 
    

       - Form Data :

            | Field               | Description                            | Optional   |
            | ------------------- | ---------------------------------------| ---------- |
            | `uploadFileData`    | file data                              | no         |
            | `uploadFileInfo`    | information of 'uploadFileData'        | no         |
            | `uploadEncodingInfo`| encoding information                   | no         |

       - 'uploadFileInfo' Json Data :

            | Field             | Description                               | Optional   |
            | ----------------- | ------------------------------------------| ---------- |
            | `fileIndex`       | index                                     | no         |
            | `courseCd`        | course code                               | no         |
            | `courseNo`        | course number(file name)                  | no         |
            | `chapterNo`       | chapter number                            | yes        |
            | `courseName`      | course name                               | yes        |
            | `courseDetailName`| course detail name                        | yes        |
            | `uploadFileName`  | filename ('courseNo' + 'uploadFileExt')   | no         |
            | `uploadFileSize`  | file size                                 | yes        |
            | `uploadFileExt`   | filename extension                        | yes        |

       - 'uploadEncodingInfo' Json Data :

            | Field             | Description                               | Optional   |
            | ----------------- | ------------------------------------------| ---------- |
            | `videoCodec`      | video codec                               | no         |
            | `videoBitrate`    | video bitrate                             | no         |
            | `audioCodec`      | audio codec                               | no         |
            | `audioBitrate`    | audio bitrate                             | no         |
            | `resolution`      | resolution                                | no         |

    * Response

            Content-Type : text/plain;

* To stop trancoding

    * Request

            DELETE /transcoding
    
       - Form Data :

            | Field          | Description                            | Optional   |
            | -------------- | ---------------------------------------| ---------- |
            | `fileIndex`    | index                                  | no         |

    * Response

             Content-Type : text/plain;
             exmaple : true / false

* To check if service program is running

    * Request

            Get /serviceCheck

    * Response

            Content-Type : text/plain;
            exmaple : 1(running) / error


## Notice

* If you want to check logs, you must use Init() method in Run() method.
* Due to the security issue(fakepath) on the web, the process will copy a stored data received from the web client, then encode it and send it to the ftp server.
 Therefore, this process must need to delete a stored data and encoded data on local.
* You can stop to encoding a file, but you can't stop during uploading.



* ftp 서버의 파일 구조

        /GoTest/과목코드/순번.mp4
        
## References

 * [ffmpeg command](https://ffmpeg.org/ffmpeg.html)
