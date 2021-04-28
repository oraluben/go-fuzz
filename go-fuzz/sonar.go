// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"sync"

	. "github.com/oraluben/go-fuzz/go-fuzz-defs"
)

type SonarSite struct {
	id  int    // unique site id (const)
	loc string // file:line.pos,line.pos (const)
	sync.Mutex
	dynamic    bool   // both operands are not constant
	takenFuzz  [2]int // number of times condition evaluated to false/true during fuzzing
	takenTotal [2]int // number of times condition evaluated to false/true in total
	val        [2][]byte
}

type SonarSample struct {
	site  *SonarSite
	flags byte
	val   [2][]byte
}

func (w *Worker) parseSonarData(sonar []byte) (res []SonarSample) {
	ro := w.hub.ro.Load().(*ROData)
	sonar = makeCopy(sonar)
	for len(sonar) > SonarHdrLen {
		id := binary.LittleEndian.Uint32(sonar)
		flags := byte(id)
		id >>= 8
		n1 := sonar[4]
		n2 := sonar[5]
		sonar = sonar[SonarHdrLen:]
		if n1 > SonarMaxLen || n2 > SonarMaxLen || len(sonar) < int(n1)+int(n2) {
			log.Fatalf("corrupted sonar data: hdr=[%v/%v/%v] data=%v", flags, n1, n2, len(sonar))
		}
		v1 := makeCopy(sonar[:n1])
		v2 := makeCopy(sonar[n1 : n1+n2])
		sonar = sonar[n1+n2:]

		// Trim trailing 0x00 and 0xff bytes (we don't know exact size of operands).
		if flags&SonarString == 0 {
			for len(v1) > 0 || len(v2) > 0 {
				i := len(v1) - 1
				if len(v2) > len(v1) {
					i = len(v2) - 1
				}
				var c1, c2 byte
				if i < len(v1) {
					c1 = v1[i]
				}
				if i < len(v2) {
					c2 = v2[i]
				}
				if (c1 == 0 || c1 == 0xff) && (c2 == 0 || c2 == 0xff) {
					if i < len(v1) {
						v1 = v1[:i]
					}
					if i < len(v2) {
						v2 = v2[:i]
					}
					continue
				}
				break
			}
		}

		res = append(res, SonarSample{&ro.sonarSites[id], flags, [2][]byte{v1, v2}})
	}
	return res
}

func (w *Worker) processSonarData(data SqlWrap, sonar []byte, depth int, smash bool) {
	// Sonar is not necessary, for we guarantee every input is valid.
	// However I'm not 100% sure that my understanding of sonar is correct.
	panic("sonar disabled")
}

var dumpMu sync.Mutex

func (site *SonarSite) update(sam SonarSample, smash, resb bool) (updated, skip bool) {
	res := 0
	if resb {
		res = 1
	}
	site.Lock()
	defer site.Unlock()
	v1 := sam.val[0]
	v2 := sam.val[1]
	if !site.dynamic && sam.flags&SonarConst1+sam.flags&SonarConst2 == 0 {
		if site.val[0] == nil {
			site.val[0] = makeCopy(v1)
		}
		if site.val[1] == nil {
			site.val[1] = makeCopy(v2)
		}
		if !bytes.Equal(site.val[0], v1) && !bytes.Equal(site.val[1], v2) {
			site.val[0] = nil
			site.val[1] = nil
			site.dynamic = true
		}
	}
	if site.takenTotal[res] == 0 {
		updated = true
	}
	site.takenTotal[res]++
	if !smash {
		site.takenFuzz[res]++
	}
	if !site.dynamic && !smash {
		// Skip this site if it has at least one const operand
		// and is taken both ways enough times.
		// Check sites that don't have const operands always,
		// because it can be a CRC-like verification, which
		// won't be cracked otherwise.
		if site.takenFuzz[0] > 10 && site.takenFuzz[1] > 10 || site.takenFuzz[0]+site.takenFuzz[1] > 100 {
			skip = true
			return
		}
	}
	return
}

func (sam *SonarSample) evaluate() bool {
	v1 := sam.val[0]
	v2 := sam.val[1]
	if sam.flags&SonarString != 0 {
		s1 := string(v1)
		s2 := string(v2)
		switch sam.flags & SonarOpMask {
		case SonarEQL:
			return s1 == s2
		case SonarNEQ:
			return s1 != s2
		case SonarLSS:
			return s1 < s2
		case SonarGTR:
			return s1 > s2
		case SonarLEQ:
			return s1 <= s2
		case SonarGEQ:
			return s1 >= s2
		default:
			panic("bad")
		}
	}
	if len(v1) == 0 || len(v2) == 0 || len(v1) > 8 || len(v2) > 8 || len(v1) != len(v2) {
		return false
	}
	v1 = makeCopy(v1)
	for len(v1) < 8 {
		if int8(v1[len(v1)-1]) >= 0 {
			v1 = append(v1, 0)
		} else {
			v1 = append(v1, 0xff)
		}
	}
	v2 = makeCopy(v2)
	for len(v2) < 8 {
		if int8(v2[len(v2)-1]) >= 0 {
			v2 = append(v2, 0)
		} else {
			v2 = append(v2, 0xff)
		}
	}
	// Note: assuming le machine.
	if sam.flags&SonarSigned == 0 {
		s1 := binary.LittleEndian.Uint64(v1)
		s2 := binary.LittleEndian.Uint64(v2)
		switch sam.flags & SonarOpMask {
		case SonarEQL:
			return s1 == s2
		case SonarNEQ:
			return s1 != s2
		case SonarLSS:
			return s1 < s2
		case SonarGTR:
			return s1 > s2
		case SonarLEQ:
			return s1 <= s2
		case SonarGEQ:
			return s1 >= s2
		default:
			panic("bad")
		}
	} else {
		s1 := int64(binary.LittleEndian.Uint64(v1))
		s2 := int64(binary.LittleEndian.Uint64(v2))
		switch sam.flags & SonarOpMask {
		case SonarEQL:
			return s1 == s2
		case SonarNEQ:
			return s1 != s2
		case SonarLSS:
			return s1 < s2
		case SonarGTR:
			return s1 > s2
		case SonarLEQ:
			return s1 <= s2
		case SonarGEQ:
			return s1 >= s2
		default:
			panic("bad")
		}
	}
}

func dumpSonarData(site *SonarSite, flags byte, v1, v2 []byte) {
	// Debug output.
	op := ""
	switch flags & SonarOpMask {
	case SonarEQL:
		op = "=="
	case SonarNEQ:
		op = "!="
	case SonarLSS:
		op = "<"
	case SonarGTR:
		op = ">"
	case SonarLEQ:
		op = "<="
	case SonarGEQ:
		op = ">="
	default:
		log.Fatalf("bad")
	}
	sign := ""
	if flags&SonarSigned != 0 {
		sign = "(signed)"
	}
	isstr := ""
	if flags&SonarString != 0 {
		isstr = "(string)"
	}
	const1 := ""
	if flags&SonarConst1 != 0 {
		const1 = "c"
	}
	const2 := ""
	if flags&SonarConst2 != 0 {
		const2 = "c"
	}
	log.Printf("SONAR %v%v %v %v%v %v%v %v",
		hex.EncodeToString(v1), const1, op, hex.EncodeToString(v2), const2, sign, isstr, site.loc)
}
