package g2productclient

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	truncator "github.com/aquilax/truncate"
	pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
)

var (
	grpcAddress              = "localhost:8258"
	grpcConnection           *grpc.ClientConn
	g2productClientSingleton *G2productClient
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

func getTestObject(ctx context.Context, test *testing.T) G2productClient {
	if g2productClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2productClientSingleton = &G2productClient{
			GrpcClient: pb.NewG2ProductClient(grpcConnection),
		}
	}
	return *g2productClientSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if 1 == 0 {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2product G2productClient, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2product G2productClient, err error, messageId string) {
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2product G2productClient, err error) {
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

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2productImpl_Init(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2product, jsonErr)
	err := g2product.Init(ctx, moduleName, iniParams, verboseLogging)
	expectError(test, ctx, g2product, err, "senzing-60164002")
}

func TestG2productImpl_License(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.License(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productImpl_ValidateLicenseFile(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	licenseFilePath := "/etc/opt/senzing/g2.lic"
	actual, err := g2product.ValidateLicenseFile(ctx, licenseFilePath)
	testErrorNoFail(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productImpl_ValidateLicenseStringBase64(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	licenseString := "AQAAADgCAAAAAAAAU2VuemluZyBQdWJsaWMgVGVzdCBMaWNlbnNlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTFVBVElPTiAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjItMTEtMjkAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFDDAAAAAAAAMjAyMy0xMS0yOQAAAAAAAAAAAABNT05USExZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACQfw5e19QAHetkvd+vk0cYHtLaQCLmgx2WUfLorDfLQq15UXmOawNIXc1XguPd8zJtnOaeI6CB2smxVaj10mJE2ndGPZ1JjGk9likrdAj3rw+h6+C/Lyzx/52U8AuaN1kWgErDKdNE9qL6AnnN5LLi7Xs87opP7wbVMOdzsfXx2Xi3H7dSDIam7FitF6brSFoBFtIJac/V/Zc3b8jL/a1o5b1eImQldaYcT4jFrRZkdiVO/SiuLslEb8or3alzT0XsoUJnfQWmh0BjehBK9W74jGw859v/L1SGn1zBYKQ4m8JBiUOytmc9ekLbUKjIg/sCdmGMIYLywKqxb9mZo2TLZBNOpYWVwfaD/6O57jSixfJEHcLx30RPd9PKRO0Nm+4nPdOMMLmd4aAcGPtGMpI6ldTiK9hQyUfrvc9z4gYE3dWhz2Qu3mZFpaAEuZLlKtxaqEtVLWIfKGxwxPargPEfcLsv+30fdjSy8QaHeU638tj67I0uCEgnn5aB8pqZYxLxJx67hvVKOVsnbXQRTSZ00QGX1yTA+fNygqZ5W65wZShhICq5Fz8wPUeSbF7oCcE5VhFfDnSyi5v0YTNlYbF8LOAqXPTi+0KP11Wo24PjLsqYCBVvmOg9ohZ89iOoINwUB32G8VucRfgKKhpXhom47jObq4kSnihxRbTwJRx4o"
	actual, err := g2product.ValidateLicenseStringBase64(ctx, licenseString)
	testErrorNoFail(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productImpl_Version(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.Version(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productImpl_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	err := g2product.Destroy(ctx)
	expectError(test, ctx, g2product, err, "senzing-60164001")
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2productClient_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2product.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60164002" error.
	}
	// Output:
}

func ExampleG2productClient_License() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	result, err := g2product.License(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"customer":"Senzing Public Test License","contract":"EVALUATION - support@senzing.com","issueDate":"2022-11-29","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"MONTHLY","expireDate":"2023-11-29","recordLimit":50000}
}

func ExampleG2productClient_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	err := g2product.SetLogLevel(ctx, logger.LevelInfo)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2productClient_ValidateLicenseFile() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	licenseFilePath := "/etc/opt/senzing/g2.lic"
	result, err := g2product.ValidateLicenseFile(ctx, licenseFilePath)
	if err != nil {
		fmt.Println("Invalid license")
	} else {
		fmt.Println(result)
	}
	// Output: Success
}

func ExampleG2productClient_ValidateLicenseStringBase64() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	licenseString := "AQAAADgCAAAAAAAAU2VuemluZyBQdWJsaWMgVGVzdCBMaWNlbnNlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTFVBVElPTiAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjItMTEtMjkAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFDDAAAAAAAAMjAyMy0xMS0yOQAAAAAAAAAAAABNT05USExZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACQfw5e19QAHetkvd+vk0cYHtLaQCLmgx2WUfLorDfLQq15UXmOawNIXc1XguPd8zJtnOaeI6CB2smxVaj10mJE2ndGPZ1JjGk9likrdAj3rw+h6+C/Lyzx/52U8AuaN1kWgErDKdNE9qL6AnnN5LLi7Xs87opP7wbVMOdzsfXx2Xi3H7dSDIam7FitF6brSFoBFtIJac/V/Zc3b8jL/a1o5b1eImQldaYcT4jFrRZkdiVO/SiuLslEb8or3alzT0XsoUJnfQWmh0BjehBK9W74jGw859v/L1SGn1zBYKQ4m8JBiUOytmc9ekLbUKjIg/sCdmGMIYLywKqxb9mZo2TLZBNOpYWVwfaD/6O57jSixfJEHcLx30RPd9PKRO0Nm+4nPdOMMLmd4aAcGPtGMpI6ldTiK9hQyUfrvc9z4gYE3dWhz2Qu3mZFpaAEuZLlKtxaqEtVLWIfKGxwxPargPEfcLsv+30fdjSy8QaHeU638tj67I0uCEgnn5aB8pqZYxLxJx67hvVKOVsnbXQRTSZ00QGX1yTA+fNygqZ5W65wZShhICq5Fz8wPUeSbF7oCcE5VhFfDnSyi5v0YTNlYbF8LOAqXPTi+0KP11Wo24PjLsqYCBVvmOg9ohZ89iOoINwUB32G8VucRfgKKhpXhom47jObq4kSnihxRbTwJRx4o"
	result, err := g2product.ValidateLicenseStringBase64(ctx, licenseString)
	if err != nil {
		fmt.Println("Invalid license")
	} else {
		fmt.Println(result)
	}
	// Output: Success
}

func ExampleG2productClient_Version() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	result, err := g2product.Version(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":"3.4.0","BUILD_VERSION":"3.4.0.23005","BUILD_DATE":"2023-01-04","BUILD_NUMBER":"2023_01_04__23_02","COMPATIBILITY_VERSION":{"CONFIG_VERSION":"10"},"SCHEMA_VERSION":{"ENGINE_SCHEMA_VERSION":"3.4","MINIMUM_REQUIRED_SCHEMA_VERSION":"3.0","MAXIMUM_REQUIRED_SCHEMA_VERSION":"3.99"}}
}

func ExampleG2productClient_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2productclient/g2productclient_test.go
	grpcConnection := getGrpcConnection()
	g2product := &G2productClient{
		GrpcClient: pb.NewG2ProductClient(grpcConnection),
	}
	ctx := context.TODO()
	err := g2product.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	// Output:
}
