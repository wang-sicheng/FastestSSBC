package main

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/fastestssbc/commonconst"
	"github.com/fastestssbc/net"
	"github.com/fastestssbc/util"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Error("输入的参数有误！")
	}
	nodeID := os.Args[1]
	if nodeID == "N0" {
		//默认N0是主节点
		commonconst.IsLeader = true
	}

	//TPS测试日志重定向
	if commonconst.IsLeader{
		if util.IsExist("log"){
			os.Remove("log")
		}
		f, _ := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND,0755)
		os.Stdout = f
		defer f.Close()
	}

	if addr, ok := commonconst.NodeTable[nodeID]; ok {
		//开启监听
		net.HttpListen(addr)
	} else {
		log.Error("无此节点编号")
	}
}
