package main

import (
	"Assignment/Server"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	fmt.Println("=================================================================")
	fmt.Println("======================DSP ASSIGNMENT=============================")
	fmt.Println("=================================================================")
	fmt.Println("==================Distributed Chat System========================")
	fmt.Println("=================================================================")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Welcome to the entry screen for this demo DCS")
	fmt.Println("")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter a username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Println(username)
	fmt.Print("Please enter the port of the DCS: ")
	scanner.Scan()
	port := scanner.Text()
	fmt.Println(port)
	fmt.Print("How many members will be in this chat: ")
	scanner.Scan()
	numberOfMember := scanner.Text()

	member, err := strconv.Atoi(numberOfMember)
	if err != nil {
		fmt.Println("please use Integer to indicate the number of members.")
		return
	}
	fmt.Println(member)
	Server.StartServer(port, username, member)
}
