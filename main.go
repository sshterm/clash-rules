package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	getDomain("blacklist_full.yaml", "blacklist_full.conf", "foreign", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/blacklist_full.txt")
	getDomain("whitelist_full.yaml", "whitelist_full.conf", "domestic", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/whitelist_full.txt")
	asn_cn()
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
func asn_cn() {
	var data Data
	if res, err := http.Get("https://raw.githubusercontent.com/ncceylan/China-ASN/refs/heads/main/asn_cn.conf"); err == nil {
		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		var data Data
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			data.Payload = append(data.Payload, fmt.Sprintf("IP-ASN,%s,no-resolve", line))
		}
		//data.Payload = append(data.Payload, "+."+"chocolatey.org")

		if d, err := yaml.Marshal(data); err == nil {
			os.WriteFile("asn_cn.yaml", d, 0644)
		}
	}

	if d, err := yaml.Marshal(data); err == nil {
		os.WriteFile("github.yaml", d, 0644)
	}
}
func getDomain(file, file2, name, url string) {
	if res, err := http.Get(url); err == nil {
		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		var data Data
		var domain []string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			data.Payload = append(data.Payload, "+."+line)
			domain = append(domain, fmt.Sprintf("nameserver /%s/%s", line, name))
		}
		//data.Payload = append(data.Payload, "+."+"chocolatey.org")

		if d, err := yaml.Marshal(data); err == nil {
			os.WriteFile(file, d, 0644)
		}
		os.WriteFile(file2, []byte(strings.Join(domain, "\n")), 0644)
	}
}

type Data struct {
	Payload []string `yaml:"payload" `
}
