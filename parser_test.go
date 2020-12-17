package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"fmt"
)

type ProtobufParserTestSuite struct {
	suite.Suite
	testProtoContent string
}

func (suite *ProtobufParserTestSuite) SetupTest() {
	suite.testProtoContent = `
	syntax = "proto3";
	package sample;

	option go_package ".;main";

	import "google/api/annotations.proto";

	service SampleService {
		rpc SampleAction(SampleRequest) returns (SampleResponse) {
			option (google.api.http) = {
				post: "/v1/SampleService/SampleAction";
			};
		};
	}

	message SampleRequest {
		dummyVarialbe string = 1;
	}

	message SampleResponse {
		message ResponseInfo {
			string responseType = 1;
			int32 httpStatusCode = 2;
			string httpStatusText = 3;
			string applicationMessageType = 4;
			string applicationMessageCode = 5;
			string applicationMessageText = 6;
		}
		map<string, google.protobuf.Struct> sampleInfo = 1;
		ResponseInfo responseInfo = 1;
	}
	`
}

func (suite *ProtobufParserTestSuite) TestReadFile() {
	t := suite.T()

	t.Skip("Pending")
}

func (suite *ProtobufParserTestSuite) TestGetSyntax() {
	t := suite.T()

	parser := &Parser{syntax: "proto1"}
	assert.NotContains(t, []string{"proto2", "proto3"}, parser.GetSyntax(), `syntax should be "proto1"`)
	parser = &Parser{syntax: "proto3"}
	assert.Contains(t, []string{"proto2", "proto3"}, parser.GetSyntax(), `Cannot obtain syntax "proto2" or "proto3"`)
}

func (suite *ProtobufParserTestSuite) TestGetPackageName() {
	t := suite.T()

	parser := &Parser{packageName: "sample"}
	assert.Equal(t, "sample", parser.GetPackageName(), `Package name should be "sample"`)
}

func (suite *ProtobufParserTestSuite) TestGetOptions() {
	t := suite.T()

	parser := &Parser{options: map[string]string{"test": "abc"}}
	assert.Equal(t, map[string]string{"test": "abc"}, parser.GetOptions(), `Invalid options`)
}

func (suite *ProtobufParserTestSuite) TestImports() {
	t := suite.T()

	parser := &Parser{
		imports: []string{
			"google/api/annotations.proto",
		},
	}
	assert.Contains(t, parser.GetImports(), "google/api/annotations.proto", "Missing import in parser")
}

func (suite *ProtobufParserTestSuite) TestServices() {
	t := suite.T()

	t.Skip("Pending")
}

func (suite *ProtobufParserTestSuite) TestMessage() {
	t := suite.T()

	t.Skip("Pending")
}

func (suite *ProtobufParserTestSuite) TestProcessSyntaxLines() {
	var err error 
	t := suite.T()

	parser := &Parser{}

	err = parser.processSyntaxLines([]string{})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax proto3`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = proto3`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = "proto3"`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = 'proto3'`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = 'proto2'`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = "proto2";`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = "proto2";`, `sytnax = "proto3"`})
	assert.NotNil(t, err)

	err = parser.processSyntaxLines([]string{`syntax = "proto3";`})
	assert.Nil(t, err)
	assert.Equal(t, "proto3", parser.GetSyntax(), `Sytnax is not "proto3"`)

}

func (suite *ProtobufParserTestSuite) TestProcessPackageLines() {
	var err error 
	t := suite.T()

	parser := &Parser{}

	err = parser.processPackageLines([]string{})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package`})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package;`})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package "abc"`})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package "abc" ;`})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package abc`})
	assert.NotNil(t, err)

	err = parser.processPackageLines([]string{`package abc ;`})
	assert.Nil(t, err)
	assert.Equal(t, "abc", parser.GetPackageName(), `Package is not "abc"`)
}

func (suite *ProtobufParserTestSuite) TestProcessImportLines() {
	var err error 
	t := suite.T()

	parser := &Parser{}

	err = parser.processImportLines([]string{})
	assert.Nil(t, err)

	err = parser.processImportLines([]string{`import`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import ""`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import ''`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import ;`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import "";`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import '';`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import google/api/annotations.proto;`})
	assert.NotNil(t, err)

	err = parser.processImportLines([]string{`import "google/api/annotations.proto";`})
	assert.Nil(t, err)
	assert.Contains(t, parser.GetImports(), "google/api/annotations.proto", "Invalid import")

	err = parser.processImportLines([]string{`import 'google/api/annotations.proto';`})
	assert.Nil(t, err)
	assert.Contains(t, parser.GetImports(), "google/api/annotations.proto", "Invalid import")
}

