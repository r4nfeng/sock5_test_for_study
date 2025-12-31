/**
 * SOCKS5 协议分析工具 - Vue应用
 * 用于可视化展示SOCKS5协议过程，配合Wireshark抓包学习
 */

const { createApp } = Vue;

createApp({
  data() {
    return {
      // SOCKS5连接配置
      config: {
        host: '127.0.0.1',
        port: 1080,
        username: 'testuser',
        password: 'testpass',
        target: 'www.example.com',
        targetPort: 80
      },

      // 测试状态
      testing: false,
      testSuccess: false,
      currentStep: -1,

      // 终端输出
      terminalLines: [],

      // 步骤状态
      steps: [
        { name: '握手', failed: false },
        { name: '认证', failed: false },
        { name: '请求', failed: false },
        { name: '连接', failed: false }
      ]
    };
  },

  methods: {
    // 开始测试SOCKS5连接
    async startTest() {
      if (!this.config.host || !this.config.port) {
        this.addTerminal('❌ 请填写完整的SOCKS5服务器地址和端口', 'error');
        return;
      }

      // 重置状态
      this.testing = true;
      this.testSuccess = false;
      this.currentStep = -1;
      this.steps.forEach(s => s.failed = false);
      this.clearTerminal();

      try {
        // 步骤1: 握手
        this.addTerminal(`📡 连接到 ${this.config.host}:${this.config.port}`, 'info');
        await this.delay(500);
        this.currentStep = 0;
        this.addTerminal('🤝 握手: 发送SOCKS5初始化', 'info');
        this.addTerminal('   → 05 02 00 02', 'info');
        await this.delay(800);
        this.addTerminal('🤝 握手: 服务器选择认证方法', 'success');
        this.addTerminal('   ← 05 02 (用户名密码认证)', 'success');

        // 步骤2: 认证
        await this.delay(500);
        this.currentStep = 1;
        this.addTerminal(`🔐 认证: 用户 ${this.config.username}`, 'info');
        this.addTerminal(`   → 发送用户名和密码`, 'info');
        await this.delay(800);

        // 模拟认证（80%成功率）
        const authSuccess = Math.random() > 0.2;
        if (authSuccess) {
          this.addTerminal('✅ 认证: 成功', 'success');
          this.addTerminal('   ← 01 00', 'success');
        } else {
          this.steps[1].failed = true;
          this.addTerminal('❌ 认证: 失败', 'error');
          this.addTerminal('   ← 01 01', 'error');
          this.testing = false;
          return;
        }

        // 步骤3: 请求
        await this.delay(500);
        this.currentStep = 2;
        this.addTerminal(`📡 请求: CONNECT ${this.config.target}:${this.config.targetPort}`, 'info');
        await this.delay(800);

        // 步骤4: 连接
        await this.delay(600);
        this.currentStep = 3;

        // 模拟连接（70%成功率）
        const connectSuccess = Math.random() > 0.3;
        if (connectSuccess) {
          this.addTerminal('✅ 连接: 成功建立', 'success');
          this.addTerminal('   ← 05 00 00 01 [绑定地址] [绑定端口]', 'success');
          this.addTerminal('🎉 SOCKS5代理连接完成', 'success');
          this.testSuccess = true;
        } else {
          this.steps[3].failed = true;
          this.addTerminal('❌ 连接: 失败', 'error');
          this.addTerminal('   ← 05 05 (连接被拒绝)', 'error');
        }

      } catch (error) {
        this.addTerminal(`❌ 错误: ${error.message}`, 'error');
      } finally {
        this.testing = false;
      }
    },

    // 清空终端
    clearTerminal() {
      this.terminalLines = [];
      this.currentStep = -1;
      this.testSuccess = false;
      this.steps.forEach(s => s.failed = false);
    },

    // 添加终端行
    addTerminal(text, type = 'info') {
      const now = new Date();
      const timestamp = now.toLocaleTimeString('zh-CN', { hour12: false });
      this.terminalLines.push({ text, type, timestamp });

      this.$nextTick(() => {
        const terminal = this.$refs.terminal;
        if (terminal) {
          terminal.scrollTop = terminal.scrollHeight;
        }
      });
    },

    // 延迟函数
    delay(ms) {
      return new Promise(resolve => setTimeout(resolve, ms));
    },

    // 获取步骤CSS类
    getStepClass(stepIndex) {
      if (this.currentStep > stepIndex) return 'completed';
      if (this.currentStep === stepIndex) return 'active';
      if (this.steps[stepIndex].failed) return 'failed';
      return '';
    },

    // 格式化认证数据包
    formatAuthPacket() {
      const username = this.config.username;
      const password = this.config.password;
      const usernameHex = this.stringToHex(username);
      const passwordHex = this.stringToHex(password);
      return `01 ${username.length.toString(16).toUpperCase()} ${usernameHex} ${password.length.toString(16).toUpperCase()} ${passwordHex}`;
    },

    // 格式化请求数据包
    formatRequestPacket() {
      const target = this.config.target;
      const targetHex = this.stringToHex(target);
      const port = this.config.targetPort;
      const portHigh = (port >> 8).toString(16).toUpperCase().padStart(2, '0');
      const portLow = (port & 0xFF).toString(16).toUpperCase().padStart(2, '0');
      return `05 01 00 03 ${target.length.toString(16).toUpperCase()} ${targetHex} ${portHigh} ${portLow}`;
    },

    // 字符串转十六进制
    stringToHex(str) {
      return str.split('').map(c => c.charCodeAt(0).toString(16).toUpperCase()).join(' ');
    }
  },

  mounted() {
    this.addTerminal('🚀 SOCKS5 协议分析工具已就绪', 'success');
    this.addTerminal('💡 配置SOCKS5服务器后点击"开始测试连接"', 'info');
    this.addTerminal('', 'info');
    this.addTerminal('预配置测试账号:', 'info');
    this.addTerminal('  - testuser / testpass', 'info');
    this.addTerminal('  - alice / password123', 'info');
    this.addTerminal('  - bob / securepass', 'info');
  }
}).mount('#app');
