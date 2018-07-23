package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

type SimpleChaincode struct {

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
	args := stub.GetStringArgs()
	if len(args) != 2{
		return shim.Error("指定了错误的参数个数")
	}

	err := stub.PutState(args[0],[]byte(args[1]))
	if err !=nil{
		return shim.Error("保存数据发生错误")
	}
	fmt.Println("数据保存成功")
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fun,args:=stub.GetFunctionAndParameters()
	if fun == "query"{
		return query(stub,args)
	}
	return shim.Error("指定的功能没有定义")
}

func query(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("只能指定相应的key")
	}
	ret,err := stub.GetState(args[1])
	if err != nil{
		return shim.Error("查询数据时发生错误")
	}
	if ret == nil{
		return shim.Error("没有查询道相应的数据")
	}

	//返回查询结果
	return shim.Success(ret)
}


//main方法
func main(){

	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Println("链码启动失败",err)
	}

}
