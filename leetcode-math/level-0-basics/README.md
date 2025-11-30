# Level 0：基礎數學入門 🌱

> 只有加減乘除，但帶入『模式』與『數學觀念』  
> 這一段是入門數學思維，讓你看到暴力法 vs. 思維法的差異

## 進度追蹤

**完成：** 0/6 題

---

## ① 1480. Running Sum（前綴和 Prefix Sum）

- [ ] 完成題目
- [ ] 理解前綴和概念
- [ ] 優化：從暴力 O(n²) → 思維 O(n)

**數學點：**
- 前綴和（Prefix Sum）概念，之後 100 題以上都會用到
- 從暴力 O(n²) → 思維 O(n)

**你會學到：**
- 「重複運算如何壓到一次」
- 演算法的最基本數學：將一段區間用總和減法得到

**LeetCode 連結：** [1480. Running Sum of 1d Array](https://leetcode.com/problems/running-sum-of-1d-array/)

---

## ② 268. Missing Number（高斯求和 / XOR）

- [ ] 完成題目（加法公式版本）
- [ ] 完成題目（XOR 版本）
- [ ] 理解高斯求和公式

**數學點：**
- 公式 Sum = n(n+1)/2（高斯和）
- 或用 XOR（但先學加法版本）

**你會學到：**
- 如何用加法公式取代遍歷
- 如何用 XOR 消掉重複值（算是數學的位元等價運算）

**LeetCode 連結：** [268. Missing Number](https://leetcode.com/problems/missing-number/)

---

## ③ 136. Single Number（XOR 數學性質）

- [ ] 完成題目
- [ ] 理解 XOR 性質

**數學點：**
- `a XOR a = 0`
- `0 XOR b = b`

**你會學到：**
- XOR 的本質：「兩次相抵」的運算
- 完全是一個群論（Group Theory）的味道（但不用真的學群論）

**LeetCode 連結：** [136. Single Number](https://leetcode.com/problems/single-number/)

---

## ④ 171. Excel Sheet Column Number（26 進位 → 十進位）

- [ ] 完成題目
- [ ] 理解進位制轉換

**數學點：**
- 字元轉換、進位制
- 本質是一個數學轉換題

**你會學到：**
- 26 進位如何轉十進位
- 之後可以直接銜接各種進位題（2進位、36進位等）

**LeetCode 連結：** [171. Excel Sheet Column Number](https://leetcode.com/problems/excel-sheet-column-number/)

---

## ⑤ 1523. Count Odd Numbers in an Interval Range（區間奇數計數）

- [ ] 完成題目
- [ ] 理解數學公式推導

**數學點：**
- 區間計數公式
- 奇偶性質

**你會學到：**
- 如何用數學公式 O(1) 解決計數問題
- 避免暴力遍歷

**LeetCode 連結：** [1523. Count Odd Numbers in an Interval Range](https://leetcode.com/problems/count-odd-numbers-in-an-interval-range/)

---

## ⑥ 9. Palindrome Number（回文數判斷）

- [ ] 完成題目
- [ ] 不轉字串的數學解法

**數學點：**
- 數字反轉
- 取模運算 `% 10` 和除法 `/ 10`

**你會學到：**
- 如何用數學操作拆解數字
- 為後續數位操作題打基礎

**LeetCode 連結：** [9. Palindrome Number](https://leetcode.com/problems/palindrome-number/)

---

## 🎯 本關重點

完成這 6 題後，你應該掌握：

1. ✅ **前綴和思想** - 如何用累加避免重複計算
2. ✅ **數學公式** - 高斯求和、區間計數
3. ✅ **XOR 基礎** - 理解 XOR 的「相消」性質
4. ✅ **進位制** - 26 進位轉換
5. ✅ **數位操作** - 用 % 和 / 拆解數字

準備好了就前往 [Level 1：數列與組合](../level-1-sequences/README.md) 吧！
