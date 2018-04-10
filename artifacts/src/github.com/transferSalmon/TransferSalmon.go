/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

 package main
 
 import (
	 "encoding/json"
	 "fmt"
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 "github.com/hyperledger/fabric/protos/peer"
 )

 type TransferSalmon struct {}

 type TransferSalmonModel struct {
	 DateTime string  	`json:"timeStamp"`
	 FromSeller string `json:"from"`
	 ToBuyer string 	`json:"to"`
 }

 // Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SellSalmon) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables00000 or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *TransferSalmon) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	if function == "queryTransferSalmon" {
		return t.queryTransferSalmon(stub, args)
	} else if function == "transferSalmon" {
		return t.transferSalmon(stub, args)
	}
	return shim.Error("Received unknown function invocation")
}

func (t *TransferSalmon) transferSalmon(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 4 {
		return shim.Error("Incorrect number of Args, expection 1")
	}

	var dateTime = args[1]
	var fromSeller = args[2]
	var toBuyer = args[3]

	var transferSalmonModel = TransferSalmonModel{DateTime: dateTime, FromSeller: fromSeller, ToBuyer: toBuyer}
	transferSalmonModelAsBytes, _ := json.Marshal(transferSalmonModel)

	err := stub.PutState(args[0], transferSalmonModelAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record salmon catch %s", args[0]))
	}

	return shim.Success(nil)
}

func (t *TransferSalmon) queryTransferSalmon(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of Args, expection 1")
	}

	transferSalmonAsBytes, _ := stub.GetState(args[0])
	if transferSalmonAsBytes == nil {
		return shim.Error("Could not find the salmon")
	}

	return shim.Success(transferSalmonAsBytes)
}

func main() {
	if err := shim.Start(new(TransferSalmon)); err != nil {
		fmt.Printf("Error starting SettingPrice chaincode: %s", err)
	}
}

