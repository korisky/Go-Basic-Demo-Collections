package algorithm

type Algorithm func(items []int, lookup int) int

func CheckEveryItem(items []int, lookup int) int {
	for i := 0; i < len(items); i++ {
		if items[i] == lookup {
			return i
		}
	}
	return -1
}

func BinarySearch(items []int, lookup int) int {
	left := 0
	right := len(items) - 1

	for true {

		if left == lookup {
			return left
		}
		if right == lookup {
			return right
		}

		center := (left + right) / 2
		if items[center] == lookup {
			return center
		}
		if center > lookup {
			right = center
		}
		if center < lookup {
			left = center
		}
		if left >= right-1 {
			return -1
		}
	}
	return -1
}
