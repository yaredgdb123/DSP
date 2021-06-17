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
	addressFlag = []bool{
		false,
	}

	adress = []string{
		"169.254.6.85",
	}

	dialConns   = make(map[string]net.Conn)
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

func contactServers(port string, n int) {
	ipServer := getAddress()
	fmt.Println("The Local IP is: " + ipServer)
	count := 0
	for {
		for index, ipValue := range adress {
			if addressFlag[index] == true {
				continue
			}

			//this doesn't always work cause getServerAdress() might not send the right IP
			if ipValue == ipServer {
				continue
			}
			dialAddress := ipValue + ":" + port
			dialConn, err := net.Dial("tcp", dialAddress)
			if err == nil {
				addressFlag[index] = true
				count = count + 1
				dialConns[dialConn.RemoteAddr().String()] = dialConn
				go Handler(dialConn, &dialConns, n)
				fmt.Println(ipServer + " connecting to IP address: " + dialAddress + "--successful")
			}
		}
		if count == n-1 {
			for i, val := range addressFlag {
				if adress[i] == ipServer {
					break
				}

				if val {
					fmt.Println(adress[i])
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
	fmt.Println(ipServer + " connecting to ALL IP adresses" + "--successful")
	fmt.Println("===================================")
	fmt.Println("======The chat room is Live!=======")
	fmt.Println("===================================")

}

func StartServer(port string, username string, n int) {
	go contactServers(port, n)
	host := ":" + port
	tcpAddress, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("adress was not resolved")
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		fmt.Println("Listener was unable to listen on port:" + port + " try a different port or close the conflicting instance")
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
	ipConn := strings.Split(conn.RemoteAddr().String(), ":")[0]
	IP2Username[ipConn] = recvStr
	for {
		length, err := conn.Read(buf)
		if err != nil {
			ipDisconnected := ipConn
			usernameDisconn := IP2Username[ipDisconnected]
			fmt.Println(usernameDisconn + " has left the chat")
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
		currentTime := time.Now()
		for key, value := range holdback {
			parsedTime, err := time.Parse(layout, stayLong[key])
			if err != nil {
				fmt.Println("cannot parse time")
			}
			timeDiff := currentTime.Sub(parsedTime)
			dura := int64(timeDiff / time.Second)
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

func broadcastUsername(conns *map[string]net.Conn, username string, ipNum int) {
	for {
		msg := username
		if len(*conns) == ipNum {
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
