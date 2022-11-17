
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

 
func Test_0_Stop(t *testing.T) {
	payload := []byte(`{"Name":"STOP","Hint":"","Code":{"Bin":"00","Asm":"STOP"},"Expect":{"Stack":[],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 0, payload)
}

 
func Test_1_Push(t *testing.T) {
	payload := []byte(`{"Name":"PUSH","Hint":"Read \"Program Counter\" section of the course learning materials for an example on how to parse the bytecode","Code":{"Bin":"6001","Asm":"PUSH1 1"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 1, payload)
}

 
func Test_2_Push2(t *testing.T) {
	payload := []byte(`{"Name":"PUSH2","Hint":"PUSH2 reads the next 2 bytes, don't forget to properly increment PC","Code":{"Bin":"611122","Asm":"PUSH2 0x1122"},"Expect":{"Stack":["0x1122"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 2, payload)
}

 
func Test_3_Push4(t *testing.T) {
	payload := []byte(`{"Name":"PUSH4","Hint":"PUSH2 reads the next 4 bytes","Code":{"Bin":"6300112233","Asm":"PUSH4 0x112233"},"Expect":{"Stack":["0x112233"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 3, payload)
}

 
func Test_4_Push6(t *testing.T) {
	payload := []byte(`{"Name":"PUSH6","Hint":"PUSH6 reads the next 6 bytes. Can you implement all PUSH1...PUSH32 using the same code?","Code":{"Bin":"65112233445566","Asm":"PUSH6 0x112233445566"},"Expect":{"Stack":["0x112233445566"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 4, payload)
}

 
func Test_5_Push10(t *testing.T) {
	payload := []byte(`{"Name":"PUSH10","Hint":"SIZE = OPCODE - PUSH1 + 1, then transform take the next SIZE bytes, PC += SIZE","Code":{"Bin":"69112233445566778899aa","Asm":"PUSH10 0x112233445566778899aa"},"Expect":{"Stack":["0x112233445566778899aa"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 5, payload)
}

 
func Test_6_Push11(t *testing.T) {
	payload := []byte(`{"Name":"PUSH11","Hint":"SIZE = OPCODE - PUSH1 + 1, program.slice(pc + 1, pc + 1 + size)","Code":{"Bin":"6a112233445566778899aabb","Asm":"PUSH11 0x112233445566778899aabb"},"Expect":{"Stack":["0x112233445566778899aabb"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 6, payload)
}

 
func Test_7_Push32(t *testing.T) {
	payload := []byte(`{"Name":"PUSH32","Hint":"PUSH32 reads the next 32 bytes (256 bits)","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},"Expect":{"Stack":["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 7, payload)
}

 
func Test_8_PushTwice(t *testing.T) {
	payload := []byte(`{"Name":"PUSH (twice)","Hint":"Note the order of items on the stack. The tests expect the top of the stack to be the first element","Code":{"Bin":"60016002","Asm":"PUSH1 1\nPUSH1 2"},"Expect":{"Stack":["0x2","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 8, payload)
}

 
func Test_9_Pop(t *testing.T) {
	payload := []byte(`{"Name":"POP","Hint":"POP removes the top item from the stack and discards it","Code":{"Bin":"6001600250","Asm":"PUSH1 1\nPUSH1 2\nPOP"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 9, payload)
}

 
func Test_10_StopMidway(t *testing.T) {
	payload := []byte(`{"Name":"STOP (midway)","Hint":"Note that the 'PUSH1 2' didn't execute because the program stops after STOP opcode","Code":{"Bin":"6001006002","Asm":"PUSH1 1\nSTOP\nPUSH1 2"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 10, payload)
}

 
func Test_11_Add(t *testing.T) {
	payload := []byte(`{"Name":"ADD","Hint":"ADD takes the first 2 items from the stack, adds them together and pushes the result","Code":{"Bin":"6001600201","Asm":"PUSH1 0x01\nPUSH1 0x02\nADD"},"Expect":{"Stack":["0x3"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 11, payload)
}

 
func Test_12_AddOverflow(t *testing.T) {
	payload := []byte(`{"Name":"ADD (overflow)","Hint":"EVM operates with uint256, if you add 2 to the max possible value it overflows and wraps around","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600201","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH1 0x02\nADD"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 12, payload)
}

 
func Test_13_Mul(t *testing.T) {
	payload := []byte(`{"Name":"MUL","Hint":"","Code":{"Bin":"6002600302","Asm":"PUSH1 0x02\nPUSH1 0x03\nMUL"},"Expect":{"Stack":["0x6"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 13, payload)
}

 
func Test_14_MulOverflow(t *testing.T) {
	payload := []byte(`{"Name":"MUL (overflow)","Hint":"All math is performed with implicit [mod 2^256]","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600202","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH1 0x02\nMUL"},"Expect":{"Stack":["0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 14, payload)
}

 
func Test_15_Sub(t *testing.T) {
	payload := []byte(`{"Name":"SUB","Hint":"SUB takes the first element from the stack and subtracts the second element from the stack","Code":{"Bin":"6002600303","Asm":"PUSH1 0x02\nPUSH1 0x03\nSUB"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 15, payload)
}

 
func Test_16_SubUnderflow(t *testing.T) {
	payload := []byte(`{"Name":"SUB (underflow)","Hint":"Underflow works the same way as overflow, 3 - 2 wraps around and results in MAX_UINT256","Code":{"Bin":"6003600203","Asm":"PUSH1 0x03\nPUSH1 0x02\nSUB"},"Expect":{"Stack":["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 16, payload)
}

 
func Test_17_Div(t *testing.T) {
	payload := []byte(`{"Name":"DIV","Hint":"DIV takes the first element from the stack and divides it by the second element from the stack","Code":{"Bin":"6002600604","Asm":"PUSH1 0x02\nPUSH1 0x06\nDIV"},"Expect":{"Stack":["0x3"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 17, payload)
}

 
func Test_18_DivWhole(t *testing.T) {
	payload := []byte(`{"Name":"DIV (whole)","Hint":"Fraction part of the division is discarded","Code":{"Bin":"6006600204","Asm":"PUSH1 0x06\nPUSH1 0x02\nDIV"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 18, payload)
}

 
func Test_19_DivByZero(t *testing.T) {
	payload := []byte(`{"Name":"DIV (by zero)","Hint":"In EVM you can divide by zero! Modern Solidity protects from this by adding instructions that check for zero","Code":{"Bin":"6000600204","Asm":"PUSH1 0x00\nPUSH1 0x02\nDIV"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 19, payload)
}

 
func Test_20_Mod(t *testing.T) {
	payload := []byte(`{"Name":"MOD","Hint":"10 mod 3 = 1","Code":{"Bin":"6003600a06","Asm":"PUSH1 3\nPUSH1 10\nMOD"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 20, payload)
}

 
func Test_21_ModByLargerNumber(t *testing.T) {
	payload := []byte(`{"Name":"MOD (by larger number)","Hint":"5 mod 17 = 5","Code":{"Bin":"6011600506","Asm":"PUSH1 17\nPUSH1 5\nMOD"},"Expect":{"Stack":["0x5"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 21, payload)
}

 
func Test_22_ModByZero(t *testing.T) {
	payload := []byte(`{"Name":"MOD (by zero)","Hint":"In EVM you can divide by zero! Modern Solidity protects from this by adding instructions that check for zero","Code":{"Bin":"6000600206","Asm":"PUSH1 0\nPUSH1 2\nMOD"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 22, payload)
}

 
func Test_23_Addmod(t *testing.T) {
	payload := []byte(`{"Name":"ADDMOD","Hint":"10 + 10 mod 8 = 4","Code":{"Bin":"6008600a600a08","Asm":"PUSH1 8\nPUSH1 10\nPUSH1 10\nADDMOD"},"Expect":{"Stack":["0x4"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 23, payload)
}

 
func Test_24_AddmodWrapped(t *testing.T) {
	payload := []byte(`{"Name":"ADDMOD (wrapped)","Hint":"","Code":{"Bin":"600260027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff08","Asm":"PUSH1 2\nPUSH1 2\nPUSH32 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\nADDMOD"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 24, payload)
}

 
func Test_25_Mulmod(t *testing.T) {
	payload := []byte(`{"Name":"MULMOD","Hint":"10 * 10 mod 8 = 4","Code":{"Bin":"6008600a600a09","Asm":"PUSH1 8\nPUSH1 10\nPUSH1 10\nMULMOD"},"Expect":{"Stack":["0x4"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 25, payload)
}

 
func Test_26_MulmodWrapped(t *testing.T) {
	payload := []byte(`{"Name":"MULMOD (wrapped)","Hint":"","Code":{"Bin":"600c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff09","Asm":"PUSH1 12\nPUSH32 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\nPUSH32 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\nMULMOD"},"Expect":{"Stack":["0x9"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 26, payload)
}

 
func Test_27_Exp(t *testing.T) {
	payload := []byte(`{"Name":"EXP","Hint":"","Code":{"Bin":"6002600a0a","Asm":"PUSH1 2\nPUSH1 10\nEXP"},"Expect":{"Stack":["0x64"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 27, payload)
}

 
func Test_28_SignextendPositive(t *testing.T) {
	payload := []byte(`{"Name":"SIGNEXTEND (positive)","Hint":"Read \"Negative Numbers\" section of the course learning materials. SIGNEXTEND has no effect on \"positive\" numbers","Code":{"Bin":"607f60000b","Asm":"PUSH1 0x7F\nPUSH1 0\nSIGNEXTEND"},"Expect":{"Stack":["0x7f"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 28, payload)
}

 
func Test_29_SignextendNegative(t *testing.T) {
	payload := []byte(`{"Name":"SIGNEXTEND (negative)","Hint":"Read \"Negative Numbers\" section of the course learning materials. The first bit of 0xFF is 1, so it is a negative number and needs to be padded by 1s in front","Code":{"Bin":"60ff60000b","Asm":"PUSH1 0xFF\nPUSH1 0\nSIGNEXTEND"},"Expect":{"Stack":["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 29, payload)
}

 
func Test_30_Sdiv(t *testing.T) {
	payload := []byte(`{"Name":"SDIV","Hint":"Read \"Negative Numbers\" section of the course learning materials. SDIV works like DIV for \"positive\" numbers","Code":{"Bin":"600a600a05","Asm":"PUSH1 10\nPUSH1 10\nSDIV"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 30, payload)
}

 
func Test_31_SdivNegative(t *testing.T) {
	payload := []byte(`{"Name":"SDIV (negative)","Hint":"Read \"Negative Numbers\" section of the course learning materials. -2 / -1 = 2","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe05","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe\nSDIV"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 31, payload)
}

 
func Test_32_SdivMixOfNegativeAndPositive(t *testing.T) {
	payload := []byte(`{"Name":"SDIV (mix of negative and positive)","Hint":"Read \"Negative Numbers\" section of the course learning materials. 10 / -2 = -5","Code":{"Bin":"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe600a05","Asm":"PUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe\nPUSH1 10\nSDIV"},"Expect":{"Stack":["0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 32, payload)
}

 
func Test_33_Smod(t *testing.T) {
	payload := []byte(`{"Name":"SMOD","Hint":"Read \"Negative Numbers\" section of the course learning materials. SMOD works like MOD for \"positive\" numbers","Code":{"Bin":"6003600a07","Asm":"PUSH1 3\nPUSH1 10\nSMOD"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 33, payload)
}

 
func Test_34_SmodNegative(t *testing.T) {
	payload := []byte(`{"Name":"SMOD (negative)","Hint":"Read \"Negative Numbers\" section of the course learning materials. -10 mod -3 = -1","Code":{"Bin":"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff807","Asm":"PUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8\nSMOD"},"Expect":{"Stack":["0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 34, payload)
}

 
func Test_35_SdivByZero(t *testing.T) {
	payload := []byte(`{"Name":"SDIV (by zero)","Hint":"In EVM you can divide by zero","Code":{"Bin":"60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd05","Asm":"PUSH1 0x00\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd\nSDIV"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 35, payload)
}

 
func Test_36_SmodByZero(t *testing.T) {
	payload := []byte(`{"Name":"SMOD (by zero)","Hint":"In EVM you can divide by zero","Code":{"Bin":"60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd07","Asm":"PUSH1 0x00\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd\nSMOD"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 36, payload)
}

 
func Test_37_Lt(t *testing.T) {
	payload := []byte(`{"Name":"LT","Hint":"9 \u003c 10 = true (1)","Code":{"Bin":"600a600910","Asm":"PUSH1 10\nPUSH1 9\nLT"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 37, payload)
}

 
func Test_38_LtEqual(t *testing.T) {
	payload := []byte(`{"Name":"LT (equal)","Hint":"10 \u003c 10 = false (0)","Code":{"Bin":"600a600a10","Asm":"PUSH1 10\nPUSH1 10\nLT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 38, payload)
}

 
func Test_39_LtGreater(t *testing.T) {
	payload := []byte(`{"Name":"LT (greater)","Hint":"11 \u003c 10 = false (0)","Code":{"Bin":"600a600b10","Asm":"PUSH1 10\nPUSH1 11\nLT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 39, payload)
}

 
func Test_40_Gt(t *testing.T) {
	payload := []byte(`{"Name":"GT","Hint":"10 \u003e 9 = true (1)","Code":{"Bin":"6009600a11","Asm":"PUSH1 9\nPUSH1 10\nGT"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 40, payload)
}

 
func Test_41_GtEqual(t *testing.T) {
	payload := []byte(`{"Name":"GT (equal)","Hint":"10 \u003e 10 = false (0)","Code":{"Bin":"600a600a11","Asm":"PUSH1 10\nPUSH1 10\nGT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 41, payload)
}

 
func Test_42_GtLess(t *testing.T) {
	payload := []byte(`{"Name":"GT (less)","Hint":"10 \u003e 11 = false (0)","Code":{"Bin":"600b600a11","Asm":"PUSH1 11\nPUSH1 10\nGT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 42, payload)
}

 
func Test_43_Slt(t *testing.T) {
	payload := []byte(`{"Name":"SLT","Hint":"Same as LT but treats arguments as signed numbers. -1 \u003c 0 = true (1)","Code":{"Bin":"60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff12","Asm":"PUSH1 0\nPUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nSLT"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 43, payload)
}

 
func Test_44_SltEqual(t *testing.T) {
	payload := []byte(`{"Name":"SLT (equal)","Hint":"Same as LT but treats arguments as signed numbers. -1 \u003c -1 = false (0)","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff12","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nSLT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 44, payload)
}

 
func Test_45_SltGreater(t *testing.T) {
	payload := []byte(`{"Name":"SLT (greater)","Hint":"Same as LT but treats arguments as signed numbers. -1 \u003c -1 = false (0)","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600012","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH1 0\nSLT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 45, payload)
}

 
func Test_46_Sgt(t *testing.T) {
	payload := []byte(`{"Name":"SGT","Hint":"Same as GT but treats arguments as signed numbers. No effect on \"positive\" numbers: 10 \u003e 9 = true (1)","Code":{"Bin":"6009600a13","Asm":"PUSH1 9\nPUSH1 10\nSGT"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 46, payload)
}

 
func Test_47_SgtEqual(t *testing.T) {
	payload := []byte(`{"Name":"SGT (equal)","Hint":"Same as GT but treats arguments as signed numbers. -2 \u003e -2 = false (0)","Code":{"Bin":"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe13","Asm":"PUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe\nSGT"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 47, payload)
}

 
func Test_48_SgtGreater(t *testing.T) {
	payload := []byte(`{"Name":"SGT (greater)","Hint":"Same as GT but treats arguments as signed numbers. -2 \u003e -3 = true (1)","Code":{"Bin":"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe13","Asm":"PUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd\nPUSH32 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe\nSGT"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 48, payload)
}

 
func Test_49_Eq(t *testing.T) {
	payload := []byte(`{"Name":"EQ","Hint":"10 == 10 = true (1)","Code":{"Bin":"600a600a14","Asm":"PUSH1 10\nPUSH1 10\nEQ"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 49, payload)
}

 
func Test_50_EqNotEqual(t *testing.T) {
	payload := []byte(`{"Name":"EQ (not equal)","Hint":"10 == 9 = false (0)","Code":{"Bin":"6009600a14","Asm":"PUSH1 9\nPUSH1 10\nEQ"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 50, payload)
}

 
func Test_51_IszeroNotZero(t *testing.T) {
	payload := []byte(`{"Name":"ISZERO (not zero)","Hint":"If the top element on the stack is not zero, pushes 0","Code":{"Bin":"600915","Asm":"PUSH1 9\nISZERO"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 51, payload)
}

 
func Test_52_IszeroZero(t *testing.T) {
	payload := []byte(`{"Name":"ISZERO (zero)","Hint":"If the top element on the stack is zero, pushes 1","Code":{"Bin":"600015","Asm":"PUSH1 0\nISZERO"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 52, payload)
}

 
func Test_53_Not(t *testing.T) {
	payload := []byte(`{"Name":"NOT","Hint":"Bitwise NOT operation, flips every bit 1-\u003e0, 0-\u003e1","Code":{"Bin":"600f19","Asm":"PUSH1 0x0f\nNOT"},"Expect":{"Stack":["0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 53, payload)
}

 
func Test_54_And(t *testing.T) {
	payload := []byte(`{"Name":"AND","Hint":"Bitwise AND operation of the top 2 items on the stack","Code":{"Bin":"600e600316","Asm":"PUSH1 0xe\nPUSH1 0x3\nAND"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 54, payload)
}

 
func Test_55_Or(t *testing.T) {
	payload := []byte(`{"Name":"OR","Hint":"Bitwise OR operation of the top 2 items on the stack","Code":{"Bin":"600e600317","Asm":"PUSH1 0xe\nPUSH1 0x3\nOR"},"Expect":{"Stack":["0xf"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 55, payload)
}

 
func Test_56_Xor(t *testing.T) {
	payload := []byte(`{"Name":"XOR","Hint":"Bitwise XOR operation of the top 2 items on the stack","Code":{"Bin":"60f0600f18","Asm":"PUSH1 0xf0\nPUSH1 0x0f\nXOR"},"Expect":{"Stack":["0xff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 56, payload)
}

 
func Test_57_Shl(t *testing.T) {
	payload := []byte(`{"Name":"SHL","Hint":"Bitwise shift left, 1 \u003c\u003c 1 = 2","Code":{"Bin":"600160011b","Asm":"PUSH1 1\nPUSH1 1\nSHL"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 57, payload)
}

 
func Test_58_ShlDiscards(t *testing.T) {
	payload := []byte(`{"Name":"SHL (discards)","Hint":"Bits that end up outside MAX_UINT256 are discarded","Code":{"Bin":"7fff0000000000000000000000000000000000000000000000000000000000000060041b","Asm":"PUSH32 0xFF00000000000000000000000000000000000000000000000000000000000000\nPUSH1 4\nSHL"},"Expect":{"Stack":["0xf000000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 58, payload)
}

 
func Test_59_ShlTooLarge(t *testing.T) {
	payload := []byte(`{"Name":"SHL (too large)","Hint":"When shift amount is too large, returns zero","Code":{"Bin":"600163ffffffff1b","Asm":"PUSH1 1\nPUSH4 0xFFFFFFFF\nSHL"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 59, payload)
}

 
func Test_60_Shr(t *testing.T) {
	payload := []byte(`{"Name":"SHR","Hint":"Bitwise shift right, 2 \u003e\u003e 1 = 1","Code":{"Bin":"600260011c","Asm":"PUSH1 2\nPUSH1 1\nSHR"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 60, payload)
}

 
func Test_61_ShrDiscards(t *testing.T) {
	payload := []byte(`{"Name":"SHR (discards)","Hint":"Bits that end up outside are discarded","Code":{"Bin":"60ff60041c","Asm":"PUSH1 0xFF\nPUSH1 4\nSHR"},"Expect":{"Stack":["0xf"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 61, payload)
}

 
func Test_62_ShrTooLarge(t *testing.T) {
	payload := []byte(`{"Name":"SHR (too large)","Hint":"When shift amount is too large, returns zero","Code":{"Bin":"600163ffffffff1c","Asm":"PUSH1 1\nPUSH4 0xFFFFFFFF\nSHR"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 62, payload)
}

 
func Test_63_Sar(t *testing.T) {
	payload := []byte(`{"Name":"SAR","Hint":"Like SHR but treats the argument as signed number. No effect on \"positive\" numbers, 2 \u003e\u003e 1 = 1","Code":{"Bin":"600260011d","Asm":"PUSH1 2\nPUSH1 1\nSAR"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 63, payload)
}

 
func Test_64_SarFills1s(t *testing.T) {
	payload := []byte(`{"Name":"SAR (fills 1s)","Hint":"Note that unlike SHR, it fills the empty space with 1s","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060041d","Asm":"PUSH32 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00\nPUSH1 4\nSAR"},"Expect":{"Stack":["0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 64, payload)
}

 
func Test_65_SarTooLarge(t *testing.T) {
	payload := []byte(`{"Name":"SAR (too large)","Hint":"When shift amount is too large and the first bit is 1, fills the whole number with 1s","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0063ffffffff1d","Asm":"PUSH32 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00\nPUSH4 0xFFFFFFFF\nSAR"},"Expect":{"Stack":["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 65, payload)
}

 
func Test_66_SarPositiveTooLarge(t *testing.T) {
	payload := []byte(`{"Name":"SAR (positive, too large)","Hint":"When shift amount is too large and the first bit is 0, fills the whole number with 0s","Code":{"Bin":"7f0fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0063ffffffff1d","Asm":"PUSH32 0x0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00\nPUSH4 0xFFFFFFFF\nSAR"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 66, payload)
}

 
func Test_67_Byte(t *testing.T) {
	payload := []byte(`{"Name":"BYTE","Hint":"The value on the stack is treated as 32 bytes, take 31st (counting from the most significant one)","Code":{"Bin":"60ff601f1a","Asm":"PUSH1 0xff\nPUSH1 31\nBYTE"},"Expect":{"Stack":["0xff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 67, payload)
}

 
func Test_68_Byte30th(t *testing.T) {
	payload := []byte(`{"Name":"BYTE (30th)","Hint":"The value on the stack is treated as 32 bytes, take 30st (counting from the most significant one)","Code":{"Bin":"61ff00601e1a","Asm":"PUSH2 0xff00\nPUSH1 30\nBYTE"},"Expect":{"Stack":["0xff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 68, payload)
}

 
func Test_69_Byte29th(t *testing.T) {
	payload := []byte(`{"Name":"BYTE (29th)","Hint":"Try to generalize your code to work with any argument","Code":{"Bin":"62ff0000601d1a","Asm":"PUSH3 0xff0000\nPUSH1 29\nBYTE"},"Expect":{"Stack":["0xff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 69, payload)
}

 
func Test_70_ByteOutOfRange(t *testing.T) {
	payload := []byte(`{"Name":"BYTE (out of range)","Hint":"Treat other elements as zeros","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602a1a","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPUSH1 42\nBYTE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 70, payload)
}

 
func Test_71_Dup1(t *testing.T) {
	payload := []byte(`{"Name":"DUP1","Hint":"Duplicate the first element from the stack and push it onto the stack","Code":{"Bin":"60018001","Asm":"PUSH1 1\nDUP1\nADD"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 71, payload)
}

 
func Test_72_Dup3(t *testing.T) {
	payload := []byte(`{"Name":"DUP3","Hint":"Duplicate the 3rd element from the stack and push it onto the stack","Code":{"Bin":"60016002600382","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nDUP3"},"Expect":{"Stack":["0x1","0x3","0x2","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 72, payload)
}

 
func Test_73_Dup5(t *testing.T) {
	payload := []byte(`{"Name":"DUP5","Hint":"Try to implement your code to handle any DUP1...DUP16","Code":{"Bin":"6001600260036004600584","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nPUSH1 4\nPUSH1 5\nDUP5"},"Expect":{"Stack":["0x1","0x5","0x4","0x3","0x2","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 73, payload)
}

 
func Test_74_Dup8(t *testing.T) {
	payload := []byte(`{"Name":"DUP8","Hint":"No seriously try to implement your code to handle any DUP1...DUP16 generically. You can do OPCODE - DUP1 + 1 to learn which item to take from the stack","Code":{"Bin":"6001600260036004600560066007600887","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nPUSH1 4\nPUSH1 5\nPUSH1 6\nPUSH1 7\nPUSH1 8\nDUP8"},"Expect":{"Stack":["0x1","0x8","0x7","0x6","0x5","0x4","0x3","0x2","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 74, payload)
}

 
func Test_75_Swap(t *testing.T) {
	payload := []byte(`{"Name":"SWAP","Hint":"Swap the top item from the stack with the 1st one after that","Code":{"Bin":"6001600290","Asm":"PUSH1 1\nPUSH1 2\nSWAP1"},"Expect":{"Stack":["0x1","0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 75, payload)
}

 
func Test_76_Swap3(t *testing.T) {
	payload := []byte(`{"Name":"SWAP3","Hint":"Swap the top item from the stack with the 3rd one after that","Code":{"Bin":"600160026003600492","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nPUSH1 4\nSWAP3"},"Expect":{"Stack":["0x1","0x3","0x2","0x4"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 76, payload)
}

 
func Test_77_Swap5(t *testing.T) {
	payload := []byte(`{"Name":"SWAP5","Hint":"Swap the top item from the stack with the 5th one after that. Try to implement SWAP1..SWAP16 with the same code","Code":{"Bin":"60016002600360046005600694","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nPUSH1 4\nPUSH1 5\nPUSH1 6\nSWAP5"},"Expect":{"Stack":["0x1","0x5","0x4","0x3","0x2","0x6"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 77, payload)
}

 
func Test_78_Swap7(t *testing.T) {
	payload := []byte(`{"Name":"SWAP7","Hint":"No seriously try to implement your code to handle any SWAP1...SWAP16 generically. You can do OPCODE - SWAP1 + 2 to learn which item to take from the stack","Code":{"Bin":"6001600260036004600560066007600896","Asm":"PUSH1 1\nPUSH1 2\nPUSH1 3\nPUSH1 4\nPUSH1 5\nPUSH1 6\nPUSH1 7\nPUSH1 8\nSWAP7"},"Expect":{"Stack":["0x1","0x7","0x6","0x5","0x4","0x3","0x2","0x8"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 78, payload)
}

 
func Test_79_Invalid(t *testing.T) {
	payload := []byte(`{"Name":"INVALID","Hint":"Invalid instruction. Note that your code is expected to return success = false, not throw exceptions","Code":{"Bin":"fe","Asm":"INVALID"},"Expect":{"Stack":[],"Success":false,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 79, payload)
}

 
func Test_80_Pc(t *testing.T) {
	payload := []byte(`{"Name":"PC","Hint":"Read \"Program Counter\" section of the course learning materials","Code":{"Bin":"58","Asm":"PC"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 80, payload)
}

 
func Test_81_PcMoreCode(t *testing.T) {
	payload := []byte(`{"Name":"PC (more code)","Hint":"'PUSH1 0' is counted as 2 bytes (even though it is a single instruction)","Code":{"Bin":"60005058","Asm":"PUSH1 0\nPOP\nPC"},"Expect":{"Stack":["0x3"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 81, payload)
}

 
func Test_82_Gas(t *testing.T) {
	payload := []byte(`{"Name":"GAS","Hint":"In this version of the tests, GAS is not supported yet and is always expected to return MAX_UINT256","Code":{"Bin":"5a","Asm":"GAS"},"Expect":{"Stack":["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 82, payload)
}

 
func Test_83_Jump(t *testing.T) {
	payload := []byte(`{"Name":"JUMP","Hint":"Set the Program Counter (PC) to the top value from the stack","Code":{"Bin":"60055660015b6002","Asm":"PUSH1 5\nJUMP\nPUSH1 1\nJUMPDEST\nPUSH1 2"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 83, payload)
}

 
func Test_84_JumpNotJumpdest(t *testing.T) {
	payload := []byte(`{"Name":"JUMP (not JUMPDEST)","Hint":"Offset 4 is not a valid JUMPDEST instruction","Code":{"Bin":"6003566001","Asm":"PUSH1 3\nJUMP\nPUSH1 1"},"Expect":{"Stack":[],"Success":false,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 84, payload)
}

 
func Test_85_JumpBadInstructionBoundry(t *testing.T) {
	payload := []byte(`{"Name":"JUMP (bad instruction boundry)","Hint":"See \"9.4.3. Jump Destination Validity\" of the Yellow Paper https://ethereum.github.io/yellowpaper/paper.pdf","Code":{"Bin":"600456605b60ff","Asm":"PUSH1 4\nJUMP\nPUSH1 0x5b\nPUSH1 0xff"},"Expect":{"Stack":[],"Success":false,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 85, payload)
}

 
func Test_86_JumpiNoJump(t *testing.T) {
	payload := []byte(`{"Name":"JUMPI (no jump)","Hint":"Conditional JUMP, second argument is 0, not jumping","Code":{"Bin":"600060075760015b600250","Asm":"PUSH1 0\nPUSH1 7\nJUMPI\nPUSH1 1\nJUMPDEST\nPUSH1 2\nPOP"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 86, payload)
}

 
func Test_87_JumpiJump(t *testing.T) {
	payload := []byte(`{"Name":"JUMPI (jump)","Hint":"Conditional JUMP, second argument is not 0, jumping","Code":{"Bin":"600160075760015b6002","Asm":"PUSH1 1\nPUSH1 7\nJUMPI\nPUSH1 1\nJUMPDEST\nPUSH1 2"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 87, payload)
}

 
func Test_88_Mstore(t *testing.T) {
	payload := []byte(`{"Name":"MSTORE","Hint":"Read \"Memory\" section of the course learning materials before implementing memory opcodes","Code":{"Bin":"7f0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20600052600051","Asm":"PUSH32 0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20\nPUSH1 0\nMSTORE\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 88, payload)
}

 
func Test_89_MstoreTail(t *testing.T) {
	payload := []byte(`{"Name":"MSTORE (tail)","Hint":"MLOAD starts from byte offset 31 and picks up the last byte (0x20), the rest of the memory is 00","Code":{"Bin":"7f0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20600052601f51","Asm":"PUSH32 0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20\nPUSH1 0\nMSTORE\nPUSH1 31\nMLOAD"},"Expect":{"Stack":["0x2000000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 89, payload)
}

 
func Test_90_Mstore8(t *testing.T) {
	payload := []byte(`{"Name":"MSTORE8","Hint":"Store a single byte at the given offset","Code":{"Bin":"60ff601f53600051","Asm":"PUSH1 0xff\nPUSH1 31\nMSTORE8\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0xff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 90, payload)
}

 
func Test_91_Msize(t *testing.T) {
	payload := []byte(`{"Name":"MSIZE","Hint":"No memory has been accessed, so the memory size is 0","Code":{"Bin":"59","Asm":"MSIZE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 91, payload)
}

 
func Test_92_Msize0x20(t *testing.T) {
	payload := []byte(`{"Name":"MSIZE (0x20)","Hint":"The first 32-byte section has been accessed, so the memory size is 32 (0x20)","Code":{"Bin":"6000515059","Asm":"PUSH1 0\nMLOAD\nPOP\nMSIZE"},"Expect":{"Stack":["0x20"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 92, payload)
}

 
func Test_93_Msize0x60(t *testing.T) {
	payload := []byte(`{"Name":"MSIZE (0x60)","Hint":"Memory is measured in 32-byte chunks","Code":{"Bin":"6039515059","Asm":"PUSH1 0x39\nMLOAD\nPOP\nMSIZE"},"Expect":{"Stack":["0x60"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 93, payload)
}

 
func Test_94_MsizeAfterMstore8(t *testing.T) {
	payload := []byte(`{"Name":"MSIZE (after MSTORE8)","Hint":"Any opcode touching memory should update MSIZE, including the future ones. Implement memory access in a way that automatically updates MSIZE no matter which opcode used it","Code":{"Bin":"60ff60ff5359","Asm":"PUSH1 0xff\nPUSH1 0xff\nMSTORE8\nMSIZE"},"Expect":{"Stack":["0x100"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 94, payload)
}

 
func Test_95_Sha3(t *testing.T) {
	payload := []byte(`{"Name":"SHA3","Hint":"Use an existing library for your programming language. Note that even though the opcode is called SHA3, the algorythm used is keccak256","Code":{"Bin":"7fffffffff000000000000000000000000000000000000000000000000000000006000526004600020","Asm":"PUSH32 0xffffffff00000000000000000000000000000000000000000000000000000000\nPUSH1 0\nMSTORE\nPUSH1 4\nPUSH1 0\nSHA3"},"Expect":{"Stack":["0x29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c238"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 95, payload)
}

 
func Test_96_Address(t *testing.T) {
	payload := []byte(`{"Name":"ADDRESS","Hint":"Read \"Transaction\" section of the course learning materials. Change your evm function parameters list to include transaction data","Code":{"Bin":"30","Asm":"ADDRESS"},"Expect":{"Stack":["0x1000000000000000000000000000000000000aaa"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 96, payload)
}

 
func Test_97_Caller(t *testing.T) {
	payload := []byte(`{"Name":"CALLER","Hint":"Solidity calls this msg.sender","Code":{"Bin":"33","Asm":"CALLER"},"Expect":{"Stack":["0x1e79b045dc29eae9fdc69673c9dcd7c53e5e159d"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 97, payload)
}

 
func Test_98_Origin(t *testing.T) {
	payload := []byte(`{"Name":"ORIGIN","Hint":"Solidity calls this tx.origin","Code":{"Bin":"32","Asm":"ORIGIN"},"Expect":{"Stack":["0x1337"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 98, payload)
}

 
func Test_99_Gasprice(t *testing.T) {
	payload := []byte(`{"Name":"GASPRICE","Hint":"","Code":{"Bin":"3a","Asm":"GASPRICE"},"Expect":{"Stack":["0x99"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 99, payload)
}

 
func Test_100_Basefee(t *testing.T) {
	payload := []byte(`{"Name":"BASEFEE","Hint":"","Code":{"Bin":"48","Asm":"BASEFEE"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 100, payload)
}

 
func Test_101_Coinbase(t *testing.T) {
	payload := []byte(`{"Name":"COINBASE","Hint":"Do not hardcode these numbers, pull them from the test cases","Code":{"Bin":"41","Asm":"COINBASE"},"Expect":{"Stack":["0x777"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 101, payload)
}

 
func Test_102_CoinbaseDifferentOne(t *testing.T) {
	payload := []byte(`{"Name":"COINBASE (different one)","Hint":"Do not hardcode these numbers, pull them from the test cases","Code":{"Bin":"41","Asm":"COINBASE"},"Expect":{"Stack":["0x888"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 102, payload)
}

 
func Test_103_Timestamp(t *testing.T) {
	payload := []byte(`{"Name":"TIMESTAMP","Hint":"Solidity calls this block.timestamp","Code":{"Bin":"42","Asm":"TIMESTAMP"},"Expect":{"Stack":["0xe4e1c1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 103, payload)
}

 
func Test_104_Number(t *testing.T) {
	payload := []byte(`{"Name":"NUMBER","Hint":"Solidity calls this block.number","Code":{"Bin":"43","Asm":"NUMBER"},"Expect":{"Stack":["0x1000001"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 104, payload)
}

 
func Test_105_Difficulty(t *testing.T) {
	payload := []byte(`{"Name":"DIFFICULTY","Hint":"Also known as PREVRANDAO, not used in these test cases yet","Code":{"Bin":"44","Asm":"DIFFICULTY"},"Expect":{"Stack":["0x20000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 105, payload)
}

 
func Test_106_Gaslimit(t *testing.T) {
	payload := []byte(`{"Name":"GASLIMIT","Hint":"","Code":{"Bin":"45","Asm":"GASLIMIT"},"Expect":{"Stack":["0xffffffffffff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 106, payload)
}

 
func Test_107_Chainid(t *testing.T) {
	payload := []byte(`{"Name":"CHAINID","Hint":"","Code":{"Bin":"46","Asm":"CHAINID"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 107, payload)
}

 
func Test_108_Blockhash(t *testing.T) {
	payload := []byte(`{"Name":"BLOCKHASH","Hint":"Not used in this test suite, can return 0","Code":{"Bin":"600040","Asm":"PUSH1 0\nBLOCKHASH"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 108, payload)
}

 
func Test_109_Balance(t *testing.T) {
	payload := []byte(`{"Name":"BALANCE","Hint":"Read \"State\" section of the course learning materials. Modify your evm function to take state as one of the arguments, or turn it into a class","Code":{"Bin":"731e79b045dc29eae9fdc69673c9dcd7c53e5e159d31","Asm":"PUSH20 0x1e79b045dc29eae9fdc69673c9dcd7c53e5e159d\nBALANCE"},"Expect":{"Stack":["0x100"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 109, payload)
}

 
func Test_110_BalanceEmpty(t *testing.T) {
	payload := []byte(`{"Name":"BALANCE (empty)","Hint":"Balance of accounts not present in state is zero","Code":{"Bin":"73af69610ea9ddc95883f97a6a3171d52165b69b0331","Asm":"PUSH20 0xaf69610ea9ddc95883f97a6a3171d52165b69b03\nBALANCE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 110, payload)
}

 
func Test_111_Callvalue(t *testing.T) {
	payload := []byte(`{"Name":"CALLVALUE","Hint":"Read \"Calls\" section of the course learning materials. Solidity calls this msg.value, it is amount of wei sent as part of this transaction","Code":{"Bin":"34","Asm":"CALLVALUE"},"Expect":{"Stack":["0x1000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 111, payload)
}

 
func Test_112_Calldataload(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATALOAD","Hint":"Read \"Calls\" section of the course learning materials. Calldata is an array of bytes sent to the evm function","Code":{"Bin":"600035","Asm":"PUSH1 0\nCALLDATALOAD"},"Expect":{"Stack":["0x102030405060708090a0b0c0d0e0f00112233445566778899aabbccddeeff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 112, payload)
}

 
func Test_113_CalldataloadTail(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATALOAD (tail)","Hint":"Overflow bytes filled with zeros","Code":{"Bin":"601f35","Asm":"PUSH1 31\nCALLDATALOAD"},"Expect":{"Stack":["0xff00000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 113, payload)
}

 
func Test_114_Calldatasize(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATASIZE","Hint":"Size (in bytes) of calldata buffer","Code":{"Bin":"36","Asm":"CALLDATASIZE"},"Expect":{"Stack":["0x20"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 114, payload)
}

 
func Test_115_CalldatasizeNoData(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATASIZE (no data)","Hint":"","Code":{"Bin":"36","Asm":"CALLDATASIZE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 115, payload)
}

 
func Test_116_Calldatacopy(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATACOPY","Hint":"Copy 32-byte chunk of calldata into memory. Do not forget to update MSIZE after CALLDATACOPY","Code":{"Bin":"60206000600037600051","Asm":"PUSH1 32\nPUSH1 0\nPUSH1 0\nCALLDATACOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x102030405060708090a0b0c0d0e0f00112233445566778899aabbccddeeff"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 116, payload)
}

 
func Test_117_CalldatacopyTail(t *testing.T) {
	payload := []byte(`{"Name":"CALLDATACOPY (tail)","Hint":"Overflow bytes filled with zeros","Code":{"Bin":"6001601f600037600051","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nCALLDATACOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0xff00000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 117, payload)
}

 
func Test_118_CodesizeSmall(t *testing.T) {
	payload := []byte(`{"Name":"CODESIZE (small)","Hint":"Size of the bytecode running in the current context","Code":{"Bin":"38","Asm":"CODESIZE"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 118, payload)
}

 
func Test_119_Codesize(t *testing.T) {
	payload := []byte(`{"Name":"CODESIZE","Hint":"","Code":{"Bin":"7300000000000000000000000000000000000000005038","Asm":"PUSH20 0\nPOP\nCODESIZE"},"Expect":{"Stack":["0x17"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 119, payload)
}

 
func Test_120_Codecopy(t *testing.T) {
	payload := []byte(`{"Name":"CODECOPY","Hint":"Copy your own code into memory. Implementing quines in EVM is really easy","Code":{"Bin":"60206000600039600051","Asm":"PUSH1 32\nPUSH1 0\nPUSH1 0\nCODECOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x6020600060003960005100000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 120, payload)
}

 
func Test_121_CodecopyTail(t *testing.T) {
	payload := []byte(`{"Name":"CODECOPY (tail)","Hint":"Overflow bytes filled with zeros","Code":{"Bin":"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5060026020600039600051","Asm":"PUSH32 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff\nPOP\nPUSH1 2\nPUSH1 32\nPUSH1 0\nCODECOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0xff50000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 121, payload)
}

 
func Test_122_ExtcodesizeEmpty(t *testing.T) {
	payload := []byte(`{"Name":"EXTCODESIZE (empty)","Hint":"","Code":{"Bin":"731e79b045dc29eae9fdc69673c9dcd7c53e5e159d3b","Asm":"PUSH20 0x1e79b045dc29eae9fdc69673c9dcd7c53e5e159d\nEXTCODESIZE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 122, payload)
}

 
func Test_123_Extcodesize(t *testing.T) {
	payload := []byte(`{"Name":"EXTCODESIZE","Hint":"Read \"State\" section of the course learning materials","Code":{"Bin":"731000000000000000000000000000000000000aaa3b","Asm":"PUSH20 0x1000000000000000000000000000000000000aaa\nEXTCODESIZE"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 123, payload)
}

 
func Test_124_Extcodecopy(t *testing.T) {
	payload := []byte(`{"Name":"EXTCODECOPY","Hint":"","Code":{"Bin":"602060006000731000000000000000000000000000000000000aaa3c600051","Asm":"PUSH1 32\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000aaa\nEXTCODECOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x6001000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 124, payload)
}

 
func Test_125_Extcodehash(t *testing.T) {
	payload := []byte(`{"Name":"EXTCODEHASH","Hint":"Use the same library you used for SHA3 opcode","Code":{"Bin":"731000000000000000000000000000000000000aaa3f","Asm":"PUSH20 0x1000000000000000000000000000000000000aaa\nEXTCODEHASH"},"Expect":{"Stack":["0x29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c238"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 125, payload)
}

 
func Test_126_ExtcodehashEmpty(t *testing.T) {
	payload := []byte(`{"Name":"EXTCODEHASH (empty)","Hint":"","Code":{"Bin":"731000000000000000000000000000000000000aaa3f","Asm":"PUSH20 0x1000000000000000000000000000000000000aaa\nEXTCODEHASH"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 126, payload)
}

 
func Test_127_Selfbalance(t *testing.T) {
	payload := []byte(`{"Name":"SELFBALANCE","Hint":"","Code":{"Bin":"47","Asm":"SELFBALANCE"},"Expect":{"Stack":["0x200"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 127, payload)
}

 
func Test_128_Sstore(t *testing.T) {
	payload := []byte(`{"Name":"SSTORE","Hint":"Read \"Storage\" section of the course learning materials","Code":{"Bin":"6001600055600054","Asm":"PUSH1 1\nPUSH1 0\nSSTORE\nPUSH1 0\nSLOAD"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 128, payload)
}

 
func Test_129_SstoreNonZeroLocation(t *testing.T) {
	payload := []byte(`{"Name":"SSTORE (non-zero location)","Hint":"","Code":{"Bin":"60026398fe5c2c556398fe5c2c54","Asm":"PUSH1 2\nPUSH4 0x98fe5c2c\nSSTORE\nPUSH4 0x98fe5c2c\nSLOAD"},"Expect":{"Stack":["0x2"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 129, payload)
}

 
func Test_130_SloadEmpty(t *testing.T) {
	payload := []byte(`{"Name":"SLOAD (empty)","Hint":"All storage is initialized to zeros","Code":{"Bin":"60ff54","Asm":"PUSH1 0xff\nSLOAD"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 130, payload)
}

 
func Test_131_Log0(t *testing.T) {
	payload := []byte(`{"Name":"LOG0","Hint":"Make evm function return array of logs, modify the testing code to assert that the logs match","Code":{"Bin":"60aa6000526001601fa0","Asm":"PUSH1 0xaa\nPUSH1 0\nMSTORE\nPUSH1 1\nPUSH1 31\nLOG0"},"Expect":{"Stack":null,"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 131, payload)
}

 
func Test_132_Log1(t *testing.T) {
	payload := []byte(`{"Name":"LOG1","Hint":"Make evm function return array of logs, modify the testing code to assert that the logs match","Code":{"Bin":"60bb6000527f11111111111111111111111111111111111111111111111111111111111111116001601fa1","Asm":"PUSH1 0xbb\nPUSH1 0\nMSTORE\nPUSH32 0x1111111111111111111111111111111111111111111111111111111111111111\nPUSH1 1\nPUSH1 31\nLOG1"},"Expect":{"Stack":null,"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 132, payload)
}

 
func Test_133_Log2(t *testing.T) {
	payload := []byte(`{"Name":"LOG2","Hint":"Use the same code to handle LOG1...LOG4 opcodes","Code":{"Bin":"60cc6000527f11111111111111111111111111111111111111111111111111111111111111117f22222222222222222222222222222222222222222222222222222222222222226001601fa2","Asm":"PUSH1 0xcc\nPUSH1 0\nMSTORE\nPUSH32 0x1111111111111111111111111111111111111111111111111111111111111111\nPUSH32 0x2222222222222222222222222222222222222222222222222222222222222222\nPUSH1 1\nPUSH1 31\nLOG2"},"Expect":{"Stack":null,"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 133, payload)
}

 
func Test_134_Log3(t *testing.T) {
	payload := []byte(`{"Name":"LOG3","Hint":"N = OPCODE - LOG0, pop N items from the stack as topics","Code":{"Bin":"60dd6000527f11111111111111111111111111111111111111111111111111111111111111117f22222222222222222222222222222222222222222222222222222222222222227f33333333333333333333333333333333333333333333333333333333333333336001601fa3","Asm":"PUSH1 0xdd\nPUSH1 0\nMSTORE\nPUSH32 0x1111111111111111111111111111111111111111111111111111111111111111\nPUSH32 0x2222222222222222222222222222222222222222222222222222222222222222\nPUSH32 0x3333333333333333333333333333333333333333333333333333333333333333\nPUSH1 1\nPUSH1 31\nLOG3"},"Expect":{"Stack":null,"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 134, payload)
}

 
func Test_135_Log4(t *testing.T) {
	payload := []byte(`{"Name":"LOG4","Hint":"Refactoring code is always a good idea. Your code will become cleaner, and the tests will catch if something breaks","Code":{"Bin":"60ee6000527f11111111111111111111111111111111111111111111111111111111111111117f22222222222222222222222222222222222222222222222222222222222222227f33333333333333333333333333333333333333333333333333333333333333337f44444444444444444444444444444444444444444444444444444444444444446001601fa4","Asm":"PUSH1 0xee\nPUSH1 0\nMSTORE\nPUSH32 0x1111111111111111111111111111111111111111111111111111111111111111\nPUSH32 0x2222222222222222222222222222222222222222222222222222222222222222\nPUSH32 0x3333333333333333333333333333333333333333333333333333333333333333\nPUSH32 0x4444444444444444444444444444444444444444444444444444444444444444\nPUSH1 1\nPUSH1 31\nLOG4"},"Expect":{"Stack":null,"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 135, payload)
}

 
func Test_136_Return(t *testing.T) {
	payload := []byte(`{"Name":"RETURN","Hint":"Read \"Calls and Returns\" section of the course learning materials","Code":{"Bin":"60a26000526001601ff3","Asm":"PUSH1 0xA2\nPUSH1 0\nMSTORE\nPUSH1 1\nPUSH1 31\nRETURN"},"Expect":{"Stack":null,"Success":true,"Return":"a2"},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 136, payload)
}

 
func Test_137_Revert(t *testing.T) {
	payload := []byte(`{"Name":"REVERT","Hint":"Note that this test expects 'success' to be false","Code":{"Bin":"60f16000526001601ffd","Asm":"PUSH1 0xF1\nPUSH1 0\nMSTORE\nPUSH1 1\nPUSH1 31\nREVERT"},"Expect":{"Stack":null,"Success":false,"Return":"f1"},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 137, payload)
}

 
func Test_138_Call(t *testing.T) {
	payload := []byte(`{"Name":"CALL","Hint":"Read \"Calls and Returns\" section of the course learning materials. Recursively call evm function from itself when handing this opcode","Code":{"Bin":"6001601f600060006000731000000000000000000000000000000000000c426000f1600051","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nCALL\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x42","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 138, payload)
}

 
func Test_139_CallReturnsAddress(t *testing.T) {
	payload := []byte(`{"Name":"CALL (returns address)","Hint":"In the inner context, the CALLER is the contract we are sending the initial transaction to","Code":{"Bin":"60206000600060006000731000000000000000000000000000000000000c426000f1600051","Asm":"PUSH1 32\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nCALL\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x1000000000000000000000000000000000000aaa","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 139, payload)
}

 
func Test_140_CallReverts(t *testing.T) {
	payload := []byte(`{"Name":"CALL (reverts)","Hint":"Reverts can also return data","Code":{"Bin":"6001601f600060006000731000000000000000000000000000000000000c426000f1600051","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nCALL\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x42","0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 140, payload)
}

 
func Test_141_ReturndatasizeEmpty(t *testing.T) {
	payload := []byte(`{"Name":"RETURNDATASIZE (empty)","Hint":"","Code":{"Bin":"3d","Asm":"RETURNDATASIZE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 141, payload)
}

 
func Test_142_Returndatasize(t *testing.T) {
	payload := []byte(`{"Name":"RETURNDATASIZE","Hint":"","Code":{"Bin":"60006000600060006000731000000000000000000000000000000000000c426000f1503d","Asm":"PUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nCALL\nPOP\nRETURNDATASIZE"},"Expect":{"Stack":["0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 142, payload)
}

 
func Test_143_Returndatacopy(t *testing.T) {
	payload := []byte(`{"Name":"RETURNDATACOPY","Hint":"","Code":{"Bin":"6001601f600060006000731000000000000000000000000000000000000c426000f1506001600060ff3e60ff51","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nCALL\nPOP\nPUSH1 1\nPUSH1 0\nPUSH1 0xff\nRETURNDATACOPY\nPUSH1 0xff\nMLOAD"},"Expect":{"Stack":["0x4200000000000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 143, payload)
}

 
func Test_144_Delegatecall(t *testing.T) {
	payload := []byte(`{"Name":"DELEGATECALL","Hint":"Like CALL, but keep the transaction data (from, origin, address) and use the code from the other account","Code":{"Bin":"600080808073dddddddddddddddddddddddddddddddddddddddd5af4600054","Asm":"PUSH1 0\nDUP1\nDUP1\nDUP1\nPUSH20 0xdddddddddddddddddddddddddddddddddddddddd\nGAS\nDELEGATECALL\nPUSH1 0\nSLOAD"},"Expect":{"Stack":["0x1000000000000000000000000000000000000aaa","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 144, payload)
}

 
func Test_145_Staticcall(t *testing.T) {
	payload := []byte(`{"Name":"STATICCALL","Hint":"Like CALL, but disable state modifications","Code":{"Bin":"6001601f60006000731000000000000000000000000000000000000c426000fa600051","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nSTATICCALL\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0x42","0x1"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 145, payload)
}

 
func Test_146_StaticcallRevertsOnWrite(t *testing.T) {
	payload := []byte(`{"Name":"STATICCALL (reverts on write)","Hint":"Use a flag to tell the evm function whenever the context is writeable (CALL) or not (STATICCALL)","Code":{"Bin":"6001601f60006000731000000000000000000000000000000000000c426000fa","Asm":"PUSH1 1\nPUSH1 31\nPUSH1 0\nPUSH1 0\nPUSH20 0x1000000000000000000000000000000000000c42\nPUSH1 0\nSTATICCALL"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 146, payload)
}

 
func Test_147_CreateEmpty(t *testing.T) {
	payload := []byte(`{"Name":"CREATE (empty)","Hint":"Read \"Creating new contracts\" section of the course learning materials. This code creates a new empty account with balance 9","Code":{"Bin":"600060006009f031","Asm":"PUSH1 0\nPUSH1 0\nPUSH1 9\nCREATE\nBALANCE"},"Expect":{"Stack":["0x9"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 147, payload)
}

 
func Test_148_CreateWith4xFf(t *testing.T) {
	payload := []byte(`{"Name":"CREATE (with 4x FF)","Hint":"Read \"Creating new contracts\" section of the course learning materials. CALL with the given code, store the returned bytes as new contracts bytecode","Code":{"Bin":"6020600060006c63ffffffff6000526004601cf3600052600d60136000f03c600051","Asm":"PUSH1 32\nPUSH1 0\nPUSH1 0\nPUSH13 0x63FFFFFFFF6000526004601CF3\nPUSH1 0\nMSTORE\nPUSH1 13\nPUSH1 19\nPUSH1 0\nCREATE\nEXTCODECOPY\nPUSH1 0\nMLOAD"},"Expect":{"Stack":["0xffffffff00000000000000000000000000000000000000000000000000000000"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 148, payload)
}

 
func Test_149_CreateReverts(t *testing.T) {
	payload := []byte(`{"Name":"CREATE (reverts)","Hint":"No address when constructor code reverts","Code":{"Bin":"6c63ffffffff6000526004601cfd600052600d60136000f0","Asm":"PUSH13 0x63FFFFFFFF6000526004601CFD\nPUSH1 0\nMSTORE\nPUSH1 13\nPUSH1 19\nPUSH1 0\nCREATE"},"Expect":{"Stack":["0x0"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 149, payload)
}

 
func Test_150_Selfdestruct(t *testing.T) {
	payload := []byte(`{"Name":"SELFDESTRUCT","Hint":"Note that for simplicity, this opcode should delete the account from the state. In the real EVM this happens only after the transaction has been processed, but that would overcomplicate these tests","Code":{"Bin":"60008080808073dead00000000000000000000000000000000dead5af15073a1c300000000000000000000000000000000a1c33173dead00000000000000000000000000000000dead3b","Asm":"PUSH1 0\nDUP1\nDUP1\nDUP1\nDUP1\nPUSH20 0xdead00000000000000000000000000000000dead\nGAS\nCALL\nPOP\nPUSH20 0xa1c300000000000000000000000000000000a1c3\nBALANCE\nPUSH20 0xdead00000000000000000000000000000000dead\nEXTCODESIZE"},"Expect":{"Stack":["0x0","0x7"],"Success":true,"Return":""},"FnName":"","Index":0,"Payload":""}`)
	runTest(t, 150, payload)
}



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

