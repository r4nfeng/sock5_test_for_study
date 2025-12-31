@echo off
echo ===================================
echo SOCKS5 学习助手 - 启动脚本
echo ===================================
echo.

echo [1/2] 启动Web服务器...
start "SOCKS5 Learning Server" go run simple-server.go

echo.
echo [2/2] 等待服务器启动...
timeout /t 2 /nobreak > nul

echo.
echo ✅ 服务器启动成功！
echo.
echo 📖 打开浏览器访问: http://localhost:3000
echo.
echo 💡 使用提示：
echo    1. 先启动SOCKS5服务器: go run main/main.go
echo    2. 在浏览器中打开上面的地址
echo    3. 点击"开始测试"按钮
echo.
echo 按任意键关闭此窗口...
pause > nul
