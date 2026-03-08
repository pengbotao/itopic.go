```
{
    "url": "openclaw-installation-and-usage-guide",
    "time": "2026/03/07 20:15",
    "tag": "OpenClaw,AI,聊天机器人,安装,使用"
}
```

# OpenClaw 安装与使用指南

> **基于实际安装经验编写** - 记录了在 Rocky Linux 10 系统上的完整部署过程

## 什么是 OpenClaw

OpenClaw 是一个强大的个人AI助手框架，提供以下核心功能：

- 🚀 **多渠道支持**：飞书、Telegram、WhatsApp、Discord、Signal等20+消息平台
- 🧠 **智能对话**：基于大语言模型的自然语言交互
- 🛠️ **技能系统**：丰富的技能生态（飞书多维表格、日历、文档等）
- 🔒 **本地部署**：完全本地运行，保护数据隐私
- 💾 **记忆管理**：长期记忆和上下文管理
- ⚡ **工具集成**：文件操作、网络请求、系统命令等

## 系统要求

### 硬件要求
- **CPU**：x86_64 架构，推荐 4 核心以上
- **内存**：最低 2GB，推荐 4GB 以上
- **存储**：至少 10GB 可用空间
- **网络**：稳定的互联网连接

### 软件要求
- **操作系统**：Linux (Ubuntu 20.04+, CentOS 8+, Rocky Linux 8+)
- **Node.js**：v24.14.0 或更高版本
- **包管理器**：npm 或 nvm

## 第一部分：实际安装步骤（Rocky Linux 10）

### 1.1 系统环境检查

```
# 查看系统版本
cat /etc/os-release

# 检查 Node.js 版本
node --version

# 如果没有 Node.js，安装 Node.js
curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -
dnf install -y nodejs
```

### 1.2 安装 OpenClaw

#### 方法一：使用 npm 安装（推荐）
```
# 全局安装 OpenClaw
npm install -g openclaw

# 验证安装
openclaw --version
```

#### 方法二：使用 nvm 管理（推荐开发者）
```
# 安装 nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 加载 nvm
source ~/.bashrc

# 安装最新的 LTS Node.js
nvm install --lts
nvm use --lts

# 安装 OpenClaw
npm install -g openclaw
```

### 1.3 验证安装

```
# 检查版本
openclaw --version

# 查看帮助信息
openclaw --help

# 检查系统状态
openclaw status
```

## 第二部分：AI模型配置

### 2.1 配置智谱AI模型（主要使用）

```
# 设置智谱AI API 密钥
openclaw config set auth.profiles.zai:default.apiKey your-zai-api-key

# 设置默认模型
openclaw config set models.default zai:default

# 验证配置
openclaw config validate
```

### 2.2 配置OpenAI模型（备选）

```
# 设置 OpenAI API 密钥
openclaw config set auth.profiles.openai:default.apiKey your-openai-key

# 添加 OpenAI 模型
openclaw models add --provider openai --model gpt-4
```

### 2.3 查看当前配置

```
# 查看完整配置
openclaw config show

# 查看可用模型
openclaw models list

# 查看API配置
openclaw config get auth
```

## 第三部分：Gateway 服务配置

### 3.1 启动 Gateway 服务

```
# 启动 Gateway 服务
openclaw gateway

# 后台启动
openclaw gateway --daemon

# 指定端口启动
openclaw gateway --port 18789
```

### 3.2 检查服务状态

```
# 查看服务状态
openclaw status

# 查看健康状态
openclaw health

# 查看实时日志
openclaw logs --follow
```

### 3.3 服务管理

```
# 停止 Gateway
openclaw gateway --stop

# 重启 Gateway
openclaw gateway --restart

# 查看配置文件路径
openclaw config file
```

## 第四部分：飞书渠道配置

### 4.1 创建飞书应用

