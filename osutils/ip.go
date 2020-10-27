package osutils

import (
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IpArr []string

func (ip IpArr) Len() int {
	return len(ip)
}

func (ip IpArr) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

func (ip IpArr) Less(i, j int) bool {
	m, _ := strconv.ParseFloat(strings.Split(ip[i], ".")[0], 32)
	n, _ := strconv.ParseFloat(strings.Split(ip[j], ".")[0], 32)
	return m < n
}

func GetIPs() (ip []string) {
	//TODO:针对docker ip->172.17.0.1特征的处理
	addrs, _ := net.InterfaceAddrs()
	var iparr []string
	for _, address := range addrs {
		//ingore 127.0.0.1
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				iparr = append(iparr, ipnet.IP.To4().String())
			}
		}
	}
	if len(iparr) == 0 {
		//when machine started ip not bound ready
		//avoid restarted at high frequencies by supervisor
		time.Sleep(5 * time.Second)
		panic("Can not get local ip")
	}

	//High priority to real ip and the virtual ip followed
	//[ip,vip1,vip2,vip3......]
	sort.Sort(IpArr(iparr))
	return iparr
}
