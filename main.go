package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Server struct {
	IP         string
	AllowedIPs []string
	Type       string
	Metadata   map[string]string
	PublicKey  string
	PrivateKey string
}

type WireguardConfig struct {
	IP string
	PrivateKey string
	AllowedIPs string
}

type ServerRequest struct {
	Type     string
	Metadata map[string]string
}

func getPublicIP() string {
	res, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return "unknown"
	}

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "unknown"
	}

	return string(ip)
}

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: ./central-client [endpoint] [secret] [type] [template config file]")
	}
	templateString, err := ioutil.ReadFile(os.Args[4])
	if err != nil {
		log.Fatalf("failed to read config template\n%s", err.Error())
	}

	wgTemplate, err := template.New("wg").Parse(string(templateString))

	if err != nil {
		log.Fatalf("failed to parse config template\n%s", err.Error())
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	reqData, _ := json.Marshal(ServerRequest{
		Type:     os.Args[3],
		Metadata: map[string]string {
			"Hostname": hostname,
			"IP": getPublicIP(),
		},
	})

	req, _ := http.NewRequest("POST", os.Args[1], bytes.NewBuffer(reqData))
	req.Header.Set("Authorization", "Bearer " + os.Args[2])

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Fatalf("failed to send POST request to create server\n%s", err.Error())
	}

	resJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read server response\n%s", err.Error())
	}

	var server Server
	err = json.Unmarshal(resJson, &server)

	pretty, _ := json.MarshalIndent(server, "", "\t");

	_, _ = fmt.Fprintf(os.Stderr, "server info: %s\n", string(pretty))

	if err != nil {
		log.Fatalf("failed to parse server response\n%s\n%s", resJson, err.Error())
	}

	err = wgTemplate.Execute(os.Stdout, WireguardConfig{
		IP:         server.IP,
		PrivateKey: server.PrivateKey,
		AllowedIPs: strings.Join(server.AllowedIPs, ", "),
	})

	if err != nil {
		log.Fatalf("failed to template config\n%s", err.Error())
	}

}