/**
 * EVM From Scratch
 * Go template
 *
 * To work on EVM From Scratch in Go:
 *
 * - Install Golang: https://golang.org/doc/install
 * - Go to the `go` directory: `cd go`
 * - Edit `evm.go` (this file!), see TODO below
 * - Run `go test -v` to run the tests
 * - Run `go test -run TestName` to run one test
 */

package main

//go:generate go run testgen.go

import (
	"fmt"

	"github.com/holiman/uint256"
)

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

	opPop      = 0x50
	opMLoad    = 0x51
	opMStore   = 0x52
	opMStore8  = 0x53
	opJump     = 0x56
	opJumpI    = 0x57
	opPC       = 0x58
	opMSize    = 0x59
	opGas      = 0x5a
	opJumpDest = 0x5b
	opPush1    = 0x60
	opPush32   = 0x7f
	opDup1     = 0x80
	opDup16    = 0x8f
	opSwap1    = 0x90
	opSwap16   = 0x9f

	opInvalid = 0xfe
)

func evm(code []byte) (success bool, stack []uint256.Int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERR: %v\n", err)
			success = false
		}
	}()

	pc := uint64(0)
	mem := NewMemory()

	for pc < uint64(len(code)) {
		op := code[pc]
		pc++

		if op >= opPush1 && op <= opPush32 {
			pushLen := uint64(op-opPush1) + 1
			if pushLen > 32 || uint64(len(code)) < pc+pushLen {
				return false, stack
			}
			bytes := code[pc:(pc + pushLen)]
			stack = push(stack, uint256.NewInt(0).SetBytes(bytes))
			pc += pushLen
			continue
		}

		if op >= opDup1 && op <= opDup16 {
			pos := uint64(op - opDup1)
			if uint64(len(stack)) < pos {
				return false, stack
			}
			stack = push(stack, &stack[pos])
			continue
		}

		if op >= opSwap1 && op <= opSwap16 {
			pos := uint64(op-opSwap1) + 1
			if uint64(len(stack)) < pos {
				return false, stack
			}
			stack[0], stack[pos] = stack[pos], stack[0]
			continue
		}

		switch op {
		case opStop:
			return true, stack
		case opPop:
			stack, _ = pop(stack, 1)
		case opAdd:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Add(&x, &y))
		case opMul:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Mul(&x, &y))
		case opSub:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Sub(&x, &y))
		case opDiv:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Div(&x, &y))
		case opMod:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Mod(&x, &y))
		case opAddMod:
			var x, y, z uint256.Int
			stack, x, y, z = pop3(stack)
			stack = push(stack, uint256.NewInt(0).AddMod(&x, &y, &z))
		case opMulMod:
			var x, y, z uint256.Int
			stack, x, y, z = pop3(stack)
			stack = push(stack, uint256.NewInt(0).MulMod(&x, &y, &z))
		case opExp:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Exp(&x, &y))
		case opSignExtend:
			var b, x uint256.Int
			stack, b, x = pop2(stack)
			stack = push(stack, uint256.NewInt(0).ExtendSign(&x, &b))
		case opSDiv:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SDiv(&x, &y))
		case opSMod:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SMod(&x, &y))
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
			stack = push(stack, uint256.NewInt(0).And(&x, &y))
		case opOr:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Or(&x, &y))
		case opXor:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Xor(&x, &y))
		case opNot:
			var x uint256.Int
			stack, x = pop1(stack)
			stack = push(stack, uint256.NewInt(0).Not(&x))
		case opByte:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, y.Byte(&x))
		case opShl:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Lsh(&y, uint(x.Uint64())))
		case opShr:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).Rsh(&y, uint(x.Uint64())))
		case opSar:
			var x, y uint256.Int
			stack, x, y = pop2(stack)
			stack = push(stack, uint256.NewInt(0).SRsh(&y, uint(x.Uint64())))
		case opInvalid:
			return false, stack
		case opPC:
			stack = push(stack, uint256.NewInt(pc-1))
		case opGas:
			// TODO: gas is not supported by this version of the course
			stack = push(stack, uint256.NewInt(0).Not(uint256.NewInt(0)))
		case opJump:
			var dest uint256.Int
			stack, dest = pop1(stack)
			dest64, overflow := dest.Uint64WithOverflow()
			if overflow || !validJumpDest(code, dest64) {
				// overflow = dest is more than MaxUint64, and Go can't handle that
				return false, stack
			}
			pc = dest64
		case opJumpDest: // noop
		case opJumpI:
			var dest, doJump uint256.Int
			stack, dest, doJump = pop2(stack)
			dest64, overflow := dest.Uint64WithOverflow()
			if overflow || !validJumpDest(code, dest64) {
				// overflow = dest is more than MaxUint64, and Go can't handle that
				return false, stack
			}
			if !doJump.IsZero() {
				pc = dest64
			}
		case opMLoad:
			var offset uint256.Int
			stack, offset = pop1(stack)
			stack = push(stack, mem.Get(offset.Uint64()))
		case opMStore:
			var offset, val uint256.Int
			stack, offset, val = pop2(stack)
			mem.Put(offset.Uint64(), &val)
		case opMStore8:
			var offset, val uint256.Int
			stack, offset, val = pop2(stack)
			mem.PutByte(offset.Uint64(), byte(val.Uint64()))
		case opMSize:
			stack = push(stack, uint256.NewInt(mem.Len()))
		}
	}

	return true, stack
}

