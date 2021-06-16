package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestZTsd(t *testing.T) {
	Compress()
}

func TestGzip(t *testing.T)  {
	//size := 1
	//datasSlice := make([][]byte, size)
	//for i := 0; i < size; i++ {
	//	datasSlice[i], _ = ioutil.ReadFile("images/" + strconv.Itoa(i+1) + ".jpg")
	//	fmt.Println("raw size:", strconv.Itoa(i)+".jpg :", len(datasSlice[i]))
	//}
	//in,_:=os.Open("1.mp4")
	//info,_:=in.Stat()
	//
	//buff:=make([]byte,info.Size())
	//in.Read(buff)
	//
	hash:="dsdadfafdff121e3edeeejfehu4ru4fneffferf"
	buff:=make([]byte,0)
	for i:=0;i<10000;i++{
		buff=append(buff,[]byte(hash)...)
	}
	fmt.Println("len buff",len(buff))

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	w.Write(buff)

	w.Flush()
	fmt.Println("gzip size:", len(b.Bytes()))

	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	fmt.Println("ungzip size:", len(undatas))

}
