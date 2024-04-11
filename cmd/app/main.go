package main 

import(
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"net/http"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/repository/adapter"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/repository/instance"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/routes"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/rules"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/internal/rules/product"
	"github.com/akuppa9/Golang-DynamoDB-CRUD-API/utils/logger"
)

func main(){
	configs := Config.GetConfig()
	connection := instance.GetConnection()
	repository := adapter.NewAdapter(connection)
	logger.INFO("waiting for the service to start...", nil)
	errors := Migrate(connection)
	if len(errors) > 0{
		for _, err := range errors{
			logger.PANIC("error on migration:...", err)
		}
	}
	logger.PANIC("", checkTables(connection))
	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("service is running on port ", port)
	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error{
	var errors []error
	callMigrateAndAppendError(&errors, connection, &RulesProduct.Rules{})
	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, rule rules.Interface){
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTables(&dynamodb.ListTablesInput{}))
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found: ", nil)
		}
		for _, tableName := range response.TableNames{
			logger.INFO("Table found: ", *tableName)
		}
	}
	return err
}