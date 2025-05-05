package helper_test

import (
	"fmt"
	"testing"

	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testCases = []struct {
	expectedType        error
	expectedTypes       []error
	falseTypes          []error
	gRPCCode            codes.Code
	senzingErrorMessage string
	name                string
}{
	{
		name:          "helper-szerror-0023",
		expectedType:  szerror.ErrSzBadInput,
		expectedTypes: []error{szerror.ErrSzBadInput},
		falseTypes:    []error{szerror.ErrSzRetryable},
		gRPCCode:      codes.Unknown,
		senzingErrorMessage: `{
			"time": "2023-03-27T20:34:11.451202917Z",
			"level": "ERROR",
			"id": "senzing-60044001",
			"text": "Call to Sz_addRecord(CUSTOMERS, 1002, {\"DATA_SOURCE\": \"BOB\", \"RECORD_ID\": \"1002\", \"RECORD_TYPE\": \"PERSON\", \"PRIMARY_NAME_LAST\": \"Smith\", \"PRIMARY_NAME_FIRST\": \"Bob\", \"DATE_OF_BIRTH\": \"11/12/1978\", \"ADDR_TYPE\": \"HOME\", \"ADDR_LINE1\": \"1515 Adela Lane\", \"ADDR_CITY\": \"Las Vegas\", \"ADDR_STATE\": \"NV\", \"ADDR_POSTAL_CODE\": \"89111\", \"PHONE_TYPE\": \"MOBILE\", \"PHONE_NUMBER\": \"702-919-1300\", \"DATE\": \"3/10/17\", \"STATUS\": \"Inactive\", \"AMOUNT\": \"200\"}, G2Engine_test) failed. Return code: -2",
			"reason": "SENZ0023E|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'",
			"duration": 518591,
			"location": "In AddRecord() at szengineserver.go:66"
		}`,
	},
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestConvertGrpcError(test *testing.T) {
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			originalError := status.Error(testCase.gRPCCode, testCase.senzingErrorMessage)
			actual := helper.ConvertGrpcError(originalError)
			require.ErrorIs(test, actual, testCase.expectedType)

			for _, szerrorTypeID := range testCase.expectedTypes {
				require.ErrorIs(test, actual, szerrorTypeID)
			}

			for _, szerrorTypeID := range testCase.falseTypes {
				assert.NotErrorIs(test, actual, szerrorTypeID)
			}
		})
	}
}

func TestConvertGrpcError_wrapped(test *testing.T) {
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			originalError := status.Error(testCase.gRPCCode, testCase.senzingErrorMessage)
			wrappedError := fmt.Errorf("Wrap %w", originalError)
			actual := helper.ConvertGrpcError(wrappedError)
			require.ErrorIs(test, actual, testCase.expectedType)

			for _, szerrorTypeID := range testCase.expectedTypes {
				require.ErrorIs(test, actual, szerrorTypeID)
			}

			for _, szerrorTypeID := range testCase.falseTypes {
				assert.NotErrorIs(test, actual, szerrorTypeID)
			}
		})
	}
}

func TestConvertGrpcError_nil(test *testing.T) {
	actual := helper.ConvertGrpcError(nil)
	require.NoError(test, actual)
}

func TestConvertGrpcError_badParse(test *testing.T) {
	jsonMessage := `{"time": 12345}`
	gRPCError := status.Error(codes.Unknown, jsonMessage)
	actual := helper.ConvertGrpcError(gRPCError)
	require.Error(test, actual)
}

func TestConvertGrpcError_badReason(test *testing.T) {
	jsonMessage := `{"reason": "bad"}`
	gRPCError := status.Error(codes.Unknown, jsonMessage)
	actual := helper.ConvertGrpcError(gRPCError)
	require.Error(test, actual)
}

func TestConvertGrpcError_badReasonCode(test *testing.T) {
	jsonMessage := `{"reason": "SENZabcd | bad text"}`
	gRPCError := status.Error(codes.Unknown, jsonMessage)
	actual := helper.ConvertGrpcError(gRPCError)
	require.Error(test, actual)
}
