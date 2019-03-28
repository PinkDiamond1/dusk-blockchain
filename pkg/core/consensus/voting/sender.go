package voting

import (
	"bytes"

	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/committee"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/msg"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/user"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire/topics"
)

// LaunchVotingComponent will start a sender, which initializes its own
// signers, and will proceed to listen on the event bus for relevant
// messages.
func LaunchVotingComponent(eventBus *wire.EventBus, keys *user.Keys,
	committee committee.Committee) *sender {

	sender := newSender(eventBus, keys, committee)
	go sender.listen()
	return sender
}

// sender is responsible for signing and sending out consensus related
// voting messages to the wire.
type sender struct {
	eventBus               *wire.EventBus
	blockReductionChannel  chan *bytes.Buffer
	blockAgreementChannel  chan *bytes.Buffer
	sigSetReductionChannel chan *bytes.Buffer
	sigSetAgreementChannel chan *bytes.Buffer
}

// newSender will return an initialized sender struct. It will also spawn all
// the needed signers, and their channels get connected to the sender.
func newSender(eventBus *wire.EventBus, keys *user.Keys, committee committee.Committee) *sender {
	blockReductionChannel := initCollector(eventBus, msg.OutgoingBlockReductionTopic,
		unmarshalBlockReduction, newBlockReductionSigner(keys, committee))
	sigSetReductionChannel := initCollector(eventBus, msg.OutgoingSigSetReductionTopic,
		unmarshalSigSetReduction, newSigSetReductionSigner(keys, committee))
	blockAgreementChannel := initCollector(eventBus, msg.OutgoingBlockAgreementTopic,
		unmarshalBlockAgreement, newBlockAgreementSigner(keys, committee))
	sigSetAgreementChannel := initCollector(eventBus, msg.OutgoingSigSetAgreementTopic,
		unmarshalSigSetAgreement, newSigSetAgreementSigner(keys, committee))

	return &sender{
		eventBus:               eventBus,
		blockReductionChannel:  blockReductionChannel,
		sigSetReductionChannel: sigSetReductionChannel,
		blockAgreementChannel:  blockAgreementChannel,
		sigSetAgreementChannel: sigSetAgreementChannel,
	}
}

// Listen will set the sender up to listen for incoming requests to vote.
func (v *sender) listen() {
	for {
		select {
		case m := <-v.blockReductionChannel:
			message, _ := addTopic(m, topics.BlockReduction)
			v.eventBus.Publish(string(topics.Gossip), message)
		case m := <-v.sigSetReductionChannel:
			message, _ := addTopic(m, topics.SigSetReduction)
			v.eventBus.Publish(string(topics.Gossip), message)
		case m := <-v.blockAgreementChannel:
			message, _ := addTopic(m, topics.BlockAgreement)
			v.eventBus.Publish(string(topics.Gossip), message)
		case m := <-v.sigSetAgreementChannel:
			message, _ := addTopic(m, topics.SigSetAgreement)
			v.eventBus.Publish(string(topics.Gossip), message)
		}
	}
}

func addTopic(m *bytes.Buffer, topic topics.Topic) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	topicBytes := topics.TopicToByteArray(topic)
	if _, err := buffer.Write(topicBytes[:]); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(m.Bytes()); err != nil {
		return nil, err
	}

	return buffer, nil
}