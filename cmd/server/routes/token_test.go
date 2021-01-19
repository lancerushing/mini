package routes

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestToken_Good(t *testing.T) {

	t.Parallel()
	check := is.New(t)

	testMsg := []byte("Test Message")
	token, err := tokenCreate(testMsg)
	check.NoErr(err)

	msg, err := tokenExtractMessage(token)
	check.NoErr(err)

	check.Equal(testMsg, msg)

}

func TestToken_Bad(t *testing.T) {

	t.Parallel()
	check := is.New(t)

	testMsg := []byte("Test Message")
	token, err := tokenCreate(testMsg)
	check.NoErr(err)

	token = append(token, byte(0))

	_, err = tokenExtractMessage(token)
	if err == nil {
		check.Fail()
	} else {
		check.Equal("token is invalid", err.Error())
	}

}

func TestToken_Expired(t *testing.T) {

	t.Parallel()
	check := is.New(t)

	message := []byte("Test Message")
	salt := make([]byte, saltBytes)
	_, _ = rand.Read(salt)
	expires, _ := time.Now().Add(-1 * time.Second).MarshalBinary()
	sum, err := computeSum(salt, expires, message)
	check.NoErr(err)

	result := bytes.NewBuffer(nil)
	result.Write(salt)
	result.Write(expires)
	result.Write(message)
	result.Write(sum)

	token := result.Bytes()

	_, err = tokenExtractMessage(token)
	if err == nil {
		check.Fail()
	} else {
		check.Equal("token is expired", err.Error())
	}

}
