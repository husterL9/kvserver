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

}

func (tx *Tx) rollBack() {

}
func (tx *Tx) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	tx.genRV()
	currentTxId := tx.txID
	item, ok := tx.db.get(args)
	itemTxId := item.TxID
	// 遍历版本链找到合适的版本,RR
	for item != nil {
		if itemTxId < tx.readView.upLimitID {
			return kvstore.GetResponse{Value: item.Value}, true
		}
		if itemTxId == currentTxId {
			return kvstore.GetResponse{Value: item.Value}, true
		}
		//如果介于up和low之间，说明在创建ReadView时生成该版本的事务仍处于活跃状态，因此该版本不能被访问
		if itemTxId < tx.readView.lowLimitID && !Contains(tx.readView.activeTxIDs, itemTxId) {
			return kvstore.GetResponse{Value: item.Value}, true
		}
		item = item.Next
	}
	//没有找到已提交的版本
	return kvstore.GetResponse{Value: nil}, ok
}
func (tx *Tx) Set(args kvstore.SetArgs) {
	key := args.Key
	clientId := args.ClientId
	opId := args.OpId
	oldItem, ok := tx.db.get(kvstore.GetArgs{
		Key:      key,
		ClientId: clientId,
		OpId:     opId,
	})
	if ok {
		//插入到链表头部
		newItem := &kvstore.Item{
			Value: args.Value,
			TxID:  tx.txID,
			Next:  oldItem,
		}
		tx.db.store.Set(key, newItem)
	} else {

	}

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
