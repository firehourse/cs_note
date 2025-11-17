
# LeetCode 第27題 Remove Element

Given an integer array nums and an integer val, remove all occurrences of val in nums in-place. The order of the elements may be changed. Then return the number of elements in nums which are not equal to val.

Consider the number of elements in nums which are not equal to val be k, to get accepted, you need to do the following things:

Change the array nums such that the first k elements of nums contain the elements which are not equal to val. The remaining elements of nums are not important as well as the size of nums.
Return k.

題目大概的意思就是要操作一個陣列，那他會給定一個val，他希望妳回傳撇除val元素的長度且不需要排序，必須要實際操作這個nums陣列他會檢查

那我這邊的思路其實很簡單，先看LeetCode給的範例
Input: nums = [3,2,2,3], val = 3
Output: 2, nums = [2,2,_,_]

既然val都要再最後，那我就維護兩個指針，一個一直指著最後，直到最後這個值變成了val他才移動，另一個指針則從後面遍歷回來
```Go
func removeElement(nums []int, val int) int {
    // 慢指針
    slow := len(nums)-1
    // 遍歷指針
    for i := len(nums)-1; i>=0;i--{
        if(nums[i] == val ){
            temp := nums [i]
            nums[i] = nums[slow]
            nums[slow] = temp
            slow --
        }
    }
// 回傳長度
return slow +1 ;
}
```