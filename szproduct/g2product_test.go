package szproduct

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2product"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2productSingleton g2api.G2product
	grpcAddress        = "localhost:8261"
	grpcConnection     *grpc.ClientConn
)

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

func getTestObject(ctx context.Context, test *testing.T) g2api.G2product {
	if g2productSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2productSingleton = &G2product{
			GrpcClient: g2pb.NewG2ProductClient(grpcConnection),
		}
	}
	return g2productSingleton
}

func getG2Product(ctx context.Context) g2api.G2product {
	if g2productSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2productSingleton = &G2product{
			GrpcClient: g2pb.NewG2ProductClient(grpcConnection),
		}
	}
	return g2productSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2product g2api.G2product, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2product g2api.G2product, err error, messageId string) {
	if err != nil {
		errorMessage := err.Error()[strings.Index(err.Error(), "{"):]
		var dictionary map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(errorMessage), &dictionary)
		if unmarshalErr != nil {
			test.Log("Unmarshal Error:", unmarshalErr.Error())
		}
		assert.Equal(test, messageId, dictionary["id"].(string))
	} else {
		assert.FailNow(test, "Should have failed with", messageId)
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2product g2api.G2product, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
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

func TestG2product_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
}

func TestG2product_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	actual := g2product.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2product_Init(test *testing.T) {
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2product.Init(ctx, moduleName, iniParams, verboseLogging)
	expectError(test, ctx, g2product, err, "senzing-60164002")
}

func TestG2product_License(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.License(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2product_Version(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.Version(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2product_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	err := g2product.Destroy(ctx)
	expectError(test, ctx, g2product, err, "senzing-60164001")
}