func (suite *ProtobufParserTestSuite) TestProcessOptionLines() {
  var err error 
	t := suite.T()

	parser := &Parser{}

	err = parser.processOptionLines([]string{})
	assert.Nil(t, err)

	err = parser.processOptionLines([]string{`option`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option ""`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option ''`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option ;`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option "";`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option '';`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option abc;`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option "abc";`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option 'abc';`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option 'abc';`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option "abc" "123";`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option "abc" '123';`})
	assert.NotNil(t, err)

	err = parser.processOptionLines([]string{`option abc "123";`})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"abc": "123"}, parser.GetOptions(), "Invalid option statements")

	err = parser.processOptionLines([]string{`option abc '123';`})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"abc": "123"}, parser.GetOptions(), "Invalid option statements")
}

func (suite *ProtobufParserTestSuite) TestProcessServiceLines() {
  var err error 
	t := suite.T()

	parser := &Parser{}

	err = parser.processServiceLines([]string{`service;`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service {};`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service{}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service {}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {};`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `};`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde;`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde();`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg);`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg)fgh;`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) fgh;`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns (fgh) {`, `option (x) = {`, `xyz: "cde";`, `}`,`}`, `}`})
	assert.NotNil(t, err)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns (fgh) {};`, `}`})
	assert.Nil(t, err)
	services := parser.GetServices()
	assert.Equal(t, 1, len(services), fmt.Sprintf(`Not return single service: %v`, services))
	assert.Equal(t, "abc", services[0].GetServiceName(), `Cannot get service name "abc"`)
	rpcs := services[0].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "cde", rpcs[0].GetRpcName(), `Cannot get rpc "cde"`)
	assert.Equal(t, "efg", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "efg"`)
	assert.Equal(t, "fgh", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "fgh"`)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns fgh;`, `}`})
	assert.Nil(t, err)
	services = parser.GetServices()
	assert.Equal(t, 1, len(services), `Not return single service`)
	assert.Equal(t, "abc", services[0].GetServiceName(), `Cannot get service name "abc"`)
	rpcs = services[0].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "cde", rpcs[0].GetRpcName(), `Cannot get rpc "cde"`)
	assert.Equal(t, "efg", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "efg"`)
	assert.Equal(t, "fgh", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "fgh"`)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns (fgh);`, `}`})
	assert.Nil(t, err)
	services = parser.GetServices()
	assert.Equal(t, 1, len(services), `Not return single service`)
	assert.Equal(t, "abc", services[0].GetServiceName(), `Cannot get service name "abc"`)
	rpcs = services[0].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "cde", rpcs[0].GetRpcName(), `Cannot get rpc "cde"`)
	assert.Equal(t, "efg", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "efg"`)
	assert.Equal(t, "fgh", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "fgh"`)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns (fgh);`, `}`, `service ghi {`, `rpc hij(ijk) returns jkl;`, `}`})
	assert.Nil(t, err)
	services = parser.GetServices()
	assert.Equal(t, 2, len(services), `Not return single service`)
	assert.Equal(t, "abc", services[0].GetServiceName(), `Cannot get service name "abc"`)
	assert.Equal(t, "ghi", services[1].GetServiceName(), `Cannot get service name "ghi"`)
	rpcs = services[0].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "cde", rpcs[0].GetRpcName(), `Cannot get rpc "cde"`)
	assert.Equal(t, "efg", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "efg"`)
	assert.Equal(t, "fgh", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "fgh"`)
	rpcs = services[1].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "hij", rpcs[0].GetRpcName(), `Cannot get rpc "hij"`)
	assert.Equal(t, "ijk", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "ijk"`)
	assert.Equal(t, "jkl", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "jkl"`)

	err = parser.processServiceLines([]string{`service abc {`, `rpc cde(efg) returns (fgh) {`, `option (x) = {`, `xyz: "cde";`, `}`,`};`, `}`})
	assert.Nil(t, err)
	services = parser.GetServices()
	assert.Equal(t, 1, len(services), `Not return single service`)
	assert.Equal(t, "abc", services[0].GetServiceName(), `Cannot get service name "abc"`)
	rpcs = services[0].GetRpcs()
	assert.Equal(t, 1, len(rpcs), "Not return single rpc")
	assert.Equal(t, "cde", rpcs[0].GetRpcName(), `Cannot get rpc "cde"`)
	assert.Equal(t, "efg", rpcs[0].GetRpcRequestName(), `Cannot get rpc request "efg"`)
	assert.Equal(t, "fgh", rpcs[0].GetRpcResponseName(), `Cannot get rpc response "fgh"`)
}

func (suite *ProtobufParserTestSuite) TestProcessMessageLines() {

}


func (suite *ProtobufParserTestSuite) TestParse() {
	t := suite.T()

	parser := &Parser{}
	err := parser.Parse(suite.testProtoContent)

	assert.Nil(t, err)
	assert.Contains(t, []string{"proto3"}, parser.GetSyntax(), `Cannot obtain syntax "proto2" or "proto3"`)
	assert.Equal(t, "sample", parser.GetPackageName(), `Package name should be "sample"`)
	assert.Contains(t, parser.GetImports(), "google/api/annotations.proto", "Missing import in parser")
	assert.Equal(t, map[string]string{"go_package": ".;main"}, parser.GetOptions(), "Options parsing error")
}

func TestProtobufParserTestSuite(t *testing.T) {
	suite.Run(t, new(ProtobufParserTestSuite))
}