package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

//该链码要实现的功能
//1.增加贷款功能
//2.增加还款功能
//3.根据账户名称查询相应的信息(包含该账户所有的历史纪录信息)

type TraceChaincode struct {

}

func (t *TraceChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

	return shim.Success(nil)
}
func (t *TraceChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{

	fun,args := stub.GetFunctionAndParameters()


	if fun == "loan"{
		return loan(stub,args)
	}else  if fun == "repayment"{
		return repayment(stub,args)
	}else if fun == "queryAccountByIdCard"{
		return queryAccountByIdCard(stub,args)
	}
	return shim.Success(nil)
}




func main(){

	error := shim.Start(new(TraceChaincode))
	if error!= nil{
		fmt.Println("启动链码失败")
	}
}
