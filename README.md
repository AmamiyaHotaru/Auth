# Euthenticator

<div align="center">
  <img src="build/appicon.png" alt="Euthenticator Logo" width="150">
  <h3>纯本地的 TOTP 身份验证器</h3>
</div>

## 简介

Euthenticator 是一个纯本地的 TOTP (基于时间的一次性密码) 安全令牌生成器，为您的在线账户提供额外的安全保障。所有数据均存储在本地，无需任何网络连接，保证您的账户安全。

## 主要特性

- 🔒 **纯本地存储**：所有验证码数据均加密存储在本地设备上
- 🔄 **支持多种导入**：兼容 Google Authenticator 和 Steam 令牌导入
- 👁️ **验证码管理**：可以一键显示/隐藏所有验证码，保证账户安全
- 📋 **便捷复制**：点击即可复制验证码到剪贴板
- 🖥️ **跨平台**：支持 Windows 和 Linux 系统
- 🎨 **简洁界面**：桌面端友好的用户界面设计

## 下载安装

访问 [GitHub Releases](https://github.com/AmamiyaHotaru/Auth/releases) 下载最新版本。

### Windows

1. 下载最新的 `Euthenticator-windows-vX.X.X.zip` 文件
2. 解压到任意位置
3. 运行 `Euthenticator.exe`

### Linux

1. 下载最新的 `Euthenticator-linux-vX.X.X.tar.gz` 文件
2. 解压到任意位置：`tar -xzvf Euthenticator-linux-vX.X.X.tar.gz -C /path/to/destination`
3. 运行应用程序

## 使用说明

### 添加新账户

Euthenticator 支持多种方式添加新的验证账户：

1. **解析二维码**
   - 点击右上角"添加账户"按钮
   - 选择"解析二维码"选项
   - 上传或粘贴包含 TOTP 信息的二维码图片

2. **手动输入**
   - 点击右上角"添加账户"按钮
   - 选择"手动输入"选项
   - 填写账户名称、服务名称和密钥信息

### 管理验证码

- **显示/隐藏单个验证码**：点击每个账户卡片上的眼睛图标
- **显示/隐藏所有验证码**：点击顶部工具栏的"显示所有"/"隐藏所有"按钮
- **复制验证码**：点击验证码或复制图标将验证码复制到剪贴板

### 删除账户

1. 长按账户卡片进入选择模式
2. 选择要删除的一个或多个账户
3. 点击"删除选中"按钮
4. 确认删除操作

## 手动构建

如果您想自己构建 Euthenticator，请按照以下步骤操作：

### 前置要求

- [Go](https://golang.org/dl/) 1.23 或更高版本
- [Node.js](https://nodejs.org/) 20 或更高版本
- [Wails](https://wails.io/) v2
- 适用于您平台的构建工具
  - **Windows**：gcc (可通过 MinGW 或 TDM-GCC 获取)
  - **Linux**：gcc、libgtk-3-dev、libwebkit2gtk-4.0-dev 或 libwebkit2gtk-4.1-dev

### 安装 Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 克隆代码库

```bash
git clone https://github.com/AmamiyaHotaru/Auth.git
cd Auth
```

### 配置环境变量

创建 `.env` 文件，设置主密码（用于加密保存的密钥）：

```
MASTER_PASSWORD=你的安全密码
```

### 构建应用

#### Windows

```powershell
# 使用 PowerShell 构建
./build.ps1
```

#### Linux

```bash
# 确保脚本有执行权限
chmod +x ./build.sh
# 执行构建脚本
./build.sh
```

构建完成后，可执行文件将位于 `build/bin` 目录中。

## 技术栈

Euthenticator 使用以下技术构建：

- **前端**：Vue.js、Vite
- **后端**：Go
- **桌面应用框架**：Wails
- **加密**：Go 标准库加密功能

## 贡献

欢迎提交 Issue 或 Pull Request 来改进 Euthenticator。

## 联系方式

开发者：AmamiyaHotaru  
邮箱：eunie@eunie.online  
项目仓库：[https://github.com/AmamiyaHotaru/Auth](https://github.com/AmamiyaHotaru/Auth)
