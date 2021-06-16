package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cfssl/log"
	"github.com/deckarep/golang-set"
	"github.com/fastestssbc/chain"
	"github.com/fastestssbc/commonconst"
	"github.com/fastestssbc/meta"
	"github.com/fastestssbc/util"
	//"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type HClient struct {
	Client *http.Client
	Url    string
}

var (
	round       = 0
	httpClients = make([]*HClient, 0)
	//交易map
	transMap = make(map[int]map[string][]string)
	//第一轮投票统计
	//round1Map = make(map[string][]meta.Vote)
	round1Map sync.Map
	//第二轮投票统计
	//round2Map = make(map[string][]meta.ReVote)  //高并发场景下会出现协程并发写panic，故使用协程安全的
	round2Map sync.Map
	StartTime time.Time
	EndTime   time.Time
	Flag      bool
)

type Server struct {
}

//初始化通讯客户端且保持alive
func init() {
	for _, u := range commonconst.Urls {
		hc := &HClient{
			Url: u,
		}
		hc.Client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				DisableKeepAlives:   false,
			},
		}
		httpClients = append(httpClients, hc)
	}
}

func HttpListen(addr string) {
	r := gin.Default()
	//r.Use(gzip.Gzip(gzip.DefaultCompression))

	//触发性能测试
	r.GET("/speedTest", speedTest)
	//触发其他节点
	r.POST("/inform", inform)
	//广播交易
	r.POST("/broadcastTrans", broadcastTrans)
	//处理新区块
	r.POST("/newBlock", newBlock)
	//处理第一轮投票结果
	r.POST("/blockVoteRound1", blockVoteRound1)
	//获取到指定高度区块的交易列表
	r.POST("/blockVoteRound2", blockVoteRound2)

	err := r.Run(addr)
	if err != nil {
		log.Error(err)
	}

}

func speedTest(ctx *gin.Context) {
	//主节点触发，通知其他所有的节点
	Broadcast("inform", []byte("inform"))
	ctx.JSON(http.StatusOK,"Start Speed Test")
}
func inform(ctx *gin.Context) {
	//所有节点收到触发，开始广播交易
	//触发使用redis进行交易生成
	if !Flag {
		chain.GenerateTrans()
	}
	StartTime=time.Now()
	//开始广播交易
	SendTrans()
	ctx.JSON(http.StatusOK, "ok")
}

//节点收到其他节点广播的交易集进行处理
func broadcastTrans(ctx *gin.Context) {
	//如果是主节点的话,接收所有的节点的交易，并且查找公共集
	if commonconst.IsLeader {
		tH := meta.TransHash{}
		err := ctx.ShouldBindJSON(&tH)
		if err != nil {
			log.Error(err)
		}
		//记录收到的交易
		if tMap, exists := transMap[round]; exists {
			tMap[ctx.Request.RemoteAddr] = tH.TransHashs
			transMap[round] = tMap
		} else {
			temMap := make(map[string][]string)
			temMap[ctx.Request.RemoteAddr] = tH.TransHashs
			transMap[round] = temMap
		}
		//判断是否收到了足量的请求
		c := len(transMap[round])
		if c == commonconst.Nodes {
			//说明收到了所有节点的广播
			//使用性能强悍的golang-set寻找公共集
			commonTrans := findCommonTrans(transMap[round])
			newBlock := meta.Block{TX: commonTrans}
			newBlock = chain.GenerateNewBlock(chain.CurrentBlock, newBlock)
			//开始广播区块
			newBlockB, _ := json.Marshal(newBlock)
			//log.Info("生成新区块")
			Broadcast("newBlock", newBlockB)

		} else {
			//说明没有收到足够量的节点的广播
			log.Info("Has Received TransHash Count:", c)
			return
		}
	} else {
		log.Info("Not Leader")
	}

}

func findCommonTrans(m map[string][]string) []meta.Transaction {
	//使用性能强悍的golang-set寻找公共集
	commonTranHashs := mapset.NewSet()
	for _, hashs := range m {
		for _, h := range hashs {
			commonTranHashs.Add(h)
		}
	}

	commonTranHashsSlice := commonTranHashs.ToSlice()

	commonTrans := make([]meta.Transaction, 0)
	//基于公共交易集生成新区块
	for _, t := range commonTranHashsSlice {
		//基于hash值与交易信息的映射来获取到交易信息
		if trans,exists:=chain.TransHashDataMap[t.(string)];exists{
			commonTrans = append(commonTrans, trans)
		}else {
			log.Warning("交易漏失,hash=",t)
		}
	}
	return commonTrans
}

//节点接收到主节点广播的区块的处理
func newBlock(ctx *gin.Context) {
	nB := meta.Block{}
	err := ctx.ShouldBindJSON(&nB)
	if err != nil {
		log.Error(err)
	}


	//nBB,_:=json.Marshal(nB)
	//log.Info("接收到新区块",string(nBB))

	//验证主节点打包的区块是否合法
	vote := chain.VerifyBlock(nB)
	//设置新区块
	if vote {
		chain.NewBlock = nB
	}
	//进行第一轮投票
	v := &meta.Vote{Sender: ctx.Request.Host, Hash: nB.Hash, Vote: vote}
	vB, err := json.Marshal(v)
	if err != nil {
		log.Error(err)
	}
	//广播第一轮投票
	Broadcast("blockVoteRound1", vB)
}

