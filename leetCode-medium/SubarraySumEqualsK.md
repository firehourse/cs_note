# SubarraySumEqualsK

Given an array of integers nums and an integer k, return the total number of subarrays whose sum equals to k.

A subarray is a contiguous non-empty sequence of elements within an array.

 

Example 1:

Input: nums = [1,1,1], k = 2
Output: 2

給定一個整數陣列跟一個整數k，返回陣列中和為k的子陣列數量
題目中明確定義了subarray 是連續非空的序列

```go
// 現在可以知道整個陣列是整數的，那因為又是子問題，但因為有可能有負數
// 題目說要等於K，考慮負數的情況下，有點難用快慢雙指針
// 有負數的情況下感覺要用dp來做，那我把0~x位存進map裡，那用two sum的思路來解
func subarraySum(nums []int, k int) int {
    result := 0
    // k 存 sum v 存 出現次數 而如果target跟sum相減的值存在 就增加他的次數
    hmap := make(map[int]int)
    // 如果index從0開始 會有邊界問題
    hmap[0] = 1
    sum := 0
    for _, v := range nums {
        sum += v
        if _, ok := hmap[sum-k]; ok{
            result += hmap[sum-k]
        }
        hmap[sum] += 1
    }

    return result
}