package helper_test

import (
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultTruncation = 76
	printErrors       = false
	printResults      = false
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestConvertGrpcError(test *testing.T) {
	testCases := getTestCasesForConvertGrpcError()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			originalError := status.Error(testCase.gRPCCode, testCase.expectedErrMessage)
			err := helper.ConvertGrpcError(originalError)
			printDebug(test, err)

			require.ErrorIs(test, err, testCase.expectedErr)
			require.JSONEq(test, testCase.expectedErrMessage, err.Error())

			for _, szerrorTypeID := range testCase.acceptableErrs {
				require.ErrorIs(test, err, szerrorTypeID)
			}

			for _, szerrorTypeID := range testCase.unacceptableErrs {
				require.NotErrorIs(test, err, szerrorTypeID)
			}
		})
	}
}

func TestConvertGrpcErrorAnomolies(test *testing.T) {
	testCases := getTestCasesForConvertGrpcErrorAnomolies()
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			gRPCError := status.Error(codes.Unknown, testCase.originalErrMessage)
			err := helper.ConvertGrpcError(gRPCError)
			printDebug(test, err)
			require.Error(test, err)
			require.Equal(test, testCase.expectedErrMessage, err.Error())
		})
	}
}

func TestConvertGrpcError_nil(test *testing.T) {
	actual := helper.ConvertGrpcError(nil)
	require.NoError(test, actual)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		}
	}

	if printResults {
		for _, item := range items {
			outLine := truncator.Truncate(fmt.Sprintf("%v", item), defaultTruncation, "...", truncator.PositionEnd)
			t.Logf("Result: %s\n", outLine)
		}
	}
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

type TestMetadataForConvertGrpcError struct {
	acceptableErrs     []error
	expectedErr        error
	expectedErrMessage string
	gRPCCode           codes.Code
	name               string
	unacceptableErrs   []error
}

type TestMetadataForConvertGrpcErrorAnomolies struct {
	expectedErrMessage string
	originalErrMessage string
	name               string
}

func getTestCasesForConvertGrpcError() []TestMetadataForConvertGrpcError {
	result := []TestMetadataForConvertGrpcError{
		{
			name:           "szerror-0023",
			acceptableErrs: []error{szerror.ErrSzBadInput, szerror.ErrSz},
			expectedErr:    szerror.ErrSzBadInput,
			expectedErrMessage: `{
			"time": "2023-03-27T20:34:11.451202917Z",
			"level": "ERROR",
			"id": "senzing-60044001",
			"text": "Call to Sz_addRecord(CUSTOMERS, 1002, {\"DATA_SOURCE\": \"BOB\", \"RECORD_ID\": \"1002\", \"RECORD_TYPE\": \"PERSON\", \"PRIMARY_NAME_LAST\": \"Smith\", \"PRIMARY_NAME_FIRST\": \"Bob\", \"DATE_OF_BIRTH\": \"11/12/1978\", \"ADDR_TYPE\": \"HOME\", \"ADDR_LINE1\": \"1515 Adela Lane\", \"ADDR_CITY\": \"Las Vegas\", \"ADDR_STATE\": \"NV\", \"ADDR_POSTAL_CODE\": \"89111\", \"PHONE_TYPE\": \"MOBILE\", \"PHONE_NUMBER\": \"702-919-1300\", \"DATE\": \"3/10/17\", \"STATUS\": \"Inactive\", \"AMOUNT\": \"200\"}, G2Engine_test) failed. Return code: -2",
			"reason": "SENZ0023E|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'",
			"duration": 518591,
			"location": "In AddRecord() at szengineserver.go:66"}`,
			gRPCCode:         codes.Unknown,
			unacceptableErrs: []error{szerror.ErrSzRetryable},
		},
		{
			name:               "szerror-0037",
			acceptableErrs:     []error{szerror.ErrSzNotFound, szerror.ErrSzBadInput, szerror.ErrSz},
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyEntities","error":{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			gRPCCode:           codes.Unknown,
			unacceptableErrs:   []error{szerror.ErrSzRetryable},
		},
	}

	return result
}

func getTestCasesForConvertGrpcErrorAnomolies() []TestMetadataForConvertGrpcErrorAnomolies {
	result := []TestMetadataForConvertGrpcErrorAnomolies{
		{
			name:               "badParse",
			expectedErrMessage: `{"function": "helper.createErrorFromReason", "text": "errorMessage: {\"time\": 12345}; reason: {\"time\":12345}", "error": "strconv.Atoi: parsing \"me\\\":\": invalid syntax"}`,
			originalErrMessage: `{"time": 12345}`,
		},
		{
			name:               "badReason",
			expectedErrMessage: `{"function": "helper.createErrorFromReason", "text": "errorMessage: {\"reason\": \"bad\"}; reason: bad", "error": "helper"}`,
			originalErrMessage: `{"reason": "bad"}`,
		},
		{
			name:               "badReasonCode",
			expectedErrMessage: `{"function": "helper.createErrorFromReason", "text": "errorMessage: {\"reason\": \"SENZabcd | bad text\"}; reason: SENZabcd | bad text", "error": "strconv.Atoi: parsing \"abcd\": invalid syntax"}`,
			originalErrMessage: `{"reason": "SENZabcd | bad text"}`,
		},
	}

	return result
}
