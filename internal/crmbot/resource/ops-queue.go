package resource

import "github.com/sirupsen/logrus"

type OperationType string

type step int

func (s step) Int() int {
	return int(s)
}

const (
	OperationType_CategoryAdd  OperationType = "CategoryAdd"
	OperationType_CategoryEdit OperationType = "CategoryEdit"

	OperationType_SupplierAdd  OperationType = "SupplierAdd"
	OperationType_SupplierEdit OperationType = "SupplierEdit"

	OperationType_ProductAdd  OperationType = "ProductAdd"
	OperationType_ProductEdit OperationType = "ProductEdit"

	OperationType_QuantityAdd OperationType = "QuantityAdd"
	OperationType_QuantitySet OperationType = "QuantitySet"
	OperationType_QuantityAll OperationType = "QuantityAll"

	OperationType_TransactionsGetAll OperationType = "TransactionsGetAll"

	OperationType_CashAdd OperationType = "CashAdd"

	OperationType_ProductGetByCategory OperationType = "ProductGetByCategory"

	OperationType_RevisionProcess OperationType = "RevisionProcess"
)

type Step int

type Operation struct {
	Name OperationType
	Step Step
	Data interface{}
}

func deleteQueuesOfUser(userID int) {
	if _, found := OpsQueue[userID]; found {
		logrus.Printf("Workflow %s - was rejected by user | user_id=%d", OpsQueue[userID].Name, userID)
		delete(OpsQueue, userID)
	}
}

// OpsQueue - map to store all queues for operations
// map[userID]Operation
var OpsQueue = make(map[int]*Operation)
