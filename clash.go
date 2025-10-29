package main

import (
	"bufio"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	getDomain("blacklist_full.yaml", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/blacklist_full.txt")
	getDomain("whitelist_full.yaml", "https://raw.githubusercontent.com/hezhijie0327/GFWList2AGH/refs/heads/main/gfwlist2domain/whitelist_full.txt")
}

func getDomain(file string, url string) {
	if res, err := http.Get(url); err == nil {
		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		var data Data
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			data.Payload = append(data.Payload, "+."+line)
		}
		if d, err := yaml.Marshal(data); err == nil {
			os.WriteFile(file, d, 0644)
		}
	}
}

type Data struct {
	Payload []string `yaml:"payload" `
}
