package main

import (
	a2 "assignment02IBC"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	// "bufio"
)

type addrs struct {
	Name string
	Addr string
}
type partner struct {
	Name string
	Addr string
}

var blockChain *a2.Block

// var addresses []addrs
var addres []addrs
var partners []partner
var noOfPartners int

func addMyPartnerInBlockchain(myName string, newPartner partner) {
	conn, err := net.Dial("tcp", "localhost:7000")
	if err != nil {
		//handle error
	}
	// log.Println("A client has connected with satoshi", conn.RemoteAddr())
	// fmt.Printf(" client connected with satoshi ")

	conn.Write([]byte("addMyPartner"))
	time.Sleep(1 * time.Second)
	conn.Write([]byte(newPartner.Name + ":" + newPartner.Addr + ":" + myName))

}

func handleValidationAndDistribution(transact a2.TransSend, addres []addrs) { // myName string,
	conn, err := net.Dial("tcp", "localhost:7000")
	if err != nil {
		//handle error
	}
	conn.Write([]byte("validate"))
	time.Sleep(1 * time.Second)
	gobEncoder := gob.NewEncoder(conn)
	err = gobEncoder.Encode(transact)

	buf := make([]byte, 4096)
	n, err := conn.Read(buf) //get validation
	if err != nil || n == 0 {
		conn.Close()
		fmt.Println("Closing connection")

	}
	// "yes"
	// "no"
	recievedValidation := string(buf[0:n])
	if recievedValidation == "yes" {
		// var transact a2.TransSend
		// fmt.Println(transact.FreeCoin)
		// fmt.Println(transact.Transactions)
		// fmt.Println(transact.Sender)
		var transaction a2.Trans
		transaction.Transactions = append(transaction.Transactions, transact.Transactions)
		transaction.FreeCoin = append(transaction.FreeCoin, transact.FreeCoin)
		transaction.Transactions = append(transaction.Transactions, transact.Sender)
		transaction.FreeCoin = append(transaction.FreeCoin, 75) //CoinBase transaction

		transaction.NoOfTrans = 2
		// AddMoneyToValidStore(validStore ,"satoshi" , 100 )
		var hashVal string

		hashVal, blockChain = a2.InsertBlock(transaction, blockChain)

		//   fmt.Print(addres[i].Addr)
		//   fmt.Print(addres[i].Name)
		// addres [] addrs
		conn, err := net.Dial("tcp", "localhost:7000") //sendin blockChain for satoshi
		if err != nil {
			//handle error
		}

		conn.Write([]byte("broadcast"))
		time.Sleep(1 * time.Second)
		gobEncoder := gob.NewEncoder(conn)
		err = gobEncoder.Encode(blockChain)
		if err != nil {
			log.Println(err)
		}

		for o := 0; o < len(addres); o++ { //lets broadcast blockchain

			conn, err := net.Dial("tcp", "localhost:"+addres[o].Addr)
			if err != nil {
				//handle error
			}

			conn.Write([]byte("broadcast"))
			time.Sleep(1 * time.Second)
			conn.Write([]byte(hashVal))
			time.Sleep(1 * time.Second)

			gobEncoder := gob.NewEncoder(conn)
			err = gobEncoder.Encode(blockChain)
			if err != nil {
				log.Println(err)
			}

		}
		fmt.Println("Block Broadcated from -> " + transact.Sender)

	} else if recievedValidation == "no" {
		fmt.Println("Invalid Transation recieved from -> " + transact.Sender)

	}

} // function ens here
func ListenForTransactions(addres []addrs, name string, Address string) {
	ln, err := net.Listen("tcp", ":"+Address)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("A client has connected with "+name+" : ", conn.RemoteAddr())

		buf := make([]byte, 4096)

		n, err := conn.Read(buf) //get address of the peer
		if err != nil || n == 0 {
			conn.Close()
			fmt.Println("Closing connection")

		}
		fmt.Println(string(buf[0:n]))
		choice := string(buf[0:n])
		if choice == "makePartner" {

			bufe := make([]byte, 4096)

			n, err := conn.Read(bufe) //get address of the peer
			if err != nil || n == 0 {
				conn.Close()
				fmt.Println("Closing connection")

			}

			fmt.Println(string(bufe[0:n]))
			s := strings.Split(string(bufe[0:n]), ":")
			nd := addrs{Name: s[0], Addr: s[1]}
			nd2 := partner{Name: s[0], Addr: s[1]}

			// Add partner code goes here

			gobEncoder := gob.NewEncoder(conn)
			err = gobEncoder.Encode(blockChain)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
			err = gobEncoder.Encode(addres)
			if err != nil {
				log.Println(err)
			}
			// adding address of the partner
			addres = append(addres, nd)
			partners = append(partners, nd2)
			noOfPartners += 1
			go addMyPartnerInBlockchain(name, nd2)
			// go addMyPartnerInBlockchain()

		} else if choice == "mine" {

			var transact a2.TransSend
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&transact)
			if err != nil {
				//handle error
			}

			// "validate"
			// var transact a2.TransSend
			// fmt.Println(transact.FreeCoin)
			// fmt.Println(transact.Transactions)
			// fmt.Println(transact.Sender)
			go handleValidationAndDistribution(transact, addres) // myName,

		} else if choice == "broadcast" {
			//COde of stopin validation

			bufe := make([]byte, 40960)

			n2, err := conn.Read(bufe) //get address of the peer
			if err != nil || n == 0 {
				conn.Close()
				fmt.Println("Closing connection")

			}
			hashVal := string(bufe[0:n2])

			existHash := a2.CheckHashExists(blockChain, hashVal)
			var tmpBlock *a2.Block
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&tmpBlock)
			a2.ListBlocks(tmpBlock)

			if err != nil {
				//handle error
			}

			if existHash == false { //check that block does not exist beffore this
				fmt.Println("Previous blockchain")
				a2.ListBlocks(blockChain)

				blockChain = tmpBlock
				fmt.Println("new block chain is ")
				a2.ListBlocks(blockChain)
				fmt.Println("new block chain is ")
				//Now lets broadCast to the connections

				for o := 0; o < len(addres); o++ { //lets broadcast blockchain

					conn, err := net.Dial("tcp", "localhost:"+addres[o].Addr)
					if err != nil {
						//handle error
					}

					conn.Write([]byte("broadcast"))
					time.Sleep(1 * time.Second)

					conn.Write([]byte(hashVal))
					time.Sleep(1 * time.Second)

					gobEncoder := gob.NewEncoder(conn)
					err = gobEncoder.Encode(blockChain)
					if err != nil {
						log.Println(err)
					}

				}

			} else {
				//Just recieve and dont broadCast

				fmt.Println("Block Received but not broadCasted")
				// a2.ListBlocks(blockChain)

			}

		} // Broad ends here

		// s := strings.Split(string(buf[0:n]), ":")
		// nd := addrs{Name: s[0], Addr:s[1]}
		// addresses = append(addresses,nd)

	}

}

