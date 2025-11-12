package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	getIP()
}

var ip map[string][]string

func getIP() {
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
	if d, err := yaml.Marshal(data); err == nil {
		os.WriteFile("github.yaml", d, 0644)
	}
}

type Data struct {
	Payload []string `yaml:"payload" `
}
