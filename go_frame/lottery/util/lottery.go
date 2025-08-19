package util

import (
	"cmp"
	"math/rand/v2"
)

// 抽奖。给定每个奖品被抽中的概率（无需要做归一化，但概率必须大于0），返回被抽中的奖品下标
func Lottery(probs []float64) int {
	if len(probs) == 0 {
		return -1
	}
	sum := 0.0
	acc := make([]float64, 0, len(probs)) //累积概率
	for _, prob := range probs {
		sum += prob
		acc = append(acc, sum)
	}

	// 获取(0,sum] 随机数
	r := rand.Float64() * sum
	index := BinarySearch4Section(acc, r)
	return index
}

func BinarySearch4Section[T cmp.Ordered](arr []T, target T) int {
	if len(arr) == 0 {
		return -1
	}
	begin, end := 0, len(arr)-1

	for {
		//arr[begin]在target后面
		if arr[begin] >= target {
			return begin
		}
		//arr[end]在target前面
		if arr[end] < target {
			return end + 1
		}

		//二分查找法
		middle := (begin + end) / 2
		if arr[middle] > target {
			end = middle - 1 //arr[end]可能会跑到target前面
		} else if arr[middle] < target {
			begin = middle + 1 //arr[begin]可能会跑到target后面
		} else {
			return middle
		}
	}
}
