package commonconst

const (
	Nodes            = 4
	TransInRedis     = 300000
	TransInBlockStep = 100 //每次区块中交易数的增幅调动
)

var (
	TransInBlock    = 100000
	MaxTransInBlock = 1000
	Rounds          = 10
	IsLeader        = false
)

//redis key
const (
	//生成待测交易集
	GenerateTrans = "generate_trans"
)

var Urls = []string{
	"http://127.0.0.1:8000",
	"http://127.0.0.1:8001",
	"http://127.0.0.1:8002",
	"http://127.0.0.1:8003",
}

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
