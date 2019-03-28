package notary

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.dusk.network/dusk-core/dusk-go/mocks"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/committee"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/msg"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/crypto"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
)

func TestSimpleBlockCollection(t *testing.T) {
	committeeMock := mockCommittee(2, true, nil)
	bc := newBlockCollector(committeeMock)
	bc.UpdateRound(1)

	blockHash := []byte("pippo")
	bc.Unmarshaller = mockBEUnmarshaller(blockHash, 1, 1)
	bc.Collect(bytes.NewBuffer([]byte{}))
	bc.Collect(bytes.NewBuffer([]byte{}))

	select {
	case res := <-bc.BlockChan:
		assert.Equal(t, blockHash, res)
		// testing that we clean after collection
		assert.Equal(t, 0, len(bc.StepEventCollector))
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Collection did not complete")
	}
}

func TestNoQuorumCollection(t *testing.T) {
	committeeMock := mockCommittee(3, true, nil)
	bc := newBlockCollector(committeeMock)
	bc.UpdateRound(1)

	blockHash := []byte("pippo")
	bc.Unmarshaller = mockBEUnmarshaller(blockHash, 1, 1)
	bc.Collect(bytes.NewBuffer([]byte{}))
	bc.Collect(bytes.NewBuffer([]byte{}))

	select {
	case <-bc.BlockChan:
		assert.Fail(t, "Collection was not supposed to complete since Quorum should not be reached")
	case <-time.After(100 * time.Millisecond):
		// testing that we still have collected for 1 step
		assert.Equal(t, 1, len(bc.StepEventCollector))
		// testing that we collected 2 messages
		assert.Equal(t, 2, len(bc.StepEventCollector[string(1)]))
	}
}

func TestSkipNoMember(t *testing.T) {
	committeeMock := mockCommittee(1, false, nil)
	bc := newBlockCollector(committeeMock)
	bc.UpdateRound(1)

	blockHash := []byte("pippo")
	bc.Unmarshaller = mockBEUnmarshaller(blockHash, 1, 1)
	bc.Collect(bytes.NewBuffer([]byte{}))

	select {
	case <-bc.BlockChan:
		assert.Fail(t, "Collection was not supposed to complete since Quorum should not be reached")
	case <-time.After(50 * time.Millisecond):
		// test successfull
	}
}

func TestBlockNotary(t *testing.T) {
	bus := wire.New()
	committee := mockCommittee(1, true, nil)
	notary := LaunchBlockNotary(bus, committee)
	blockHash := []byte("pippo")
	notary.blockCollector.Unmarshaller = mockBEUnmarshaller(blockHash, 1, 1)

	blockChan := make(chan *bytes.Buffer)
	bus.Subscribe(msg.PhaseUpdateTopic, blockChan)
	bus.Publish(msg.BlockAgreementTopic, bytes.NewBuffer([]byte("test")))

	result := <-blockChan
	assert.Equal(t, result.String(), "pippo")

	// test that after a round update messages for previous rounds get ignored
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(2))
	bus.Publish(msg.RoundUpdateTopic, bytes.NewBuffer(b))

	// we need to wait for the round update to be propagated before publishing other round related messages. This is what this timeout is about
	<-time.After(100 * time.Millisecond)

	notary.blockCollector.Unmarshaller = mockBEUnmarshaller(blockHash, 1, 2)
	bus.Publish(msg.BlockAgreementTopic, bytes.NewBuffer([]byte("test")))

	select {
	case <-blockChan:
		assert.FailNow(t, "Previous round messages should not trigger a phase update")
	case <-time.After(100 * time.Millisecond):
		// all well
		return
	}

}

type MockBEUnmarshaller struct {
	event *BlockEvent
	err   error
}

func (m *MockBEUnmarshaller) Unmarshal(b *bytes.Buffer, e wire.Event) error {
	if m.err != nil {
		return m.err
	}

	blsPub, _ := crypto.RandEntropy(32)
	ev := e.(*BlockEvent)
	ev.Step = m.event.Step
	ev.Round = m.event.Round
	ev.BlockHash = m.event.BlockHash
	ev.PubKeyBLS = blsPub
	return nil
}

func mockBEUnmarshaller(blockHash []byte, round uint64, step uint8) wire.EventUnmarshaller {
	ev := committee.NewNotaryEvent()
	ev.BlockHash = blockHash
	ev.Step = step
	ev.Round = round

	return &MockBEUnmarshaller{
		event: ev,
		err:   nil,
	}
}

func mockCommittee(quorum int, isMember bool, verification error) committee.Committee {
	committeeMock := &mocks.Committee{}
	committeeMock.On("Quorum").Return(quorum)
	committeeMock.On("VerifyVoteSet",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(verification)
	committeeMock.On("IsMember", mock.AnythingOfType("[]uint8")).Return(isMember)
	return committeeMock
}