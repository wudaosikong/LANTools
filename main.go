package main

import (
	"LANTools/Manager"
	"fmt"
	"net"
)

var ch = make(chan []byte, 10)

func main() {
	fmt.Println("开始启动！")
	LocalIps := GetIntranetIp()
	fmt.Print("你的ID是：")
	for _, LocalIp := range LocalIps {
		fmt.Println(LocalIp)
	}
	fmt.Println("输入 help 以获取更多帮助")

	gui := Manager.GUI{}
	gui.LocalIP = LocalIps
	gui.Render()
}

func GetIntranetIp() []string {
	address, _ := net.InterfaceAddrs()
	// fmt.Printf("%#v", address)
	var res []string
	for _, addr := range address {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				ip4 := ip.IP.To4().String()
				if ip4[:7] != "169.254" {
					res = append(res, ip4)
				}
			}
		}
	}
	return res
}
