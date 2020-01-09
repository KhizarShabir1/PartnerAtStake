package main

import (
	a1 "assignment02IBC"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type validationStore struct {
	Names []string
	Coins []int
}
type addrs struct {
	Name string
	Addr string
}
type partnerStore struct { //
	Name         string
	WarningCount int
	Names        []string
	Addreses     []string
}

var partnerList []partnerStore
var addresses []addrs
var partnerLent int
var blacklist []string

func handleSatoshiConnection(c net.Conn, totalCons int, blockChain *a1.Block) {
	// (conn,totalCons,addresses, blockChain )
	log.Println("A client has connected with satoshi ", c.RemoteAddr())
	fmt.Printf(" client connected with satoshi-  ")
	buf := make([]byte, 4096)

	n, err := c.Read(buf) //get address of the peer
	if err != nil || n == 0 {
		c.Close()
		fmt.Println("Closing connection")

	}

	fmt.Println(string(buf[0:n]))
	s := strings.Split(string(buf[0:n]), ":")
	nd := addrs{Name: s[0], Addr: s[1]}
	addresses = append(addresses, nd)

}

func randArray(lent int, min int, max int, notIn int) []int {
	var arr []int

	for {

		if len(arr) == 0 {

			rand.Seed(time.Now().UnixNano())
			n := min + rand.Intn(max-min+1)
			if n != notIn {
				arr = append(arr, n)
				if len(arr) >= lent {
					break
				}
			}

		} else {

			check := true
			rand.Seed(time.Now().UnixNano())
			n := min + rand.Intn(max-min+1)
			for i := 0; i < len(arr); i++ {
				if n == arr[i] {
					check = false
				}
			}

			if check == true && n != notIn {
				arr = append(arr, n)
				if len(arr) >= lent {
					break
				}
			}
		}

	}
	return arr
}
func AddMoneyToValidStore(validStore validationStore, name string, coins int) {
	for i := 0; i < len(validStore.Names); i++ {
		if validStore.Names[i] == name {
			validStore.Coins[i] += coins
			break
		}
	}
}
func printValid(validStore validationStore) {
	for i := 0; i < len(validStore.Names); i++ {
		fmt.Print(validStore.Names[i])
		fmt.Print(" : ")
		fmt.Print(validStore.Coins[i])
		fmt.Print("\n")

	}
}
func handleMinerDistribution(transact a1.TransSend, index int) {
	transact.Miner = addresses[index].Name
	conn, err := net.Dial("tcp", "localhost:"+addresses[index].Addr)
	if err != nil {
		fmt.Print("error occured in mining connection ")
		//handle error
	}

	conn.Write([]byte("mine"))
	fmt.Print("sending mine message to -> " + addresses[index].Addr + "\n")
	time.Sleep(1 * time.Second)

	gobEncoder := gob.NewEncoder(conn)
	err = gobEncoder.Encode(transact)
	if err != nil {
		log.Println(err)
	}
}

func checkValidity(validStore validationStore, name string, coins int) bool {

	for i := 0; i < len(validStore.Names); i++ {
		if validStore.Names[i] == name {
			if validStore.Coins[i] >= coins {
				return true
			}

		}

	}
	return false
}

func deductStore(validStore validationStore, name string, coins int) {

	for i := 0; i < len(validStore.Names); i++ {
		if validStore.Names[i] == name {
			validStore.Coins[i] = validStore.Coins[i] - coins
		}

	}

}
func addToStore(validStore validationStore, name string, coins int) {

	for i := 0; i < len(validStore.Names); i++ {
		if validStore.Names[i] == name {
			validStore.Coins[i] = validStore.Coins[i] + coins
		}

	}

}
func sendValidationAnswer(c net.Conn, chk bool) {
	if chk == true {
		c.Write([]byte("yes"))
		time.Sleep(1 * time.Second)
	} else {
		c.Write([]byte("no"))
		time.Sleep(1 * time.Second)
	}

}

func checkIfAlreadyPresent(s string) bool {
	check := false
	for i := 0; i < partnerLent; i++ {
		if partnerList[i].Name == s {
			check = true
		}
	}
	return check
}

func isPartner(Name string, Sender string) bool {

	check := false
	for i := 0; i < partnerLent; i++ {
		if partnerList[i].Name == Name {
			for j := 0; j < len(partnerList[i].Names); j++ {
				if partnerList[i].Names[j] == Sender {
					check = true
				}
			}
		}
	}
	return check
}

func main() {

	partnerLent = 0

	var totalCons int
	totalCons = 3
	consNo := totalCons
	var blockChain *a1.Block
	blockChain = nil
	var validStore validationStore
	validStore.Names = append(validStore.Names, "satoshi")
	validStore.Coins = append(validStore.Coins, 0)

	var cons []net.Conn

	ln, err := net.Listen("tcp", ":7000")

	if err != nil {
		log.Fatal(err)
	}

	for totalCons > 0 {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		handleSatoshiConnection(conn, totalCons, blockChain)

		// Satoshi performs his own transaction with the addition of every new handleSatoshiConnection
		var transaction a1.Trans
		transaction.Transactions = append(transaction.Transactions, "satoshi")
		transaction.FreeCoin = append(transaction.FreeCoin, 100)
		transaction.NoOfTrans = 1
		AddMoneyToValidStore(validStore, "satoshi", 100)
		_, blockChain = a1.InsertBlock(transaction, blockChain)
		totalCons = totalCons - 1
		cons = append(cons, conn)
	}
	for i := 0; i < consNo; i++ {
		validStore.Names = append(validStore.Names, addresses[i].Name)
		validStore.Coins = append(validStore.Coins, 100) // Giving every new connection 100 FreeCoins as welcome gift
	}
	fmt.Println("FreeCoin Store updated with welcome transactions to all the connections ")
	printValid(validStore)

	for i := 0; i < consNo; i++ {

		gobEncoder := gob.NewEncoder(cons[i])
		err = gobEncoder.Encode(blockChain)
		if err != nil {
			log.Println(err)
		}

	}
	for i := 0; i < consNo; i++ {
		var addres []addrs
		var lengthOfRandomCons int
		if i == consNo-1 {
			lengthOfRandomCons = consNo - 1
		} else {
			lengthOfRandomCons = i + 1
		}
		arr := randArray(lengthOfRandomCons, 0, consNo-1, i) // return random values indexes of connections
		for j := 0; j < len(arr); j++ {

			fmt.Print(arr[j])
			addres = append(addres, addresses[arr[j]])
		}

		gobEncoder := gob.NewEncoder(cons[i])
		err = gobEncoder.Encode(addres)
		if err != nil {
			log.Println(err)
		}

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		buf := make([]byte, 4096)
		fmt.Println("New connection in satoshi transaction")

		// buffer :=bytes.NewBuffer([]byte 4096)

		n, err := conn.Read(buf) //get address of the peer
		if err != nil || n == 0 {
			conn.Close()
			fmt.Println("Closing connection")

		}

		choice := string(buf[0:n])

		fmt.Println(choice)
		if choice == "addMyPartner" {

			bufe := make([]byte, 4096)

			n, err := conn.Read(bufe) //get address of the peer
			if err != nil || n == 0 {
				conn.Close()
				fmt.Println("Closing connection")

			}

			fmt.Println(string(bufe[0:n]))
			s := strings.Split(string(bufe[0:n]), ":")
			nd := addrs{Name: s[0], Addr: s[1]}
			addresses = append(addresses, nd)

			// makin acount of partner
			validStore.Names = append(validStore.Names, s[0])
			validStore.Coins = append(validStore.Coins, 100)
			if !checkIfAlreadyPresent(s[2]) {

				var names []string
				var addresd []string
				addresd = append(addresd, s[1])
				names = append(names, s[0])
				nd5 := partnerStore{Name: s[2], WarningCount: 0, Names: names, Addreses: addresd}
				partnerList = append(partnerList, nd5)
				log.Println("partner Added "+s[0]+" : ", conn.RemoteAddr())
				partnerLent++

			} else {
				for i := 0; i < partnerLent; i++ {

					if partnerList[i].Name == s[2] {
						partnerList[i].Names = append(partnerList[i].Names, s[0])
						partnerList[i].Addreses = append(partnerList[i].Addreses, s[1])

					}
				} //for ends here

			} //else ends check

		} // addMyPartner ends here

		if choice == "chooseMiner" {

			var transact a1.TransSend
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&transact)
			fmt.Println("Mining transaction with following credentials ")
			fmt.Println("FreeCoins")
			fmt.Println(transact.FreeCoin)
			fmt.Println("Reciever")
			fmt.Println(transact.Transactions)
			fmt.Println("Sender")
			fmt.Println(transact.Sender)
			var index int
			var index2 int
			for k := 0; k < consNo; k++ {
				if addresses[k].Name == transact.Sender { //don't choose sender and reciever as miners
					index = k
				} //addresses[k].Name==transact.Transactions
				if addresses[k].Name == transact.Transactions { //don't choose sender and reciever as miners
					index2 = k
				} //addresses[k].Name==transact.Transactions

			}
			var n int
			for {
				rand.Seed(time.Now().UnixNano())
				n = 0 + rand.Intn((consNo-1)-0+1)
				if n != index && n != index2 {
					break
				}
			}

			// go handleMinerDistribution(transact a1.TransSend, index int)

			go handleMinerDistribution(transact, n)

		} //ChooseMiner ends here
		if choice == "validate" {

			var transact a1.TransSend
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&transact)
			fmt.Println("Validating transaction with following credentials ")
			fmt.Println("FreeCoins")
			fmt.Println(transact.FreeCoin)
			fmt.Println("Reciever")
			fmt.Println(transact.Transactions)
			fmt.Println("Sender")
			fmt.Println(transact.Sender)
			// func checkValidity(validStore validationStore, name string, coins int ) bool {
			// func deductStore(validStore validationStore, name string, coins int ){
			//   func addToStore(validStore validationStore, name string, coins int ){
			var chk bool
			chk = checkValidity(validStore, transact.Sender, transact.FreeCoin)
			chk2 := true
			// if transation is bein sent by/to blacklisted candidate block it
			for b := 0; b < len(blacklist); b++ {
				if blacklist[b] == transact.Sender || blacklist[b] == transact.Transactions {
					chk2 = false
				}
			}

			if chk == true && chk2 == true { //transaction is valid
				deductStore(validStore, transact.Sender, transact.FreeCoin)
				addToStore(validStore, transact.Transactions, transact.FreeCoin)
				addToStore(validStore, transact.Miner, 75) //coinbase

			} else {
				fmt.Println("length of partner ", partnerLent)

				for k := 0; k < partnerLent; k++ {
					fmt.Println("partner list name  " + partnerList[k].Name)
					//partnerList[k].Name == transact.Sender
					if isPartner(partnerList[k].Name, transact.Sender) {

						partnerList[k].WarningCount++

						//if ver exceeds warnin count  blacklist sender and all its partners

						if partnerList[k].WarningCount >= 2 {
							fmt.Println("// WARNING: exeeded from 1 of : " + transact.Sender)

							blacklist = append(blacklist, transact.Sender)
							for l := 0; l < len(partnerList[k].Names); l++ {
								blacklist = append(blacklist, partnerList[k].Names[l])

							}

						} else {

							fmt.Println("Warning being issued to the : ", transact.Sender)
						}

					}

				}
			}
			fmt.Println("FreeCoin Store ")
			printValid(validStore)
			if chk == true && chk2 == true {
				chk = true
			} else {
				chk = false
			}
			go sendValidationAnswer(conn, chk)

		} //validate ends here
		if choice == "broadcast" {
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&blockChain)

			if err != nil {
				//handle error
			}
			a1.ListBlocks(blockChain)
		}

	}

}
