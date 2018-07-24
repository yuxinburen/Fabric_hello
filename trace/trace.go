package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"encoding/json"
)

//保存账户数据
func saveAccount(stub shim.ChaincodeStubInterface, account Account) bool {
	acc, err := json.Marshal(account)
	if err != nil {
		return false
	}
	err = stub.PutState(account.CardNo, acc)
	if err != nil {
		return false
	}
	return true
}

//实现贷款功能
//-c '{"Args":["loan","idcard","bankname","loanvalue"]}'
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	am, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("给定的贷款金额错误")
	}
	bank := Bank{
		BankName:  args[1],
		Flag:      Bank_Flag_Loan,
		Amount:    am,
		StartDate: "20100901",
		EndDate:   "20101201",
	}
	account := Account{
		CardNo: args[0],
		Aname:  "alice",
		Age:    39,
		Gender: "男",
		Mobil:  "13167582311",
		Bank:   bank,
	}
	//将要保存的对象进行序列化处理
	saveResult := saveAccount(stub, account)
	if !saveResult {
		return shim.Error("保存贷款纪录失败")
	}
	return shim.Success([]byte("贷款成功"))
}

func repayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	am, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("给定的贷款金额错误")
	}
	bank := Bank{
		BankName:  args[1],
		Flag:      Bank_Flag_Repayment,
		Amount:    am,
		StartDate: "20100901",
		EndDate:   "20101201",
	}
	account := Account{
		CardNo: args[0],
		Aname:  "alice",
		Age:    39,
		Gender: "男",
		Mobil:  "13167582311",
		Bank:   bank,
	}
	//将要保存的对象进行序列化处理
	saveResult := saveAccount(stub, account)
	if !saveResult {
		return shim.Error("保存还款纪录失败")
	}
	return shim.Success([]byte("还款成功"))
}

//根据身份证号码查询账户信息
func GetAccountByNo(stub shim.ChaincodeStubInterface, cardNo string) (Account, bool) {
	var account Account
	result, err := stub.GetState(cardNo)
	if err != nil {
		return account, false
	}
	err = json.Unmarshal(result, &account)
	if err != nil {
		return account, false
	}
	return account, true
}

//根据账户身份证号查询相应的信息(包含该账户的历史纪录信息)
//－c '{"Args":[]}'
func queryAccountByIdCard(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("功能调用参数设置错误，请核对后重新调用")
	}
	account, exit := GetAccountByNo(stub, args[0])
	if !exit {
		return shim.Error("查询账户出错,请重试")
	}
	//历史纪录信息
	historIterator, erro := stub.GetHistoryForKey(account.CardNo)
	if erro != nil {
		return shim.Error("溯源历史纪录查询失败")
	}
	defer historIterator.Close()

	var historys []HistoryItem
	//处理查询到的历史纪录消息迭代对象
	var acc Account
	for historIterator.HasNext() {
		//依次获取迭代器中的元素
		hisData, err := historIterator.Next()
		if err != nil {
			return shim.Error("处理迭代器数据出错")
		}
		var histItem HistoryItem
		histItem.TxId = hisData.TxId
		error := json.Unmarshal(hisData.Value, &acc)
		if error != nil {
			return shim.Error("查询历史数据出错")
		}

		if hisData.Value == nil { //处理当前纪录为nil的情况
			var empty Account
			histItem.Account = empty
		} else {
			histItem.Account = acc
		}
		//将当前处理完毕的历史状态保存至数组中
		historys = append(historys, histItem)
	}
	accByte, err := json.Marshal(historys)
	if err != nil {
		return shim.Error("账户序列化错误,请重试")
	}
	return shim.Success(accByte)
}
