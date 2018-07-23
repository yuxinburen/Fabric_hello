package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

//查询账户余额
type PaymentChaincode struct {
}

func (p *PaymentChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 4 {
		return shim.Error("必须指定两个账户名称及相应的初始余额")
	}
	var a = args[0]
	var avalstr = args[1]
	var b = args[2]
	var bvalstr = args[3]

	if len(a) < 2 {
		return shim.Error("账户名称不能少于两个字符长度")
	}

	if len(b) < 2 {
		return shim.Error("账户名称不能少于两个字符长度")
	}

	_, err := strconv.Atoi(avalstr)
	if err != nil {
		return shim.Error("指定的账户初始值错误," + avalstr)
	}

	_, error := strconv.Atoi(bvalstr)

	if error != nil {
		return shim.Error("指定的账户余额初始值错误," + bvalstr)
	}
	err = stub.PutState(a, []byte(avalstr))
	if err != nil {
		return shim.Error("保存a账户数据时发生错误,")
	}
	err = stub.PutState(b, []byte(bvalstr))
	if err != nil {
		return shim.Error("保存b账户数据时发生错误")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *PaymentChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//获取用户意图
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "find" {
		return find(stub, args)
	} else if funcName == "payment" {
		return payment(stub, args)
	} else if funcName == "deleteAccount" {
		return del(stub, args)
	} else if funcName == "" {
	}
	return shim.Error("调用的功能未定义")
}

func find(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定要查询的账户名称")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("查询账户余额失败")
	}
	if result == nil {
		return shim.Error("根据指定账户查询，没有查询道结果")
	}
	return shim.Success(result)
}

//实现转账功能
// -c '{"Args",["from","to","value"]}'
func payment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("参数调用不合法，请检查")
	}
	var source, target string
	var value string
	source = args[0]
	target = args[1]
	value = args[2]

	//转账，from账户扣除value,to账户增加value
	fval, err := stub.GetState(source)
	if err != nil {
		return shim.Error("查询原账户失败")
	}

	tval, error := stub.GetState(target)
	if error != nil {
		return shim.Error("查询转账目标账户失败")
	}

	//实现转账
	_, err = strconv.Atoi(value)
	if err != nil {
		return shim.Error("指定的转账金额错误")
	}

	svi, err := strconv.Atoi(string(fval))
	if err != nil {
		return shim.Error("转出账户余额处理错误")
	}

	tvi, err := strconv.Atoi(string(tval))
	if err != nil {
		return shim.Error("转行目标账户余额处理错误")
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return shim.Error("转账金额转换失败")
	}
	if svi < val {
		return shim.Error("指定的转出账户余额不足，无法执行转账操作")
	}

	svi -= val
	tvi += val

	//将修改后的原账户和目标账户数据保存至账本中
	err = stub.PutState(source, []byte(strconv.Itoa(svi)))
	if err != nil {
		return shim.Error("转账后原账户数据更新错误")
	}

	err = stub.PutState(target, []byte(strconv.Itoa(tvi)))
	if err != nil {
		return shim.Error("目标账户数据更新错误")
	}

	return shim.Success([]byte("转账成功"))
}

//根据指定账户名称删除相应的信息
// -c '{"Args":["del","账户名称"]}'
func del(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定要删除的账户名称")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("查询失败")
	}
	if result != nil {
		return shim.Error("根据指定的账户名，没有查找道相应的余额")
	}
	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除指定的账户失败")
	}
	return shim.Success([]byte("删除指定账户成功" + args[0]))
}

func main() {
	err := shim.Start(new(PaymentChaincode))
	if err != nil {
		fmt.Println("启动链码失败")
	}
}
