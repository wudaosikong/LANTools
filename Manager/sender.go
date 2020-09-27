package Manager

import (
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
)

func Send(filename string, ip string) bool {
	host, _ := net.ResolveTCPAddr("tcp4", ip+filePort)
	client, err := net.DialTCP("tcp", nil, host)
	if err != nil {
		color.Red("连接对方主机失败", err)
		return false
	}
	defer client.Close()

	// 成功分割线----------------------
	if fileInfo, _ := os.Stat(filename); fileInfo.IsDir() {
		filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				tmp := filepath.ToSlash(path)
				sendFile(tmp, true, client)
			} else {
				tmp := filepath.ToSlash(path)
				sendFile(tmp, false, client)
			}
			return nil
		})
	} else {
		return sendFile(filename, false, client)
	}
	return true
}

func sendFile(filename string, isDir bool, client *net.TCPConn) bool {
	if isDir {
		size, _ := DirSize(filename)
		if !sendDirInfo(filename, client) {
			return false
		}
		if !sendName(filename, client) {
			return false
		}
		if !sendSize(size, client) {
			return false
		}
	} else {
		file, err := os.Open(filename)
		if err != nil {
			color.Red("发送前文件打开失败", err)
			return false
		}
		info, _ := file.Stat()
		size := info.Size()
		if !sendFileInfo(filename, client) {
			return false
		}
		if !sendName(filename, client) {
			return false
		}
		if !sendSize(size, client) {
			return false
		}
		file.Close()
		readerResult := make(chan bool)
		senderResult := make(chan bool)
		counter := make(chan int64)
		data := make(chan []byte, 1024)
		go func() {
			readerResult <- FileReader(filename, data)
		}()
		go func() {
			senderResult <- Sender(client, data, true, counter)
		}()

		go DisplayCounter(size, counter)

		if <-readerResult && <-senderResult {
			// cost := time.Since(start)
			// end := cost.Seconds()
			// if end == 0 {
			// 	end = 1
			// }
			// fmt.Println(end)
			// color.Yellow("%s发送成功，用时：%v\n速度：%d Mb//s", filename, cost, size/int64(end)/1024/1024)
			color.Yellow("%s 发送成功", filename)

		} else {
			color.Red("发送失败")
			return false
		}
	}
	return true
}

func sendDirInfo(filename string, client *net.TCPConn) bool {
	tmpDir := []byte("isDir")
	_, err := client.Write(tmpDir)
	if err != nil {
		color.Red("发送文件夹信息失败", err)
		return false
	}
	tmp := make([]byte, 7)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件夹信息失败")
		return false
	}
	return true
}

func sendFileInfo(filename string, client *net.TCPConn) bool {
	tmpFile := []byte("isFile")
	_, err := client.Write(tmpFile)
	if err != nil {
		color.Red("发送文件信息失败", err)
		return false
	}
	tmp := make([]byte, 200)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件信息失败")
		return false
	}
	return true
}

func sendName(filename string, client *net.TCPConn) bool {
	tmpName := []byte(filename)
	_, err := client.Write(tmpName)
	if err != nil {
		color.Red("发送文件(夹)名失败", err)
		return false
	}
	tmp := make([]byte, 7)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件(夹)名失败")
		return false
	}
	return true
}

func sendSize(size int64, client *net.TCPConn) bool {
	tmpSize := make([]byte, 200)
	tmpSize = []byte(strconv.FormatInt(size, 10))
	_, err := client.Write(tmpSize)
	if err != nil {
		color.Red("发送文件大小失败", err)
		return false
	}
	tmp := make([]byte, 7)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件大小失败")
		return false
	}
	return true
}
