package pref

import (
	"fmt"
	"sort"
)

type User struct {
	ID      uint64 // hot field, 8 bytes
	Score   int    // hot field, 8 bytes
	Name    string
	Email   string
	Address string
	Bio     string
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
	Name    string
	Email   string
	Address string
	Bio     string
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

const numUsers = 10000

func makeUserData(i int) (uint64, int, string, string, string, string) {
	return uint64(i),
		i % 10000,
		fmt.Sprintf("User%d", i),
		fmt.Sprintf("user%d@example.com", i),
		fmt.Sprintf("%d Main Street", i),
		fmt.Sprintf("Biometrix of user %d with text......", i)
}

func setupMixedUsers() []User {
	users := make([]User, numUsers)
	for i := range users {
		id, score, name, email, addr, bio := makeUserData(i)
		users[i] = User{
			ID:      id,
			Score:   score,
			Name:    name,
			Email:   email,
			Address: addr,
			Bio:     bio,
		}
	}
	return users
}

func setupSeparatedUsers() []UserHot {
	users := make([]UserHot, numUsers)
	for i := range users {
		id, score, name, email, addr, bio := makeUserData(i)
		users[i] = UserHot{id, score,
			&UserCold{name, email, addr, bio},
		}
	}
	return users
}
