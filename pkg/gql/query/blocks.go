package query

import (
	"encoding/hex"
	"errors"

	"github.com/dusk-network/dusk-blockchain/pkg/core/block"
	"github.com/dusk-network/dusk-blockchain/pkg/core/database"
	"github.com/graphql-go/graphql"

	core "github.com/dusk-network/dusk-blockchain/pkg/wallet/transactions"
)

// File purpose is to define all arguments and resolvers relevant to "blocks" query only

type blocks struct {
}

func (b blocks) getQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(Block),
		Args: graphql.FieldConfigArgument{
			"hash": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"hashes": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
			"height": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"range": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
			"last": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: b.resolve,
	}
}
func (b blocks) resolve(p graphql.ResolveParams) (interface{}, error) {

	// Retrieve DB conn from context
	db, ok := p.Context.Value("database").(database.DB)
	if !ok {
		return nil, errors.New("context does not store database conn")
	}

	// resolve argument hash (single block)
	hash, ok := p.Args["hash"].(interface{})
	if ok {
		hashes := make([]interface{}, 0)
		hashes = append(hashes, hash)
		return b.fetchBlocksByHashes(db, hashes)
	}

	// resolve argument hashes (multiple blocks)
	hashes, ok := p.Args["hashes"].([]interface{})
	if ok {
		return b.fetchBlocksByHashes(db, hashes)
	}

	// resolve argument height (single block)
	// Chain height type is uint64 whereas `resolve` can handle height values up to MaxUInt
	height, ok := p.Args["height"].(int)
	if ok {
		return b.fetchBlocksByHeights(db, int64(height), int64(height))
	}

	// resolve argument range (range of blocks)
	heightRange, ok := p.Args["range"].([]interface{})
	if ok && len(heightRange) == 2 {
		from, ok := heightRange[0].(int)
		if !ok {
			return nil, errors.New("range `from` value not int64")
		}

		to, ok := heightRange[1].(int)
		if !ok {
			return nil, errors.New("range `to` value not int64")
		}

		return b.fetchBlocksByHeights(db, int64(from), int64(to))
	}

	offset, ok := p.Args["last"].(int)
	if ok {
		if offset <= 0 {
			return nil, errors.New("invalid offset")
		}
		return b.fetchBlocksByHeights(db, int64(offset)*-1, -1)
	}

	return nil, nil
}

func resolveTxs(p graphql.ResolveParams) (interface{}, error) {

	var txs []*core.Standard
	b, ok := p.Source.(*block.Block)
	if ok {

		// Retrieve DB conn from context
		db, ok := p.Context.Value("database").(database.DB)
		if !ok {
			return nil, errors.New("context does not store database conn")
		}

		err := db.View(func(t database.Transaction) error {
			fetched, err := t.FetchBlockTxs(b.Header.Hash)
			if err != nil {
				return err
			}

			for _, tx := range fetched {
				sTx := tx.StandardTx()
				sTx.TxID, _ = tx.CalculateHash()
				txs = append(txs, sTx)
			}
			return nil
		})

		return txs, err
	}
	return nil, errors.New("invalid source block")
}

// Fetch block headers by a list of hashes
func (b blocks) fetchBlocksByHashes(db database.DB, hashes []interface{}) ([]*block.Block, error) {

	blocks := make([]*block.Block, 0)
	err := db.View(func(t database.Transaction) error {
		for _, v := range hashes {
			encodedHash, ok := v.(string)
			if !ok {
				continue
			}

			hash, err := hex.DecodeString(encodedHash)
			if err != nil {
				continue
			}

			header, err := t.FetchBlockHeader(hash)
			if err != nil {
				continue
			}

			// Reconstructing block with header only
			b := &block.Block{
				Header: header,
				Txs:    nil,
			}

			blocks = append(blocks, b)
		}
		return nil
	})

	return blocks, err
}

// Fetch block headers by a range of heights
func (b blocks) fetchBlocksByHeights(db database.DB, from, to int64) ([]*block.Block, error) {

	blocks := make([]*block.Block, 0)
	err := db.View(func(t database.Transaction) error {

		var tip uint64
		if from < 0 || to < 0 {
			var err error
			tip, err = t.FetchCurrentHeight()
			if err != nil {
				return err
			}
		}

		if from == -1 {
			from = int64(tip)
		}

		// For cases where last N blocks are required
		if from < 0 {
			from = int64(tip) + from + 1
			if from < 0 {
				from = 0
			}
		}

		if to == -1 {
			to = int64(tip)
		}

		for height := from; height <= to; height++ {
			hash, err := t.FetchBlockHashByHeight(uint64(height))
			if err != nil {
				return err
			}

			header, err := t.FetchBlockHeader(hash)
			if err != nil {
				continue
			}

			// Reconstructing block with header only
			b := &block.Block{
				Header: header,
				Txs:    nil,
			}

			blocks = append(blocks, b)
		}
		return nil
	})

	return blocks, err
}