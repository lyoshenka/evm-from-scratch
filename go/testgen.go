//go:build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type testCase struct {
	Name string
	Hint string
	Code struct {
		Bin string
		Asm string
	}
	Expect struct {
		Stack   []string
		Success bool
		Return  string
	}

	FnName  string
	Index   int
	Payload string
}

var tmpl = template.Must(template.New("").Parse(`
package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/holiman/uint256"
)

type testCase struct {
	Name string
	Hint string
	Code struct {
		Bin string
		Asm string
	}
	Expect struct {
		Stack   []string
		Success bool
		Return  string
	}
}

{{ range . }} 
func Test_{{ .Index }}_{{ .FnName }}(t *testing.T) {
	payload := []byte(` + "`" + `{{ .Payload }}` + "`" + `)
	runTest(t, {{ .Index }}, payload)
}

{{ end }}

func runTest(t *testing.T, index int, payload []byte) {
	var test testCase
	err := json.Unmarshal(payload, &test)
	if err != nil {
		t.Fatal("Error during json.Unmarshal(): ", err)
	}

	bin, err := hex.DecodeString(test.Code.Bin)
	if err != nil {
		log.Fatal("Error during hex.DecodeString(): ", err)
	}

	var expectedStack []uint256.Int
	for _, s := range test.Expect.Stack {
		i, err := uint256.FromHex(s)
		if err != nil {
			log.Fatal("Error during big.Int.SetString(): ", err)
		}
		expectedStack = append(expectedStack, *i)
	}

	// Note: as the test cases get more complex, you'll need to modify this
	// to pass down more arguments to the evm function and return more than
	// just the stack.
	success, stack := evm(bin)

	match := len(stack) == len(expectedStack)
	if match {
		for i, s := range stack {
			match = match && (s.Cmp(&expectedStack[i]) == 0)
		}
	}
	match = match && (success == test.Expect.Success)

	if !match {
		fmt.Printf("Instructions: \n%v\n", test.Code.Asm)
		fmt.Printf("Expected: success=%v, stack=%v\n", test.Expect.Success, toStrings(expectedStack))
		fmt.Printf("Got:      success=%v, stack=%v\n\n", success, toStrings(stack))
		fmt.Printf("Hint: %v\n\n", test.Hint)
		fmt.Printf("Progress: %v/%v\n\n", index, len(payload))
		log.Fatal("Stack mismatch")
	}
}

func toStrings(stack []uint256.Int) []string {
	var strings []string
	for _, s := range stack {
		strings = append(strings, s.String())
	}
	return strings
}

`))

func main() {
	content, err := ioutil.ReadFile("../evm.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var cases []testCase
	err = json.Unmarshal(content, &cases)
	if err != nil {
		log.Fatal("Error during json.Unmarshal(): ", err)
	}

	out, err := os.Create("evm_test.go")
	if err != nil {
		log.Fatal("Error creating file: ", err)
	}

	for i, t := range cases {
		cases[i].FnName = niceFnName(t.Name)
		p, err := json.Marshal(t)
		if err != nil {
			log.Fatal("Error re-marshalling: ", err)
		}
		cases[i].Index = i
		cases[i].Payload = strings.ReplaceAll(string(p), "`", "'")
	}

	tmpl.Execute(out, cases)

	out.Close()
}

func niceFnName(s string) string {
	fixed := ""
	re := regexp.MustCompile(`\W+`)
	for _, value := range strings.Fields(re.ReplaceAllString(s, " ")) {
		fixed += strings.Title(strings.ToLower(value))
	}
	return fixed
}
