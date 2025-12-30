package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	githubIP()
}

var ip map[string][]string

func githubIP() {
	var data Data
	if res, err := http.Get("https://api.github.com/meta"); err == nil {
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&ip)
		for _, v := range ip {
			for _, v2 := range v {
				_, _, err := net.ParseCIDR(v2)
				if err != nil {
					continue
				}
				data.Payload = append(data.Payload, v2)
			}
		}
	}
	// 对IP地址进行排序
	sort.Slice(data.Payload, func(i, j int) bool {
		return compareCIDR(data.Payload[i], data.Payload[j])
	})
	if d, err := yaml.Marshal(data); err == nil {
		os.WriteFile("github.yaml", d, 0644)
	}
}

// compareCIDR 比较两个CIDR地址
func compareCIDR(a, b string) bool {
	// 分离IP和前缀长度
	ipA, _, _ := net.ParseCIDR(a)
	ipB, _, _ := net.ParseCIDR(b)

	// 比较IP地址
	if ipA.To4() != nil && ipB.To4() != nil {
		// 都是IPv4，逐字节比较
		for i := 0; i < 4; i++ {
			if ipA.To4()[i] != ipB.To4()[i] {
				return ipA.To4()[i] < ipB.To4()[i]
			}
		}
	} else if ipA.To16() != nil && ipB.To16() != nil {
		// 都是IPv6，逐字节比较
		for i := 0; i < 16; i++ {
			if ipA.To16()[i] != ipB.To16()[i] {
				return ipA.To16()[i] < ipB.To16()[i]
			}
		}
	}

	// 如果IP相同，比较前缀长度
	aParts := strings.Split(a, "/")
	bParts := strings.Split(b, "/")
	if len(aParts) > 1 && len(bParts) > 1 {
		prefixA, _ := strconv.Atoi(aParts[1])
		prefixB, _ := strconv.Atoi(bParts[1])
		return prefixA < prefixB
	}

	return a < b
}

type Data struct {
	Payload []string `yaml:"payload" `
}
