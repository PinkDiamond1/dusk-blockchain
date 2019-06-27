package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "gitlab.dusk.network/dusk-core/dusk-go/pkg/core/database/lite"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/tests/helper"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
)

func TestDemoSaveFunctionality(t *testing.T) {

	eb := wire.NewEventBus()
	rpc := wire.NewRPCBus()
	chain, err := New(eb, rpc)

	assert.Nil(t, err)

	defer chain.Close()

	for i := 1; i < 5; i++ {
		nextBlock := helper.RandomBlock(t, 200, 10)
		nextBlock.Header.PrevBlockHash = chain.prevBlock.Header.Hash
		nextBlock.Header.Height = uint64(i)
		err = chain.AcceptBlock(*nextBlock)
		assert.NoError(t, err)
	}

	err = chain.AcceptBlock(chain.prevBlock)
	assert.Error(t, err)

}
