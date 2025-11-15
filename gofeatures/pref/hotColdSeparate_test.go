package pref

import (
	"fmt"
	"sort"
	"testing"
)

// User struct足够大时, 拆分冷热数据就变得非常有必要
type User struct {
	ID      uint64 // hot field, 8 bytes
	Score   int    // hot field, 8 bytes
	Name    string // cold field, 16bytes, headers(pointer8bytes + length8bytes) only
	Email   string
	Address string
	Bio     string
	extra1  string
	extra2  string
	extra3  string
	extra4  string
	extra5  string
	extra6  string
	extra7  string
	extra8  string
	extra9  string
	extra10 string
	extra11 string
	extra12 string
	extra13 string
	extra14 string
	extra15 string
	extra16 string
	extra17 string
	extra18 string
	extra19 string
	extra20 string
}

// TopUsersMixed 获取排名前n的User的ID信息
func TopUsersMixed(users []User, n int) []uint64 {
	// 对slice进行排序
	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score
	})
	// remain top10
	result := make([]uint64, n)
	for i := 0; i < n && i < len(users); i++ {
		result[i] = users[i].ID
	}
	return result
}

// GetUserDetailsMixed 配合TopUsersMixed, 当需要某个User的所有数据(包括冷数据)
func GetUserDetailsMixed(users []User, id uint64) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

// UserHot 将冷数据直接抽取另一个struct, 并在热数据中带一个对应冷数据的Ref
type UserHot struct {
	ID    uint64    // hot field, 8 bytes
	Score int       // hot field, 8 bytes
	Cold  *UserCold // cold fields
}

type UserCold struct {
	Name    string // for string in golang, 8bytes of pointer, 8bytes about length
	Email   string // so -> 16bytes per string in a struct
	Address string
	Bio     string
	extra1  string
	extra2  string
	extra3  string
	extra4  string
	extra5  string
	extra6  string
	extra7  string
	extra8  string
	extra9  string
	extra10 string
	extra11 string
	extra12 string
	extra13 string
	extra14 string
	extra15 string
	extra16 string
	extra17 string
	extra18 string
	extra19 string
	extra20 string
}

// TopUsersSeparated 由于直接操作热数据, 而每个struct包含的冷数据都不需要load到内存
func TopUsersSeparated(users []UserHot, n int) []uint64 {
	// 对slice进行排序
	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score
	})
	// remain top10
	result := make([]uint64, n)
	for i := 0; i < n && i < len(users); i++ {
		result[i] = users[i].ID
	}
	return result
}

func GetUserDetailsSeparated(users []UserHot, id uint64) (*UserHot, *UserCold) {
	for i := range users {
		if users[i].ID == id {
			return &users[i], users[i].Cold
		}
	}
	return nil, nil
}

const numUsers = 100000

func makeUserData(i int) (uint64, int, string, string, string, string, string,
	string, string, string, string, string, string, string, string, string,
	string, string, string, string, string, string, string, string, string,
	string) {
	return uint64(i),
		i % 10000,
		fmt.Sprintf("User%d", i),
		fmt.Sprintf("user%d@example.com", i),
		fmt.Sprintf("%d Main Street", i),
		fmt.Sprintf("Biometrix of user %d with text......", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i),
		fmt.Sprintf("extra%d", i)
}

func setupMixedUsers() []User {
	users := make([]User, numUsers)
	for i := range users {
		id, score, name, email, addr, bio, e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14, e15, e16, e17, e18, e19, e20 := makeUserData(i)
		users[i] = User{id, score, name, email, addr, bio, e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14, e15, e16, e17, e18, e19, e20}
	}
	return users
}

func setupSeparatedUsers() []UserHot {
	users := make([]UserHot, numUsers)
	for i := range users {
		id, score, name, email, addr, bio, e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14, e15, e16, e17, e18, e19, e20 := makeUserData(i)
		users[i] = UserHot{id, score,
			&UserCold{name, email, addr, bio, e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14, e15, e16, e17, e18, e19, e20},
		}
	}
	return users
}

// BenchmarkComparisonHotCodeSeparation 使用冷热分离的时候需要考虑
//  1. 主要: struct的大小有多少, 当struct大致计算已经接近/超过 CPU的cache-line(m1pro-128bytes),
//     使用冷热分离则非常有必要, 已经不care具体操作的数据量
//  2. 次要: 当struct不是特别夸张, 但数据量非常离谱的时候, 那进行冷热分离则也有一定的bene
//
// 反例: 当struct也好操作的size也好都不是夸张的时候, 简单一个struct反而能更快
func BenchmarkComparisonHotCodeSeparation(b *testing.B) {

	mixUsers := setupMixedUsers()
	sepUsers := setupSeparatedUsers()

	b.Run("Mixed-TopUsers", func(b *testing.B) {
		for b.Loop() {
			usersCopy := append([]User(nil), mixUsers...)
			_ = TopUsersMixed(usersCopy, numUsers)
		}
	})

	b.Run("Separated-TopUsers", func(b *testing.B) {
		for b.Loop() {
			usersCopy := append([]UserHot(nil), sepUsers...)
			_ = TopUsersSeparated(usersCopy, numUsers)
		}
	})

	b.Run("Mixed-GetDetails", func(b *testing.B) {
		i := 0
		for b.Loop() {
			_ = GetUserDetailsMixed(mixUsers, uint64(i%numUsers))
			i++
		}
	})

	b.Run("Separated-GetDetails", func(b *testing.B) {
		i := 0
		for b.Loop() {
			_, _ = GetUserDetailsSeparated(sepUsers, uint64(i%numUsers))
			i++
		}
	})

}