func blockVoteRound1(ctx *gin.Context) {
	vote := meta.Vote{}
	err := ctx.ShouldBindJSON(&vote)
	if err != nil {
		log.Error(err)
	}
	r := strconv.Itoa(round)
	//投票数统计
	key := r + "_" + "round1_" + vote.Hash
	if val,ok:=round1Map.Load(key);ok{
		votes:=val.([]meta.Vote)
		votes=append(votes,vote)
		round1Map.Store(key,votes)
		val,_=round1Map.Load(key)
		votes=val.([]meta.Vote)
		if len(votes)==commonconst.Nodes{
			//说明广播已收齐
			log.Info("Received all votes,start round2")
			//投票鉴别是否投同意
			count := 0
			for _, v := range votes {
				if v.Vote {
					//同意票+1
					count++
				}
			}

			re := false
			if float64(count) > float64(commonconst.Nodes)*0.75 {
				re = true
			}
			//开始第二轮投票
			rv := meta.ReVote{
				Sender: ctx.Request.Host,
				Vote:   votes,
				Hash:   vote.Hash,
				V:      re,
			}
			rvB, _ := json.Marshal(rv)
			Broadcast("blockVoteRound2", rvB)
		}else {
			//说明广播还没收齐
			log.Info("Round1 Not receive all votes:", votes)
			return
		}
	}else {
		votes:=make([]meta.Vote,0)
		votes=append(votes,vote)
		round1Map.Store(key,votes)
		return
	}
}

//第二轮投票处理
func blockVoteRound2(ctx *gin.Context) {
	reVote := meta.ReVote{}
	err := ctx.ShouldBindJSON(&reVote)
	if err != nil {
		log.Error(err)
	}
	log.Info("收到第二轮投票:",reVote)
	//投票数统计
	r := strconv.Itoa(round)
	key := r + "_" + "round2_" + reVote.Hash

	if val,ok:=round2Map.Load(key);ok{
		votes:=val.([]meta.ReVote)
		votes=append(votes,reVote)
		round2Map.Store(key,votes)
		val,_=round2Map.Load(key)
		votes=val.([]meta.ReVote)
		if len(votes)==commonconst.Nodes{
			//说明投票已收齐
			//先进行区块hash核验
			if reVote.Hash == chain.NewBlock.Hash && reVote.Hash != chain.CurrentBlock.Hash {
				count := 0
				for _, v := range votes {
					if v.V {
						//同意票+1
						count++
					}
				}
				if float64(count) > float64(commonconst.Nodes)*0.75 {
					//第二轮投票过3/4,新区块固化
					chain.StoreNewBlock()
					//本轮交易共识上链流程TPS估算
					CalCulTPS()
					//判断是否开启新的一轮交易共识上链流程
					JudgeNextRound()
				}
			} else {
				log.Error("区块Hash校验失败")
				return
			}
		}else {
			//说明票还没收齐
			//说明广播还没收齐
			log.Info("Round2 Not receive all votes:", votes)
			return
		}
	}else {
		votes:=make([]meta.ReVote,0)
		votes=append(votes,reVote)
		round2Map.Store(key,votes)
		return
	}
}

//TPS计算
func CalCulTPS() {
	if commonconst.IsLeader {
		EndTime = time.Now()
		dura := EndTime.Sub(StartTime).Seconds()
		//log.Info("duration: ",t2.Sub(t1))
		//totalTrans := float64((round + 1) * commonconst.TransInBlock)
		tps := float64(commonconst.TransInBlock) / dura
		fmt.Println("TransCount=", commonconst.TransInBlock, ", Sequence=", round, ", Duration=", dura, ", TPS=", tps)
	}
}

func JudgeNextRound() {
	Flag = true
	if round+1 < commonconst.Rounds {
		//未达到设置的轮次，继续
		round++
		if commonconst.IsLeader {
			//主节点开启新一轮的触发
			go Broadcast("inform", []byte("inform"))
		}
	} else {
		//说明在区块指定交易数情况下跑到了指定的轮次，更改区块内交易数参数
		//Step1:轮次重置
		round = 0
		//Step2:交易数按固定增幅增加
		commonconst.TransInBlock = commonconst.TransInBlock + commonconst.TransInBlockStep
		if commonconst.IsLeader && commonconst.TransInBlock<=commonconst.MaxTransInBlock{
			//主节点开启新一轮的触发
			go Broadcast("inform", []byte("inform"))
		}
	}
}

func SendTrans() {
	//先从redis中取交易
	log.Infof("第%d轮开始", round)
	trans := chain.PullTrans()
	//transB,_:=util.Fastestjson.Marshal(trans)

	//广播的不是交易，而是hash
	tT := meta.TransHash{}
	//先将每个交易进行hash
	hashs := make([]string, 0)
	for _, t := range trans {
		tB, _ := json.Marshal(t)
		h := util.CalHash(tB)
		//将hash值与交易数据的映射关系进行保存
		chain.TransHashDataMap[h] = t
		hashs = append(hashs, h)
	}

	tT.BlockHash = chain.CurrentBlock.Hash
	tT.TransHashs = hashs
	tTB, _ := json.Marshal(tT)

	Broadcast("broadcastTrans", tTB)

}

//广播
func Broadcast(s string, reqBody []byte)  {
	for _, client := range httpClients {
		//go send(client,s,reqBody) //并发发送
		send(client,s,reqBody)
	}
}

func send(c *HClient,s string,reqBody []byte)  {
	endPoint := c.Url + "/" + s
	req, err := NewPost(endPoint, reqBody)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")
	err = c.SendReq(req, nil)
	if err != nil {
		log.Error(err)
	}
}

func NewPost(endPoint string, reqBody []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(reqBody))
	if err != nil {
		log.Error("[NewPost] err:",err)
		return nil, errors.New("Failed posting to " + endPoint)
	}
	return req, nil
}
func (c *HClient) SendReq(req *http.Request, result interface{}) (err error) {
	_, err = c.Client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "%s failure of request", req.Method)
	}
	return nil
}