1. 访问 [飞书开放平台](https://open.feishu.cn/)
2. 创建企业自建应用
3. 配置应用权限：
   - 获取与发送单聊、群聊消息
   - 读取用户与群组信息
   - 多维表格操作权限
   - 日历读写权限
   - 文档编辑权限

### 4.2 配置Webhook

#### 飞书平台配置
- **请求网址**：`http://your-server-ip:18789/feishu/webhook`
- **验证 Token**：自定义字符串
- **消息加密密钥**：自动生成

#### OpenClaw 中配置
```
# 交互式配置飞书
openclaw channels login feishu

# 或手动配置
openclaw config set channels.feishu.appId your-feishu-app-id
openclaw config set channels.feishu.appSecret your-feishu-app-secret
```

### 4.3 验证连接

```
# 查看渠道状态
openclaw channels status

# 测试消息发送
openclaw message send --channel feishu --target "test-group" --message "Hello from OpenClaw!"
```

## 第五部分：常用命令大全

### 5.1 基础命令

```
# 版本信息
openclaw --version
openclaw --help

# 状态检查
openclaw status
openclaw health
openclaw doctor

# 配置管理
openclaw config show
openclaw config get <key>
openclaw config set <key> <value>
openclaw config validate
```

### 5.2 Gateway 管理

```
# 启动服务
openclaw gateway
openclaw gateway --daemon
openclaw gateway --port 19001

# 服务管理
openclaw gateway --stop
openclaw gateway --restart
openclaw gateway --force

# 日志管理
openclaw logs
openclaw logs --follow
openclaw logs --tail 100
openclaw logs --level info
```

### 5.3 渠道管理

```
# 渠道操作
openclaw channels list
openclaw channels status
openclaw channels login <channel>
openclaw channels logout <channel>

# 消息发送
openclaw message send --channel <channel> --target <target> --message <message>
openclaw message reply --session <session> --message <message>
```

### 5.4 技能和插件

```
# 技能管理
openclaw skills list
openclaw plugins list
openclaw plugins install <plugin>
openclaw plugins enable <plugin>
openclaw plugins disable <plugin>
```

### 5.5 会话管理

```
# 会话操作
openclaw sessions list
openclaw sessions show <session-id>
openclaw agent --message "test message"
```

## 第六部分：配置文件详解

### 6.1 配置文件位置

```
# 主配置文件
~/.openclaw/openclaw.json

# 工作区配置
~/.openclaw/workspace/

# 技能目录
~/.openclaw/skills/

# 日志目录
~/.openclaw/logs/
```

### 6.2 完整配置示例

```json
{
  "meta": {
    "lastTouchedVersion": "2026.3.2",
    "lastTouchedAt": "2026-03-07T20:15:00.000Z"
  },
  "auth": {
    "profiles": {
      "zai:default": {
        "provider": "zai",
        "mode": "api_key",
        "apiKey": "your-zai-api-key-here"
      },
      "openai:default": {
        "provider": "openai", 
        "mode": "api_key",
        "apiKey": "your-openai-api-key-here"
      }
    }
  },
  "models": {
    "mode": "merge",
    "providers": {
      "zai": {
        "baseUrl": "https://open.bigmodel.cn/api/coding/paas/v4",
        "api": "openai-completions",
        "models": [
          {
            "id": "glm-4.6",
            "name": "GLM-4.6"
          }
        ]
      },
      "openai": {
        "baseUrl": "https://api.openai.com/v1",
        "api": "openai-completions", 
        "models": [
          {
            "id": "gpt-4",
            "name": "GPT-4"
          }
        ]
      }
    },
    "default": "zai:default",
    "fallback": "zai:default"
  },
  "channels": {
    "feishu": {
      "enabled": true,
      "appId": "your-feishu-app-id",
      "appSecret": "your-feishu-app-secret",
      "verifierToken": "your-verifier-token",
      "encryptKey": "your-encrypt-key"
    }
  },
  "gateway": {
    "port": 18789,
    "host": "0.0.0.0"
  }
}
```

## 第七部分：故障排除

### 7.1 常见问题及解决方案

#### 问题1：API Rate Limit 错误
```
# 检查当前模型配置
openclaw config get models

# 切换到备用模型
openclaw config set models.fallback zai:default

# 重启服务
openclaw gateway --restart
```

#### 问题2：端口被占用
```
# 查看端口占用
sudo netstat -tlnp | grep :18789

# 强制启动
openclaw gateway --force

# 或更换端口
openclaw gateway --port 19001
```

#### 问题3：飞书连接失败
```
# 检查网络连通性
curl -I http://your-server-ip:18789

# 重新配置飞书
openclaw channels logout feishu
openclaw channels login feishu

# 查看详细日志
openclaw logs --follow
```

#### 问题4：权限不足
```
# 检查工作区权限
ls -la ~/.openclaw/

# 修复权限
chmod -R 755 ~/.openclaw/

# 检查系统权限
openclaw doctor
```

### 7.2 日志分析

```
# 查看错误日志
openclaw logs --level error

# 查看最近错误
openclaw logs --tail 50 | grep ERROR

# 查看飞书相关日志
openclaw logs --follow | grep feishu
```

### 7.3 重置配置

```
# 重置所有配置（谨慎使用）
openclaw reset

# 备份配置
cp ~/.openclaw/openclaw.json ~/.openclaw/openclaw.json.backup

# 重置特定配置
openclaw config unset channels.feishu
```

## 第八部分：高级配置

### 8.1 环境变量配置

```
# 在 ~/.bashrc 中添加
export OPENAI_API_KEY=your-openai-key
export ZAI_API_KEY=your-zai-key
export OPENCLAW_GATEWAY_PORT=18789
export OPENCLAW_LOG_LEVEL=info

# 重新加载
source ~/.bashrc
```

### 8.2 系统服务配置（生产环境）

创建 systemd 服务文件 `/etc/systemd/system/openclaw.service`：

```
[Unit]
Description=OpenClaw Gateway Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/.openclaw
ExecStart=/usr/local/bin/openclaw gateway --daemon
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

启用服务：
```
sudo systemctl daemon-reload
sudo systemctl enable openclaw
sudo systemctl start openclaw
sudo systemctl status openclaw
```

### 8.3 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:18789;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 第九部分：监控和维护

### 9.1 性能监控

```
# 查看系统资源使用
top -p $(pgrep openclaw)

# 查看内存使用
ps aux | grep openclaw

# 查看连接状态
netstat -an | grep :18789
```

### 9.2 定期维护

```
# 查看磁盘使用
du -sh ~/.openclaw/

# 清理旧日志
find ~/.openclaw/logs -name "*.log" -mtime +7 -delete

# 检查配置文件
openclaw config validate
```

### 9.3 备份和恢复

```
# 备份配置
tar -czf openclaw-backup-$(date +%Y%m%d).tar.gz ~/.openclaw/

# 恢复配置
tar -xzf openclaw-backup-20260307.tar.gz -C ~/
```

## 第十部分：开发扩展

### 10.1 创建自定义技能

```
# 创建技能目录
mkdir -p ~/.openclaw/skills/my-custom-skill

# 创建技能文件
cat > ~/.openclaw/skills/my-custom-skill/SKILL.md << 'EOF'
# My Custom Skill

## Description
This is a custom skill for OpenClaw.

## Usage
When to use this skill...

## Implementation
Add your implementation details here.
EOF
```

### 10.2 插件开发

```
# 列出插件开发命令
openclaw plugins --help

# 安装本地插件进行测试
openclaw plugins install ./my-plugin

# 查看插件状态
openclaw plugins list
```

## 实际使用技巧

### 1. 日常管理命令
```
# 快速重启服务
openclaw gateway --restart

# 查看实时状态
watch -n 5 'openclaw status'

# 快速查看错误
openclaw logs --tail 20 --level error
```

### 2. 配置最佳实践
```
# 使用版本控制管理配置
git init ~/.openclaw
git add ~/.openclaw/openclaw.json
git commit -m "Initial OpenClaw config"
```

### 3. 监控脚本
```bash
#!/bin/bash
# 监控 OpenClaw 服务状态
if ! pgrep -f "openclaw gateway" > /dev/null; then
    echo "OpenClaw Gateway is down, restarting..."
    openclaw gateway --daemon
fi
```

## 总结

通过本指南，你应该能够：

✅ **成功安装** OpenClaw 在 Rocky Linux 系统上  
✅ **配置AI模型** 智谱AI和OpenAI  
✅ **启动Gateway服务** 并进行基本管理  
✅ **配置飞书渠道** 实现企业级集成  
✅ **掌握常用命令** 进行日常维护  
✅ **解决常见问题** 保证服务稳定运行  
✅ **进行高级配置** 满足生产环境需求  

### 快速参考

| 命令 | 功能 |
|------|------|
| `openclaw status` | 查看服务状态 |
| `openclaw gateway --restart` | 重启Gateway |
| `openclaw logs --follow` | 查看实时日志 |
| `openclaw config show` | 查看完整配置 |
| `openclaw channels status` | 查看渠道状态 |
| `openclaw doctor` | 运行诊断检查 |

### 相关资源

- [OpenClaw 官方文档](https://docs.openclaw.ai)
- [OpenClaw GitHub](https://github.com/openclaw/openclaw)
- [飞书开放平台](https://open.feishu.cn/)
- [社区 Discord](https://discord.gg/clawd)
- [技能市场](https://clawhub.com)

---

**最后更新：** 2026年3月7日  
**文档版本：** v3.0.0 (实际安装经验版)  
**测试环境：** Rocky Linux 10.0 + OpenClaw 2026.3.2