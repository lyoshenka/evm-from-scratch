/**
 * EVM From Scratch
 * Go template
 *
 * To work on EVM From Scratch in Go:
 *
 * - Install Golang: https://golang.org/doc/install
 * - Go to the `go` directory: `cd go`
 * - Edit `evm.go` (this file!), see TODO below
 * - Run `go run evm.go` to run the tests
 */

package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/holiman/uint256"
)

type code struct {
	Bin string
	Asm string
}

type expect struct {
	Stack   []string
	Success bool
	Return  string
}

type TestCase struct {
	Name   string
	Hint   string
	Code   code
	Expect expect
}

const (
	opStop   = 0x00
	opPush1  = 0x60
	opPush32 = 0x7f
	opPop    = 0x50
)

func evm(code []byte) (bool, []uint256.Int) {
	var stack []uint256.Int
	var err error
	pc := 0

	for pc < len(code) {
		op := code[pc]
		pc++

		// cleaner than listing all push opcodes
		if op >= opPush1 && op <= opPush32 {
			len := int(op-opPush1) + 1
			stack = push(stack, code[pc:(pc+len)]...)
			pc += len
		}

		switch op {
		case opStop:
			return true, stack
		case opPop:
			stack, _, err = pop(stack, 1)
			if err != nil {
				return false, stack
			}
			pc++
		}

	}

	return true, stack
}

func push(stack []uint256.Int, data ...byte) []uint256.Int {
	front := make([]uint256.Int, len(data))
	for i, b := range data {
		f := uint256.NewInt(uint64(b))
		front[i] = *f
	}
	stack = append(front, stack...)
	return stack
}

func pop(stack []uint256.Int, n int) ([]uint256.Int, []uint256.Int, error) {
	if n > len(stack) {
		return stack, nil, errors.New("stack smaller than n")
	}
	vals := make([]uint256.Int, n)
	copy(vals, stack[:n])
	return stack[n:], vals, nil
}

func main() {
	content, err := ioutil.ReadFile("../evm.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []TestCase
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during json.Unmarshal(): ", err)
	}

	for index, test := range payload {
		fmt.Printf("Test #%v of %v: %v\n", index+1, len(payload), test.Name)

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
}

func toStrings(stack []uint256.Int) []string {
	var strings []string
	for _, s := range stack {
		strings = append(strings, s.String())
	}
	return strings
}
