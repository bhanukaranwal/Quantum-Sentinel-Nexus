// services/audit-svc/chaincode/audit_contract.go
package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an audit log.
type SmartContract struct {
	contractapi.Contract
}

// AuditEvent describes the basic structure of an audit log entry.
// The JSON tags are used for serialization.
type AuditEvent struct {
	DocType   string    `json:"docType"` // Used to distinguish different document types in CouchDB
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"` // The service that generated the event (e.g., "scoring-svc")
	UserID    string    `json:"userID"`  // The user or system principal associated with the action
	Action    string    `json:"action"`  // The action being logged (e.g., "ModelUpdate", "HighValueTxnScored")
	Details   string    `json:"details"` // A JSON string containing arbitrary details about the event
}

// CreateAuditEvent records a new event to the audit log on the ledger.
// This function is invoked by a client application submitting a transaction proposal.
func (s *SmartContract) CreateAuditEvent(
	ctx contractapi.TransactionContextInterface,
	id string,
	service string,
	userID string,
	action string,
	details string,
) error {
	// 1. Check if an event with this ID already exists.
	exists, err := s.AuditEventExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the audit event %s already exists", id)
	}

	// 2. Get the transaction timestamp from the proposal.
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	// 3. Create the event struct.
	event := AuditEvent{
		DocType:   "AuditEvent",
		ID:        id,
		Timestamp: time.Unix(txTimestamp.GetSeconds(), int64(txTimestamp.GetNanos())),
		Service:   service,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	// 4. Serialize the event to JSON.
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 5. Put the state on the ledger.
	// The endorsement policy for this chaincode would require signatures from
	// multiple organizations, ensuring no single entity can tamper with the log.
	err = ctx.GetStub().PutState(id, eventJSON)
	if err != nil {
		return err
	}

	// 6. (Optional) Set an event to notify off-chain listeners.
	ctx.GetStub().SetEvent("AuditEventCreated", eventJSON)

	return nil
}

// GetAuditEvent retrieves an audit event from the ledger by its ID.
// This is a read-only operation and does not result in a transaction.
func (s *SmartContract) GetAuditEvent(ctx contractapi.TransactionContextInterface, id string) (*AuditEvent, error) {
	eventJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if eventJSON == nil {
		return nil, fmt.Errorf("the audit event %s does not exist", id)
	}

	var event AuditEvent
	err = json.Unmarshal(eventJSON, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// AuditEventExists returns true if an audit event with the given ID exists in the world state.
func (s *SmartContract) AuditEventExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	eventJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return eventJSON != nil, nil
}
