package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/DataDog/zstd"
	"github.com/cloudflare/cfssl/log"
	"github.com/fastestssbc/meta"
	jsoniter "github.com/json-iterator/go"
	//"github.com/klauspost/compress/zstd"
	//jsoniter "github.com/json-iterator/go"
	"os"
	//"compress/gzip"
	"strconv"
)

//使用速度最快的json工具
var FastestJson = jsoniter.ConfigCompatibleWithStandardLibrary
//go build -tags=jsoniter

//计算hash值
func CalHash(data []byte) string {
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func CalBlockHash(block meta.Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.PrevHash + block.Merkle + block.Signature
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//判断文件或文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		log.Info(err)
		return false
	}
	return true
}

//func Compress(in io.Reader, out io.Writer) error {
//	enc, err := zstd.NewWriter(out)
//	if err != nil {
//		return err
//	}
//	_, err = io.Copy(enc, in)
//	if err != nil {
//		enc.Close()
//		return err
//	}
//	return enc.Close()
//}


//func Decompress(in io.Reader, out io.Writer) error {
//	d, err := zstd.NewReader(in)
//	if err != nil {
//		return err
//	}
//	defer d.Close()
//	// Copy content...
//	_, err = io.Copy(out, d)
//	return err
//}

func Compress()  {
	file,err:=os.Open("1.mp4")

	fmt.Println(err)
	defer file.Close()
	fileInfo,err:=file.Stat()
	fmt.Println(err)
	fmt.Println("开始大小：",fileInfo.Size())
	buffer := make([]byte, fileInfo.Size())
	_,err=file.Read(buffer)
	fmt.Println("传入大小:",len(buffer))
	ret,err:=zstd.CompressLevel(nil,buffer,5)
	fmt.Println(err)
	fmt.Println("传出大小:",len(ret))
}



