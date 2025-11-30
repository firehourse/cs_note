package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	addr     string
	listener net.Listener
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}

}

func (s *Server) Start() error {
	// 1. 用 OS 建立 TCP 監聽 socket
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln

	fmt.Println("Server started on", s.addr)

	// 2. 進入接受連線的無限迴圈
	for {
		// Accept() 在「沒有新的連線」時會阻塞
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		// 每條連線用 goroutine 處理
		go s.handleConn(conn)
	}
}

// 看到這邊可能會疑惑，為什麼我把資料讀出來又寫回去，因為這只是一個最小化的測試能應用而已，那後續我會逐漸完善整個server
func (s *Server) handleConn(conn net.Conn) {
	// 確保關閉連線的時候會執行defer
	defer conn.Close()
	// 宣告一個buffer 來承接字符
	buf := make([]byte, 1024)
	// 開一個While true 來運行
	for {
		// n 來進行讀取
		n, err := conn.Read(buf)
		// 當 回傳異常的時候
		if err != nil {
			// 如果宣告end of file 代表斷開連結
			if err == io.EOF {
				fmt.Println("client disconnected")
			} else {
				// 否則打印異常
				fmt.Println("read error:", err)
			}
			// 提早結束
			return
		}
		// 用一個data進行值複製 然後打印
		data := buf[:n]
		fmt.Println("recv:", string(data))
		// 這邊進行寫進 os send buffer 傳送緩衝區
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}

}
