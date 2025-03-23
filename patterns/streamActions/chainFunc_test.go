package streamActions

import (
	"fmt"
	"testing"
)

func Test_ChainFrom(t *testing.T) {
	From([]int{1, 2, 3, 4}).Each(func(a int) { println(a) })
}

func Test_Reverse(t *testing.T) {
	From([]int{1, 2, 3, 4}).Reverse().Each(func(a int) { println(a) })
}

func Test_IteratorStructs(t *testing.T) {
	type User struct {
		Id   int
		Name string
		Hash int
	}

	users := []User{
		{0, "abc", 0},
		{1, "bbbsdfas", 0},
		{2, "cadsd", 0},
	}

	From(users).
		Map(func(u User) User {
			h := 0
			for i, r := range u.Name {
				h += int(r)*31 ^ (len(u.Name) - i - 1)
			}
			u.Hash = h
			return u
		}).
		Each(func(u User) {
			fmt.Printf("%#+v\n", u)
		})
}
