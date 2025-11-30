## server

這個專案主要是想從底層理解網路通訊框架是怎麼運作的，類似 Netty 那種東西。

很多人可能會覺得 server 就是開個 port 然後收資料、送資料，但實際上一個可用的 server framework 要處理的東西遠比想像中複雜。

這個專案不是要重做 Netty，而是透過實作來理解：

- OS Socket 跟網路底層到底在幹嘛
- Go runtime 的調度器、IO 多工、M:N 模型
- TCP 黏包/拆包問題（framing）
- pipeline/middleware 設計模式
- event loop 跟連線管理
- graceful shutdown 怎麼做
- 雲原生的 networking 底層（Gateway、Sidecar、Proxy 的基礎）

## 目標

最終目標是做出一個可以支援：
- TCP server（長連線）
- pipeline handler（類似 middleware）
- framing（decoder/encoder）
- event loop（多 worker 或 goroutine 模式）
- graceful shutdown
- 可擴展的協議（JSON、Protobuf、自訂 binary）

最後可以拿來做 WebSocket / RPC server 的核心。

不過一開始的 MVP 很簡單：
- 開一個 TCP server
- 能 Accept 連線
- 能收資料、回資料（Echo）
- 能正常關閉

## 專案結構

```
server/
    README.md
    
    /basics
        socket.md         ← OS Socket 底層
        tcp.md            ← TCP 協議基礎
        
    /phase-1-minimal
        README.md         ← 最小 server 實作
        echo.go           ← Echo server 範例
        
    /phase-2-connection
        README.md         ← Connection 包裝
        conn.go           ← Conn struct
        
    /phase-3-pipeline
        README.md         ← Pipeline 設計
        handler.go        ← Handler interface
        
    /phase-4-framing
        README.md         ← 黏包/拆包
        decoder.go        ← Decoder 實作
        
    /phase-5-eventloop
        README.md         ← Event loop 模型
        
    /phase-6-lifecycle
        README.md         ← Graceful shutdown
```

## 學習路線

### Phase 0 — 基礎知識

先搞懂 OS Socket 跟 TCP 的基本概念，不然後面會看不懂。

### Phase 1 — 最小 Server

重點：
- `net.Listen` 背後做了什麼（OS socket）
- `Accept()` 為什麼會阻塞（accept queue）
- `go handleConn(conn)` 為什麼要開 goroutine
- Read / Write 的 buffer 機制
- 怎麼正常關閉 server

### Phase 2 — Connection 包裝

把 `net.Conn` 包一層，加上：
- 每個連線的唯一 ID
- Write / Close 的抽象化
- connection table（管理所有連線）

這樣 server 就不再是「裸運行」，而是有 framework 的雛形。

### Phase 3 — Pipeline

實作 Handler Chain，類似 middleware：
- Handler interface
- Pipeline = Handler slice
- OnRead → 經過多層 handler
- OnWrite → outbound handler

這樣就可以做到可擴展的設計。

### Phase 4 — Framing

這是最重要的部分，要處理 TCP 黏包/拆包問題。

實作幾種常見的 framing：
- LengthFieldBasedFrameDecoder（長度前綴）
- DelimiterBasedFrameDecoder（分隔符，例如 \n）
- FixedLengthFrameDecoder（固定長度）

這層做好了，WebSocket / RPC 都可以建在上面。

### Phase 5 — Event Loop

這是進階部分，要理解 Netty、Nginx、Node.js 的核心架構：
- 單 loop（goroutine per conn）
- 多 loop（類 Netty Reactor）
- 每個連線固定綁定 loop（避免 race）

### Phase 6 — Lifecycle

實作 graceful shutdown：
- context with cancel
- server shutdown hook
- connection draining
- connection idle timeout

到這邊就是一個可以用在 production 的 server framework 了。

### Phase 7 — Protocol

選一個協議來實作：
- JSON protocol
- Protobuf protocol
- 自訂 binary 協議

### Phase 8 — 擴展

加上一些實用的功能：
- builder 模式
- 內建 middleware
- 錯誤處理
- logging
- metrics

### Phase 9 — 進階應用

可以選擇：
- WebSocket server
- RPC server（類 gRPC）
- Gateway / Proxy
