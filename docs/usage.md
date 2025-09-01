# Relay Panel 使用文档

本文档旨在帮助用户快速上手使用 Relay Panel。

## 1. 系统部署

### 1.1. 环境要求

*   一台具有公网 IP 的服务器
*   Go (1.18 或更高版本)
*   Node.js (16 或更高版本)
*   pnpm
*   Nginx

### 1.2. 编译前端

1.  **进入前端目录**：
    ```bash
    cd relay-panel/frontend
    ```

2.  **安装依赖**：
    ```bash
    pnpm install
    ```

3.  **编译打包**：
    ```bash
    pnpm build
    ```

    编译后的静态文件将位于 `dist` 目录下。

### 1.3. 编译后端

1.  **进入后端目录**：
    ```bash
    cd relay-panel/backend
    ```

2.  **编译**：
    ```bash
    go build -o relay-panel-backend .
    ```

    编译后的可执行文件名为 `relay-panel-backend`。

### 1.4. 部署

1.  将 `relay-panel-backend` 可执行文件和前端 `dist` 目录上传到服务器的同一目录下，例如 `/var/www/relay-panel`。

2.  **设置 JWT_KEY 环境变量**：
    为了保证安全，请设置一个复杂的 JWT 密钥。
    ```bash
    export JWT_KEY="your_custom_secure_secret_key"
    ```

3.  **启动后端服务**：
    建议使用 `systemd` 或 `supervisor` 等工具来管理后端进程。
    ```bash
    ./relay-panel-backend
    ```
    服务将在 `http://localhost:8088` 启动。

### 1.5. 配置 Nginx 反向代理

以下是一个示例 Nginx 配置：

```nginx
server {
    listen 80;
    server_name your_domain.com;

    location / {
        root /var/www/relay-panel/dist;
        try_files $uri /index.html;
    }

    location /api {
        proxy_pass http://localhost:8088;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 2. 功能介绍

### 2.1. 用户注册与登录

访问 `http://your_domain.com` 即可看到登录页面。首次使用请先注册账号。

### 2.2. 节点管理

节点是运行 `gost` 服务的服务器。在 Relay Panel 中，你需要先添加节点，然后才能创建隧道和转发规则。

1.  在左侧菜单栏中，点击 “节点管理”。
2.  点击 “新增节点” 按钮。
3.  填写节点信息，包括节点名称、地址、API 端口、用户名和密码。
4.  保存后，你可以在节点列表中看到新添加的节点。
5.  点击 “安装” 按钮可以获取该节点的 `gost` 安装和配置命令。

### 2.3. 隧道管理

隧道是 `gost` 中的一个概念，它可以在两个节点之间建立一个加密的、可靠的通信通道。

1.  在左侧菜单栏中，点击 “隧道管理”。
2.  点击 “新增隧道” 按钮。
3.  填写隧道信息，选择入口节点和出口节点。
4.  保存后，隧道将自动在相应的节点上创建。

### 2.4. 转发配置

转发规则定义了如何将流量从一个端口转发到另一个地址。

1.  在左侧菜单栏中，点击 “转发管理”。
2.  点击 “新增转发” 按钮。
3.  填写转发规则，包括选择隧道、监听端口、目标地址等。
4.  保存后，转发规则将立即生效。

### 2.5. 速度限制

你可以为特定的隧道或转发规则设置速度限制。

1.  在左侧菜单栏中，点击 “限速管理”。
2.  点击 “新增限速” 按钮。
3.  填写限速规则，设置上传和下载的速率限制。
4.  将限速规则关联到隧道或转发规则。
