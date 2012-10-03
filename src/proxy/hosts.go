package proxy

import (
	"bufio"
	"common"
	"github.com/yinqiwen/godns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//var useGlobalProxy bool
var repoUrls []string
var hostMapping = make(map[string]string)
var hostsEnable bool
var trustedDNS = []string{"8.8.8.8", "8.8.4.4", "208.67.222.222", "208.67.220.220"}
var injectCRLFPatterns = []*regexp.Regexp{}

//type AutoHostConnection struct {
//	http_client  net.Conn
//	https_client net.Conn
//	forwardChan  chan int
//	manager *AutoHost
//}
//
//func (conn *AutoHostConnection) Close() error {
//	if nil != conn.http_client {
//		conn.http_client.Close()
//		conn.http_client = nil
//	}
//	if nil != conn.https_client {
//		conn.https_client.Close()
//		conn.https_client = nil
//	}
//	if nil != conn.forwardChan {
//		close(conn.forwardChan)
//		conn.forwardChan = nil
//	}
//	return nil
//}
//
//func (conn *AutoHostConnection) initHttpsClient(host string) {
//	if nil != conn.https_client {
//		return
//	}
//	host = strings.Split(host, ":")[0]
//	addr, ok := hostMapping[host]
//	if !ok {
//		addr = host
//	}
//	log.Printf("AutoHost use mapping:%s fost host:%s\n", addr, host)
//	var err error
//	conn.https_client, err = net.DialTimeout("tcp", addr+":443", connTimeoutSecs)
//	if nil != err {
//		log.Printf("Failed to dial address:%s for reason:%s\n", addr, err.Error())
//	}
//}
//
//func (conn *AutoHostConnection) initHttpClient(host string) {
//	if nil != conn.http_client {
//		return
//	}
//	host = strings.Split(host, ":")[0]
//	addr, ok := hostMapping[host]
//	if !ok {
//		addr = host
//	}
//	log.Printf("AutoHost use mapping:%s fost host:%s\n", addr, host)
//	var err error
//	conn.http_client, err = net.DialTimeout("tcp", addr+":80", connTimeoutSecs)
//	if nil != err {
//		log.Printf("Failed to dial address:%s for reason:%s\n", addr, err.Error())
//		return
//	}
//}
//
//func (conn *AutoHostConnection) GetConnectionManager() RemoteConnectionManager {
//	return conn.manager
//}
//
//func (conn *AutoHostConnection) writeHttpRequest(req *http.Request) error {
//	var err error
//	index := 0
//	for {
//		err = req.Write(conn.http_client)
//		if nil != err {
//			log.Printf("Resend request since error:%s occured.\n", err.Error())
//			conn.Close()
//			conn.initHttpClient(req.Host)
//		} else {
//			return nil
//		}
//		index++
//		if index == 2 {
//			return err
//		}
//	}
//	return nil
//}
//
//func (auto *AutoHostConnection) Request(conn *SessionConnection, ev event.Event) (err error, res event.Event) {
//	//c := make(chan int)
//	//defer close(c)
//	f := func(local, remote net.Conn) {
//	    io.Copy(remote, local)
//		auto.forwardChan <- 1
//	}
//	switch ev.GetType() {
//	case event.HTTP_REQUEST_EVENT_TYPE:
//		req := ev.(*event.HTTPRequestEvent)
//		if conn.Type == HTTPS_TUNNEL {
//			auto.initHttpsClient(req.RawReq.Host)
//			//try again
//			if nil == auto.https_client {
//				auto.initHttpsClient(req.RawReq.Host)
//			}
//			//log.Printf("Host is %s\n", req.RawReq.Host)
//			log.Printf("Session[%d]Request URL:%s %s\n", ev.GetHash(), req.RawReq.Method, req.RawReq.RequestURI)
//			if nil != auto.https_client {
//				conn.LocalRawConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
//			} else {
//				conn.LocalRawConn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
//				return io.EOF, nil
//			}
//			go f(conn.LocalRawConn, auto.https_client)
//			go f(auto.https_client, conn.LocalRawConn)
//			<-auto.forwardChan
//			<-auto.forwardChan
//			auto.Close()
//			conn.State = STATE_SESSION_CLOSE
//		} else {
//			auto.initHttpClient(req.RawReq.Host)
//			//try again
//			if nil == auto.http_client {
//				auto.initHttpClient(req.RawReq.Host)
//			}
//			if nil == auto.http_client {
//				log.Printf("Failed to connect google http site.\n")
//				conn.LocalRawConn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
//				return nil, nil
//			}
//			log.Printf("Session[%d]Request URL:%s %s\n", ev.GetHash(), req.RawReq.Method, req.RawReq.RequestURI)
//			err := auto.writeHttpRequest(req.RawReq)
//			if nil != err {
//				return err, nil
//			}
//			resp, err := http.ReadResponse(bufio.NewReader(auto.http_client), req.RawReq)
//			if err != nil {
//				return err, nil
//			}
//			resp.Write(conn.LocalRawConn)
//			if resp.Close {
//				conn.LocalRawConn.Close()
//				auto.Close()
//				conn.State = STATE_SESSION_CLOSE
//			} else {
//				conn.State = STATE_RECV_HTTP
//			}
//
//		}
//	default:
//	}
//	return nil, nil
//}

//type AutoHost struct {
//	auths      *util.ListSelector
//	idle_conns chan RemoteConnection
//}
//
//func (manager *AutoHost) GetName() string {
//	return AUTOHOST_NAME
//}
//
////func (manager *AutoHost) GetArg() string {
////	return ""
////}
//func (manager *AutoHost) RecycleRemoteConnection(conn RemoteConnection) {
//	select {
//	case manager.idle_conns <- conn:
//		// Buffer on free list; nothing more to do.
//	default:
//		// Free list full, just carry on.
//	}
//}
//
//func (manager *AutoHost) GetRemoteConnection(ev event.Event) (RemoteConnection, error) {
//	var b RemoteConnection
//	// Grab a buffer if available; allocate if not.
//	select {
//	case b = <-manager.idle_conns:
//		// Got one; nothing more to do.
//	default:
//		// None free, so allocate a new one.
//		g := new(AutoHostConnection)
//		g.manager = manager
//		b = g
//		//b.auth = 
//	} // Read next message from the net.
//	b.Close()
//	return b, nil
//}

func loadHostFile() {
	hostMapping = make(map[string]string)
	os.Mkdir(common.Home+"hosts/", 0755)
	for index, urlstr := range repoUrls {
		resp, err := http.DefaultClient.Get(urlstr)
		if err != nil {
			http_proxy := os.Getenv("http_proxy")
			https_proxy := os.Getenv("https_proxy")
			if addr, exist := common.Cfg.GetProperty("LocalServer", "Listen"); exist {
				_, port, _ := net.SplitHostPort(addr)
				os.Setenv("http_proxy", "http://"+net.JoinHostPort("127.0.0.1", port))
				os.Setenv("https_proxy", "http://"+net.JoinHostPort("127.0.0.1", port))
			}

			defer func() {
				os.Setenv("http_proxy", http_proxy)
				os.Setenv("https_proxy", https_proxy)
			}()
			resp, err = http.DefaultClient.Get(urlstr)
		}
		if err != nil || resp.StatusCode != 200 {
			log.Printf("Failed to fetch host from %s\n", urlstr)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if nil == err {
				hf := common.Home + "hosts/" + "hosts_" + strconv.Itoa(index) + ".txt"
				ioutil.WriteFile(hf, body, 0755)
			}
		}
	}
	files, err := ioutil.ReadDir(common.Home + "hosts/")
	if nil == err {
		for _, file := range files {
			//file.
			content, err := ioutil.ReadFile(common.Home + "hosts/" + file.Name())
			if nil == err {

				reader := bufio.NewReader(strings.NewReader(string(content)))
				for {
					line, _, err := reader.ReadLine()
					if nil != err {
						break
					}
					str := string(line)
					str = strings.TrimSpace(str)

					if strings.HasPrefix(str, "#") || len(str) == 0 {
						continue
					}
					ss := strings.Split(str, " ")
					if len(ss) == 1 {
						ss = strings.Split(str, "\t")
					}
					if len(ss) == 2 {
						k := strings.TrimSpace(ss[1])
						v := strings.TrimSpace(ss[0])
						hostMapping[strings.TrimSpace(k)] = strings.TrimSpace(v)
					}
				}
			}
		}
	}
}

func getOnlineMappingHost(host string) (string, bool) {
	v, exist := hostMapping[host]
	return v, exist
}

func lookupHostPort(hostport string) string {
	if !hostsEnable {
		return hostport
	}
	host, port, err := net.SplitHostPort(hostport)
	if nil != err {
		return hostport
	}
	if v, exist := hostMapping[host]; exist {
		return net.JoinHostPort(v, port)
	}
	addrs, err := godns.LookupHost(host, &godns.LookupOptions{DNSServers: trustedDNS, Cache: true})
	if nil == err && len(addrs) > 0 {
		return net.JoinHostPort(addrs[0], port)
	}
	return hostport
}

func needInjectCRLF(host string) bool {
	for _, regex := range injectCRLFPatterns {
		if regex.MatchString(host) {
			return true
		}
	}
	return false
}

func InitHosts() error {
	if enable, exist := common.Cfg.GetIntProperty("Hosts", "Enable"); exist {
		if enable == 0 {
			hostsEnable = false
			return nil
		}
	}
	hostsEnable = true
	log.Println("Init AutoHost.")
	if dnsserver, exist := common.Cfg.GetProperty("Hosts", "TrustedDNS"); exist {
		trustedDNS = strings.Split(dnsserver, "|")
	}

	if pattern, exist := common.Cfg.GetProperty("Hosts", "InjectCRLF"); exist {
		ps := strings.Split(pattern, "|")
		for _, p := range ps {
			originrule := p
			p = strings.TrimSpace(p)
			p = strings.Replace(p, ".", "\\.", -1)
			p = strings.Replace(p, "*", ".*", -1)
			reg, err := regexp.Compile(p)
			if nil != err {
				log.Printf("[ERROR]Invalid pattern:%s for reason:%v\n", originrule, err)
			} else {
				injectCRLFPatterns = append(injectCRLFPatterns, reg)
			}
		}
	}

	repoUrls = make([]string, 0)
	index := 0
	for {
		v, exist := common.Cfg.GetProperty("Hosts", "HostsRepo["+strconv.Itoa(index)+"]")
		if !exist || len(v) == 0 {
			break
		}
		repoUrls = append(repoUrls, v)
		index++
	}
	if index > 0 {
		go loadHostFile()
	}
	return nil
}