package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type UserRating struct {
	User    string  `json:"user"`
	Average float64 `json:"average"`
	Pnum    int     `json:"projectNum"`
	Ptitle  string  `json:"ptitle"`
	Pscore  float64 `json:"pscore"`
	Pstate  string  `json:"pstate"` // 지원, 선정, 시작, 종료(rated), 프리(free)
}

type HistoryQueryResult struct {
	Record    	*UserRating    `json:"record"`
	TxId     	string    		`json:"txId"`
	Timestamp 	time.Time 		`json:"timestamp"`
	IsDelete  	bool      		`json:"isDelete"`
}

//func (s *SmartContract) addUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, username string) error {

	var user = UserRating{User: username, Average: 0, Pnum: 0, Pstate: "free"}
	
	userAsBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(username, userAsBytes)

}

//func (s *SmartContract) addRating(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
func (s *SmartContract) AddRating(ctx contractapi.TransactionContextInterface, username string, projectname string, projectscore string) error {

	// getState User
	userAsBytes, err := ctx.GetStub().GetState(username)
	if err != nil {
		return fmt.Errorf("Failed to get state for %s", err.Error())
	} else if userAsBytes == nil { // no State! error
		return fmt.Errorf("User does not exist: %s", username)
	}
	// state ok
	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return err
	}
	// create rate structure
	newRate, _ := strconv.ParseFloat(projectscore, 64)

	rateCount := float64(user.Pnum)

	user.Ptitle = projectname
	user.Pscore = newRate

	user.Average = (rateCount*user.Average + newRate) / (rateCount + 1)

	user.Pnum = user.Pnum + 1
	user.Pstate = "rated"

	// update to User World state
	userAsBytes, err = json.Marshal(user)

	return ctx.GetStub().PutState(username, userAsBytes)

}

// func (s *SmartContract) ReadRating(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
func (s *SmartContract) readRating(ctx contractapi.TransactionContextInterface, username string) (*UserRating, error) {

	UserAsBytes, err := ctx.GetStub().GetState(username)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if UserAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", username)
	}

	user := new(UserRating)
	_ = json.Unmarshal(UserAsBytes, user)

	return user, nil
	
}
// func (s *SmartContract) getHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {
func (s *SmartContract) getHistory(ctx contractapi.TransactionContextInterface, username string)([]HistoryQueryResult, error) {

	log.Printf("getHistory: ID %v", username)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(username)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user UserRating
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &user)
			if err != nil {
				return nil, err
			}
		} else {
			user = UserRating{
				User: username,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &user,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})

	if err != nil {
		log.Panicf("Error creating teamate chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting teamate chaincode: %v", err)
	}
}
