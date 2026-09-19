# AR 网课代刷聚合管理系统

一站式网课代刷下单管理平台，支持多平台网课代刷、实习打卡、盖章病历聚合管理。

## 功能特性

- **多源台对接**：支持 29 网课源台、实习打卡源台、盖章病历源台
- **商品管理**：分类管理、商品管理、一键同步货源商品
- **订单管理**：网课订单、打卡订单、盖章/病历订单统一管理
- **用户系统**：注册登录、QQ 头像自动拉取、邮箱验证、忘记密码
- **系统配置**：站点名称/Logo/Favicon 自定义、SMTP 邮件配置
- **安装引导**：一键安装向导，自动建表

## 技术栈

- **后端**：Go + Gin + GORM + MySQL
- **前端**：Vue 3 + Vite + Element Plus + UnoCSS
- **部署**：Linux + Nginx + 宝塔面板

## 目录结构

```
├── web/          # 前端（Vue 3）
├── server/       # 后端（Go）
└── README.md
```

## 本地开发

### 后端

```bash
cd server
go mod download
go run ./cmd
```

配置文件：`config/config.yaml`

### 前端

```bash
cd web
npm install
npm run dev
```

## 生产部署

### 编译

```bash
# 前端
cd web && npm run build
# 输出到 dist/，上传到服务器 web 目录

# 后端 Linux 编译
cd server
GOOS=linux GOARCH=amd64 go build -o AR ./cmd
```

### 服务器目录

```
/www/wwwroot/你的域名/
├── AR              # 后端二进制
├── config/config.yaml
├── web/            # 前端 dist
└── uploads/        # 上传目录
```

### Nginx 配置（完整版）

```nginx
server
{
    listen 80;
    listen 443 ssl;
    http2 on;
    server_name your.domain.com;
    index index.html index.htm;
    root /www/wwwroot/your.domain.com/web;

    # SSL 证书（宝塔申请后自动替换路径）
    ssl_certificate    /www/server/panel/vhost/cert/你的域名/fullchain.pem;
    ssl_certificate_key    /www/server/panel/vhost/cert/你的域名/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers EECDH+CHACHA20:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HTTP 跳转 HTTPS
    if ($server_port != 443) {
        rewrite ^(/.*)$ https://$host$1 permanent;
    }

    # API 反代到后端
    location /api/ {
        proxy_pass http://127.0.0.1:8091;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_connect_timeout 30s;
        proxy_read_timeout 300s;
    }

    # 上传文件反代
    location /uploads/ {
        proxy_pass http://127.0.0.1:8091;
        proxy_set_header Host $host;
    }

    # SPA 路由 fallback
    location / {
        try_files $uri $uri/ /index.html;
    }

    access_log  /www/wwwlogs/your.domain.com.log;
    error_log  /www/wwwlogs/your.domain.com.error.log;
}
```

## 默认账号

首次访问会自动进入安装引导，按提示配置数据库和管理员账号即可。
