package env

import "os"

func GetEnv(env, defaultValue string) string{
	environment := os.GetEnv(env)
	if environment == ""{
		return defaultValue
	}
	return environment
}