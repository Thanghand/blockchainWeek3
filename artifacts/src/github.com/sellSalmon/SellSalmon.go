/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

var countSalmon = 0

// SellSalmon implements a simple chaincode to manage an assetuh
type SellSalmon struct {
}

// Salmon struct is salmon
type Salmon struct {
	Vessel   string `json:"vessel"`
	DateTime string `json:"datetime"`
	Location string `json:"location"`
	Holder   string `json:"holder"`
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
func (t *SellSalmon) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	if function == "querySalmon" {
		return t.querySalmon(stub, args)
	} else if function == "recordSalmon" {
		return t.recordSalmon(stub, args)
	} else if function == "changeSalmonHolder" {
		return t.changeSalmonHolder(stub, args)
	} else if function == "queryAllSalmon" {
		return t.queryAllSalmon(stub, args)
	}
	return shim.Error("Received unknown function invocation")
}

func (t *SellSalmon) recordSalmon(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments, Expecting")
	}

	var vessel = args[1]
	var dateTime = args[2]
	var location = args[3]
	var holder = args[4]
	var salmon = Salmon{Vessel: vessel, DateTime: dateTime, Location: location, Holder: holder}

	salmonAsBytes, _ := json.Marshal(salmon)
	err := stub.PutState(args[0], salmonAsBytes)
	countSalmon++
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record salmon catch %s", args[0]))
	}

	return shim.Success(nil)
}

func (t *SellSalmon) changeSalmonHolder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of Args, expection 1")
	}

	salmonAsBytes, _ := stub.GetState(args[0])
	salmon := Salmon{}
	err := json.Unmarshal(salmonAsBytes, salmon)
	if err != nil {
		return shim.Error("Bad Data was stored for catch")
	}
	salmon.Holder = args[1]

	salmonAsBytes, _ = json.Marshal(salmon)

	errorResult := stub.PutState(args[0], salmonAsBytes)

	if errorResult != nil {
		return shim.Error("Failed to changeSalmonHolder")
	}

	return shim.Success(nil)
}

func (t *SellSalmon) querySalmon(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of Args, expection 1")
	}

	salmonAsBytes, _ := stub.GetState(args[0])
	if salmonAsBytes == nil {
		return shim.Error("Could not find the salmon")
	}

	return shim.Success(salmonAsBytes)
}

func (t *SellSalmon) queryAllSalmon(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	startKey := "0"
	endKey := strconv.Itoa(countSalmon - 1)

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func main() {
	if err := shim.Start(new(SellSalmon)); err != nil {
		fmt.Printf("Error starting SellSalmon chaincode: %s", err)
	}
}