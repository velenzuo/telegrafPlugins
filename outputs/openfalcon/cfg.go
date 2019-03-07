package openfalcon

import (
	"fmt"
	"log"
	"os"
	"sync"
)
type TransferResponse struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}

func (this *TransferResponse) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
		this.Total,
		this.Invalid,
		this.Latency,
		this.Message,
	)
}

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Addrs    string `json:"addr"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
	MountPoint  []string `json:"mountPoint"`
}

type GlobalConfig struct {
	Transfer      *TransferConfig   `json:"transfer"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	return config
}

func Hostname() (string, error) {
	hostname := ""
	if hostname != "" {
		return hostname, nil
	}

	if os.Getenv("FALCON_ENDPOINT") != "" {
		hostname = os.Getenv("FALCON_ENDPOINT")
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := ""
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return ip
}


