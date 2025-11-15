package pref

import (
	"fmt"
	"sort"
	"testing"
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

// Not comment out -> causing compile issue
// 重点在: print.go 的672行 -> case Stringer:
//				handled = true
//				defer p.catchPanic(p.arg, verb, "String")  <= 'safety net'
//				p.fmtString(v.String(), verb)
//				return
//			}
// 这里的 catchPanic 是作为保护网, 当出现OOM堆栈溢出等问题直接通过该方法拦截
// 这里通过使用defer + catchPanic + recover来抓取异常, 一般lib才会这样使用
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

func BenchmarkSequenceComparison(b *testing.B) {
	b.Run("Slow", func(b *testing.B) {
		s := SequenceSlow{5, 3, 8, 1, 9, 2, 7, 4, 6}
		for b.Loop() {
			_ = s.String()
		}
	})
	b.Run("Normal", func(b *testing.B) {
		s := Sequence{5, 3, 8, 1, 9, 2, 7, 4, 6}
		for b.Loop() {
			_ = s.String()
		}
	})
}
