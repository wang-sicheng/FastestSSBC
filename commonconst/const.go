package commonconst

import "github.com/fastestssbc/meta"

const (
	Nodes            = 4
	TransInRedis     = 10000
	TransInBlockStep = 300 //每次区块中交易数的增幅调动
)

var (
	TransInBlock    = 100
	MaxTransInBlock = 10000 //先不用跑完十轮后增加再跑，先直接跑十轮看结果
	Rounds          = 2
	IsLeader        = false
	Ready			= make(chan string)
	PrivateKey 		[]byte //私钥
	PublicKey  		[]byte //公钥
	//公共交易集
	CommonTransList = make([]meta.Transaction,0)
	TotalRound = 0
)

//redis key
const (
	//生成待测交易集
	GenerateTrans = "generate_trans"
	BlockChain    = "block_chain"
)

var Urls = []string{
	"http://127.0.0.1:8000",
	"http://127.0.0.1:8001",
	"http://127.0.0.1:8002",
	"http://127.0.0.1:8003",
}
//var Urls = []string{
	//"http://192.168.1.101:8000",
	//"http://192.168.1.155:8001",
	//"http://192.168.1.102:8002",
	//"http://192.168.1.104:8003",
//}

//生成全局变量-节点池
var NodeTable map[string]string

func init() {
	NodeTable = make(map[string]string)
	NodeTable = map[string]string{
		"N0": ":8000",
		"N1": ":8001",
		"N2": ":8002",
		"N3": ":8003",
	}
}
