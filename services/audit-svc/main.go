// services/audit-svc/main.go
package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/quantum-sentinel-nexus/audit-svc/chaincode" // Import the chaincode package
)

func main() {
	// Create a new instance of our smart contract.
	auditChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating audit-svc chaincode: %v", err)
	}

	// Start the chaincode, making it ready to listen for invocations.
	if err := auditChaincode.Start(); err != nil {
		log.Panicf("Error starting audit-svc chaincode: %v", err)
	}
}
