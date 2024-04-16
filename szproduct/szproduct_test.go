package szproduct

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szProductSingleton *Szproduct
	grpcAddress        = "localhost:8261"
	grpcConnection     *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzProduct_GetLicense(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetLicense(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzProduct_GetVersion(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetVersion(ctx)
	testError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzProduct_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
}

func TestSzProduct_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzProduct_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzProduct_Initialize(test *testing.T) {
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	instanceName := "Test name"
	settings, err := getSettings()
	testError(test, err)
	verboseLogging := sz.SZ_NO_LOGGING
	err = szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	testError(test, err)
}

func TestSzProduct_Destroy(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	testError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getSettings() (string, error) {
	return "{}", nil
}

func getSzProduct(ctx context.Context) *Szproduct {
	_ = ctx
	if szProductSingleton == nil {
		grpcConnection := getGrpcConnection()
		szProductSingleton = &Szproduct{
			GrpcClient: szpb.NewSzProductClient(grpcConnection),
		}
	}
	return szProductSingleton
}

func getSzProductAsInterface(ctx context.Context) sz.SzProduct {
	return getSzProduct(ctx)
}

func getTestObject(ctx context.Context, test *testing.T) *Szproduct {
	_ = test
	return getSzProduct(ctx)
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
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
