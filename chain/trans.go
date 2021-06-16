package chain

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/fastestssbc/commonconst"
	"github.com/fastestssbc/meta"
	"github.com/fastestssbc/redis"
	"strconv"
)

var (
	//保存交易与hash值一一对应的映射
	TransHashDataMap = make(map[string]meta.Transaction)
)

//生成待测交易存至redis
func GenerateTrans() {
	//接受交易，验证 存入redis
	trans := generateTx()
	//将交易存至redis
	transB, _ := json.Marshal(trans)
	err := redis.SetIntoRedis(commonconst.GenerateTrans, string(transB))
	if err != nil {
		log.Error(err)
	}
}

func generateTx() []meta.Transaction {
	//生成待测交易集
	txs := make([]meta.Transaction, 0)
	for i := 0; i <= commonconst.TransInRedis; i++ {
		//cur := time.Now()
		tmp := meta.Transaction{
			From:      strconv.Itoa(i), //int(cur.Unix())+
			To:        "To",
			Timestamp: strconv.Itoa(i), //cur.String(),
			Signature: "Signature",
			Message:   "Message",
		}
		txs = append(txs, tmp)
	}
	return txs
}

//从redis中取交易（取一个区块中的交易数）
func PullTrans() []meta.Transaction {
	transS, err := redis.GetFromRedis(commonconst.GenerateTrans)
	if err != nil {
		log.Error(err)
	}
	trans := make([]meta.Transaction, 0)
	err = json.Unmarshal([]byte(transS), &trans)
	if err != nil {
		log.Error("[PullTrans] json err:",err)
	}

	trans = trans[:commonconst.TransInBlock]
	return trans
}
