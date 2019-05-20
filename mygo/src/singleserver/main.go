package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"net"
	"os"
	"thriftexample/gen-go/echo"
	"time"
)

func main() {
	serverStart()
	//send(nil,0)
}

func serverStart() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(3 * time.Second)) // set 2 minutes timeout

	defer conn.Close() // close connection before exit
	for {
		request := make([]byte, 128) // set maxium request length to 128B to prevent flood attack
		read_len, err := conn.Read(request)
		fmt.Println(request[:read_len])
		tmbClient(request, read_len)
		if err != nil {
			checkError(err)
		}

		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finishe
		break
	}
}

func tmbClient(request []byte, read_len int){

	tmb := thrift.NewTMemoryBuffer()
	tmb.Write(request[:read_len-1])
	readcount:= len([]byte("fdsfdasfsafsfaf"))+1+1
	protocolFactory := thrift.NewTCompactProtocolFactory()
	tp := protocolFactory.GetProtocol(tmb)
	tp.WriteFieldBegin("name",thrift.STRING,int16(1))
	tp.WriteString("fdsfdasfsafsfaf")

	tp.WriteFieldEnd()
	tp.WriteFieldStop()
	godata := make([]byte,1024)

	tmb.Read(godata)
	fmt.Println("new data ",godata)
	fmt.Println("new data ",godata[:read_len+readcount])


	service := "127.0.0.1:9898"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write(godata[:read_len+readcount])
	checkError(err)
	defer conn.Close()

	r := bufio.NewReader(conn)
	//reader.Buffered()
	//r := bufio.NewReaderSize(reader,32)

	result := make([]byte, 0)
	tmp := make([]byte, 12)
	for {
		len, _ := r.Read(tmp)
		if r.Buffered()==0{
			break
		}
		fmt.Println("tmp", tmp[:len])
		result = append(result,tmp[:len]...)
	}


	fmt.Println("result",result)
	fmt.Println("result",string(result))
}


func tcpClient(request []byte, read_len int) {
	service := "127.0.0.1:9898"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write(request[:read_len])
	checkError(err)

	r := bufio.NewReader(conn)
	//result := make([]byte, 12)
	tmp := make([]byte, 128)

	len,_ := r.Read(tmp)
	fmt.Println("tmp", tmp[:len])
	fmt.Println("tmp", string(tmp[:len]))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func send(request []byte, read_len int) {
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()

	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9898"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9898", " ", err)
		os.Exit(1)
	}
	defer transport.Close()
	ctx := context.Background()
	useTransport, err := transportFactory.GetTransport(transport)
	client := thrift.NewTStandardClient(protocolFactory.GetProtocol(useTransport), protocolFactory.GetProtocol(useTransport))

	req := &echo.EchoReq{Msg: "You are welcome."}
	var _args0 echo.EchoEchoArgs
	_args0.Req = req
	temp, err := thrift.NewTSerializer().Write(ctx, &_args0)
	fmt.Println(temp)

	var _result1 echo.EchoEchoResult
	if err = client.Call(ctx, "echo", &_args0, &_result1); err != nil {
		return
	}

	log.Println("response:", _result1.GetSuccess())
	fmt.Println("well done")
}

func sendByte(data []byte, data_len int) {

	ctx := context.Background()
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9898"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9898", " ", err)
		os.Exit(1)
	}
	defer transport.Close()
	transport.Write(data[:data_len])
	transport.Flush(ctx)

	res := make([]byte, 1024)
	if transport.IsOpen() {
		len, err := transport.Read(res)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
		}
		fmt.Println(len)
		fmt.Println("dsdsddsa", BytesToInt(res[0:3]))
		fmt.Println(res[:len])
	}

	fmt.Println("thrift end ... ")
}

func sendTMessage(data []byte, data_len int) {
	ctx := context.Background()
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9898"))
	protocolFactory := thrift.NewTCompactProtocolFactory()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9898", " ", err)
		os.Exit(1)
	}
	defer transport.Close()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	useTransport, err := transportFactory.GetTransport(transport)

	useTransport.Write(data[:data_len])
	useTransport.Flush(ctx)
	tt := protocolFactory.GetProtocol(useTransport)
	method, typeid, seqid, _ := tt.ReadMessageBegin()
	fmt.Println("read head info : ", method, typeid, seqid)

	result := echo.NewEchoEchoResult()
	var ts thrift.TStruct
	ts = result
	result.Read(tt)
	tt.ReadMessageEnd()
	fmt.Println("read end : ", ts)

}

func processFun(data []byte, data_len int) {
	ctx := context.Background()
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9898"))
	protocolFactory := thrift.NewTCompactProtocolFactory()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9898", " ", err)
		os.Exit(1)
	}
	defer transport.Close()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	useTransport, err := transportFactory.GetTransport(transport)
	tt := protocolFactory.GetProtocol(useTransport)
	useTransport.Write(data[:data_len])
	useTransport.Flush(ctx)

	tclient := thrift.NewTStandardClient(tt, tt)

	method, _, seqId, _ := tt.ReadMessageBegin()
	result := echo.NewEchoEchoResult()
	tclient.Recv(tt, seqId, method, result)
	fmt.Println(result)
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
