package pref

import "sort"

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
