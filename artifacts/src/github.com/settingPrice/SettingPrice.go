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
	"strconv"
)

type SettingPrice struct {
}

type SalmonPrice struct {
	Price int    `json:"price"`
	Buyer string `json:"buyer"`
	Seller string `json:"seller"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SettingPrice) Init(stub shim.ChaincodeStubInterface) peer.Response {
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
func (t *SettingPrice) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	if function == "setupPriceToBuyer" {
		return t.setupPriceToBuyer(stub, args)
	} else if function == "querySettingPrice" {
		return t.querySettingPrice(stub, args)
	}
	return shim.Error("Received unknown function invocation")
}


func (t *SettingPrice) setupPriceToBuyer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 4 {
		return shim.Error("Incorrect number of Args, expection 1")
	}
	var id = args[0]
	price, _ := strconv.Atoi(args[1])
	var buyerName = args[2]
	var sellerName = args[3]
	
	var salmonPrice = SalmonPrice{Price: price, Buyer: buyerName, Seller: sellerName}
	salmonPriceByte, _ := json.Marshal(salmonPrice)
	err := stub.PutState(id, salmonPriceByte)

	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record salmon catch %s", id))
	}

	return shim.Success(nil)
}

func (t *SettingPrice) querySettingPrice(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of Args, expection 1")
	}

	salmonPriceAsBytes, _ := stub.GetState(args[0])
	if salmonPriceAsBytes == nil {
		return shim.Error("Could not find the salmon")
	}

	return shim.Success(salmonPriceAsBytes)
}

func main() {
	if err := shim.Start(new(SettingPrice)); err != nil {
		fmt.Printf("Error starting SettingPrice chaincode: %s", err)
	}
}