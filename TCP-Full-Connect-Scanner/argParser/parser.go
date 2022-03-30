package argParser

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
)

func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	ipList := addressList.Expand()
	return ipList, err

}

func GetPorts(ports string) ([]int, error) {
	parsedPorts := []int{}
	if ports == "" {
		return parsedPorts, nil
	}
	ranges := strings.Split(ports, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port selection segment:'%s'", r)
			}
			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port number:'%s'", parts[0])
			}
			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port number:'%s'", parts[1])
			}
			if p1 > p2 {
				return nil, fmt.Errorf("invalid port ranges:'%d-%d'", p1, p2)
			}
			for i := p1; i <= p2; i++ {
				parsedPorts = append(parsedPorts, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("invalid port number:'%s'", r)
			} else {
				parsedPorts = append(parsedPorts, port)
			}
		}
	}
	return parsedPorts, nil

}
