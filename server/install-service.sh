#!/bin/bash
# 网课代刷管理系统 wk-go systemd 服务注册脚本
# 使用方法：bash install-service.sh [部署目录] [服务名]
# 示例：bash install-service.sh /www/wwwroot/wk.iyys.ltd wk-go

DEPLOY_DIR=${1:-/www/wwwroot/wk}
SERVICE_NAME=${2:-wk-go}
SERVER_BIN="$DEPLOY_DIR/wk-go"
WORK_DIR="$DEPLOY_DIR"
ENV_FILE="$DEPLOY_DIR/.env"

echo "部署目录: $DEPLOY_DIR"
echo "服务名称: $SERVICE_NAME"
echo "二进制路径: $SERVER_BIN"

# 检查二进制文件是否存在
if [ ! -f "$SERVER_BIN" ]; then
    echo "错误: 未找到 $SERVER_BIN，请先上传后端二进制文件"
    exit 1
fi

# 赋予执行权限
chmod +x "$SERVER_BIN"

# 写入 systemd service 文件
cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=网课代刷管理系统
After=network.target

[Service]
Type=simple
WorkingDirectory=${WORK_DIR}
ExecStart=${SERVER_BIN}
Restart=always
RestartSec=5
EnvironmentFile=-${ENV_FILE}

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"

echo ""
echo "服务注册完成！"
echo "查看状态: systemctl status $SERVICE_NAME"
echo "查看日志: journalctl -u $SERVICE_NAME -f"
