package hoststat

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

type ClientInfo struct {
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
	Result struct {
		Balance         float64 `json:"balance"`
		Blocks          int     `json:"blocks"`
		Connections     int     `json:"connections"`
		Difficulty      float64 `json:"difficulty"`
		Errors          string  `json:"errors"`
		Keypoololdest   int     `json:"keypoololdest"`
		Keypoolsize     int     `json:"keypoolsize"`
		Network         string  `json:"network"`
		Paytxfee        float64 `json:"paytxfee"`
		Protocolversion int     `json:"protocolversion"`
		Proxy           string  `json:"proxy"`
		Relayfee        float64 `json:"relayfee"`
		Timeoffset      int     `json:"timeoffset"`
		Version         int     `json:"version"`
		Walletversion   int     `json:"walletversion"`
	} `json:"result"`
}

type MasterNodeStatus struct {
	txid          string
	vout          int
	status        string
	Protocol      string
	Payee         string
	Lastseen      int64
	Activeseconds int64
	Lastpaidtime  int64
	LastPaidBlock int
	IPAddress     string
}

type MasterNodeSum struct {
	Total   int
	Enabled int
	Qualify int
}

func getBlockHeight() int {
	var clientInfo ClientInfo
	cmd := exec.Command("vds-cli", "getinfo")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd error in getNodeStatus,Err:%v", err)
	}

	err = json.Unmarshal(out, &clientInfo)
	if err != nil {
		log.Fatalf("Can't Unmarshal data from chainInfo,ERR:%v", err)
	}
	return clientInfo.Result.Blocks
}

func getMasterNodeSum() MasterNodeSum {
	var masterNodeSum MasterNodeSum
	cmd := exec.Command("vds-cli", "masternode", "count", "all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd error in getNodeStatus,Err:%v", err)
	}

	reg := regexp.MustCompile(`[0-9]+`)
	data := reg.FindAllString(string(out), -1)
	masterNodeSum.Total, _ = strconv.Atoi(data[0])
	masterNodeSum.Enabled, _ = strconv.Atoi(data[1])
	masterNodeSum.Qualify, _ = strconv.Atoi(data[2])

	return masterNodeSum
}

func getMasterNodeStatus(address string) MasterNodeStatus {
	var masterNodeStatus MasterNodeStatus
	c1 := exec.Command("vds-cli", "masternodelist", "full")
	c2 := exec.Command("grep", address)
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	ourStr := b2.String() //输出执行结果
	reg := regexp.MustCompile(`[0-9A-Za-z]+`)
	data := reg.FindAllString(ourStr, -1)
	masterNodeStatus.txid = data[0]
	masterNodeStatus.vout, _ = strconv.Atoi(data[1])
	masterNodeStatus.status = data[2]
	masterNodeStatus.Protocol = data[3]
	masterNodeStatus.Payee = data[4]
	masterNodeStatus.Lastseen, _ = strconv.ParseInt(data[5], 10, 64)
	masterNodeStatus.Activeseconds, _ = strconv.ParseInt(data[6], 10, 64)
	masterNodeStatus.Lastpaidtime, _ = strconv.ParseInt(data[7], 10, 64)
	masterNodeStatus.LastPaidBlock, _ = strconv.Atoi(data[8])
	masterNodeStatus.IPAddress = data[9]
	return masterNodeStatus
}
