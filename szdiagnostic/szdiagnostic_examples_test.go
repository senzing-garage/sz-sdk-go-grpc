//go:build linux

package szdiagnostic_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

const AllLines = -1

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_CheckRepositoryPerformance() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}

	secondsToRun := 1

	result, err := szDiagnostic.CheckRepositoryPerformance(ctx, secondsToRun)
	if err != nil {
		handleError(err)
	}

	redactKeys := []string{"numRecordsInserted"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(result, AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//     "insertTime": 1000
	// }
}

func ExampleSzdiagnostic_GetRepositoryInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}

	result, err := szDiagnostic.GetRepositoryInfo(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "dataStores": [
	//         {
	//             "id": "CORE",
	//             "type": "sqlite3",
	//             "location": "nowhere"
	//         }
	//     ]
	// }
}

func ExampleSzdiagnostic_GetFeature() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}

	featureID := int64(1)

	result, err := szDiagnostic.GetFeature(ctx, featureID)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "LIB_FEAT_ID": 1,
	//     "FTYPE_CODE": "NAME",
	//     "ELEMENTS": [
	//         {
	//             "FELEM_CODE": "FULL_NAME",
	//             "FELEM_VALUE": "Robert Smith"
	//         },
	//         {
	//             "FELEM_CODE": "SUR_NAME",
	//             "FELEM_VALUE": "Smith"
	//         },
	//         {
	//             "FELEM_CODE": "GIVEN_NAME",
	//             "FELEM_VALUE": "Robert"
	//         },
	//         {
	//             "FELEM_CODE": "CULTURE",
	//             "FELEM_VALUE": "ANGLO"
	//         },
	//         {
	//             "FELEM_CODE": "CATEGORY",
	//             "FELEM_VALUE": "PERSON"
	//         },
	//         {
	//             "FELEM_CODE": "TOKENIZED_NM",
	//             "FELEM_VALUE": "ROBERT|SMITH"
	//         }
	//     ]
	// }
}

func ExampleSzdiagnostic_PurgeRepository() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}

	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)

	err := szDiagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzdiagnostic_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzdiagnostic_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	result := szDiagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}
