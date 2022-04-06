// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this
// file, you can obtain one at https://opensource.org/licenses/MIT.
//
// Copyright (c) DUSK NETWORK. All rights reserved.

package reduction

import (
	"time"

	"github.com/dusk-network/dusk-blockchain/pkg/config"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus"
	"github.com/dusk-network/dusk-blockchain/pkg/core/consensus/header"
	"github.com/dusk-network/dusk-blockchain/pkg/core/data/block"
	"github.com/dusk-network/dusk-blockchain/pkg/p2p/wire/message"
	"github.com/dusk-network/dusk-blockchain/pkg/p2p/wire/topics"
	log "github.com/sirupsen/logrus"
)

// EmptyHash ...
var EmptyHash [32]byte

// EmptyStepVotes ...
var EmptyStepVotes = message.StepVotes{}

// EmptyResult ...
var EmptyResult *Result

func init() {
	EmptyResult = &Result{
		EmptyHash[:],
		EmptyStepVotes,
	}
}

// Result of the Reduction steps.
type Result struct {
	Hash []byte
	SV   message.StepVotes
}

// IsEmpty tests if the result of the aggregation is empty.
func (r *Result) IsEmpty() bool {
	return r == EmptyResult
}

// Reduction is a struct to be embedded in the reduction steps.
type Reduction struct {
	*consensus.Emitter
	TimeOut time.Duration

	// candidate is the block proposed by Block Generator
	candidate block.Block

	// VerifyFn verifies candidate block
	VerifyFn consensus.CandidateVerificationFunc
}

// IncreaseTimeout is used when reduction does not reach the quorum or
// converges over an empty block.
func (r *Reduction) IncreaseTimeout(round uint64) {
	// if we converged on an empty block hash, we increase the timeout
	r.TimeOut = r.TimeOut * 2
	if r.TimeOut > 60*time.Second {
		lg.
			WithField("timeout", r.TimeOut).
			WithField("round", round).
			Error("max_timeout_reached")

		r.TimeOut = 60 * time.Second
	}
}

// SendReduction propagates a signed vote for the candidate block, if block is
// fully valid.
func (r *Reduction) SendReduction(round uint64, step uint8, candidate *block.Block) {
	voteHash := EmptyHash[:]

	if candidate != nil {
		// TODO: EmptyBlock handled
		if err := r.VerifyFn(*candidate); err != nil {
			log.
				WithError(err).
				WithField("round", round).
				WithField("step", step).
				Error("verifyfn failed")
		} else {
			// Candidate block is fully valid.
			// Vote for it.
			voteHash, _ = candidate.CalculateHash()
		}
		// TODO: Throttle
	}

	// Generate Reduction message to propagate my vote.
	hdr := header.Header{
		Round:     round,
		Step:      step,
		BlockHash: voteHash,
		PubKeyBLS: r.Keys.BLSPubKey,
	}

	sig, err := r.Sign(hdr)
	if err != nil {
		panic(err)
	}

	red := message.NewReduction(hdr)
	red.SignedHash = sig

	m := message.NewWithHeader(topics.Reduction, *red, config.KadcastInitHeader)
	if err := r.Republish(m); err != nil {
		panic(err)
	}
}

// ShouldProcess checks whether a message is consistent with the current round
// and step. If it is not, it either discards it or stores it for later. The
// function potentially mutates the consensus.Queue.
func ShouldProcess(m message.Message, round uint64, step uint8, queue *consensus.Queue) bool {
	msg := m.Payload().(consensus.InternalPacket)
	hdr := msg.State()

	cmp := hdr.CompareRoundAndStep(round, step)
	if cmp == header.Before {
		lg.
			WithFields(log.Fields{
				"topic":          m.Category(),
				"round":          hdr.Round,
				"step":           hdr.Step,
				"expected round": round,
			}).
			Debugln("discarding obsolete event")
		return false
	}

	// Only store events up to 10 rounds into the future
	// XXX: According to protocol specs, we should abandon consensus
	// if we notice valid messages from far into the future.
	if cmp == header.After && hdr.Round-round < 10 {
		lg.
			WithFields(log.Fields{
				"topic":          m.Category(),
				"round":          hdr.Round,
				"step":           hdr.Step,
				"expected round": round,
			}).
			Debugln("storing future event for later")
		queue.PutEvent(hdr.Round, hdr.Step, m)
		return false
	}

	if m.Category() != topics.Reduction {
		lg.
			WithFields(log.Fields{
				"topic": m.Category(),
				"round": hdr.Round,
				"step":  hdr.Step,
			}).
			Warnln("unexpected topic for this step")
		return false
	}

	return true
}
