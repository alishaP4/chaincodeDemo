package main

import(
	"fmt"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math"
	"math/rand"
	//"code.google.com/p/go.crypto/ripemd160"
	"crypto/rsa"
)

const   PREFIX_PATIENT =  "patient"

type SimpleChaincode struct {
}

type Prescription  struct {
	Disease            	string 		`json:"disease"`
	Medication          string 		`json:"medication"`
	Duration      		string 		`json:"duration"`
}

type Patient struct {
	Name            	string 			`json:"name"`
	dob           		string 			`json:"dob"`
	CurrentProblem          string 			`json:"currentproblem"`
	allergies      		string 			`json:"currentproblem"`
	Role    		string			`json:"role"`
	Prescriptions 		[]Prescription 	        `json:"prescriptions"`
	Lab_Details             []Lab_Details           `json:"lab_details"`
}
type Lab_Details struct{
	Name_Lab  	  	string  `json:"name_lab"`
	Report_Type 	  	string  `json:"report_type"`
	Date       	   	string  `json:"date"`
	Impressions 		string  `json:"impressions"`
	Findings      	 	string  `json:"findings"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}
 
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "get_patient_details" {
            		return t.get_patient_details(stub, args)
	}

	return nil, errors.New("Received unknown function invocation " + function)

}
 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "enter_patient_details" {
		return t.enter_patient_details(stub, args)
	}else if function == "enter_patient_prescription" {
		return t.enter_patient_prescription(stub, args)
	}else if function == "enter_lab_details" {
		return t.enter_lab_details(stub, args)
	}
	

	return nil, errors.New("Received unknown function invocation " + function)
   
}
 
func GetPrivateKey() *rsa.PrivateKey{
	PrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil{
		fmt.Println(err.Error)	
	}
	return privateKey
}

func GetPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey)
*rsa.PublicKey {
	publicKey := &privateKey.PublicKey
	return publicKey
}

func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
    fmt.Println("Entering GetCertAttribute")
    attr, err := stub.ReadCertAttribute(attributeName)
    if err != nil {
        return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
    }
    attrString := string(attr)
    return attrString, nil
}
 

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Println("Could not start SampleChaincode")
    } else {
        fmt.Println("SampleChaincode successfully started")
    }
 
}




func (t *SimpleChaincode) enter_patient_details(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 5 { 
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}


	patient := Patient{}
	patient.Name = args[0]
	patient.dob = args[1]
	patient.CurrentProblem = args[2]
	patient.allergies =  args[3]
	patient.Role =  args[4]
	bytes, err := json.Marshal(&patient)

	if err != nil { 
		return nil, errors.New("Error converting Patient record")
	}

	err = stub.PutState(PREFIX_PATIENT + args[0], bytes)

	if err != nil { 
		return nil, errors.New("Error storing Patient record") 
	}

	return nil, nil
}


func (t *SimpleChaincode) get_patient_details(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 { 
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	patient, err := stub.GetState(PREFIX_PATIENT + args[0])
	if err != nil {
		return nil, errors.New("Failed to get state for " + args[0])
	}

	return patient, nil
}

func (t *SimpleChaincode) enter_lab_details(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//	    0              1		   2		  3             4            5
  	//	   Name       Name_Lab        Report_Type        Date      Impressions    Findings

	if len(args) != 6 { 
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	bytes, err := stub.GetState(PREFIX_PATIENT + args[0])
	if err != nil { 
		return nil, errors.New("No Patient with name " + args[0])
	}


lab_details := Lab_Details{}
	lab_details.Name_Lab = args[1]
	lab_details.Report_Type = args[2]
	lab_details.Date = args[3]
	lab_details.Impressions = args[4]
	lab_details.Findings = args[5]

	var patient Patient
 	err = json.Unmarshal(bytes,&patient)
 	patient.Lab_Details = append(patient.Lab_Details, lab_details)

	bytes, err = json.Marshal(&patient)
	if err != nil { 
		return nil, errors.New("Error converting Patient record") 
	}

	err = stub.PutState(PREFIX_PATIENT + args[0], bytes)
	if err != nil { 
		return nil, errors.New("Error storing Patient record") 
	}

	return nil, nil
}
func (t *SimpleChaincode) enter_patient_prescription(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//		0		 1		   2		  3
	//	       Name  	      Disease 	      Medication       Duration

	if len(args) != 4 { 
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	bytes, err := stub.GetState(PREFIX_PATIENT + args[0])
	if err != nil { 
		return nil, errors.New("No Patient with name " + args[0])
	}


	prescription := Prescription{}
	prescription.Disease = args[1]
	prescription.Medication = args[2]
	prescription.Duration = args[3]

	var patient Patient
 	err = json.Unmarshal(bytes,&patient)
 	patient.Prescriptions = append(patient.Prescriptions, prescription)

	bytes, err = json.Marshal(&patient)
	if err != nil { 
		return nil, errors.New("Error converting Patient record") 
	}

	err = stub.PutState(PREFIX_PATIENT + args[0], bytes)
	if err != nil { 
		return nil, errors.New("Error storing Patient record") 
	}

	return nil, nil
}
