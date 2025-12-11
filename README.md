# DeepSeek Demo 项目

这是一个使用 Golang 开发的简单 Web 应用程序，演示了如何与 DeepSeek API 集成，提供了一个聊天界面让用户与 AI 进行交互。

## 功能特性

- 基于 Web 的聊天界面
- 与 DeepSeek API 集成
- 实时逐字响应显示（打字机效果）
- 美观的用户界面

## 技术栈

- 后端：Golang + Gin 框架
- 前端：HTML + CSS + JavaScript
- 环境管理：godotenv

## 安装和运行

1. 克隆项目到本地：
   ```
   git clone <repository-url>
   cd deepseek-demo
   ```

2. 安装 Go 依赖：
   ```
   go mod tidy
   ```

3. 配置环境变量：
   - 编辑 `.env` 文件
   - 设置你的 DeepSeek API 密钥：
     ```
     DEEPSEEK_API_KEY=your_actual_api_key_here
     ```

4. 运行应用程序：
   ```
   go run main.go
   ```

5. 在浏览器中访问 `http://localhost:8082`

## 项目结构

```
.
├── main.go              # 主程序文件
├── go.mod               # Go 模块定义
├── go.sum               # Go 依赖校验和
├── .env                 # 环境变量配置文件
├── README.md            # 项目说明文档
└── static/              # 静态资源目录
    ├── index.html       # 主页
    ├── style.css        # 样式表
    └── script.js        # JavaScript 脚本
```

## 配置说明

在 `.env` 文件中，你可以配置以下参数：

- `DEEPSEEK_API_KEY`: 你的 DeepSeek API 密钥（必需）
- `DEEPSEEK_API_URL`: DeepSeek API 的基础 URL（可选，默认为 https://api.deepseek.com/v1）
- `SERVER_PORT`: 服务器监听端口（可选，默认为 8082）

## 使用方法

1. 启动应用后，在浏览器中打开 `http://localhost:8082`
2. 在输入框中输入你想与 AI 交流的内容
3. 点击"发送"按钮或按回车键发送消息
4. 等待 AI 回复显示在聊天窗口中（具有逐字显示的打字机效果）

## 注意事项

- 你需要有一个有效的 DeepSeek API 密钥才能正常使用此应用
- 请妥善保管你的 API 密钥，不要泄露给他人
- 本项目仅供学习和演示用途