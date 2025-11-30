## single number

Given a non-empty array of integers nums, every element appears twice except for one. Find that single one.

You must implement a solution with a linear runtime complexity and use only constant extra space.

找出數組中唯一沒有重複的數字
如果用一般的方式其實也挺簡單的，直接用個map，有數字的時候將他輸入map中，如果已經有了就刪除，最後就會是答案

但因為我打算用XOR的形式做看看，那好像其實也很簡單，直接遍歷應該就可以去重得到答案了

```go
func singleNumber(nums []int) int {
    if (len(nums) == 0){
        return 0
    }
    xor := 0
    for _,v := range nums {
        xor ^=v
    }
    return xor
}
```