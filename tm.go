package db

import "sync"

// 定义隔离级别常量
const (
	ReadCommitted  = "READ COMMITTED"
	RepeatableRead = "REPEATABLE READ"
)

type TxManager struct {
	// 隔离级别
	isolation   string
	mu          sync.Mutex
	nextTxID    int64
	activeTxIDs []int64
}

func NewTxManager() *TxManager {
	return &TxManager{
		nextTxID:    0,
		activeTxIDs: make([]int64, 0),
		// 在隔离级别默认是RR
		isolation: RepeatableRead,
	}
}

func (tm *TxManager) BeginTx() (*Tx, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.nextTxID++
	txID := tm.nextTxID

	tm.activeTxIDs = append(tm.activeTxIDs, txID)

	return &Tx{
		txID:         txID,
		readViewUsed: false,
	}, nil
}

func (tm *TxManager) EndTransaction(trxID int64) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 从活动事务列表中移除事务ID
	for i, id := range tm.activeTxIDs {
		if id == trxID {
			tm.activeTxIDs = append(tm.activeTxIDs[:i], tm.activeTxIDs[i+1:]...)
			break
		}
	}
}
