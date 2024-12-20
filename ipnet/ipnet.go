package ipnet

import (
	"fmt"
	"github.com/allape/gogger"
	"net"
	"strings"
)

var l = gogger.New("ipnet")

func DescriptAddress(addr string) ([]string, error) {
	if strings.HasPrefix(addr, ":") {
		interfaces, err := net.Interfaces()
		if err != nil {
			return nil, fmt.Errorf("failed to get interfaces: %v", err)
		}

		bindableAddrs := make([]string, 0, len(interfaces))

		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err != nil {
				l.Warn().Println("Error getting addresses for interface:", iface.Name, err)
				continue
			}

			for _, address := range addrs {
				ipAddr, ok := address.(*net.IPNet)
				if !ok {
					continue
				} else if ipAddr.IP.IsMulticast() || ipAddr.IP.IsLinkLocalMulticast() || ipAddr.IP.IsLinkLocalUnicast() {
					continue
				}

				if ipAddr.IP.To16() == nil {
					bindableAddrs = append(bindableAddrs, fmt.Sprintf("%s%s", ipAddr.IP, address))
				} else {
					bindableAddrs = append(bindableAddrs, fmt.Sprintf("[%s]%s", ipAddr.IP, address))
				}
			}
		}

		return bindableAddrs, nil
	}
	return []string{addr}, nil
}
