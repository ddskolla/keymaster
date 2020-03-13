package main

import (
	"context"
	"github.com/bsycorp/keymaster/keymaster/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestHandler(t *testing.T) {
//	res, err :=
//}

func TestHandler(t *testing.T) {
	//var err error

	// Expect a JSON unmarshal error
	//res, err = Handler(context.Background(), []byte(""))
	//if assert.NoError(t, err) && assert.NotNil(t, res) {
	//	assert.False(t, res.Success)
	//	assert.Contains(t, res.Message, "bad request")
	//}


	// No operation
	//res, err = Handler(context.Background(), []byte("{}"))
	//if assert.NoError(t, err) && assert.NotNil(t, res) {
	//	assert.False(t, res.Success)
	//	assert.Contains(t, res.Message, "bad request: unknown operation")
	//}

	// Dumb operation
	//res, err = Handler(context.Background(), []byte("{ \"op\": \"foo\" }"))
	//assert.Nil(t, res)
	//assert.NotNil(t, err)
	//assert.Contains(t, err.Error(), "unknown operation")
}

func TestHandlePing(t *testing.T) {
	resp, err := Handler(context.Background(), []byte(`{ "op": "ping", "message": {} }`))
	if assert.NoError(t, err) {
		if pingResponse, ok := resp.(api.PingResponseMessage); assert.True(t, ok) {
			assert.Equal(t, pingResponse.Message, "OK")
		}
	}

	//assert.True(t, ok)
	//assert.True(t, pingResponse.Success)
	//assert.Equal(t, pingResponse.Message, "OK")
}

func TestHandleAuth(t *testing.T) {

}
