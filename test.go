package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

//简单的资产管理
//两个功能：查询余额，修改余额 get／set
//初始化时jack账户默认100

type TestChaincode struct {

}

func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 2{
		return shim.Error("给定参数个数错误")
	}
	err :=stub.PutState(args[0],[]byte(args[1]))
	if err!= nil{
		return shim.Error("保存数据时发生错误")
	}
	fmt.Println("保存数据成功")
	return shim.Success(nil)
}

func (t *TestChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fun,args :=stub.GetFunctionAndParameters()
	if fun == "get"{
		return get(stub,args)
	}else if fun == "set" {
		return set(stub,args)
	}
	return shim.Error("非法操作,没有定义的功能")
}

func  get(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args) != 1{
		return shim.Error("必须指定要查询的账户名称")
	}
	result ,err :=stub.GetState(args[0])
	if err != nil{
		fmt.Println("根据指定的账户名查询数据失败")
		return shim.Error("根据指定的账户名称查询数据失败")
	}
	if result == nil{
		return shim.Error("查询结果为空")
	}
	return shim.Success(result)

}

func set(stub shim.ChaincodeStubInterface,args []string)peer.Response{
	if len(args) !=2 {
		return shim.Error("给定的参数个数错误")
	}
	//数据验证
	err :=stub.PutState(args[0],[]byte(args[1]))
	if err != nil{
		return shim.Error("保存数据发生错误")
	}
	fmt.Println("保存数据成功")
	return shim.Success([]byte("保存数据成功"))
}

func main(){
	err := shim.Start(new(TestChaincode))
	if err!= nil{
		fmt.Errorf("启动链码失败,%v",err)
	}
}
