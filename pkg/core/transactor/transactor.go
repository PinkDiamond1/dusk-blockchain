package transactor

import (
	"bytes"
	"fmt"
	"math/big"

	ristretto "github.com/bwesterb/go-ristretto"
	cfg "github.com/dusk-network/dusk-blockchain/pkg/config"
	"github.com/dusk-network/dusk-blockchain/pkg/core/block"
	"github.com/dusk-network/dusk-blockchain/pkg/core/database"
	"github.com/dusk-network/dusk-blockchain/pkg/core/database/heavy"
	"github.com/dusk-network/dusk-blockchain/pkg/core/transactions"
	"github.com/dusk-network/dusk-blockchain/pkg/wallet"
	"github.com/dusk-network/dusk-wallet/key"
	log "github.com/sirupsen/logrus"
)

var l = log.WithField("process", "transactor")

// TODO: rename
type Transactor struct {
	w *wallet.Wallet
}

// Instantiate a new Transactor struct.
func New(w *wallet.Wallet) *Transactor {
	return &Transactor{
		w: w,
	}
}

func (t *Transactor) CreateStandardTx(amount uint64, address string) transactions.Transaction {
	if err := t.syncWallet(); err != nil {
		l.WithError(err).Warnln("error syncing wallet")
		return nil
	}

	// Create a new standard tx
	// TODO: customizable fee
	tx, err := t.w.NewStandardTx(cfg.MinFee)
	if err != nil {
		l.WithError(err).Warnln("error creating transaction")
		return nil
	}

	// Turn amount into a scalar
	amountScalar := ristretto.Scalar{}
	amountScalar.SetBigInt(big.NewInt(0).SetUint64(amount))

	// Send amount to address
	tx.AddOutput(key.PublicAddress(address), amountScalar)

	// Sign tx
	err = t.w.Sign(tx)
	if err != nil {
		l.WithError(err).Warnln("error signing transaction")
		return nil
	}

	// Convert wallet-tx to wireTx
	// TODO: Unification of tx structures makes this obsolete. Remove when
	// unified tx structure is merged
	wireTx, err := tx.WireStandardTx()
	if err != nil {
		l.WithError(err).Warnln("error converting transaction")
		return nil
	}

	return wireTx
}

func (t *Transactor) sendStake(amount, lockTime uint64) transactions.Transaction {
	if err := t.syncWallet(); err != nil {
		l.WithError(err).Warnln("error syncing wallet")
		return nil
	}

	// Turn amount into a scalar
	amountScalar := ristretto.Scalar{}
	amountScalar.SetBigInt(big.NewInt(0).SetUint64(amount))

	// Create a new stake tx
	tx, err := t.w.NewStakeTx(cfg.MinFee, lockTime, amountScalar)
	if err != nil {
		l.WithError(err).Warnln("error creating stake")
		return nil
	}

	// Sign tx
	err = t.w.Sign(tx)
	if err != nil {
		l.WithError(err).Warnln("error signing stake")
		return nil
	}

	// Convert wallet-tx to wireTx and encode into buffer
	wireTx, err := tx.WireStakeTx()
	if err != nil {
		l.WithError(err).Warnln("error converting stake")
		return nil
	}

	return wireTx
}

func (t *Transactor) sendBid(amount, lockTime uint64) transactions.Transaction {
	if err := t.syncWallet(); err != nil {
		l.WithError(err).Warnln("error syncing wallet")
		return nil
	}

	// Turn amount into a scalar
	amountScalar := ristretto.Scalar{}
	amountScalar.SetBigInt(big.NewInt(0).SetUint64(amount))

	// Create a new bid tx
	tx, err := t.w.NewBidTx(cfg.MinFee, lockTime, amountScalar)
	if err != nil {
		l.WithError(err).Warnln("error creating bid")
		return nil
	}

	// Sign tx
	err = t.w.Sign(tx)
	if err != nil {
		l.WithError(err).Warnln("error signing bid")
		return nil
	}

	// Convert wallet-tx to wireTx and encode into buffer
	wireTx, err := tx.WireBid()
	if err != nil {
		l.WithError(err).Warnln("error converting bid")
		return nil
	}

	return wireTx
}

func (t *Transactor) syncWallet() error {
	var totalSpent, totalReceived uint64
	_, db := heavy.CreateDBConnection()
	// keep looping until tipHash = currentBlockHash
	for {
		// Get Wallet height
		walletHeight, err := t.w.GetSavedHeight()
		if err != nil {
			t.w.UpdateWalletHeight(0)
		}

		// Get next block using walletHeight and tipHash of the node
		blk, tipHash, err := fetchBlockHeightAndState(db, walletHeight)
		if err != nil {
			return fmt.Errorf("error fetching block from node db: %v\n", err)
		}

		// call wallet.CheckBlock
		spentCount, receivedCount, err := t.w.CheckWireBlock(*blk)
		if err != nil {
			return fmt.Errorf("error checking block: %v\n", err)
		}

		totalSpent += spentCount
		totalReceived += receivedCount

		// check if state is equal to the block that we fetched
		if bytes.Equal(tipHash, blk.Header.Hash) {
			break
		}
	}

	l.WithFields(log.Fields{
		"spends":   totalSpent,
		"receives": totalReceived,
	}).Debugln("finished wallet sync")
	return nil
}

func fetchBlockHeightAndState(db database.DB, height uint64) (*block.Block, []byte, error) {
	var blk *block.Block
	var state *database.State
	err := db.View(func(t database.Transaction) error {
		hash, err := t.FetchBlockHashByHeight(height)
		if err != nil {
			return err
		}
		state, err = t.FetchState()
		if err != nil {
			return err
		}

		blk, err = t.FetchBlock(hash)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return blk, state.TipHash, nil
}

func (t *Transactor) Balance() (float64, error) {
	if err := t.syncWallet(); err != nil {
		l.WithError(err).Warnln("error syncing wallet")
		return 0.0, err
	}

	balance, err := t.w.Balance()
	if err != nil {
		l.WithError(err).Warnln("error fetching balance")
		return 0.0, err
	}

	return balance, nil
}
