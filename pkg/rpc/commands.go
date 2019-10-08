package rpc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dusk-network/dusk-blockchain/pkg/p2p/wire/protocol"
	"github.com/dusk-network/dusk-blockchain/pkg/p2p/wire/topics"
)

// handler defines a method bound to an RPC command.
type handler func(*Server, []string) (string, error)

var (

	// rpcCmd maps method names to their actual functions.
	rpcCmd = map[string]handler{
		"version": version,
		"ping":    pong,
		"uptime":  uptime,
		// Publish Topic (experimental). Injects an event directly into EventBus system.
		// Would be useful on E2E testing. Mind the supportedTopics list when sends it
		"publishTopic": publishTopic,
		"sendBidTx":    sendBidTx,
		"loadWallet":   loadWallet,
	}

	// rpcAdminCmd holds all admin methods.
	rpcAdminCmd = map[string]bool{}

	// supported topics for injection into EventBus
	supportedTopics = [2]topics.Topic{
		topics.Tx,
		topics.Block,
	}
)

// version will return the version of the client.
var version = func(s *Server, params []string) (string, error) {
	return protocol.NodeVer.String(), nil
}

// pong simply returns "pong" to let the caller know the server is up.
var pong = func(s *Server, params []string) (string, error) {
	return "pong", nil
}

// uptime returns the server uptime.
var uptime = func(s *Server, params []string) (string, error) {
	return strconv.FormatInt(time.Now().Unix()-s.startTime, 10), nil
}

var publishTopic = func(s *Server, params []string) (string, error) {

	if len(params) < 2 {
		return "", errors.New("expects always two input params - topics and payload bytes")
	}

	// Validate topic parameter.
	jsonrpcTopic := params[0]

	supported := false
	for _, topic := range supportedTopics {
		if topic.String() == jsonrpcTopic {
			supported = true
			break
		}
	}

	if !supported {
		return "", fmt.Errorf("%s is not supported by publishTopic API", jsonrpcTopic)
	}

	payload, _ := hex.DecodeString(params[1])
	rpcTopic := topics.StringToTopic(jsonrpcTopic)
	s.eventBus.Publish(rpcTopic, bytes.NewBuffer(payload))

	result :=
		`{ 
			"result": "published"
		}`
	return result, nil
}

var sendBidTx = func(s *Server, params []string) (string, error) {

	if len(params) < 2 {
		return "", fmt.Errorf("missing parameters: amount/locktime")
	}

	amount, err := strconv.Atoi(params[0])
	if err != nil {
		return "", fmt.Errorf("converting amount string to an integer: %v", err)
	}

	lockTime, err := strconv.Atoi(params[1])
	if err != nil {
		return "", fmt.Errorf("converting locktime string to an integer: %v", err)
	}

	txid, err := s.rpcBus.SendBidTx(uint64(amount), uint64(lockTime))
	result := fmt.Sprintf("{ \"txid\": \"%s\"}", hex.EncodeToString(txid))
	return result, err
}

var loadWallet = func(s *Server, params []string) (string, error) {

	if len(params) < 1 {
		return "", fmt.Errorf("missing parameter: password")
	}

	pubKey, err := s.rpcBus.LoadWallet(params[0])
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("{ \"pubkey\": \"%s\"}", pubKey)
	return result, err
}
