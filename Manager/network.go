package Manager

import (
	"io"
	"net"

	"github.com/fatih/color"
)

func Sender(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(counter)

	for tmp := range data {
		_, err := conn.Write(tmp)
		if err != nil {
			color.Red("发送失败", err)
			return false
		}
		if isDisplay {
			counter <- int64(len(tmp))
		}
	}

	return true
}

var totleSize int64

func Receiver(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(data)
	defer close(counter)

	for {
		tmp := make([]byte, blockSize)
		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			color.Red("接收失败", err)
			return false
		} else if err == io.EOF {
			return true
		}
		if string(tmp[:n]) == "EOF" {
			return true
		} else {
			data <- tmp[:n]
		}
		if isDisplay && string(tmp[:n]) != "EOF" {
			totleSize += int64(n)
			counter <- int64(n)
		}
	}
}
