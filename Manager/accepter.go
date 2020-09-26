package Manager

import (
	"LANTools/tools"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strconv"
)

func Accept() bool {
	host, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0"+filePort)
	// fmt.Printf("%#v", host)
	fmt.Println("监听：", host.IP, host.Port)
	listener, err := net.ListenTCP("tcp", host)
	if err != nil {
		color.Red("监听失败", err)
		return false
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		color.Red("接收客户端失败", err)
		return false
	}
	defer conn.Close()

	// 成功分割线---------------------------------

	filename := acceptName(conn)
	if len(filename) == 0 {
		color.Red("接收文件名有误")
		return false
	} else if filename == "isDir" {
		filename = acceptName(conn)
		fileInfo, _ := os.Stat(filename)
		for n, tmp := 1, filename; IsExit(filename); {
			if fileInfo.IsDir() {
				filename = tmp + "-副本" + strconv.Itoa(n)
			}
			n++
		}
		os.Mkdir(filename, os.ModePerm)
	} else if filename == "isFile" {
		filename = acceptName(conn)
		if len(filename) == 0 {
			color.Red("接收文件名有误2")
			return false
		}
		size := acceptSize(conn)
		diskFree := tools.GetFree()
		if size > diskFree {
			color.Red("磁盘空间不足，请清理磁盘，需要空间：%dGB", size)
			return false
		} else if size == 0 {
			color.Red("接收文件大小有误")
			return false
		}
		fileReceive(filename, conn, size)
	}
	return true
}

func fileReceive(filename string, conn *net.TCPConn, size int64) bool {
	data := make(chan []byte, blockSize)
	writerResult := make(chan bool)
	receiveResult := make(chan bool)
	counter := make(chan int64)
	go func() {
		writerResult <- FileWriter(filename, data)
	}()
	go func() {
		receiveResult <- Receiver(conn, data, true, counter)
	}()

	go DisplayCounter(size, counter)

	if <-writerResult && <-receiveResult {
		color.Yellow("接收文件成功")
	} else {
		color.Red("接收文件失败")
		return false
	}
	return true
}

func acceptName(conn *net.TCPConn) string {
	tmp := make([]byte, 200)
	n, err := conn.Read(tmp)
	if err != nil {
		color.Red("接收文件名失败", err)
		tmp = []byte("fail")
		_, _ = conn.Write(tmp)
		return ""
	}
	res := string(tmp[:n])
	tmp = []byte("success")
	_, _ = conn.Write(tmp)
	return res
}

func acceptSize(conn *net.TCPConn) int64 {
	tmp := make([]byte, 200)
	n, err := conn.Read(tmp)
	if err != nil {
		color.Red("接收文件大小失败", err)
		tmp = []byte("fail")
		_, _ = conn.Write(tmp)
		return 0
	}
	res, _ := strconv.ParseInt(string(tmp[:n]), 10, 64)
	tmp = []byte("success")
	_, _ = conn.Write(tmp)
	return res
}

func DisplayCounter(size int64, counter chan int64) {
	totle := int64(0)
	green := color.New(color.FgGreen)
	for tmp := range counter {
		totle += tmp
		_, _ = green.Printf("总进度：%d%%\r", int(float64(totle)/float64(size)*100))
	}
	fmt.Println("")
}
