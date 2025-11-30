package code

/*
Given a non-empty array of integers nums, every element appears twice except for one. Find that single one.

You must implement a solution with a linear runtime complexity and use only constant extra space.
*/
func SingleNumber(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	xor := 0
	for _, v := range nums {
		xor ^= v
	}
	return xor
}
