# Longest Substring Without Repeating Characters

Given a string s, find the length of the longest substring without duplicate characters.

 

Example 1:

Input: s = "abcabcbb"
Output: 3
Explanation: The answer is "abc", with the length of 3. Note that "bca" and "cab" are also correct answers.


```go
// 首先因為這算是有順序的，所以雙指針，二分查找這種概念絕對有
// 我們同時可以記錄一個長度，用max函數
// 接下來就是思考指針要怎麼擺放，因為他是不能重複
// 直覺一點的作法是我永遠有兩個指針，然後有個哈希紀錄現在某個值他在第幾個位置？
// 然後當遇到重複的，慢指針就移動到重複的下一個位置這樣
func lengthOfLongestSubstring(s string) int {
    // 最後的長度返回
    result := 0
    // map來記錄byte 裡面存位置所以是int
    mem := make(map[byte]int)
    // 初始化兩個指針
    slow := 0
    // 先不考慮什麼優化，先讓迴圈能進行
    for fast := 0; fast < len(s); fast ++ {
        // 判斷slow邊界
        if lastIndex, ok := mem[s[fast]]; ok && lastIndex >= slow {
                slow = lastIndex + 1
        }
        mem[s[fast]] = fast
        result = max(result, fast-slow+1)

    }
    return result
}
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
  