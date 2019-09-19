package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

var saltBytes = 16
var timeBytes = 15
var sumBytes = sha512.Size

// @todo how to configure? ENV? "Config" struct?
var shaSecret = []byte("need-to-configure-secret")

func computeSum(salt []byte, expires []byte, message []byte) ([]byte, error) {
	// verify input
	if len(salt) != saltBytes {
		return nil, errors.New("salt is unexpected length")
	}

	// verify input
	if len(expires) != timeBytes {
		return nil, errors.New("expires is unexpected length")
	}

	h := hmac.New(sha512.New, shaSecret)
	h.Write(salt)
	h.Write(expires)
	h.Write(message)

	return h.Sum(nil), nil
}

func tokenExtractMessage(token []byte) ([]byte, error) {

	sumStart := len(token) - sumBytes

	extractedSalt := token[0:saltBytes]
	extractedExpires := token[saltBytes : saltBytes+timeBytes]
	extractedMessage := token[saltBytes+timeBytes : sumStart]
	extractedSum := token[sumStart:]

	// verify sum
	computedSum, err := computeSum(extractedSalt, extractedExpires, extractedMessage)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(extractedSum, computedSum) {
		return nil, errors.New("token is invalid")
	}

	// verify  expiration
	var expires time.Time
	err = expires.UnmarshalBinary(extractedExpires)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expires) {
		return nil, errors.New("token is expired")
	}

	return extractedMessage, nil
}

func tokenCreate(message []byte) ([]byte, error) {

	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	expires, err := time.Now().Add(2 * time.Hour).MarshalBinary()
	if err != nil {
		return nil, err
	}

	sum, err := computeSum(salt, expires, message)
	if err != nil {
		return nil, err
	}

	result := bytes.NewBuffer(nil)
	result.Write(salt)
	result.Write(expires)
	result.Write(message)
	result.Write(sum)

	return result.Bytes(), nil

}
