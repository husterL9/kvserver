package db

import (
	"fmt"

	"github.com/husterL9/kvserver/internal/kvstore"
)

// 开启事务
func (db *KVStore) Tx(f func(tx *Tx) error) error {
	tx, err := db.tm.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.rollBack()
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	err = tx.commit()

	return err
}

func (db *KVStore) loadRecord(record *Record) error {
	op := OpFromRecord(record)
	switch record.Op {
	case uint16(ModifyOp):
		db.rollVersion(op)
	default:
		return fmt.Errorf("invalid operation %d", record.Op)
	}
	return nil
}

func (db *KVStore) rollVersion(op *Op) {
	key := op.key
	txId := op.txID
	item, ok := db.store.Get(key)
	if !ok {
		return
	}
	//找到xID为op.txID的version
	var prev *kvstore.Version
	var current = item.Version
	for current != nil {
		if current.TxID == txId {
			//删除这个version
			if prev == nil {
				item.Version = current.Next
			} else {
				prev.Next = current.Next
			}
		}
		prev = current
		current = current.Next
	}
	if current == nil {
		return
	}
	if prev == nil {
		item.Version = current.Next
	} else {
		prev.Next = current.Next
	}
}
