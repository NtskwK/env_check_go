# env_checkGo

[![Build and Release](https://github.com/Ntskwk/env_check_go/actions/workflows/build.yaml/badge.svg)](https://github.com/Ntskwk/env_check_go/actions/workflows/build.yaml)
[![License](https://img.shields.io/badge/license-LGPL%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://go.dev/)

一个跨平台的系统环境检测工具，用于快速获取系统硬件信息和运行时环境。

## ✨ 功能特性

- 🖥️ **系统信息检测**

  - 当前时间
  - 操作系统版本
  - CPU 型号
  - 内存大小
  - GPU 信息（支持多显卡）

- 🔧 **运行时环境检测**

  - .NET Runtime 版本检测（跨平台）
  - Visual C++ Runtime 检测（Windows）

- 🌍 **跨平台支持**
  - Windows (amd64/arm64)
  - Linux (amd64/arm64)
  - macOS (amd64/arm64)

## 📦 安装

### 方式一：下载预编译二进制文件

从 [Releases](https://github.com/YOUR_USERNAME/env_check_go/releases) 页面下载适合您系统的预编译版本。

### 方式二：从源码构建

**前置要求：**

- Go 1.21 或更高版本

**克隆并构建：**

```bash
# 克隆仓库
git clone https://github.com/YOUR_USERNAME/env_check_go.git
cd env_check_go

# 构建当前平台
go build -o env_checkGo .

# 或使用 Makefile 构建所有平台
make build
```

## 🚀 使用方法

直接运行可执行文件：

```bash
# Windows
.\env_checkGo-windows-amd64.exe

# Linux
./env_checkGo-linux-amd64

# macOS
./env_checkGo-darwin-amd64
```

### 输出示例

```
Time: 2025-12-02 00:23:54
OS: Microsoft Windows 11 Pro
CPU: Intel(R) Core(TM) i7-12700K CPU @ 3.60GHz
Memory: 32.00 GB
GPU:
  - NVIDIA GeForce RTX 4090

=== .NET Runtimes ===
Found 3 .NET Runtime(s):
  - Microsoft.NETCore.App 6.0.25
  - Microsoft.NETCore.App 7.0.14
  - Microsoft.NETCore.App 8.0.0

=== Visual C++ Runtimes (Windows Only) ===
Found 4 VC++ Runtime(s):
  - Microsoft Visual C++ 2015-2022 Redistributable (x64) - 14.38.33135
  - Microsoft Visual C++ 2015-2022 Redistributable (x86) - 14.38.33135
  - Microsoft Visual C++ 2013 Redistributable (x64) - 12.0.40664
  - Microsoft Visual C++ 2010 Redistributable (x64) - 10.0.40219
```

## 🛠️ 开发

### 项目结构

```
env_check_go/
├── main.go           # 主程序入口
├── hardware.go       # 硬件信息检测
├── checker.go        # 运行时环境检测
├── go.mod            # Go 模块依赖
├── Makefile          # 构建脚本
└── .github/
    └── workflows/
        └── build.yaml # CI/CD 工作流
```

### 本地开发

```bash
# 安装依赖
go mod download

# 运行程序
go run .

# 运行测试
go test ./...

# 格式化代码
go fmt ./...
```

### 构建所有平台

```bash
make build
```

构建产物将输出到 `release/` 目录。

## 🔄 CI/CD

本项目使用 GitHub Actions 进行自动化构建和发布：

- **自动构建触发条件：**

  - Push 到 `main` 分支
  - Pull Request 到 `main` 分支
  - 手动触发

- **自动发布：**
  - 推送以 `v` 开头的 tag（如 `v1.0.0`）会自动创建 GitHub Release
  - 所有平台的二进制文件会自动上传到 Release

**创建新版本发布：**

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 📋 依赖项

- [gopsutil](https://github.com/shirou/gopsutil) - 跨平台系统信息库

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📧 联系方式

如有问题或建议，请通过 [Issues](https://github.com/YOUR_USERNAME/env_check_go/issues) 联系我们。

---

⭐ 如果这个项目对您有帮助，请给它一个 Star！
