package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	getDomain("blacklist_full.yaml", "blacklist_full.conf", "foreign", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/blacklist_full.txt")
	getDomain("whitelist_full.yaml", "whitelist_full.conf", "domestic", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/whitelist_full.txt")
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
		if d, err := yaml.Marshal(data); err == nil {
			os.WriteFile(file, d, 0644)
		}
		os.WriteFile(file2, []byte(strings.Join(domain, "\n")), 0644)
	}
}

type Data struct {
	Payload []string `yaml:"payload" `
}
