# Relay Panel 开发文档

本文档旨在帮助开发者快速上手 Relay Panel 的开发工作。

## 1. 项目概述

Relay Panel 是一个基于 `flux-panel` 重构的流量转发管理面板。它提供了用户管理、节点管理、隧道管理、流量转发和速度限制等功能，旨在为用户提供一个稳定、高效、易于使用的流量管理解决方案。

## 2. 技术栈

*   **后端**：
    *   Go
    *   Gin (Web 框架)
    *   GORM (ORM 框架)
    *   JWT (认证)
*   **前端**：
    *   React
    *   TypeScript
    *   Vite (构建工具)
    *   @heroui (UI 组件库)
*   **数据库**：
    *   SQLite (默认)
    *   MySQL (可配置)

## 3. 环境准备

在开始开发之前，请确保你的开发环境已经安装并配置好以下工具：

*   Go (1.18 或更高版本)
*   Node.js (16 或更高版本)
*   pnpm (推荐使用的包管理器)

## 4. 后端开发

### 4.1. 目录结构

```
backend/
├── controllers/    # 控制器，处理业务逻辑
├── database/       # 数据库连接和初始化
├── gost_api/       # 与 gost 服务交互的 API
├── middlewares/    # 中间件，如认证
├── models/         # 数据模型 (GORM)
└── main.go         # 应用入口
```

### 4.2. 本地开发

1.  **进入后端目录**：
    ```bash
    cd relay-panel/backend
    ```

2.  **安装依赖**：
    ```bash
    go mod tidy
    ```

3.  **启动服务**：
    ```bash
    go run main.go
    ```

    后端服务将在 `http://localhost:8088` 启动。

4.  **环境变量**：
    为了安全起见，JWT 密钥应通过环境变量进行配置。在启动后端服务之前，请设置 `JWT_KEY` 环境变量：
    ```bash
    export JWT_KEY="your_custom_secure_secret_key"
    ```
    如果不设置，系统将在每次启动时生成一个随机密钥，这会导致用户在服务重启后需要重新登录。


## 5. 前端开发

### 5.1. 目录结构

```
frontend/
├── public/           # 静态资源
├── src/
│   ├── api/          # API 请求模块
│   ├── assets/       # 静态资源 (图片、样式等)
│   ├── components/   # 可复用组件
│   ├── layouts/      # 页面布局
│   ├── pages/        # 页面组件
│   ├── utils/        # 工具函数
│   └── main.tsx      # 应用入口
└── package.json
```

### 5.2. 本地开发

1.  **进入前端目录**：
    ```bash
    cd relay-panel/frontend
    ```

2.  **安装依赖**：
    ```bash
    pnpm install
    ```

3.  **启动服务**：
    ```bash
    pnpm dev
    ```

    前端开发服务将在 `http://localhost:5173` (或 Vite 指定的其他端口) 启动。

4.  **代理配置**：
    前端使用 Vite 的代理功能将 API 请求转发到后端。配置文件为 `vite.config.ts`。默认配置下，所有 `/api` 开头的请求都会被转发到 `http://localhost:8088`。

## 6. API 接口

所有 API 的基础路径为 `/api`。

### 6.1. 公共接口

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/register` | `POST` | 用户注册 |
| `/login` | `POST` | 用户登录 |
| `/flow/upload` | `POST` | 上报流量数据 |
| `/captcha/generate` | `GET` | 生成验证码 |
| `/captcha/:captchaId` | `GET` | 获取验证码图片 |
| `/captcha/verify` | `POST` | 校验验证码 |

### 6.2. 认证接口

以下接口需要`Authorization: Bearer <token>` 请求头。

#### 6.2.1. 用户管理

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/users` | `GET` | 获取所有用户 |
| `/users/:id` | `PUT` | 更新用户信息 |
| `/users/:id` | `DELETE` | 删除用户 |
| `/users/update_password`| `POST` | 修改密码 |
| `/users/:id/reset_traffic`| `POST` | 重置用户流量 |
| `/users/:id/tunnels`| `GET` | 获取用户的隧道列表 |
| `/users/:id/tunnels/:tunnel_id` | `DELETE`| 移除用户的隧道 |

#### 6.2.2. 节点管理

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/nodes` | `POST` | 创建节点 |
| `/nodes` | `GET` | 获取所有节点 |
| `/nodes/:id` | `GET` | 获取单个节点 |
| `/nodes/:id/install` | `POST` | 获取节点安装命令 |
| `/nodes/:id` | `PUT` | 更新节点 |
| `/nodes/:id` | `DELETE` | 删除节点 |

#### 6.2.3. 隧道管理

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/tunnels` | `POST` | 创建隧道 |
| `/tunnels` | `GET` | 获取所有隧道 |
| `/tunnels/:id` | `GET` | 获取单个隧道 |
| `/tunnels/:id` | `PUT` | 更新隧道 |
| `/tunnels/:id` | `DELETE` | 删除隧道 |
| `/tunnels/assign` | `POST` | 为用户分配隧道 |
| `/tunnels/:id/diagnose` | `GET` | 诊断隧道 |

#### 6.2.4. 转发管理

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/forwards` | `POST` | 创建转发规则 |
| `/forwards` | `GET` | 获取所有转发规则 |
| `/forwards/:id` | `GET` | 获取单个转发规则 |
| `/forwards/:id` | `PUT` | 更新转发规则 |
| `/forwards/:id` | `DELETE` | 删除转发规则 |
| `/forwards/:id/pause`| `POST` | 暂停转发 |
| `/forwards/:id/resume`| `POST` | 恢复转发 |
| `/forwards/:id/diagnose`| `GET` | 诊断转发 |
| `/forwards/reorder` | `POST` | 重新排序转发规则 |

#### 6.2.5. 速度限制

| 接口路径 | HTTP 方法 | 功能描述 |
| :--- | :--- | :--- |
| `/speedlimits` | `POST` | 创建速度限制 |
| `/speedlimits` | `GET` | 获取所有速度限制 |
| `/speedlimits/:id`| `GET` | 获取单个速度限制 |
| `/speedlimits/:id`| `PUT` | 更新速度限制 |
| `/speedlimits/:id`| `DELETE`| 删除速度限制 |
