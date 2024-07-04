package gorpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

// /////////////////////////////////////////////////////////
// gorpc Log区
// /////////////////////////////////////////////////////////
// 控制台日志输出
func Log(level string, format string, args ...interface{}) {
	time := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(fmt.Sprintf("[%v] [%v]", time, level), fmt.Sprintf(format, args...))
}
func LogInfo(format string, args ...interface{}) {
	Log("INFO ", format, args...)
}
func LogWarn(format string, args ...interface{}) {
	Log("WARN ", format, args...)
}
func LogError(format string, args ...interface{}) {
	Log("ERROR", format, args...)
}

// /////////////////////////////////////////////////////////
// gorpc 默认配置区
// /////////////////////////////////////////////////////////
const DefaultRpcPort = 9797

// /////////////////////////////////////////////////////////
// gorpc server 区
// /////////////////////////////////////////////////////////
func GetDefaultRpcServer() *RpcServer {
	ret := &RpcServer{
		Server:   rpc.NewServer(),
		Services: map[string]interface{}{},
		Port:     DefaultRpcPort,
	}

	return ret
}
func GetRpcServer(port int) *RpcServer {
	ret := &RpcServer{
		Server:   rpc.NewServer(),
		Services: map[string]interface{}{},
		Port:     port,
	}

	return ret
}

type RpcServer struct {
	Server   *rpc.Server
	Services map[string]interface{}
	Port     int
}

func (svr *RpcServer) AddService(name string, svc interface{}) *RpcServer {
	svr.Services[name] = svc
	return svr
}

func (svr *RpcServer) Run() *sync.WaitGroup {
	if svr.Port == 0 {
		svr.Port = DefaultRpcPort
	}
	LogInfo("rpc server start run at port: %v", svr.Port)

	wg := sync.WaitGroup{}
	for key, val := range svr.Services {
		err := svr.Server.RegisterName(key, val)
		LogInfo("register rpc service: %v", key)
		if err != nil {
			LogWarn("register rpc service err: %v", err)
		}
	}

	addrStr := fmt.Sprintf(":%v", svr.Port)
	LogInfo("rpc server listen %v", addrStr)
	listener, err := net.Listen("tcp", addrStr)
	if err != nil {
		LogError("rpc listen tcp err: %v\n", err)
		return &wg
	}

	LogInfo("rpc server run in goroutine")
	wg.Add(1)
	go func() {
		for {
			conn, err := listener.Accept()
			LogInfo("rpc server accept client: %v", conn)
			if err != nil {
				LogWarn("accept rpc client err: %v\n", err)
			}
			svr.Server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
		wg.Done()
		LogInfo("rpc server shutdown")
	}()

	LogInfo("rpc server run ok.")
	return &wg
}

func (svr *RpcServer) Stop() {
	if svr.Server != nil {

	}
}

// /////////////////////////////////////////////////////////
// gorpc client 区
// /////////////////////////////////////////////////////////
func GetDefaultRpcClient(host string) *RpcClient {
	ret := &RpcClient{
		Host: host,
		Port: DefaultRpcPort,
	}

	return ret
}

func GetRpcClient(host string, port int) *RpcClient {
	ret := &RpcClient{
		Host: host,
		Port: port,
	}

	return ret
}

type RpcClient struct {
	Host   string
	Port   int
	Conn   net.Conn
	Client *rpc.Client
}

func (cli *RpcClient) Run() {
	if cli.Port == 0 {
		cli.Port = DefaultRpcPort
	}
	if cli.Host == "" {
		cli.Host = "127.0.0.1"
	}
	addrStr := fmt.Sprintf("%v:%v", cli.Host, cli.Port)
	LogInfo("rpc client connect to: %v", addrStr)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		LogError("connect rpc server err: %v\n", err)
		return
	}

	cli.Conn = conn
	LogInfo("rpc client connect conn: %v", conn)

	LogInfo("rpc client has built ok.")
	cli.Client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(cli.Conn))
}

func (cli *RpcClient) Stop() {
	if cli.Client != nil {
		err := cli.Client.Close()
		if err != nil {
			LogWarn("close rpc client err: %v", err)
		}
		cli.Client = nil
	}
}

func (cli *RpcClient) Call(method string, req interface{}, resp interface{}) error {
	if cli.Client == nil {
		cli.Run()
	}
	LogInfo("rpc invoke %v @ %v:%v", method, cli.Host, cli.Port)
	err := cli.Client.Call(method, req, resp)
	return err
}
