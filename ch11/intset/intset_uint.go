//copy here for test
//use uint32 type, for e6.5
package intset

import (
	"bytes"
	"fmt"
)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint32
}

const WordSize = 32 << (^uint32(0) >> 63) //32 or 64
// Has reports whether the set contains the non-negative value x.
//get panic if x < 0,
//painic, if s is nil or not initial pointer
func (s *IntSet) Has(x int) bool {
	// if x < 0 {
	// 	return false
	// }
	word, bit := x/WordSize, uint32(x%WordSize)
	return word < len(s.words) &&
		s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	// if x < 0 {
	// 	return
	// }
	word, bit := x/WordSize, uint32(x%WordSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < WordSize; j++ {
			if word&(1<<uint32(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", WordSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// !-string

func (s *IntSet) Len() (len int) {
	for _, word := range s.words {
		c := 0
		for word > 0 {
			word &= word - 1
			c++
		}
		len += c
	}
	return
}

func (s *IntSet) Remove(x int)  {
	word, bit := x/WordSize, uint32(x % WordSize)
	if word < len(s.words){
		s.words[word] &^= (1 << bit)
	}
}

func (s *IntSet) AddAll(nums ...int)  {
	for _, x := range nums {
		s.Add(x)
	}
}

// func NewSet(n int) *IntSet {
// //
// }

func (s *IntSet) Copy() *IntSet {
	t := new(IntSet)
	t.words = make([]uint32, len(s.words))
	copy(t.words, s.words)
	return t
}

func (s *IntSet) Clear() {
	s.words = nil
}

func (s * IntSet) Elems() (elems []int) {
	for i, word := range s.words {
		for bit := 0; bit < WordSize; bit++ {
			if word & uint32(1 << uint32(bit)) != 0 {
				elems = append(elems, i * WordSize + bit)
			}
		}
	}
	return
}

func (s *IntSet) IntersectWith(t *IntSet)  {
	// for _, v := range s.Elems() {
	// 	if !t.Has(v) {
	// 		s.Remove(v)
	// 	}
	// }
	for i := 0; i < len(s.words); i++ {
		if i < len(t.words) {
			s.words[i] &= t.words[i]
		}else {
			s.words[i] = 0
		}
	}

}

func (s *IntSet) DifferenceWith(t *IntSet)  {
	// for _, v := range s.Elems() {
	// 	if t.Has(v) {
	// 		s.Remove(v)
	// 	}
	// }

	for i, word := range t.words {
		if i < len(s.words) {
			s.words[i] &^= word
		}
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet)  {
	// tc := t.Copy()
	// tc.DifferenceWith(s)
	// s.DifferenceWith(t)
	// s.UnionWith(tc)

	for i, word := range t.words {
		if i < len(s.words) {
			s.words[i] ^= word
		} else{
			s.words = append(s.words, word)
		}
	}
}

func (s *IntSet) Equal(t *IntSet) bool {

	if len(s.words) != len(t.words) {
		return false
	}
	for i, word := range t.words {
		if s.words[i] != word {
			return false
		}
	}
	return true
}
