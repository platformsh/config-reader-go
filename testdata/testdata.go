package testdata

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "reflect"
  "runtime"
  "testing"
  psh "github.com/platformsh/config-reader-go"
)

// This function produces a getter of the same signature as os.Getenv() that
// always returns an empty string, simulating a non-Platform environment.
func NonPlatformEnv() func(string) string {
	return func(key string) string {
		return ""
	}
}

// This function produces a getter of the same signature as os.gGetenv()
// that returns test values to simulate a build environment.
func BuildEnv(env psh.EnvList) func(string) string {

  var (
    _, b, _, _ = runtime.Caller(0)
    basepath   = filepath.Dir(b)
  )

	// Create build time env.
	vars := LoadJsonFile(basepath + "/ENV.json")
	env = MergeMaps(vars, env)
	env["PLATFORM_VARIABLES"] = EncodeJsonFile(basepath + "/PLATFORM_VARIABLES.json")
	env["PLATFORM_APPLICATION"] = EncodeJsonFile(basepath + "/PLATFORM_APPLICATION.json")

	return func(key string) string {
		if val, ok := env[key]; ok {
			return val
		} else {
			return ""
		}
	}
}

// This function produces a getter of the same signature as os.gGetenv()
// that returns test values to simulate a runtime environment.
func RuntimeEnv(env psh.EnvList) func(string) string {

  var (
    _, b, _, _ = runtime.Caller(0)
    basepath   = filepath.Dir(b)
  )

	// Create runtimeVars env.
	vars := LoadJsonFile(basepath + "/ENV.json")
	env = MergeMaps(vars, env)
	env["PLATFORM_VARIABLES"] = EncodeJsonFile(basepath + "/PLATFORM_VARIABLES.json")
	env["PLATFORM_APPLICATION"] = EncodeJsonFile(basepath + "/PLATFORM_APPLICATION.json")
	env["PLATFORM_RELATIONSHIPS"] = EncodeJsonFile(basepath + "/PLATFORM_RELATIONSHIPS.json")
	env["PLATFORM_ROUTES"] = EncodeJsonFile(basepath + "/PLATFORM_ROUTES.json")

	vars = LoadJsonFile(basepath + "/ENV_runtime.json")
	env = MergeMaps(vars, env)

	return func(key string) string {
		if val, ok := env[key]; ok {
			return val
		} else {
			return ""
		}
	}
}

func GetKeys(data psh.EnvList) []string {
	keys := make([]string, 0)
	for key := range data {
		keys = append(keys, key)
	}

	return keys
}

func MergeMaps(a psh.EnvList, b psh.EnvList) psh.EnvList {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func EncodeJsonFile(file string) string {
	jsonFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	val := base64.StdEncoding.EncodeToString(byteValue)
	return val
}

func LoadJsonFile(file string) psh.EnvList {
	jsonFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result psh.EnvList
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

// These utilities copied with permission from:
// https://github.com/benbjohnson/testing

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
