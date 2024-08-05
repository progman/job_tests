//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// 1.0.0
// Alexey Potehin <gnuplanet@gmail.com>, http://www.overtask.org/doc/cv
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
package main
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
import (
	"crypto/sha256"
	"crypto/rand"
r	"math/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"errors"
)
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// convert uint64 to []
// 0 -> [ 0 ], 255 -> [ 255 ], 256 -> [ 0, 1 ]
func Uint64ToByteSlice(src uint64) (dst []byte) {
	if src == 0 {
		dst = append(dst, 0)
		return
	}

	var tmp []byte
	for i := 0; i < 8; i++ {
		tmp = append([]byte{byte(src & 255)}, tmp...)
		src >>= 8
	}

	var flagSkip bool = true
	for i := 0; i < 8; i++ {
		if flagSkip == true {
			if tmp[i] != 0 {
				flagSkip = false
			}
		}
		if flagSkip == false {
			dst = append([]byte{byte(tmp[i])}, dst...)
		}
	}

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// convert bin to hex
// byte 0x00 -> "00"
func Bin2Hex(src []byte) (dst string) {
	var tmp []byte = make([]byte, len(src) * 2)

	hex.Encode(tmp, src)

	dst = string(tmp)

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// convert hex to bin
// "00" -> byte 0x00
func Hex2Bin(src string) (dst []byte) {
	dst = make([]byte, len(src) / 2)

	hex.Decode(dst, []byte(src))

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// check string, maybe it uint
func IsUint(str string) (flagOk bool) {
	str_size := len(str)

	if str_size == 0 {
		return
	}

	for i := 0; i < str_size; i++ {
		if (str[i] < byte('0')) || (str[i] > byte('9')) {
			return
		}
	}
	flagOk = true
	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// convert string to uint64
func StrToUint(str string) (value uint64, err error) {

	if IsUint(str) == false {
		err = errors.New("this is not uint")
		return
	}

	value, err = strconv.ParseUint(str, 10, 64)
	if err != nil {
		return
	}

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
/// convert uint64 to string
func UintToStr(value uint64) (str string) {

	str = fmt.Sprintf("%d", value)

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// make challange
func GetChallenge(size int) (challenge string, err error) {
	src := make([]byte, size)
	_, err = rand.Read(src)
	if err != nil {
		return
	}

	challenge = Bin2Hex(src)

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// check solve
func CheckSolve(solve []byte, difficulty int) (flagOk bool) {
	tmp := solve[:difficulty]

	flagOk = true
	for i := 0; i < len(tmp); i++ {
		if tmp[i] != 0 {
			flagOk = false
			break
		}
	}

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// solve challenge
func ChallengeSolve(challenge string, difficulty int) (nonce uint64) {
	for {
		solve := sha256.Sum256(append(Hex2Bin(challenge), Uint64ToByteSlice(nonce)...))
//		fmt.Printf("%x\n", solve)

		flagOk := CheckSolve([]byte(solve[:]), difficulty)
		if flagOk == true {
			break
		}
		nonce++
	}


	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// verify challenge
func ChallengeVerify(challenge string, difficulty int, nonce uint64) (flagOk bool) {

	solve := sha256.Sum256(append(Hex2Bin(challenge), Uint64ToByteSlice(nonce)...))
//	fmt.Printf("%x\n", solve)

	flagOk = CheckSolve([]byte(solve[:]), difficulty)

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// send line (message with '\n')
func SendMessage(conn net.Conn, message string) (err error) {
	_, err = fmt.Fprintf(conn, message + "\n")
	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// read line (message with '\n')
func ReadMessage(conn net.Conn) (message string, err error) {
	buf := make([]byte, 1)
	var n int

	for {
		n, err = conn.Read(buf)
		if err != nil {
			return
		}


		if n == 1 {
			if buf[0] == '\n' {
				break
			}
			message += string(buf)
		}
	}

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
// get quote
func GetQuote() (quote string) {
	quoteList := []string{
		"Imagination is more important than knowledge. For knowledge is limited, whereas imagination embraces the entire world. - Albert Einstein",
		"Be the change that you wish to see in the world. - Mahatma Gandhi",
		"Success is not final, failure is not fatal: It is the courage to continue that counts. - Winston Churchill",
		"Your work is going to fill a large part of your life, and the only way to be truly satisfied is to do what you believe is great work. - Steve Jobs",
		"I cannot live without books. - Thomas Jefferson",
	}

	quote = quoteList[r.Intn(len(quoteList))]

	return
}
//-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------//
