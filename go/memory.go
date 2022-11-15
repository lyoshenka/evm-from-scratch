package main

import (
	"github.com/holiman/uint256"
)

type Memory struct {
	data []byte
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Put(offset uint64, val *uint256.Int) {
	m.expandIfNeeded(offset, 32)
	val.WriteToSlice(m.data[offset : offset+32])
}

func (m *Memory) PutByte(offset uint64, val byte) {
	m.expandIfNeeded(offset, 1)
	m.data[offset] = val
}

func (m *Memory) Get(offset uint64) *uint256.Int {
	m.expandIfNeeded(offset, 32)
	return uint256.NewInt(0).SetBytes(m.data[offset : offset+32])
}

func (m *Memory) expandIfNeeded(offset, size uint64) {
	lastByte := offset + size
	if lastByte%32 != 0 {
		lastByte += 32 - (lastByte % 32)
	}
	if uint64(len(m.data)) < lastByte {
		m.data = append(m.data, make([]byte, lastByte-uint64(len(m.data)))...)
	}
}

func (m *Memory) Len() uint64 {
	return uint64(len(m.data))
}
