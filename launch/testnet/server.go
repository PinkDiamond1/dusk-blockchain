package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/chain"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/msg"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/mempool"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/peer"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/peer/dupemap"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/peer/processing/chainsync"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire/protocol"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/rpc"

	cfg "gitlab.dusk.network/dusk-core/dusk-go/pkg/config"
)

var timeOut = 3 * time.Second

type Server struct {
	eventBus *wire.EventBus
	rpcBus   *wire.RPCBus
	chain    *chain.Chain
	dupeMap  *dupemap.DupeMap
	counter  *chainsync.Counter
}

// Setup creates a new EventBus, generates the BLS and the ED25519 Keys, launches a new `CommitteeStore`, launches the Blockchain process and inits the Stake and Blind Bid channels
func Setup() *Server {
	// creating the eventbus
	eventBus := wire.NewEventBus()

	// creating the rpcbus
	rpcBus := wire.NewRPCBus()

	// generating the keys
	// TODO: this should probably lookup the keys on a local storage before recreating new ones
	// keys, _ := user.NewRandKeys()

	m := mempool.NewMempool(eventBus, nil)
	m.Run()

	// creating and firing up the chain process
	chain, err := chain.New(eventBus, rpcBus)
	if err != nil {
		panic(err)
	}
	go chain.Listen()

	// Setting up a dupemap
	dupeBlacklist := launchDupeMap(eventBus)

	if cfg.Get().RPC.Enabled {
		rpcServ, err := rpc.NewRPCServer(eventBus, rpcBus)
		if err != nil {
			log.Errorf("RPC server error: %s", err.Error())
		}

		if err := rpcServ.Start(); err != nil {
			log.Errorf("RPC server error: %s", err.Error())
		}
	}

	// creating the Server
	srv := &Server{
		eventBus: eventBus,
		rpcBus:   rpcBus,
		chain:    chain,
		dupeMap:  dupeBlacklist,
		counter:  chainsync.NewCounter(eventBus),
	}

	// Connecting to the general monitoring system
	// ConnectToMonitor(eventBus, d)

	return srv
}

func launchDupeMap(eventBus wire.EventBroker) *dupemap.DupeMap {
	roundChan := consensus.InitRoundUpdate(eventBus)
	dupeBlacklist := dupemap.NewDupeMap(1)
	go func() {
		for {
			round := <-roundChan
			// NOTE: do we need locking?
			dupeBlacklist.UpdateHeight(round)
		}
	}()
	return dupeBlacklist
}

func (s *Server) StartConsensus(round uint64) {
	roundBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(roundBytes[:8], round)
	s.eventBus.Publish(msg.InitializationTopic, bytes.NewBuffer(roundBytes))
}

func (s *Server) OnAccept(conn net.Conn) {
	responseChan := make(chan *bytes.Buffer, 100)
	peerReader, err := peer.NewReader(conn, protocol.TestNet, s.dupeMap, s.eventBus, s.rpcBus, s.counter, responseChan)
	if err != nil {
		panic(err)
	}

	if err := peerReader.Accept(); err != nil {
		log.WithFields(log.Fields{
			"process": "server",
			"error":   err,
		}).Warnln("problem performing handshake")
		return
	}
	log.WithFields(log.Fields{
		"process": "server",
		"address": peerReader.Addr(),
	}).Debugln("connection established")

	go peerReader.ReadLoop()
	peerWriter := peer.NewWriter(conn, protocol.TestNet, s.eventBus)
	peerWriter.Subscribe(s.eventBus)
	go peerWriter.WriteLoop(responseChan)
}

func (s *Server) OnConnection(conn net.Conn, addr string) {
	responseChan := make(chan *bytes.Buffer, 100)
	peerWriter := peer.NewWriter(conn, protocol.TestNet, s.eventBus)

	if err := peerWriter.Connect(s.eventBus); err != nil {
		log.WithFields(log.Fields{
			"process": "server",
			"error":   err,
		}).Warnln("problem performing handshake")
		return
	}
	log.WithFields(log.Fields{
		"process": "server",
		"address": peerWriter.Addr(),
	}).Debugln("connection established")

	peerReader, err := peer.NewReader(conn, protocol.TestNet, s.dupeMap, s.eventBus, s.rpcBus, s.counter, responseChan)
	if err != nil {
		panic(err)
	}

	go peerReader.ReadLoop()
	go peerWriter.WriteLoop(responseChan)
}

func (s *Server) Close() {
	s.chain.Close()
	s.rpcBus.Close()
}
