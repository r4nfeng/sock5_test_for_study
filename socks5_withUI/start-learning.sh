#!/bin/bash

echo "==================================="
echo "SOCKS5 学习助手 - 启动脚本"
echo "==================================="
echo ""

echo "[1/2] 启动Web服务器..."
go run simple-server.go &
SERVER_PID=$!

echo ""
echo "[2/2] 等待服务器启动..."
sleep 2

echo ""
echo "✅ 服务器启动成功！"
echo ""
echo "📖 打开浏览器访问: http://localhost:3000"
echo ""
echo "💡 使用提示："
echo "   1. 先启动SOCKS5服务器: go run main/main.go"
echo "   2. 在浏览器中打开上面的地址"
echo "   3. 点击'开始测试'按钮"
echo ""
echo "按 Ctrl+C 停止服务器"

# 等待服务器进程
wait $SERVER_PID
