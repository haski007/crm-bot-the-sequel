package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/google/uuid"
	"time"
)

func (bot *CrmBotService) changeCashAmount(diff model.Money, author string, comment string) error {
	var txType model.TxType
	if diff > 0 {
		txType = model.TxAddCash
	} else {
		txType = model.TxGetCash
	}

	if err := bot.TransactionRepository.Add(model.Transaction{
		ID:        uuid.New().String(),
		Author:    author,
		Amount:    diff,
		Type:      txType,
		CreatedAt: time.Now(),
		Comment:   comment,
	}); err != nil {
		return fmt.Errorf("TransactionRepository.Add | err: %s", err)
	}

	if err := bot.CashRepository.ChangeAmount(diff); err != nil {
		return fmt.Errorf("TransactionRepository.Add | err: %s", err)
	}
	return nil
}
