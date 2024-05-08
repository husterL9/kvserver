package db

import "github.com/husterL9/kvserver/internal/kvstore"

type ReadView struct {
	//Read View 创建时其他未提交的活跃事务 ID 列表。
	activeTxIDs []int64
	//创建该 Read View 的事务 ID。
	creatorTxID int64
	//不可见下限，目前出现过的最大的事务 ID+1，即下一个将被分配的事务 ID。大于等于这个 ID 的数据版本均不可见。
	lowLimitID int64
	//可见上限 ，activeTxIDs中最小的事务 ID，如果 activeTxIDs 为空，则 lowLimitID 为 upLimitID 。小于这个 ID 的数据版本均可见。
	upLimitID int64
}

type Tx struct {
	db           *KVStore
	txID         int64
	commits      []*Record
	undos        []*Record
	readView     ReadView
	readViewUsed bool // 标志位，表示ReadView是否已经被生成并使用
}

// 生成readView
func (tx *Tx) genRV() error {
	if tx.readViewUsed && tx.db.tm.isolation == RepeatableRead {
		return nil // 对于RR, 如果已生成ReadView则重用
	}
	tx.readViewUsed = true
	var upLimitID int64
	// 如果 activeTxIDs 为空，则 lowLimitID 为 upLimitID
	if len(tx.db.tm.activeTxIDs) == 0 {
		upLimitID = tx.db.tm.nextTxID + 1
	} else {
		upLimitID = tx.db.tm.activeTxIDs[0]

	}
	tx.readView = ReadView{
		activeTxIDs: tx.db.tm.activeTxIDs,
		creatorTxID: tx.txID,
		lowLimitID:  tx.db.tm.nextTxID + 1,
		upLimitID:   upLimitID,
	}

	return nil
}

func (tx *Tx) commit() error {

}

func (tx *Tx) rollBack() {

}

func (tx *Tx) Set(key string, value []byte, meta kvstore.MetaData) {
	tx.db.store.Set(key, value, meta)
}

func (tx *Tx) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	val, ok := tx.db.store.Get(args)
	return kvstore.GetResponse{
		Value: val,
	}, ok
}

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
