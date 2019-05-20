package reduction

import (
	"bytes"

	"github.com/stretchr/testify/mock"
	"gitlab.dusk.network/dusk-core/dusk-go/mocks"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/committee"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/header"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/selection"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/user"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/crypto"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/crypto/bls"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/util/nativeutils/sortedset"
	"golang.org/x/crypto/ed25519"
)

func MockReduction(keys *user.Keys, hash []byte, round uint64, step uint8) (*user.Keys, *Reduction) {
	if keys == nil {
		keys, _ = user.NewRandKeys()
	}

	reduction := MockOutgoingReduction(hash, round, step)
	reduction.PubKeyBLS = keys.BLSPubKeyBytes

	r := new(bytes.Buffer)
	_ = header.MarshalSignableVote(r, reduction.Header)
	sigma, _ := bls.Sign(keys.BLSSecretKey, keys.BLSPubKey, r.Bytes())
	reduction.SignedHash = sigma.Compress()
	return keys, reduction
}

func MockOutgoingReduction(hash []byte, round uint64, step uint8) *Reduction {
	reduction := New()
	reduction.Round = round
	reduction.Step = step
	reduction.BlockHash = hash
	return reduction
}

func MockOutgoingReductionBuf(hash []byte, round uint64, step uint8) *bytes.Buffer {
	reduction := MockOutgoingReduction(hash, round, step)
	b := new(bytes.Buffer)
	_ = header.MarshalSignableVote(b, reduction.Header)
	return b
}

func MockReductionBuffer(keys *user.Keys, hash []byte, round uint64, step uint8) (*user.Keys, *bytes.Buffer) {
	k, ev := MockReduction(keys, hash, round, step)
	marshaller := NewUnMarshaller()
	buf := new(bytes.Buffer)
	_ = marshaller.Marshal(buf, ev)
	return k, buf
}

func mockSelectionEventBuffer(hash []byte) *bytes.Buffer {
	// 32 bytes
	score, _ := crypto.RandEntropy(32)
	// Var Bytes
	proof, _ := crypto.RandEntropy(1477)
	// 32 bytes
	z, _ := crypto.RandEntropy(32)
	// Var Bytes
	bidListSubset, _ := crypto.RandEntropy(32)
	// BLS is 33 bytes
	seed, _ := crypto.RandEntropy(33)
	se := &selection.ScoreEvent{
		Round:         uint64(23),
		Score:         score,
		Proof:         proof,
		Z:             z,
		Seed:          seed,
		BidListSubset: bidListSubset,
		VoteHash:      hash,
	}

	b := make([]byte, 0)
	r := bytes.NewBuffer(b)
	_ = selection.MarshalScoreEvent(r, se)
	return r
}

func mockBlockEventBuffer(round uint64, step uint8, hash []byte) *bytes.Buffer {
	keys, _ := user.NewRandKeys()
	signedHash, _ := bls.Sign(keys.BLSSecretKey, keys.BLSPubKey, hash)
	marshaller := NewUnMarshaller()

	bev := &Reduction{
		Header: &header.Header{
			PubKeyBLS: keys.BLSPubKeyBytes,
			Round:     round,
			Step:      step,
			BlockHash: hash,
		},
		SignedHash: signedHash.Compress(),
	}

	buf := new(bytes.Buffer)
	_ = marshaller.Marshal(buf, bev)
	edSig := ed25519.Sign(*keys.EdSecretKey, buf.Bytes())
	completeBuf := bytes.NewBuffer(edSig)
	completeBuf.Write(keys.EdPubKeyBytes)
	completeBuf.Write(buf.Bytes())
	return completeBuf
}

func mockCommittee(quorum int, isMember bool, amMember bool) committee.Committee {
	committeeMock := &mocks.Committee{}
	committeeMock.On("Quorum").Return(quorum)
	committeeMock.On("ReportAbsentees", mock.Anything,
		mock.Anything, mock.Anything).Return(nil)
	committeeMock.On("IsMember",
		mock.AnythingOfType("[]uint8"),
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType("uint8")).Return(isMember)
	committeeMock.On("AmMember",
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType("uint8")).Return(amMember)
	committeeMock.On("Pack",
		mock.AnythingOfType("sortedset.Set"),
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType("uint8")).Return(sortedset.All)
	return committeeMock
}
