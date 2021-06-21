package util

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/DataDog/zstd"
	"github.com/cloudflare/cfssl/log"
	"github.com/fastestssbc/meta"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"

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

//生成本节点的公私钥
func GetKeyPair() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

func CompressV2()  {
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


//压缩
func Compress(in []byte) []byte  {
	fmt.Println("origin size",len(in))
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(in)
	w.Flush()
	fmt.Println("gzip size:", len(b.Bytes()))
	return b.Bytes()


}

//解压
func DeCompress(in []byte)[]byte  {
	var b bytes.Buffer
	b.Write(in)
	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	fmt.Println("ungzip size:", len(undatas))
	return undatas
}


//数字签名
func Sign(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Info("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//签名验证
func VerifySign(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		//panic(err)
		log.Info("验签不通过！")
		return false
	}
	return true
}