func main() {

	noOfPartners = 0
	name := os.Args[1]
	Address := os.Args[2]
	if len(os.Args) > 3 {
		// new partner
		partnerName := os.Args[3]
		partnerAddress := os.Args[4]

		conn, err := net.Dial("tcp", "localhost"+":"+partnerAddress)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(" A new user wants to be the partner of : ", partnerAddress)

		conn.Write([]byte("makePartner"))
		time.Sleep(1 * time.Second)

		conn.Write([]byte(name + ":" + Address))
		time.Sleep(1 * time.Second)

		dec := gob.NewDecoder(conn)
		err = dec.Decode(&blockChain)
		if err != nil {
			//handle error
		}

		// dec := gob.NewDecoder(conn)
		// gob.Register(net.TCPConn)
		err = dec.Decode(&addres)
		if err != nil {
			//handle error
		}
		nd3 := partner{Name: partnerName, Addr: partnerAddress}
		partners = append(partners, nd3)
		noOfPartners += 1

		a2.ListBlocks(blockChain)
		fmt.Println("My connections are following Broadcated from -> ")
		for i := 0; i < len(addres); i++ {

			fmt.Print(addres[i].Addr)
			fmt.Print(" : ")
			fmt.Print(addres[i].Name)
			fmt.Println("\n")

		}
		go ListenForTransactions(addres, name, Address)
		for {

			var option int
			option = -1
			fmt.Println("Enter 1 for transaction : ")
			_, err := fmt.Scan(&option)
			if err != nil {
				//handle error
			}
			if option == 1 {
				// reader := bufio.NewReader(os.Stdin)

				var coins int
				fmt.Println("Please enter no of FreeCoins do you want to send : ")
				_, err = fmt.Scan(&coins)

				for i := 0; i < len(addres); i++ {
					fmt.Println("Enter " + strconv.Itoa(i) + " to send to " + addres[i].Name)
				}
				var tranC int

				_, err = fmt.Scan(&tranC)
				// conn, err := net.Dial("tcp", "localhost:7000")
				// if err != nil {
				//   //handle error
				// }
				// conn.Write([]byte("chooseMiner"))
				// time.Sleep(1 * time.Second)
				//Encodin transaction
				var transact a2.TransSend
				transact.Transactions = addres[tranC].Name
				transact.FreeCoin = coins
				transact.Sender = name
				fmt.Println("sender name is ", name)
				fmt.Println("reciever Name is  ", addres[tranC].Name)

				// gobEncoder := gob.NewEncoder(conn)
				// err = gobEncoder.Encode(transact)
				// if err != nil {
				//   log.Println(err)
				// }
				var index int
				var index2 int
				index = -1
				index2 = -1

				fmt.Println("no off partners", noOfPartners)

				for k := 0; k < noOfPartners; k++ {

					fmt.Print("partner  : ", partners[k].Name)
					if partners[k].Name == transact.Sender { //don't choose sender and reciever as miners
						index = k
					} //addresses[k].Name==transact.Transactions
					if partners[k].Name == transact.Transactions { //don't choose sender and reciever as miners

						index2 = k
					} //addresses[k].Name==transact.Transactions

				}

				fmt.Println("index", index)

				fmt.Println("index2", index2)
				// Randoml choosin one partner from mining
				var n int
				for {
					rand.Seed(time.Now().UnixNano())
					n = 0 + rand.Intn((noOfPartners-1)-0+1)

					if n != index && n != index2 {
						fmt.Println("breaking now")
						break
					}
				}
				transact.Miner = partners[n].Name
				fmt.Println("miner name ", transact.Miner)
				fmt.Println("miner address ", partners[n].Addr)
				conn, err := net.Dial("tcp", "localhost:"+partners[n].Addr)
				if err != nil {
					fmt.Println("error occured in mining connection ")
					//handle error
				}

				conn.Write([]byte("mine"))
				fmt.Println("sending mine message to -> " + partners[n].Addr + "\n")
				time.Sleep(1 * time.Second)

				gobEncoder := gob.NewEncoder(conn)
				err = gobEncoder.Encode(transact)
				if err != nil {
					log.Println(err)
				}

			} //partner based transaction

		}

	} else { // if of new partner ends here
		conn, err := net.Dial("tcp", "localhost:7000")
		if err != nil {
			//handle error
		}
		log.Println("A client has connected with satoshi", conn.RemoteAddr())
		fmt.Printf(" client connected with satoshi ")

		conn.Write([]byte(name + ":" + Address))
		// var recvdBlock * a2.Block
		dec := gob.NewDecoder(conn)
		err = dec.Decode(&blockChain)

		if err != nil {
			//handle error
		}
		var addres []addrs
		// dec := gob.NewDecoder(conn)
		// gob.Register(net.TCPConn)
		err = dec.Decode(&addres)
		if err != nil {
			//handle error
		}

		a2.ListBlocks(blockChain)
		fmt.Println("My connections are following Broadcated from -> ")
		for i := 0; i < len(addres); i++ {

			fmt.Print(addres[i].Addr)
			fmt.Print(" : ")
			fmt.Print(addres[i].Name)
			fmt.Println("\n")

		}

		go ListenForTransactions(addres, name, Address)
		for {

			var option int
			option = -1
			fmt.Println("Enter 1 for transaction : ")
			_, err = fmt.Scan(&option)
			if option == 1 {
				// reader := bufio.NewReader(os.Stdin)

				var coins int
				fmt.Println("Please enter no of FreeCoins do you want to send : ")
				_, err = fmt.Scan(&coins)
				for i := 0; i < len(addres); i++ {
					fmt.Println("Enter " + strconv.Itoa(i) + " to send to " + addres[i].Name)
				}
				var tranC int

				_, err = fmt.Scan(&tranC)
				conn, err := net.Dial("tcp", "localhost:7000")
				if err != nil {
					//handle error
				}
				conn.Write([]byte("chooseMiner"))
				time.Sleep(1 * time.Second)
				//Encodin transaction
				var transact a2.TransSend
				transact.Transactions = addres[tranC].Name
				transact.FreeCoin = coins
				transact.Sender = name

				gobEncoder := gob.NewEncoder(conn)
				err = gobEncoder.Encode(transact)
				if err != nil {
					log.Println(err)
				}
				// conn.Close()

			}

		}

	} //else of peer ends here

}
