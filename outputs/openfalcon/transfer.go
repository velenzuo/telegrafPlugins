package openfalcon

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	TransferClientsLock *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients     map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func (f *Openfalcon) SendMetrics(metrics []*MetricValue, resp *TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	addr := f.Addr

	c := getTransferClient(addr)
	if c == nil {
		c = f.initTransferClient(addr)
	}

	if f.updateMetrics(c, metrics, resp) {
		return
	} else {
		fmt.Println("发送失败")
	}
}

func (f *Openfalcon) initTransferClient(addr string) *SingleConnRpcClient {
	var c *SingleConnRpcClient = &SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   10 * time.Second,
	}
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = c

	return c
}

func (f *Openfalcon) updateMetrics(c *SingleConnRpcClient, metrics []*MetricValue, resp *TransferResponse) bool {
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Println("call Transfer.Update fail:", c, err)
		return false
	}
	return true
}

func getTransferClient(addr string) *SingleConnRpcClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()

	if c, ok := TransferClients[addr]; ok {
		return c
	}
	return nil
}

