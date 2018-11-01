package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Party struct{
	Name string `json:"name"`
	Mandatory bool `json:"mandatory"`
	Signed bool `json:"signed"`
	Weight	float64 `json:"weight"`
}

type History struct {
	User string `json:"user"`
	Stage string `json:"stage"`
	Timestamp string `json:"timestamp"`
	Details string `json:"details"`
}

type Contract struct {
	CreatedBy string `json:"createdby"`
	CreatedOn string `json:"createdon"`
	Details string `json:"details"`
	Stage string `json:"stage"`
	DocName string `json:"docname"`
	DocHash string `json:"dochash"`
	DocID string `json:"docid"`
	Condition float64 `json:"condition"`
	Progress float64 `json:"progress"`
	Validation bool `json:"validation"`
	PartyList []Party `json:"partylist"`
	HistoryList []History `json:"historylist"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCon" {
		return s.queryCon(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addCon" {
		return s.addCon(APIstub, args)
	} else if function == "addParty" {
		return s.addParty(APIstub, args)
	} else if function == "conAct" {
		return s.conAct(APIstub, args)
	} else if function == "conSign" {
		return s.conSign(APIstub, args)
	} else if function == "conDoc" {
		return s.conDoc(APIstub, args)
	} else if function == "queryAllCon" {
		return s.queryAllCon(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	conAsBytes, _ := APIstub.GetState(args[0])

	contract := Contract{}

	json.Unmarshal(conAsBytes, &contract)

	fmt.Println("Queried", contract)

	//fmt.Println("Party Length", len(contract.PartyList))

	return shim.Success(conAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	contracts := []Contract{
			Contract {
				CreatedBy: "user1",
				CreatedOn: "1530540623",
				Details: "Sample contract",
				Stage: "Contract Validation",
				DocName: "contarct.pdf",
				DocHash: "96b26f6cc52edd91cd52ac5baa1a802f4ff04daab07a308f0b2e897cc807e4bb",
				PartyList: []Party{
						Party{Name: "user2", Mandatory: true, Signed: true, Weight: 1.0 },
						Party{Name: "user3", Mandatory: false, Signed: true, Weight: 1.0 },
						Party{Name: "user4", Mandatory: true, Signed: true, Weight: 1.0 },
					},
				HistoryList: []History{
						History{Details:"New contract created",Stage:"Contract Created",Timestamp:"1530540623",User:"user1"},
						History{Details:"Added Parties",Stage:"Contract Activation",Timestamp:"1530540725",User:"user1"},
						History{Details:"Added Parties",Stage:"Contract Activation",Timestamp:"1530540750",User:"user1"},
						History{Details:"Added Parties",Stage:"Contract Activation",Timestamp:"1530540756",User:"user1"},
						History{Details:"Contarct Activated for signing",Stage:"Contract Activation",Timestamp:"1530540795",User:"user1"},
						History{Details:"Document inforamtion added",Stage:"Contract Activation",Timestamp:"1530540795",User:"user1"},
						History{Details:"Party not authorized",Stage:"Contract Signing",Timestamp:"1530540874",User:"user1"},
						History{Details:"Mandatory condition not satisfied",Stage:"Contract Validation",Timestamp:"1530540874",User:"Smart Contract"},
						History{Details:"Party Signed",Stage:"Contract Signing",Timestamp:"1530540889",User:"user2"},
						History{Details:"Mandatory condition not satisfied",Stage:"Contract Validation",Timestamp:"1530540889",User:"Smart Contract"},
						History{Details:"Party Signed",Stage:"Contract Signing",Timestamp:"1530540901",User:"user3"},
						History{Details:"Mandatory condition not satisfied",Stage:"Contract Validation",Timestamp:"1530540901",User:"Smart Contract"},
						History{Details:"Party Signed",Stage:"Contract Signing",Timestamp:"1530540917",User:"user4"},
						History{Details:"Condition satisfied",Stage:"Contract Validation",Timestamp:"1530540917",User:"Smart Contract"},						
					},	 
				Condition: 2,
				Progress: 3,
				Validation: true },			
		}

	i := 0
	for i < len(contracts) {
		fmt.Println("i is ", i)
		conAsBytes, _ := json.Marshal(contracts[i])
		APIstub.PutState("CON"+strconv.Itoa(i), conAsBytes)
		fmt.Println("Added", contracts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) addCon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	t := time.Now()

	timestamp := strconv.FormatInt(t.Unix() , 10)

	var contarct = Contract{CreatedBy: args[1], CreatedOn: timestamp, Details: args[2], Stage: "Contract Created", Validation: false,
					HistoryList: []History{ History{User: args[1], Stage: "Contract Created", Timestamp: timestamp , Details: "New contract created" }}}

	conAsBytes, _ := json.Marshal(contarct)

	err := APIstub.PutState(args[0], conAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record new contract: %s", args[0]))
	}

	return shim.Success(nil)
}

func (s *SmartContract) addParty(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	t := time.Now()

	timestamp := strconv.FormatInt(t.Unix() , 10)

	conAsBytes, _ := APIstub.GetState(args[0])

	weight, _ := strconv.ParseFloat(args[4], 64)

	contract := Contract{}

	json.Unmarshal(conAsBytes, &contract)

	fmt.Println("Queried", contract)

	b,_ := strconv.ParseBool(args[2])

	stage := "Contract Activation"

	contract.Stage = stage

	party1 := Party{Name: args[1], Mandatory: b, Signed: false, Weight: weight }

	contract.PartyList = append(contract.PartyList, party1)

	history1 := History{User: args[3], Stage: "Contract Activation", Timestamp: timestamp, Details: "Added Parties" }

	contract.HistoryList = append(contract.HistoryList, history1)

	conAsBytes, _ = json.Marshal(contract)
	APIstub.PutState(args[0], conAsBytes)

	fmt.Println("contract updated", contract)
	return shim.Success(nil)
}


func (s *SmartContract) conAct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var stage string

	t := time.Now()

	timestamp := strconv.FormatInt(t.Unix() , 10)

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	conAsBytes, _ := APIstub.GetState(args[0])

	contract := Contract{}


	json.Unmarshal(conAsBytes, &contract)

	fmt.Println("Queried", contract)

	if contract.Stage != "Contract Activation" {
		return shim.Error("Not in Activation to move to signing stage")
	}


	condition, _ := strconv.ParseFloat(args[2], 64)

	contract.Condition = condition

	stage = "Contract Signing"

	contract.Stage = stage

	contract.DocName = args[3]

	contract.DocHash = args[4] 

	history1 := History{User: args[1], Stage: "Contract Activation", Timestamp: timestamp, Details: "Document information added" }

	contract.HistoryList = append(contract.HistoryList, history1)

	history2 := History{User: args[1], Stage: "Contract Activation", Timestamp: timestamp, Details: "Contract Activated for signing" }

	contract.HistoryList = append(contract.HistoryList, history2)

	conAsBytes, _ = json.Marshal(contract)
	APIstub.PutState(args[0], conAsBytes)

	fmt.Println("updated", contract)
	return shim.Success(nil)
}


func (s *SmartContract) conSign(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var history1, history2 History

	var count float64

	t := time.Now()

	timestamp := strconv.FormatInt(t.Unix() , 10)

	conAsBytes, _ := APIstub.GetState(args[0])

	contract := Contract{}

	json.Unmarshal(conAsBytes, &contract)

	fmt.Println("Queried", contract)

	if contract.Stage != "Contract Signing" {
		return shim.Error("Contract not in signing stage")
	}

	for i,party := range contract.PartyList {
		if party.Name == args[1] {
			if party.Signed == false {
				contract.PartyList[i].Signed = true
				history1 = History{User: args[1], Stage: "Contract Signing", Timestamp: timestamp, Details: "Party Signed" }
				break
			} else {
				history1 = History{User: args[1], Stage: "Contract Signing", Timestamp: timestamp, Details: "Already signed" }
				break
			}
		} else {
			history1 = History{User: args[1], Stage: "Contract Signing", Timestamp: timestamp, Details: "Party not authorized" }	
		}
	}
	
	contract.HistoryList = append(contract.HistoryList, history1)

	manFlag := false

	for _,party := range contract.PartyList {
		if party.Mandatory == true && party.Signed == false {
			manFlag = true
		}

		if party.Signed == true {
			count = party.Weight + count
		}
	}

	contract.Progress = count

	if contract.Validation == true{
		history2 = History{User: "Smart Contract", Stage: "Contract Validation", Timestamp: timestamp, Details: "Contract valid" }		

	}else if manFlag == true {
		history2 = History{User: "Smart Contract", Stage: "Contract Validation", Timestamp: timestamp, Details: "Mandatory condition not satisfied" }	
	}else if count < contract.Condition {
		history2 = History{User: "Smart Contract", Stage: "Contract Validation", Timestamp: timestamp, Details: "Weitage condition not satisfied" }
	}else {
		contract.Validation = true
		contract.Stage = "Contract Validation"
		history2 = History{User: "Smart Contract", Stage: "Contract Validation", Timestamp: timestamp, Details: "Condition satisfied" }
	}

	contract.HistoryList = append(contract.HistoryList, history2)

	conAsBytes, _ = json.Marshal(contract)
	APIstub.PutState(args[0], conAsBytes)

	fmt.Println("updated", contract)
	return shim.Success(nil)
}


func (s *SmartContract) conDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	t := time.Now()

	timestamp := strconv.FormatInt(t.Unix() , 10)

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	conAsBytes, _ := APIstub.GetState(args[0])

	contract := Contract{}

	json.Unmarshal(conAsBytes, &contract)

	contract.DocName = args[2]

	contract.DocHash = args[3] 

	history1 := History{User: args[1], Stage: "Contract Activation", Timestamp: timestamp, Details: "Document inforamtion added" }

	contract.HistoryList = append(contract.HistoryList, history1)

	conAsBytes, _ = json.Marshal(contract)
	APIstub.PutState(args[0], conAsBytes)

	fmt.Println("updated", contract)
	return shim.Success(nil)

}

func (s *SmartContract) queryAllCon(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CON0"
	endKey := "CON999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
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

	fmt.Printf("- queryAllCon:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}