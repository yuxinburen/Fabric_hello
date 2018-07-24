package main

const(
	Bank_Flag_Loan = 1
	Bank_Flag_Repayment=2
)

type Bank struct{
	BankName string `json:"BankName"`
	Flag int `json:"Flag"`
	Amount int `json:"Amount"`
	StartDate string `json:"StartDate"`
	EndDate string `json:"EndDate"`
}

type Account struct{
	CardNo string 	`json:"CardNo"`
	Aname string `json:"Aname"`
	Age int `json:"Age"`
	Gender string `json:"Gender"`
	Mobil string `json:"Mobil"`
	Bank Bank `json:"Bank"`
	Historys []HistoryItem
}

//溯源历史信息结构体
type HistoryItem struct {
	TxId string
	Account Account
}