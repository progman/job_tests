//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// 1.0.0
// Alexey Potehin <gnuplanet@gmail.com>, http://www.overtask.org/doc/cv
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
package main
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
import (
	"log"
	"fmt"
	"net"
	"os"
)
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
type GlobalType struct {
	ServerHost string
	ServerPort string
}
var Global GlobalType
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func PowClient() (err error) {
	var message string


// get env
	Global.ServerHost = os.Getenv("SERVER_HOST")
	Global.ServerPort = os.Getenv("SERVER_PORT")


// show env
	log.Printf("SERVER_HOST: \"%s\"\n", Global.ServerHost)
	log.Printf("SERVER_PORT: \"%s\"\n", Global.ServerPort)


// connect to server
	conn, err := net.Dial("tcp", Global.ServerHost + ":" + Global.ServerPort)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()


// get challenge
	log.Printf("read challenge...\n")
	message, err = ReadMessage(conn)
	if err != nil {
		return
	}
	log.Printf("done\n")
	challenge := message
	log.Printf("challenge: %s\n", challenge)


// get difficulty
	log.Printf("read difficulty...\n")
	message, err = ReadMessage(conn)
	if err != nil {
		return
	}
	log.Printf("done\n")
	difficulty := message
	log.Printf("difficulty: \"%s\"\n", difficulty)


// convert string difficulty to int
	var difficultyUint64 uint64
	difficultyUint64, err = StrToUint(difficulty)
	if err != nil {
		return
	}
	difficultyInt := int(difficultyUint64)


// solve challenge
	var nonce uint64
	nonce = ChallengeSolve(challenge, difficultyInt)
	log.Printf("nonce: %d\n", nonce)


// send nonce
	log.Printf("send nonce...\n")
	err = SendMessage(conn, UintToStr(nonce))
	if err != nil {
		return
	}
	log.Printf("done\n")


// get quote
	log.Printf("read quote...\n")
	message, err = ReadMessage(conn)
	if err != nil {
		return
	}
	log.Printf("done\n")
	log.Printf("quote: \"%s\"\n", message)


	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func main() {
	var err error


	err = PowClient()
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}


	os.Exit(0)
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
