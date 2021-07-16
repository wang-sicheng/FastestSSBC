package chain

import (
	"encoding/json"
	"fmt"
	"github.com/fastestssbc/meta"
	"github.com/fastestssbc/util"
	"strconv"
	"testing"
	"unsafe"
)

func TestChain(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		//for i := 0; i <= commonconst.TransInRedis; i++ {
			//cur := time.Now()
		tmp := meta.Transaction{
			From:      strconv.Itoa(0), //int(cur.Unix())+
			To:        "To",
			Timestamp: strconv.Itoa(0), //cur.String(),
			Signature: "Signature",
			Message:   "927348972983749823749",
		}
			tB, _ := json.Marshal(tmp)
			h := util.CalHash(tB)
			tmp.Hash=h
			//将hash值与交易数据的映射关系进行保存
			TransHashDataMap[h] = tmp
			fmt.Println(unsafe.Sizeof(tmp))
			jsonTmp,_ := json.Marshal(tmp)

		fmt.Println(len(jsonTmp))
	})

}
