package db

import (
	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/husterL9/kvserver/internal/wal"
)

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
	tm := tx.db.tm
	if tx.readViewUsed && tm.isolation == RepeatableRead {
		return nil // 对于RR, 如果已生成ReadView则重用
	}
	tx.readViewUsed = true
	var upLimitID int64
	// 如果 activeTxIDs 为空，则 lowLimitID 为 upLimitID
	if len(tm.activeTxIDs) == 0 {
		upLimitID = tm.nextTxID + 1
	} else {
		upLimitID = tx.db.tm.activeTxIDs[0]

	}
	activeTxIDs := make([]int64, len(tm.activeTxIDs))
	copy(activeTxIDs, tm.activeTxIDs)
	tx.readView = ReadView{
		activeTxIDs: activeTxIDs,
		creatorTxID: tx.txID,
		lowLimitID:  tm.nextTxID + 1,
		upLimitID:   upLimitID,
	}

	return nil
}

func (tx *Tx) commit() error {
	if len(tx.commits) == 0 {
		return nil
	}

	batch := new(wal.Batch)
	for _, commit := range tx.commits {
		data, err := commit.Encode()
		if err != nil {
			return err
		}
		batch.Write(data)
	}

	if n, err := tx.db.wal.WriteBatch(batch, wal.WithSync, wal.WithAtomic); err != nil {
		if n > 0 {
			tx.db.wal.Truncate(n)
		}
		tx.commits = nil
		return err
	}

	// notify all commit
	// for _, commit := range tx.commits {
	// 	tx.db.notify(&Op{
	// 		key: string(commit.Key),
	// 		val: commit.Val,
	// 		op:  OpType(commit.Op),
	// 	})
	// }

	//将txID从activeTxIDs中删除
	tx.db.tm.activeTxIDs = Remove(tx.db.tm.activeTxIDs, tx.txID)
	tx.commits = nil
	tx.undos = nil

	return nil
}

func (tx *Tx) rollBack() error {
	tx.db.lock.Lock()
	defer tx.db.lock.Unlock()
	tx.commits = nil
	for _, undo := range tx.undos {
		if err := tx.db.loadRecord(undo); err != nil {
			return err
		}
	}

	tx.undos = nil
	tx.db = nil
	return nil
}

func (tx *Tx) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	tx.genRV()
	currentTxId := tx.txID
	item, ok := tx.db.get(args)
	itemTxId := item.Version.TxID
	verison := item.Version
	// 遍历版本链找到合适的版本,RR
	for verison != nil {
		if itemTxId < tx.readView.upLimitID {
			return kvstore.GetResponse{Value: item.Version.Value}, true
		}
		if itemTxId == currentTxId {
			return kvstore.GetResponse{Value: item.Version.Value}, true
		}
		//如果介于up和low之间，说明在创建ReadView时生成该版本的事务仍处于活跃状态，因此该版本不能被访问
		if itemTxId < tx.readView.lowLimitID && !Contains(tx.readView.activeTxIDs, itemTxId) {
			return kvstore.GetResponse{Value: item.Version.Value}, true
		}
		verison = verison.Next
	}
	//没有找到已提交的版本
	return kvstore.GetResponse{Value: nil}, ok
}
func (tx *Tx) Set(args kvstore.SetArgs, opts ...Option) {
	key := args.Key
	val := args.Value
	meta := args.Meta
	op := setOption(key, val, opts...)
	undo := &Op{key: key}
	err := tx.db.set(key, val, meta)
	if err == nil {
		undo.op = ModifyOp
		undo.txID = tx.txID
	}
	tx.undos = append(tx.undos, NewRecord(undo))
	tx.commits = append(tx.commits, NewRecord(op))
}
