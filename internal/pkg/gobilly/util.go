package gobilly

import (
	"bytes"
	"encoding/gob"
)

func encode(data interface{}) []byte {
	//Buffer类型实现了io.Writer接口
	var buf bytes.Buffer
	//得到编码器
	enc := gob.NewEncoder(&buf)
	//调用编码器的Encode方法来编码数据data
	enc.Encode(data)
	//编码后的结果放在buf中
	return buf.Bytes()
}

func decode(data []byte, r interface{}) {
	buf := bytes.NewReader(data)
	//获取一个解码器，参数需要实现io.Reader接口
	dec := gob.NewDecoder(buf)
	//调用解码器的Decode方法将数据解码，用Q类型的q来接收
	dec.Decode(r)
	return
}
