# Phase 1 — 最小 TCP Server（Echo）

打算先建立一個可運作的TCP server

能用來作為client 接收 socket 的請求並進行連線
能讀取資料
能回傳資料
使用 goroutine 來並行處理多個連線  
  （I/O 多工是由 Go runtime + epoll/kqueue 自動處理）


這是整個 server framework 的最底層基礎。  
後續所有功能（connection 包裝、pipeline、framing、event loop）都會建立在這個 MVP 上