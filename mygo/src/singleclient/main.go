package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	send(conn,"111111")

	buf := make([]byte,1024) //定义一个切片的长度是1024。
	for  {
		n,err := conn.Read(buf) //接收到的内容大小。
		if err == io.EOF {
			conn.Close()
			break
		}
		fmt.Println(string(buf[:n]))
	}

}

func send(conn *net.TCPConn, s string) {
	_, err := conn.Write([]byte{130 ,33, 1, 4, 101, 99, 104, 111, 28, 24, 16, 89, 111, 117, 32, 97, 114, 101, 32, 119, 101, 108, 99, 111, 109, 101, 46, 0, 0})
	checkError(err)
}

func recive(conn *net.TCPConn){

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}