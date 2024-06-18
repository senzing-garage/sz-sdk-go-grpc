package helper

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/senzing-garage/sz-sdk-go/szerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testCases = []struct {
	expectedType        error
	expectedTypes       []error
	falseTypes          []error
	gRpcCode            codes.Code
	senzingErrorMessage string
	name                string
}{
	{
		name:          "helper-szerror-0023",
		expectedType:  szerror.ErrSzBadInput,
		expectedTypes: []error{szerror.ErrSzBadInput},
		falseTypes:    []error{szerror.ErrSzRetryable},
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
	_ = test
	_ = testCases
	// TODO: Reinstate TestConvertGrpcError
	// for _, testCase := range testCases {
	// 	test.Run(testCase.name, func(test *testing.T) {
	// 		originalError := status.Error(testCase.gRpcCode, testCase.senzingErrorMessage)
	// 		actual := ConvertGrpcError(originalError)
	// 		assert.NotNil(test, actual)
	// 		assert.IsType(test, testCase.expectedType, actual)
	// 		for _, szerrorTypeID := range testCase.expectedTypes {
	// 			require.ErrorIs(test, actual, szerrorTypeID)
	// 		}
	// 		for _, szerrorTypeID := range testCase.falseTypes {
	// 			assert.False(test, errors.Is(actual, szerrorTypeID), szerrorTypeID)
	// 		}
	// 	})
	// }

}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleConvertGrpcError() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/helper/helper_test.go
	senzingErrorMessage := `{"reason": "SENZ0033|Test message"}`        // Example message from Senzing G2 engine.
	grpcStatusError := status.Error(codes.Unknown, senzingErrorMessage) // Create a gRPC *status.Error
	err := ConvertGrpcError(grpcStatusError)
	if err != nil {
		if errors.Is(err, szerror.ErrSzNotFound) {
			fmt.Println("Is an ErrSzNotFound")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Println("Is an ErrSzBadInput")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Println("Is an ErrSzRetryable.")
		}
	}
	// Output:
	// Is an ErrSzNotFound
	// Is an ErrSzBadInput
}
