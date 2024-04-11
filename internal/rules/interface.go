package rules

import(
	"io"
	"github.com/aws-sdk-go/service/dynamodb"
)

type Interface interface{
	ConvertIOReaderToStruct(data io.Reader, model interface{})(body interface{}, err error)
	GetMock() interface{}
	Migrate(connection *dynamodb.DynamoDB) error
	Validate(model interface{}) error 

}