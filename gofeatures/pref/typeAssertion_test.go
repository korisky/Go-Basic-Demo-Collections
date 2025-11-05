package pref

import (
	"fmt"
	"sort"
)

type Sequence []int
type SequenceBad []int
type SequenceSlow []int

func (s Sequence) Len() int           { return len(s) }
func (s Sequence) Less(i, j int) bool { return s[i] < s[j] }
func (s Sequence) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SequenceBad) Len() int           { return len(s) }
func (s SequenceBad) Less(i, j int) bool { return s[i] < s[j] }
func (s SequenceBad) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SequenceSlow) Len() int           { return len(s) }
func (s SequenceSlow) Less(i, j int) bool { return s[i] < s[j] }
func (s SequenceSlow) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Sequence) Copy() Sequence {
	c := make(Sequence, len(s))
	copy(c, s)
	return c
}

func (s SequenceBad) CopyBad() SequenceBad {
	c := make(SequenceBad, len(s))
	copy(c, s)
	return c
}

func (s SequenceSlow) CopySlow() SequenceSlow {
	c := make(SequenceSlow, len(s))
	copy(c, s)
	return c
}

// String -> correct impl, by converting to []int(s)
func (s Sequence) String() string {
	s = s.Copy()
	sort.Sort(s)
	return fmt.Sprint([]int(s))
}

//func (s SequenceBad) String() string {
//	s = s.CopyBad()
//	sort.Sort(s)
//	return fmt.Sprint(s) // keep calling String() again -> infinite recursion
//}

func (s SequenceSlow) String() string {
	s = s.CopySlow()
	sort.Sort(s)
	result := "["
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprint(v)
	}
	return result + "]"
}
