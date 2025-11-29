# leetcode 268 Missing Number

題目:
Given an array nums containing n distinct numbers in the range [0, n], return the only number in the range that is missing from the array.

 

Example 1:

Input: nums = [3,0,1]

Output: 2

Explanation:

n = 3 since there are 3 numbers, so all numbers are in the range [0,3]. 2 is the missing number in the range since it does not appear in nums.

給定了一個數字陣列，其中包含了從 0 到 n 的不重複的數字，返回這個陣列中缺少的數字。


那這題首先直覺的幾個作法是
1. 加總全部的數字，然後用0~n的加總去減，也就是使用高斯公式求等差數列和

2. 用set去做，先排序好接著去尋找缺失的數字

那透過ai我了解到也可以用XOR來做
那這邊主要會講解高斯公式的推導跟XOR的作法



## 等差數列求和

首先寫出等差數列的加總

(1+2+3+4+5+...n)
那倒過來
n+(n-1)+(n-2)+...+1

接著我們要讓他重新排列組合，也就是找到其中的規律
```
1 2 3 4 5

5 4 3 2 1 
```
列出了可能性之後讓他們相加

你會發現全部都是 n+1 那全部加出來之後 得到的其實是 2s 
所以你還需要除以2 
這樣就可以得到 n*(n+1)/2 這個等差數列公式

# XOR 求解

XOR其實就是真值表的概念
他比對的是兩個bit 是否一樣這件事情

不一樣的結果是 1 
一樣的結果是 0 

a ^ a = 0
```
0101 XOR 0101 = 0000
```
數字化以後

```
5 ^ 5 = 0
9 ^ 9 = 0
123 ^ 123 = 0
```
但前面有提到 他是比較bit
所以這時候我們要先用二進位來看
```
5 ^ 3 
5 = 101
3 = 011

1 0 1
0 1 1 
------
1 1 0
```
這時候我們會得到110
二進制進行轉換
1 x 2^2 + 1 x 2^1 + 0 x 2^0 = 4 + 2 + 0 = 6

那這時候我們可以得知
如果我們用一個完整的等插數列跟題目給定的nums 數組做XOR
就可以把數字都抵銷留下唯一不同的部分
那這一樣是一個O(n)的時間複雜度

```go
func missingNumber(nums []int) int {
    xor := len(nums)
    for i, v := range nums {
        xor ^= i ^ v
    }
    return xor

}
```
拆解一下上面的寫法，其實就是同時把完整的等差數列跟題目給定的nums數組做XOR，抵銷完以後就是完整的剩餘的那個數字