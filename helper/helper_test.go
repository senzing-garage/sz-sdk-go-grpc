package helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testCases = []struct {
	expectedType        error
	expectedTypes       []szerror.SzErrorTypeIds
	falseTypes          []szerror.SzErrorTypeIds
	gRpcCode            codes.Code
	senzingErrorMessage string
	name                string
}{
	{
		name:          "helper-szerror-0023",
		expectedType:  szerror.SzBadInputError{},
		expectedTypes: []szerror.SzErrorTypeIds{szerror.SzBadInput},
		falseTypes:    []szerror.SzErrorTypeIds{szerror.SzRetryable},
		gRpcCode:      codes.Unknown,
		senzingErrorMessage: `{
			"date": "2023-03-27",
			"time": "20:34:11.451202917",
			"level": "ERROR",
			"id": "senzing-60044001",
			"text": "Call to G2_addRecord(CUSTOMERS, 1002, {\"DATA_SOURCE\": \"BOB\", \"RECORD_ID\": \"1002\", \"RECORD_TYPE\": \"PERSON\", \"PRIMARY_NAME_LAST\": \"Smith\", \"PRIMARY_NAME_FIRST\": \"Bob\", \"DATE_OF_BIRTH\": \"11/12/1978\", \"ADDR_TYPE\": \"HOME\", \"ADDR_LINE1\": \"1515 Adela Lane\", \"ADDR_CITY\": \"Las Vegas\", \"ADDR_STATE\": \"NV\", \"ADDR_POSTAL_CODE\": \"89111\", \"PHONE_TYPE\": \"MOBILE\", \"PHONE_NUMBER\": \"702-919-1300\", \"DATE\": \"3/10/17\", \"STATUS\": \"Inactive\", \"AMOUNT\": \"200\"}, G2Engine_test) failed. Return code: -2",
			"duration": 518591,
			"location": "In AddRecord() at g2engineserver.go:66",
			"errors": [{
				"text": "0023E|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'"
			}],
			"details": {
				"1": "CUSTOMERS",
				"2": 1002,
				"3": {
					"DATA_SOURCE": "BOB",
					"RECORD_ID": "1002",
					"RECORD_TYPE": "PERSON",
					"PRIMARY_NAME_LAST": "Smith",
					"PRIMARY_NAME_FIRST": "Bob",
					"DATE_OF_BIRTH": "11/12/1978",
					"ADDR_TYPE": "HOME",
					"ADDR_LINE1": "1515 Adela Lane",
					"ADDR_CITY": "Las Vegas",
					"ADDR_STATE": "NV",
					"ADDR_POSTAL_CODE": "89111",
					"PHONE_TYPE": "MOBILE",
					"PHONE_NUMBER": "702-919-1300",
					"DATE": "3/10/17",
					"STATUS": "Inactive",
					"AMOUNT": "200"
				},
				"4": "G2Engine_test",
				"5": -2,
				"6": 518591
			}
		}`,
	},
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestConvertGrpcError(test *testing.T) {

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			originalError := status.Error(testCase.gRpcCode, testCase.senzingErrorMessage)
			actual := ConvertGrpcError(originalError)
			assert.NotNil(test, actual)
			assert.IsType(test, testCase.expectedType, actual)
			for _, szerrorTypeId := range testCase.expectedTypes {
				assert.True(test, szerror.Is(actual, szerrorTypeId), szerrorTypeId)
			}
			for _, szerrorTypeId := range testCase.falseTypes {
				assert.False(test, szerror.Is(actual, szerrorTypeId), szerrorTypeId)
			}
		})
	}

}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleConvertGrpcError() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/helper/helper_test.go
	senzingErrorMessage := "27E|Test message"                           // Example message from Senzing G2 engine.
	grpcStatusError := status.Error(codes.Unknown, senzingErrorMessage) // Create a gRPC *status.Error
	err := ConvertGrpcError(grpcStatusError)
	if err != nil {
		if szerror.Is(err, szerror.SzBadInput) {
			fmt.Println("Is a SzBadInputError")
		}
		if szerror.Is(err, szerror.SzUnknownDatasource) {
			fmt.Println("Is a SzUnknownDatasourceError")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Println("Is a SzRetryableError.")
		}
	}
	// Output:
	// Is a SzBadInputError
	// Is a SzUnknownDatasourceError
}
