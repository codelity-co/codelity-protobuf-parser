package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

)

type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) TestImport() {
	t := suite.T()
	assert.NotNil(t, Message)
	assert.NotNil(t, Parser)
	assert.NotNil(t, Service)
}

func TestParserTestSuitet(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}