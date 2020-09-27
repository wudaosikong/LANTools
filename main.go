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
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}
	var res []string
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip4 := ipnet.IP.String()
						res = append(res, ip4)
					}
				}
			}
		}
	}
	return res
}
