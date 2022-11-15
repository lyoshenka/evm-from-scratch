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
	opStop = 0x00

	opAdd    = 0x01
	opMul    = 0x02
	opSub    = 0x03
	opDiv    = 0x04
	opSDiv   = 0x05
	opMod    = 0x06
	opSMod   = 0x07
	opAddMod = 0x08
	opMulMod = 0x09
	opExp    = 0x0a

	opSignExtend = 0x0b

	opLT     = 0x10
	opGT     = 0x11
	opSLT    = 0x12
	opSGT    = 0x13
	opEQ     = 0x14
	opIsZero = 0x15

	opAnd = 0x16
	opOr  = 0x17
	opXor = 0x18
	opNot = 0x19

	opByte = 0x1a
	opShl  = 0x1b
	opShr  = 0x1c
	opSar  = 0x1d

	opPop    = 0x50
	opPush1  = 0x60
	opPush32 = 0x7f
)

func evm(code []byte) (success bool, stack []uint256.Int) {
	pc := 0

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERR: %v\n", err)
			success = false
		}
	}()

	for pc < len(code) {
		op := code[pc]
		pc++

		if op >= opPush1 && op <= opPush32 { // cleaner than listing all push opcodes
			len := int(op-opPush1) + 1
			stack = push(stack, code[pc:(pc+len)]...)
			pc += len
			continue
		}

		switch op {
		case opStop:
			return true, stack
		case opPop:
			stack, _ = pop(stack, 1)
			pc++
		case opAdd:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Add(&x, &y).Bytes()...)
		case opMul:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Mul(&x, &y).Bytes()...)
		case opSub:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Sub(&x, &y).Bytes()...)
		case opDiv:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Div(&x, &y).Bytes()...)
		case opMod:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Mod(&x, &y).Bytes()...)
		case opAddMod:
			var x, y, z uint256.Int
			stack, x, y, z = pop3(stack)
			stack = push(stack, uint256.NewInt(0).AddMod(&x, &y, &z).Bytes()...)
		case opMulMod:
			var x, y, z uint256.Int
			stack, x, y, z = pop3(stack)
			stack = push(stack, uint256.NewInt(0).MulMod(&x, &y, &z).Bytes()...)
		case opExp:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Exp(&x, &y).Bytes()...)
		case opSignExtend:
			var b, x uint256.Int
			stack, b, x = pop2(stack)
			stack = push(stack, uint256.NewInt(0).ExtendSign(&x, &b).Bytes()...)
		case opSDiv:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SDiv(&x, &y).Bytes()...)
		case opSMod:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SMod(&x, &y).Bytes()...)
		case opLT:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = pushBool(stack, x.Lt(&y))
		case opGT:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = pushBool(stack, x.Gt(&y))
		case opSLT:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = pushBool(stack, x.Slt(&y))
		case opSGT:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = pushBool(stack, x.Sgt(&y))
		case opEQ:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = pushBool(stack, x.Eq(&y))
		case opIsZero:
			var x uint256.Int
			stack, x = pop1(stack)
			stack = pushBool(stack, x.IsZero())
		case opAnd:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).And(&x, &y).Bytes()...)
		case opOr:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Or(&x, &y).Bytes()...)
		case opXor:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Xor(&x, &y).Bytes()...)
		case opNot:
			var x uint256.Int
			stack, x = pop1(stack)
			stack = push(stack, uint256.NewInt(0).Not(&x).Bytes()...)
		case opByte:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, y.Byte(&x).Bytes()...)
		case opShl:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Lsh(&y, uint(x.Uint64())).Bytes()...)
		case opShr:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Rsh(&y, uint(x.Uint64())).Bytes()...)
		case opSar:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SRsh(&y, uint(x.Uint64())).Bytes()...)
		}
	}

	return true, stack
}

func push(stack []uint256.Int, data ...byte) []uint256.Int {
	if len(data) > 32 {
		panic("data too long")
	}
	i := uint256.NewInt(0).SetBytes(data)
	return append([]uint256.Int{*i}, stack...)
}

func pushBool(stack []uint256.Int, b bool) []uint256.Int {
	if b {
		stack = push(stack, uint256.NewInt(1).Bytes()...)
	} else {
		stack = push(stack, uint256.NewInt(0).Bytes()...)
	}
	return stack
}

func pop(stack []uint256.Int, n int) ([]uint256.Int, []uint256.Int) {
	if n > len(stack) {
		panic(fmt.Errorf("stack len (%d) is smaller than %d", len(stack), n))
	}
	vals := make([]uint256.Int, n)
	copy(vals, stack[:n])
	return stack[n:], vals
}

func pop1(stack []uint256.Int) ([]uint256.Int, uint256.Int) {
	stack, vals := pop(stack, 1)
	return stack, vals[0]
}

func pop2(stack []uint256.Int) ([]uint256.Int, uint256.Int, uint256.Int) {
	stack, vals := pop(stack, 2)
	return stack, vals[0], vals[1]
}

func pop3(stack []uint256.Int) ([]uint256.Int, uint256.Int, uint256.Int, uint256.Int) {
	stack, vals := pop(stack, 3)
	return stack, vals[0], vals[1], vals[2]
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
