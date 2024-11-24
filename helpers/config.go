package helpers

import "github.com/joho/godotenv"

var Env = map[string]string{}

func SetupConfig() {
	var err error
	Env, err = godotenv.Read(".env")
	if err != nil {
		Logger.Fatal("failed to read env file: ", err)
	}
}

func GetEnv(key string, val string) string {
	// val is used as default value
	result := Env[key]
	if result == "" {
		result = val
	}
	return result
}
