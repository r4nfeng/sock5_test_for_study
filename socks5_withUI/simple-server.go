package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// 设置静态文件服务
	// 允许前端访问simple-frontend目录下的文件
	fs := http.FileServer(http.Dir("./simple-frontend"))

	// 处理根路径，返回index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果请求的是根路径或目录，则返回index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("./simple-frontend", "index.html"))
		} else {
			// 其他请求由文件服务器处理
			fs.ServeHTTP(w, r)
		}
	})

	// 启动HTTP服务器
	port := ":3000"
	log.Printf("🚀 SOCKS5 学习助手启动成功！")
	log.Printf("📖 访问地址: http://localhost%s", port)
	log.Printf("💡 打开浏览器访问上述地址开始学习")
	log.Printf("")
	log.Printf("提示：")
	log.Printf("  1. 确保 SOCKS5 服务器正在运行 (go run main/main.go)")
	log.Printf("  2. 在前端页面配置 SOCKS5 服务器信息")
	log.Printf("  3. 点击'开始测试'按钮查看协议过程")
	log.Printf("")
	log.Printf("按 Ctrl+C 停止服务器")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
