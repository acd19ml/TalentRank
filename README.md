# TalentRank
## 配置环境
豆包模型，到以下路径中找到app.go，修改Model使用自己的Model号
```
TalentRank/apps/user
```
引入豆包环境变量（linux系统）
```
echo 'export ARK_API_KEY=<YOUR_API_KEY>' >> ~/.bashrc
source ~/.bashrc
echo $ARK_API_KEY
```
引入Github环境变量（linux系统）
```
echo 'export GITHUB_TOKEN=<Your_GitHub_Token>' >> ~/.bashrc
source ~/.bashrc
echo $GITHUB_TOKEN
```
## 运行代码
启动前端
```
make runf
```
启动后端
```
make run
```
## 项目架构
详细文档链接（包含项目设计，分工）
https://kcnvb52mvl2f.feishu.cn/docx/BgfBdoagUojt2vxf3ecccGdcnSc?from=from_copylink
```
TalentRank
├── Makefile             # 构建脚本
├── README.md            # 项目描述文件
├── go.mod               # Go 模块配置文件
├── go.sum               # Go 依赖文件校验和
├── main.go              # Go 主程序入口
├── apps                 # 后端服务模块
│   ├── all              # 通用实现
│   │   └── impl.go      # 通用实现代码
│   ├── git              # Git 服务模块
│   │   ├── app.go       # Git 模块主文件
│   │   ├── git.pb.go    # Protobuf 生成的文件
│   │   ├── git_grpc.pb.go # gRPC 文件
│   │   ├── impl         # Git 实现文件夹
│   │   └── pb           # Protobuf 文件夹
│   ├── llm              # LLM 服务模块
│   │   ├── app.go       # LLM 模块主文件
│   │   ├── llm.pb.go    # Protobuf 生成的文件
│   │   ├── llm_grpc.pb.go # gRPC 文件
│   │   ├── impl         # LLM 实现文件夹
│   │   └── pb           # Protobuf 文件夹
│   └── user             # 用户服务模块
│       ├── app.go       # 用户服务主文件
│       ├── gpt.go       # GPT 相关文件
│       ├── http         # 用户 HTTP 实现
│       ├── impl         # 用户服务实现
│       ├── interface.go # 用户服务接口定义
│       ├── json.go      # JSON 处理文件
│       ├── json_test.go # JSON 处理测试文件
│       └── model.go     # 用户数据模型
├── cmd                  # 命令行相关代码
├── conf                 # 配置管理模块
├── etc                  # 配置文件夹 (如 demo.toml)
├── gpt                  # GPT 模块文件夹
├── middleware           # 中间件
├── protocol             # 通信协议
├── version              # 版本控制文件夹
└── my-app               # 前端应用 (可能为 React 应用)
    ├── README.md        # 前端应用说明
    ├── node_modules     # Node.js 依赖
    ├── package.json     # Node.js 依赖配置文件
    ├── package-lock.json # 锁定依赖版本
    ├── public           # 静态文件目录
    └── src              # 前端源代码
        ├── App.css
        ├── App.js       # 应用主文件
        ├── index.css
        ├── index.js     # 应用入口文件
        ├── pages        # 页面组件
        ├── reportWebVitals.js
        └── setupTests.js

```
##后端架构图
![service接口 (1)](https://github.com/user-attachments/assets/d0317a53-a925-4ddc-86f8-ad9ae5a32579)