func validJumpDest(code []byte, dest uint64) bool {
	if uint64(len(code)) < dest {
		return false // destination is past the end of the code
	}
	if code[dest] != opJumpDest {
		return false
	}

	// loop through code and make sure jump dest is not inside PUSH data
	for pc := uint64(0); pc < uint64(len(code)); {
		op := code[pc]
		pc++

		if op < opPush1 || op > opPush32 {
			continue
		}

		numBitsToPush := uint64(op) - opPush1 + 1
		if pc <= dest && dest < pc+numBitsToPush {
			return false
		}
		pc += numBitsToPush
	}

	return true
}

func push(stack []uint256.Int, i *uint256.Int) []uint256.Int {
	return append([]uint256.Int{*i}, stack...)
}

func pushBool(stack []uint256.Int, b bool) []uint256.Int {
	if b {
		stack = push(stack, uint256.NewInt(1))
	} else {
		stack = push(stack, uint256.NewInt(0))
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

//func main() {
//	content, err := ioutil.ReadFile("../evm.json")
//	if err != nil {
//		log.Fatal("Error when opening file: ", err)
//	}
//
//	var payload []testCase
//	err = json.Unmarshal(content, &payload)
//	if err != nil {
//		log.Fatal("Error during json.Unmarshal(): ", err)
//	}
//
//	for index, test := range payload {
//		fmt.Printf("Test #%v of %v: %v\n", index+1, len(payload), test.Name)
//
//		bin, err := hex.DecodeString(test.Code.Bin)
//		if err != nil {
//			log.Fatal("Error during hex.DecodeString(): ", err)
//		}
//
//		var expectedStack []uint256.Int
//		for _, s := range test.Expect.Stack {
//			i, err := uint256.FromHex(s)
//			if err != nil {
//				log.Fatal("Error during big.Int.SetString(): ", err)
//			}
//			expectedStack = append(expectedStack, *i)
//		}
//
//		// Note: as the test cases get more complex, you'll need to modify this
//		// to pass down more arguments to the evm function and return more than
//		// just the stack.
//		success, stack := evm(bin)
//
//		match := len(stack) == len(expectedStack)
//		if match {
//			for i, s := range stack {
//				match = match && (s.Cmp(&expectedStack[i]) == 0)
//			}
//		}
//		match = match && (success == test.Expect.Success)
//
//		if !match {
//			fmt.Printf("Instructions: \n%v\n", test.Code.Asm)
//			fmt.Printf("Expected: success=%v, stack=%v\n", test.Expect.Success, toStrings(expectedStack))
//			fmt.Printf("Got:      success=%v, stack=%v\n\n", success, toStrings(stack))
//			fmt.Printf("Hint: %v\n\n", test.Hint)
//			fmt.Printf("Progress: %v/%v\n\n", index, len(payload))
//			log.Fatal("Stack mismatch")
//		}
//	}
//}
