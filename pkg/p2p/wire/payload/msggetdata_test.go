package payload

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/dusk-go/crypto"
	"gitlab.dusk.network/dusk-core/dusk-go/transactions"
)

func TestMsgGetDataEncodeDecode(t *testing.T) {
	byte32 := []byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}

	// Input
	sig, _ := crypto.RandEntropy(2000)
	in := transactions.Input{
		KeyImage:  byte32,
		TxID:      byte32,
		Index:     1,
		Signature: sig,
	}

	// Output
	out := transactions.Output{
		Amount: 200,
		P:      byte32,
	}

	// Type attribute
	ta := transactions.TypeAttributes{
		Inputs:   []transactions.Input{in},
		TxPubKey: byte32,
		Outputs:  []transactions.Output{out},
	}

	R, _ := crypto.RandEntropy(32)
	s := transactions.Stealth{
		Version: 1,
		Type:    1,
		R:       R,
		TA:      ta,
	}

	msg := NewMsgGetData()
	msg.AddTx(s)

	// TODO: test Addblock function when block structure is decided
	buf := new(bytes.Buffer)
	if err := msg.Encode(buf); err != nil {
		t.Fatal(err)
	}

	msg2 := NewMsgGetData()
	if err := msg2.Decode(buf); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, msg, msg2)
}