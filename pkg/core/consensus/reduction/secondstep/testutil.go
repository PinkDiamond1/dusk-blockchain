package secondstep

import (
	"bytes"
	"sync"
	"time"

	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus/agreement"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus/header"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus/reduction"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus/user"
	"github.com/dusk-network/dusk-blockchain/pkg/p2p/wire/topics"
	"github.com/dusk-network/dusk-blockchain/pkg/util/nativeutils/eventbus"
	"github.com/dusk-network/dusk-blockchain/pkg/util/nativeutils/rpcbus"
	crypto "github.com/dusk-network/dusk-crypto/hash"
	"github.com/dusk-network/dusk-wallet/key"
)

// Helper for reducing test boilerplate
type Helper struct {
	*Factory
	Keys        []key.ConsensusKeys
	P           *user.Provisioners
	Reducer     *Reducer
	eventPlayer consensus.EventPlayer
	signer      consensus.Signer

	AgreementChan chan bytes.Buffer
	RegenChan     chan bytes.Buffer
	nr            int

	lock               sync.RWMutex
	failOnVerification bool
	Handler            *reduction.Handler
}

// NewHelper creates a Helper
func NewHelper(eb *eventbus.EventBus, rpcbus *rpcbus.RPCBus, eventPlayer consensus.EventPlayer, signer consensus.Signer, provisioners int) *Helper {
	p, keys := consensus.MockProvisioners(provisioners)
	factory := NewFactory(eb, rpcbus, keys[0], 1000*time.Millisecond)
	a := factory.Instantiate()
	red := a.(*Reducer)
	hlp := &Helper{
		Factory:            factory,
		Keys:               keys,
		P:                  p,
		Reducer:            red,
		eventPlayer:        eventPlayer,
		signer:             signer,
		AgreementChan:      make(chan bytes.Buffer, 1),
		RegenChan:          make(chan bytes.Buffer, 1),
		nr:                 provisioners,
		failOnVerification: false,
		Handler:            reduction.NewHandler(keys[0], *p),
	}
	hlp.createResultChan()
	return hlp
}

func (hlp *Helper) Verify(hash []byte, sv *agreement.StepVotes, round uint64, step uint8) error {
	vc := hlp.P.CreateVotingCommittee(round, step, hlp.nr)
	sub := vc.Intersect(sv.BitSet)
	apk, err := agreement.ReconstructApk(sub)
	if err != nil {
		return err
	}

	return header.VerifySignatures(round, step, hash, apk, sv.Signature)
}

// CreateResultChan is used by tests (internal and external) to quickly wire the StepVotes resulting from the firststep reduction to a channel to listen to
func (hlp *Helper) createResultChan() {
	agListener := eventbus.NewChanListener(hlp.AgreementChan)
	hlp.Bus.Subscribe(topics.Agreement, agListener)
	regenListener := eventbus.NewChanListener(hlp.RegenChan)
	hlp.Bus.Subscribe(topics.Regeneration, regenListener)
}

// SendBatch of consensus events to the reducer callback CollectReductionEvent
func (hlp *Helper) SendBatch(hash []byte, round uint64, step uint8) {
	batch := hlp.Spawn(hash, round, step)
	for _, ev := range batch {
		go hlp.Reducer.CollectReductionEvent(ev)
	}
}

// Spawn a number of different valid events to the Agreement component bypassing the EventBus
func (hlp *Helper) Spawn(hash []byte, round uint64, step uint8) []consensus.Event {
	evs := make([]consensus.Event, hlp.nr)
	vc := hlp.P.CreateVotingCommittee(round, step, hlp.nr)
	for i := 0; i < hlp.nr; i++ {
		ev := reduction.MockConsensusEvent(hash, round, step, hlp.Keys, vc, i)
		evs[i] = ev

	}
	return evs
}

// Initialize the reducer with a Round update
func (hlp *Helper) Initialize(ru consensus.RoundUpdate) {
	hlp.Reducer.Initialize(hlp.eventPlayer, hlp.signer, ru)
}

func (hlp *Helper) StartReduction(sv *agreement.StepVotes) error {
	buf := new(bytes.Buffer)
	if sv != nil {
		if err := agreement.MarshalStepVotes(buf, sv); err != nil {
			return err
		}
	}

	hlp.Reducer.CollectStepVotes(consensus.Event{header.Header{}, *buf})
	return nil
}

// ProduceFirstStepVotes encapsulates the process of creating and forwarding Reduction events
func ProduceFirstStepVotes(eb *eventbus.EventBus, rpcbus *rpcbus.RPCBus, eventPlayer consensus.EventPlayer, signer consensus.Signer, nr int, withTimeout bool, round uint64, step uint8) (*Helper, []byte) {
	h := NewHelper(eb, rpcbus, eventPlayer, signer, nr)
	roundUpdate := consensus.MockRoundUpdate(1, h.P, nil)
	h.Initialize(roundUpdate)
	hash, _ := crypto.RandEntropy(32)
	h.Spawn(hash, round, step)
	return h, hash
}
