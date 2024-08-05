//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// 1.0.0
// Alexey Potehin <gnuplanet@gmail.com>, http://www.overtask.org/doc/cv
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
package main
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
import (
	"log"
	"os"
	"fmt"
	"net"
	"math"
)
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
type GlobalType struct {
	ServerHost         string
	ServerPort         string
	Difficulty         string
	ChallengeLength    string

	DifficultyInt      int
	ChallengeLengthInt int
}
var Global GlobalType
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func handleConnectionInner(conn net.Conn, difficulty int, challengeLength int) (err error) {
	defer conn.Close()


// make challenge
	var challenge string
	challenge, err = GetChallenge(challengeLength)
	if err != nil {
		return
	}
	log.Printf("challenge: %s\n", challenge)
	log.Printf("challenge: %s\n", Bin2Hex(Hex2Bin(challenge)))


// send challenge
	log.Printf("send challenge...\n")
	err = SendMessage(conn, challenge)
	if err != nil {
		return
	}
	log.Printf("done\n")


// send difficulty
	log.Printf("send difficulty...\n")
	err = SendMessage(conn, fmt.Sprintf("%d", difficulty))
	if err != nil {
		return
	}
	log.Printf("done\n")


// read nonce
	log.Printf("read nonce...\n")
	var nonce string
	nonce, err = ReadMessage(conn)
	if err != nil {
		return
	}
	log.Printf("done\n")


// convert nonce to uint64
	var nonceUint64 uint64
	nonceUint64, err = StrToUint(nonce)
	if err != nil {
		return
	}


// verify challenge
	var flagOk bool
	flagOk = ChallengeVerify(challenge, difficulty, nonceUint64)
	if flagOk == false {
		err = fmt.Errorf("challenge failed")
		return
	}
	log.Printf("challenge passed\n")


// send quote
	log.Printf("send quote...\n")
	quote := GetQuote()
	err = SendMessage(conn, quote)
	if err != nil {
		return
	}
	log.Printf("done\n")


	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func handleConnection(conn net.Conn, difficulty int, challengeLength int) {
	defer conn.Close()


	var err error
	err = handleConnectionInner(conn, difficulty, challengeLength)
	if err != nil {
		log.Printf("ERROR %v\n", err)
		return
	}
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func PowServer() (err error) {
	var tmp uint64


// get env
	Global.ServerHost      = os.Getenv("SERVER_HOST")
	Global.ServerPort      = os.Getenv("SERVER_PORT")
	Global.Difficulty      = os.Getenv("DIFFICULTY")
	Global.ChallengeLength = os.Getenv("CHALLENGE_LENGTH")


// show env
	log.Printf("SERVER_HOST:      \"%s\"\n", Global.ServerHost)
	log.Printf("SERVER_PORT:      \"%s\"\n", Global.ServerPort)
	log.Printf("DIFFICULTY:       \"%s\"\n", Global.Difficulty)
	log.Printf("CHALLENGE_LENGTH: \"%s\"\n", Global.ChallengeLength)


// convert DIFFICULTY to int
	tmp, err = StrToUint(Global.Difficulty)
	if err != nil {
		return
	}
	if tmp > math.MaxInt {
		err = fmt.Errorf("DIFFICULTY is too big")
		return
	}
	Global.DifficultyInt = int(tmp)


// convert CHALLENGE_LENGTH to int
	tmp, err = StrToUint(Global.ChallengeLength)
	if err != nil {
		return
	}
	if tmp > math.MaxInt {
		err = fmt.Errorf("DIFFICULTY is too big")
		return
	}
	Global.ChallengeLengthInt = int(tmp)


// listen
	ln, err := net.Listen("tcp", Global.ServerHost + ":" + Global.ServerPort)
	if err != nil {
		return
	}
	defer ln.Close()
	log.Printf("Server started\n")


// accept new connects
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn, Global.DifficultyInt, Global.ChallengeLengthInt)
	}


	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
func main() {
	var err error


	err = PowServer()
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}


	os.Exit(0)
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
