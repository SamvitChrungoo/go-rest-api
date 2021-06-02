package utils

import "os"

//SetEnvironmentVariables -> for setting ENV variables
func SetEnvironmentVariables() {
	os.Setenv("JWT_SECRET_PHASE", "super_secret_phrase_key")
	os.Setenv("PORT_NO", "8080")
	os.Setenv("DB_PASSWORD", "samvit123")
	os.Setenv("DB_NAME", "test_db")

}

//GetEnvironmentVariable -> Get variable based on key
func GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}
