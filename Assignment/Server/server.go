package Server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	Addrs_flag = [5]bool{
		false,
	}

	Addrs = [5]string{
		"169.254.6.85",
	}

	dial_conns  = make(map[string]net.Conn)
	IP2Username = make(map[string]string)
	pNumber     = 0
	localStamp  []int
	holdback    = make(map[string]string)
	stayLong    = make(map[string]string)
)

func getAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Failed to get the IP address: %v", err)
	}

	var return_ip string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return_ip = ipnet.IP.String()
			}
		}
	}
	return return_ip
}

func connectToServers(port string, n int) {
	ip_server := getAddress()
	fmt.Println("The Local IP is: " + ip_server)
	count := 0
	for {
		for index, ip_value := range Addrs {
			if Addrs_flag[index] == true {
				continue
			}
			if ip_value == ip_server {
				continue
			}
			dial_addr := ip_value + ":" + port
			dial_conn, err := net.Dial("tcp", dial_addr)
			if err == nil {
				Addrs_flag[index] = true
				count = count + 1
				dial_conns[dial_conn.RemoteAddr().String()] = dial_conn
				go Handler(dial_conn, &dial_conns, n)
				fmt.Println(ip_server + " connecting to IP address: " + dial_addr + "--successful")
			}
		}
		if count == n-1 {
			for i, val := range Addrs_flag {
				if Addrs[i] == ip_server {
					break
				}

				if val {
					fmt.Println(Addrs[i])
					pNumber = pNumber + 1
				}
			}
			fmt.Println("pNumber: ", pNumber)

			localStamp = InitTimestamp(n, pNumber)
			fmt.Println("# initialize local stamp: ", localStamp)
			go releaseHoldback(holdback, stayLong)
			break
		}
	}
	fmt.Println(ip_server + " connecting to ALL IP adresses" + "--successful")
	fmt.Println("======================")
	fmt.Println("The chat room is Live!")
	fmt.Println("======================")

}

func StartServer(port string, username string, n int) {
	go connectToServers(port, n)
	host := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("Address was not resolved")
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Listener was unable to listen on port " + port + ": try a different port or close the conflicting instance")
		return
	}
	conns := make(map[string]net.Conn)
	go broadcastUsername(&conns, username, n-1)
	go broadcastMessages(&conns, username)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept Failed")
			continue
		}
		conns[conn.RemoteAddr().String()] = conn
	}
}

func broadcastMessages(conns *map[string]net.Conn, username string) {

	for {
		fmt.Print("Message: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		if len(input) > 0 {
			realIput := "[" + username + "]: " + input
			msg := addTimestamp(localStamp, realIput)
			for key, conn := range *conns {
				_, err := conn.Write([]byte(msg))
				if err != nil {
					fmt.Println("broadcasting message to %s failed: %v\n", key, err)
					delete(*conns, key)
				}
			}
		}
	}
}

func Handler(conn net.Conn, conns *map[string]net.Conn, n int) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Reading Client Username Failed")
		delete(*conns, conn.RemoteAddr().String())
		conn.Close()
	}
	recvStr := string(buf[0:length])
	fmt.Println(recvStr)
	ip_a := strings.Split(conn.RemoteAddr().String(), ":")[0]
	IP2Username[ip_a] = recvStr
	for {
		length, err := conn.Read(buf)
		if err != nil {
			ip_left := strings.Split(conn.RemoteAddr().String(), ":")[0]
			username_left := IP2Username[ip_left]
			fmt.Println(username_left + " has left the chat")
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		recvStr := string(buf[0:length])
		handleMsg(recvStr, localStamp, n, holdback, stayLong)
	}
}

func releaseHoldback(holdback map[string]string, stayLong map[string]string) {
	for {
		time.Sleep(1 * time.Second)
		layout := "2000-01-01 20:00:00"
		currentT := time.Now()
		for key, value := range holdback {
			whenEntered, err := time.Parse(layout, stayLong[key])
			if err != nil {
				fmt.Println("cannot parse time")
			}
			diff := currentT.Sub(whenEntered)
			dura := int64(diff / time.Second)
			if dura >= 2 {
				fmt.Println("#stay too long")
				_, ok1 := stayLong[key]

				if ok1 {
					delete(stayLong, key)
				}

				_, ok2 := holdback[key]

				if ok2 {
					delete(holdback, key)
				}

				fmt.Println(value)
			}
		}
	}
}

func broadcastUsername(conns *map[string]net.Conn, username string, ip_num int) {
	for {
		msg := username
		if len(*conns) == ip_num {
			for key, conn := range *conns {
				_, err := conn.Write([]byte(msg))
				if err != nil {
					fmt.Println("broadcasting Username to %s failed: %v\n", key, err)
				}
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}
