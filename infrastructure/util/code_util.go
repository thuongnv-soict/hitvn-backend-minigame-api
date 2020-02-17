package util

import (
	"fmt"
	"github.com/speps/go-hashids"
	"strconv"
	"strings"
)

const (
	Salt 					= "hit.vn"
	InvitedCodeMinLength 	= 8
	MobileCartMinLength 	= 20
)


func checkPhoneNumberValid(phoneNumber string) bool {
	if len(phoneNumber) != 12 {
		return false
	}

	if !strings.HasPrefix(phoneNumber, "+"){
		return false
	}

	return true
}

func EncodeInvitedCode(phoneNumber string) (string, error){
	ok := checkPhoneNumberValid(phoneNumber)
	if ok == false {
		return "", nil
	}

	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = InvitedCodeMinLength

	h, err := hashids.NewWithData(hd)
	if err != nil{
		return "", err
	}

	number, err := strconv.Atoi(phoneNumber[1:])
	if err != nil{
		return "", err
	}
	e, _ := h.Encode([]int{number})

	return e, nil
}

func EncodeMobileCard(str string) (string, error){

	var code []int
	for _, letter := range str {
		number, err := strconv.Atoi(string(letter))
		if err != nil {
			return "", err
		}
		code = append(code, number)
	}

	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = MobileCartMinLength

	h, err := hashids.NewWithData(hd)
	if err != nil{
		return "", err
	}

	e, _ := h.Encode(code)

	return e, nil
}

func DecodeMobileCard(str string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = InvitedCodeMinLength

	h, err := hashids.NewWithData(hd)
	if err != nil{
		return "", err
	}

	decodeArray, err := h.DecodeWithError(str)
	if err != nil {
		fmt.Println(err)
	}
	//if len(d) != 1{
	//	fmt.Printf("Cannot decode: %s \n", code)
	//	return "", nil
	//}
	code := ""
	for _, number := range decodeArray {
		number := strconv.Itoa(number)
		code = code + number
	}

	//if len(phoneNumber) != 11{
	//	fmt.Printf("Decode wrong phone number: %s %s\n", code, phoneNumber)
	//	return "", nil
	//}

	return code, nil
}