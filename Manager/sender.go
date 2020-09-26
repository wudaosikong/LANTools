package Manager

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
	if fileInfo,_:=os.Stat(filename);fileInfo.IsDir() {
		filepath.Walk(filename, func (path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				sendFile(path,true,client)
			}else {
				sendFile(path,false,client)
			}
			return nil
		})
	}else{
		 return sendFile(filename,false,client)
	}
	return true
}

func sendFile(filename string,isDir bool,client *net.TCPConn) bool {
	start := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		color.Red("发送前文件打开失败", err)
		return false
	}
	info, _ := file.Stat()
	size := info.Size()
	if !sendName(filename, isDir,client) {
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
		cost := time.Since(start)
		color.Yellow("%s发送成功，用时：%v\n速度：%d Mb//s",filename, cost, size/int64(cost.Seconds())/1024/1024)
	} else {
		color.Red("发送失败")
		return false
	}
	return true
}

func sendName(filename string,isDir bool, client *net.TCPConn) bool {
	var tmp []byte
	if isDir{
		tmp=[]byte("isDir")
	}else {
		tmp = []byte("isFile")
	}
	_, err := client.Write(tmp)
	if err != nil {
		color.Red("发送文件类型失败", err)
		return false
	}
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件类型失败")
		return false
	}


	tmp=[]byte(filename)
	_, err = client.Write(tmp)
	if err != nil {
		color.Red("发送文件名失败", err)
		return false
	}
	n, _ = client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件名失败")
		return false
	}
	return true
}

func sendSize(size int64, client *net.TCPConn) bool {
	tmp := []byte(strconv.FormatInt(size, 10))
	_, err := client.Write(tmp)
	if err != nil {
		color.Red("发送文件大小失败", err)
		return false
	}
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件大小失败")
		return false
	}
	return true
}
