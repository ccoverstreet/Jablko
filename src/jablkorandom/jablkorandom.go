// jablkorandom.go: Jablko Random Package
// Cale Overstreet
// 2020/11/10
// Generate cryptographically secure randomness or 
// other useful random generation methods.

package jablkorandom

import (
	"crypto/rand"
)

func GenRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

const chars = "1234567890!@#$%^&/?.,:qwertyuiopasdfghjklmnzxcvb"

func GenRandomStr(n int) (string, error) {
	randB, err := GenRandomBytes(n)
	if err != nil {
		return "", err
	} 

	for i, b := range randB {
		randB[i] = chars[b % byte(len(chars))]
	}

	return string(randB), nil
}
