// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"sort"

	. "github.com/oraluben/go-fuzz/go-fuzz-defs"
	"github.com/oraluben/go-fuzz/go-fuzz/internal/pcg"
)

type Mutator struct {
	r *pcg.Rand
}

func newMutator() *Mutator {
	return &Mutator{r: pcg.New()}
}

func (m *Mutator) rand(n int) int {
	return m.r.Intn(n)
}

// randbig generates a number in [0, 2³⁰).
func (m *Mutator) randbig() int64 {
	return int64(m.r.Uint32() >> 2)
}

func (m *Mutator) randByteOrder() binary.ByteOrder {
	if m.r.Bool() {
		return binary.LittleEndian
	}
	return binary.BigEndian
}

func (m *Mutator) generate(ro *ROData) (SqlWrap, int) {
	corpus := ro.corpus
	scoreSum := corpus[len(corpus)-1].runningScoreSum
	weightedIdx := m.rand(scoreSum)
	idx := sort.Search(len(corpus), func(i int) bool {
		return corpus[i].runningScoreSum > weightedIdx
	})
	input := &corpus[idx]
	return m.mutate(input.data, ro), input.depth + 1
}

func (m *Mutator) mutate(data SqlWrap, ro *ROData) SqlWrap {
	_ = ro.corpus
	res := data.copy()

	for {
		dml := data.getDML()
		r, err := ro.mutateConfig.Mutate(dml)
		if err != nil && r == "" {
			// we should not ignore this error because it implies there is some bug in go-squirrel
			log.Printf("error while mutating %v : %v", dml, err)
			continue
		}
		if err != nil || r == "" {
			continue
		}
		res.setDML(r)
		break
	}
	for res.len() > MaxInputSize {
		// todo: implement trunk
		panic(fmt.Sprintf("mutation result too loong: %d", res.len()))
	}
	return res
}

// chooseLen chooses length of range mutation.
// It gives preference to shorter ranges.
func (m *Mutator) chooseLen(n int) int {
	switch x := m.rand(100); {
	case x < 90:
		return m.rand(min(8, n)) + 1
	case x < 99:
		return m.rand(min(32, n)) + 1
	default:
		return m.rand(n) + 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	interesting8  = []int8{-128, -1, 0, 1, 16, 32, 64, 100, 127}
	interesting16 = []int16{-32768, -129, 128, 255, 256, 512, 1000, 1024, 4096, 32767}
	interesting32 = []int32{-2147483648, -100663046, -32769, 32768, 65535, 65536, 100663045, 2147483647}
)

func init() {
	for _, v := range interesting8 {
		interesting16 = append(interesting16, int16(v))
	}
	for _, v := range interesting16 {
		interesting32 = append(interesting32, int32(v))
	}
}
