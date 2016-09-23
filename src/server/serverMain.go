package main

import (
	"fmt"
	"net/http"
	//"time"
	"config"
	"io"
	"runtime/debug"

	zmq "github.com/pebbe/zmq4"
)

type MQMsgType map[string]string

var (
	MQMsgChan       chan MQMsgType
	MQMsgReturnChan chan string
)

func main() {
	fmt.Println("This is webserver base!")
	err := config.ReadConfigFromFile("config.json")
	if err != nil {
		fmt.Println("ReadConfigFromFile fail,", err)
		return
	}
	fmt.Println("ReadConfigFromFile ok:", config.ServerConfig.PostURL, config.ServerConfig.MQRemote)
	MQMsgChan = make(chan MQMsgType)
	MQMsgReturnChan = make(chan string)
	go MqProcess()

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/"+config.ServerConfig.PostURL, Task)

	//服务器要监听的主机地址和端口号
	err = http.ListenAndServe("127.0.0.1:7777", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func MqProcess() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("MqProcess panic: %v\n\n%s", err, debug.Stack())
		}
	}()
	fmt.Println("mqSocket init...")
	mqSocket, _ := zmq.NewSocket(zmq.REQ)
	defer mqSocket.Close()
	mqSocket.Connect(config.ServerConfig.MQRemote)
	fmt.Println("mqSocket start...")

	for {
		select {
		case msg := <-MQMsgChan:
			fmt.Println("MqProcess new MQMsgChan...")
			MQMsgReturnChan <- DoRequest(mqSocket, msg)
		}
	}
}

func DoRequest(mqSocket *zmq.Socket, msg MQMsgType) string {
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println("DoRequest panic: %v\n\n%s", err, debug.Stack())
	//	}
	//}()

	arglist := make([]interface{}, 0)
	arglist = append(arglist, config.ServerConfig.PostURL)
	arglist = append(arglist, len(msg))
	for k, v := range msg {
		arglist = append(arglist, k)
		arglist = append(arglist, v)
	}
	mqSocket.SendMessageDontwait(arglist)

	//mqSocket.SendMessage(msg)
	//mqSocket.Send("1,1,1,order1,mypay,127.0.0.1-23001", 0) //string uid, string gameAddCount, string CurrencyCount, string orderId, string channel, string gameserver=""
	result, _ := mqSocket.Recv(zmq.DONTWAIT)
	return result
}

func SendRequest(msg map[string]string) string {
	fmt.Println("SendRequest...")

	MQMsgChan <- msg
	fmt.Println("wait MQMsgReturnChan...")
	return <-MQMsgReturnChan
}

func Task(w http.ResponseWriter, req *http.Request) {
	fmt.Println("new requst from http...")
	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Task panic: %v\n\n%s", err, debug.Stack())
		}
	}()
	var msg = make(map[string]string)
	var debugMsg = ""
	for _, formkey := range config.ServerConfig.FormNames {
		formdata, found := req.Form[formkey]
		if !found {
			panic("form no found" + formkey)
		}
		msg[formkey] = formdata[0]
		debugMsg = debugMsg + formkey + ":" + formdata[0] + ","
	}

	result := SendRequest(msg) //("1,1,1,order1,mypay,127.0.0.1-23001") //string uid, string gameAddCount, string CurrencyCount, string orderId, string channel, string gameserver=""
	//s := "transdata:" + transdata[0] + ",sign:" + sign[0] + " result:" + result
	fmt.Println(debugMsg)

	io.WriteString(w, result)
}
