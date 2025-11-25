這題源自 NeetCode 最基本的 Dynamic Array。
我在學 Go 的時候，一開始不太懂為什麼大家說 Slice 跟一般 Array 不一樣。後來才發現，這其實是相對於 C++ / Java 的基本陣列來說的。

在 C++ / Java 的語境中，Array 是固定大小的，一旦建立就不能動。

但在很多高階語言中，例如 PHP / Python / JavaScript，所謂的「array」其實已經是 動態陣列的封裝物件。

Java 雖然 array 是物件，但它依然是 固定長度的 raw storage，並不具備動態擴容能力。

因此，Go 的 slice 其實不是 “特別奇怪的語法”，而是類似於：

 C++ std::vector
 Java ArrayList
  Python / PHP 的 list

Go slice 就是一個高階動態陣列的實作，只是語言內建。

下面這段 C++ 程式碼展示了一個最基本、手動實作 的動態陣列，用來說明：

底層仍然是 int* 的固長陣列

當空間不夠時進行倍增（resize）

提供 get / set / pushback / popback 等常用 API

這個實作可視為 對標 Go Slice 背後的 growth strategy。
```cpp
class DynamicArray {
private:
    int* arr;
    int size;
    int capacity;
public:

    DynamicArray(int capacity) {
        this->arr = new int[capacity];
        this->size = 0;
        this->capacity = capacity;
    }
    ~DynamicArray() {
    delete[] arr;
    }

    int get(int i) {
        if(i<0 || i >= size){
            throw std::out_of_range("Index out of range");
        }
        return arr[i];
    }

    void set(int i, int n) {
        if(i < 0){
            throw std::out_of_range("Index out of range");
        }
        while(i >= capacity){
            this->resize();
        }
        for (int j = size; j < i; j++) {
            arr[j] = 0;
        }

        arr[i] = n;
        if (i >= size) {
            size = i + 1;
        }
    }

    void pushback(int n) {
        if (size >= capacity) {
            resize();
        }
        arr[size] = n;
        size ++;
    }

    int popback() {
        if(size == 0 ){
            throw std::out_of_range("Index out of range");
        }
        size --;
        return arr[size];
    }

    void resize() {
        int newCapacity = capacity *2;
        capacity = newCapacity;
        int* arr2 = new int[capacity];
        for(int i=0;i<size;i++){
            arr2[i] = arr[i];
        }
        delete[] arr;
        arr = arr2; 
    }

    int getSize() {
        return size;
    }

    int getCapacity() {
        return capacity;
    }
};
```
# go slice 知識補充
前面列舉了我認知的動態陣列作法，那這邊有一個小知識就是ｇｏ針對擴容的優化
當cap < 1024 會變成 cap = cap * 2
當cap >= 1024 會變成 cap = cap + (cap/4)
避免大型的slice無限制的膨脹造成記憶體浪費
 
# 常見的幾個 go 宣告 slice的方法

```go
var a = []int{1, 2, 3} 
```

這宣告的是 slice。
長度是型別的一部分，因此無法自動擴容，是最基本、最原始的 array 型別。

```go
var a = make([]int, 3)
```
使用 make 宣告 slice。
這仍然是一個 slice（而不是 array），但：

len = 3
cap = 3
換句話說：
底層會分配一塊長度為 3 的陣列，並且 slice 直接引用它全部的範圍。

```go 
var b = [3]int{1, 2, 3} 
```
這宣告的是 固定大小的 array。
長度是型別的一部分，因此無法自動擴容，是最基本、最原始的 array 型別。

```go
var s = make([]int, 0, 3) 
```
這也是 slice，但：
size = 0
cap = 3
也就是說：
目前沒有元素（長度 0）
但底層已經預留了容量 3 的空間
append 的時候可以不重新配置，直到超過容量為止
這點和前面的 C++ DynamicArray 完全一致：
雖然底層會開一個固定大小（capacity）的 array，但並不會一直 free / realloc，而是用 size（Go slice 中是 len）來控制邊界、決定哪些元素是有效
的。


# go slice的幾個基本操作

```go
s := []int{1, 2, 3}
sub := s[1:3] // [2, 3]
```
上面的操作主要是依靠簡單的指標，也就是說他們實際上指到的是同個array區段，所以如果對其中一個改動另一個也會跟著發生變動

```go 
s := []int{1, 2, 3}
s = append(s, 3)
```
上面基本上就跟其他語言一樣，在最後面的時候加入
python 裡面是 append 
php 是 $arr[] = 3;
java 是 ArrayList.add()
cpp 則是  push_back()

那觸發擴容的話其實從上面cpp代碼應該可以清楚看到，就會用一個新的陣列來接，所以前面提到的sub方式要注意

```go
a := []int{1, 2, 3}
b := make([]int, len(a))
copy(b, a)
```
copy其實指的是把一個陣列複製到另一個陣列，那他會用min的方式去取，這保證了不會有邊界超出的問題

```go
a := []int{1, 2, 3}
a = a[:0]
```
這是一個清空操作，但是避免了這個slice被gc釋放
```go
i := 2
s = append(s[:i], s[i+1:]...)
```
這是一個刪除操作，這個操作是把後面的元素往前移動，然後把最後一個元素刪除
