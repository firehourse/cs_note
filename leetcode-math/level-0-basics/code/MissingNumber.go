package code

/*
Given an array nums. We define a running sum of an array as runningSum[i] = sum(nums[0]…nums[i]).

Return the running sum of nums.

Example 1:

Input: nums = [3,0,1]

Output: 2

Explanation:

n = 3 since there are 3 numbers, so all numbers are in the range [0,3]. 2 is the missing number in the range since it does not appear in nums.
*/
func MissingNumber(nums []int) int {
	// 先知道n
	n := len(nums)
	// 高斯公式
	total := n * (n + 1) / 2
	// 先初始化加總
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return total - sum
}
