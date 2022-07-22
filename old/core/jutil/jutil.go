package jutil

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
)

// Code for generating random string used as cookie value
const cookieChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Function from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func RandomString(n int) (string, error) {
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(cookieChars))))
		if err != nil {
			return "", err
		}
		res[i] = cookieChars[num.Int64()]
	}

	return string(res), nil
}

func ParseJSONBody(body io.ReadCloser, dataStruct interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, dataStruct)
	return err
}
