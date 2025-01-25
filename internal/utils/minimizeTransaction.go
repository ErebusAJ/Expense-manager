package utils

import (
	"log"
	"strconv"

	"github.com/google/uuid"
)

// convertBalanceStrToFloat
// converts netbalance from string to float64
func convertBalanceStrToFloat(balances map[uuid.UUID]string) map[uuid.UUID]float64{
	temp := make(map[uuid.UUID]float64)

	for key, item := range balances{
		val, err := strconv.ParseFloat(item, 64)
		if err != nil{
			log.Fatalf("error parsing amount %v", err)
		}
		temp[key] = val
	}

	return temp
}


type Transactions struct{
	FromUserID	uuid.UUID	`json:"from_user"`
	ToUserID	uuid.UUID	`json:"to_user"`
	Amount		string		`json:"amount"`
}
// MinimizeDebts 
// A utils function to minimize the transaction to settle debts
func MinimizeDebts(balances map[uuid.UUID]string) []Transactions{
	var creditors []uuid.UUID	// +ve balance
	var debtors []uuid.UUID	// -ve balance

	newBalances := convertBalanceStrToFloat(balances)

	for user, item := range newBalances{
		if item >= 0{
			creditors = append(creditors, user)
		}else{
			debtors = append(debtors, user)
		}
	}

	var transactions []Transactions
	
	cIdx, dIdx := 0, 0

	for cIdx < len(creditors) && dIdx < len(debtors){
		creditor := creditors[cIdx]
		debtor := debtors[dIdx]

		// determine amount settle
		amount := min(newBalances[creditor], -newBalances[debtor])

		//record payment

		transactions = append(transactions, Transactions{
			FromUserID: debtor,
			ToUserID: creditor,
			Amount: strconv.FormatFloat(amount, 'f', 2, 64),
		})

		newBalances[creditor] -= amount
		newBalances[debtor] += amount

		if newBalances[creditor] == 0{
			cIdx++
		}

		if newBalances[debtor] == 0{
			dIdx++
		}
	}



	return transactions
}