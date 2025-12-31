# SOCKS5 协议分析工具

一个简单的前端工具，用于可视化展示SOCKS5协议的实现过程，配合Wireshark抓包进行深入学习。

## 功能特点

✅ **协议可视化** - 实时展示SOCKS5的4个关键步骤
✅ **数据包详解** - 显示每个步骤的完整数据包格式
✅ **终端模拟** - 模拟命令行输出协议交互过程
✅ **Wireshark集成** - 提供详细的抓包分析指南
✅ **简单易用** - 无需安装，浏览器中直接运行

## 快速开始

### 1. 启动Web服务器

```bash
go run simple-server.go
```

或者使用启动脚本：

**Windows:**
```bash
start-learning.bat
```

**Linux/Mac:**
```bash
chmod +x start-learning.sh
./start-learning.sh
```

### 2. 打开浏览器

访问: http://localhost:3000

### 3. 配置并测试

1. 填写SOCKS5服务器信息（已预填默认值）
2. 点击"开始测试连接"
3. 观察协议过程和数据包变化

## SOCKS5协议步骤

### 步骤1: 握手 (Handshake)

客户端发送支持的认证方法列表，服务器选择并响应。

**数据包格式:**
```
05 02 00 02
```
- `05` = SOCKS5版本
- `02` = 支持2种方法
- `00` = 无需认证
- `02` = 用户名/密码认证

### 步骤2: 认证 (Authentication)

客户端发送用户名和密码，服务器验证并返回结果。

**数据包格式:**
```
01 08 74 65 73 74 75 73 65 72 08 74 65 73 74 70 61 73 73
```
- `01` = 认证版本
- `08` = 用户名长度(8)
- `74 65 73 74 75 73 65 72` = "testuser"
- `08` = 密码长度(8)
- `74 65 73 74 70 61 73 73` = "testpass"

### 步骤3: 连接请求 (Request)

客户端发送CONNECT命令，请求连接到目标服务器。

**数据包格式:**
```
05 01 00 03 0B 77 77 77 2E 65 78 61 6D 70 6C 65 2E 63 6F 6D 01 BB
```
- `05` = SOCKS5版本
- `01` = CONNECT命令
- `00` = 保留字段
- `03` = 域名类型
- `0B` = 域名长度(11)
- `77 77 77...` = "www.example.com"
- `01 BB` = 端口443

### 步骤4: 连接响应 (Response)

服务器返回连接结果，成功后开始双向转发数据。

**成功响应:**
```
05 00 00 01 7F 00 00 01 1F 90
```
- `05` = SOCKS5版本
- `00` = 成功
- `01` = IPv4地址类型
- `7F 00 00 01` = 127.0.0.1
- `1F 90` = 端口8080

## Wireshark抓包分析

### 安装Wireshark

下载地址: https://www.wireshark.org/download.html

### 抓包步骤

1. **启动Wireshark**，选择网络接口
   - Windows: 选择正在使用的网卡
   - 抓本地回环：安装Npcap时勾选"Support loopback"

2. **设置过滤器**
   ```
   tcp.port == 1080
   ```
   （替换为你的SOCKS5服务器端口）

3. **开始抓包**
   - 点击蓝色鲨鱼鳍图标
   - 或按 `Ctrl + E`

4. **触发连接**
   - 在本页面点击"开始测试连接"
   - 或使用curl命令测试：
   ```bash
   curl --socks5 127.0.0.1:1080 -U testuser:testpass http://www.example.com
   ```

5. **查看数据包**
   - Wireshark会显示完整的TCP通信过程
   - 每个数据包都可以展开查看详细信息
   - 查看SOCKS5协议层的内容

### 分析技巧

**Follow TCP Stream**
- 右键点击数据包
- 选择 "Follow" → "TCP Stream"
- 查看完整的TCP对话内容

**查看SOCKS5协议细节**
- 展开数据包
- 找到 "SOCKS" 或 "Socket Proxy" 协议层
- 查看各个字段的值和含义

**过滤器示例**
```
# 只看SOCKS5协议
socks

# 只看特定连接
tcp.stream == 0

# 只看认证包
socks.auth_name

# 只看连接请求
socks.cmd_code == 1
```

### 对比学习

将Wireshark抓到的实际数据包和本页面显示的理论数据包对比：

1. **握手包** - Wireshark显示 `05 02 00 02`
2. **认证包** - Wireshark显示用户名密码的十六进制
3. **请求包** - Wireshark显示CONNECT命令和目标地址
4. **响应包** - Wireshark显示连接结果

这样可以从理论和实践两个层面深入理解SOCKS5协议。

## 使用场景

### 1. 学习SOCKS5协议

- 理解协议的工作流程
- 掌握数据包的格式和含义
- 学习各个字段的取值和作用

### 2. 调试SOCKS5服务器

- 启动自己的SOCKS5服务器
- 使用本工具触发连接
- 在Wireshark中查看实际的数据包
- 对比预期结果，找出问题

### 3. 分析SOCKS5流量

- 使用Wireshark抓取SOCKS5流量
- 理解每个数据包的作用
- 分析认证过程和连接建立过程

## 技术实现

- **前端**: Vue 3 (CDN方式，无需构建)
- **后端**: Go HTTP服务器
- **样式**: 纯CSS，无需UI框架

## 项目结构

```
simple-frontend/
├── index.html    # 主页面（HTML + CSS + Vue模板）
├── app.js        # Vue应用逻辑
└── README.md     # 本文件
```

## 扩展阅读

- [RFC 1928 - SOCKS Protocol Version 5](https://tools.ietf.org/html/rfc1928)
- [RFC 1929 - Username/Password Authentication](https://tools.ietf.org/html/rfc1929)
- [Wireshark Wiki - SOCKS](https://wiki.wireshark.org/SOCKS)

## 许可证

MIT License
