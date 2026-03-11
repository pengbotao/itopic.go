```
{
    "url": "openclaw-cron-tutorial",
    "time": "2026/03/11 08:55",
    "tag": "OpenClaw,Cron,定时任务,自动化,运维"
}
```

# OpenClaw Cron 定时任务完全指南

> **从零到精通** - OpenClaw 定时任务的完整使用手册

## 什么是 OpenClaw Cron

OpenClaw Cron 是一个强大的定时任务调度器，集成在 Gateway 服务中。它让你能够：

- ⏰ **定时执行任务**：支持 cron 表达式和间隔执行
- 🤖 **AI 代理集成**：自动触发 Agent 执行特定任务
- 📤 **消息投递**：自动发送消息到各种聊天渠道
- 📊 **执行监控**：查看任务历史和执行状态
- 🔄 **动态管理**：运行时添加、修改、删除任务

## 第一部分：基础概念

### 1.1 Cron 表达式

OpenClaw 支持 5 字段和 6 字段 cron 表达式：

#### 5 字段格式（分钟级精度）
```
* * * * *
│ │ │ │ │
│ │ │ │ └─ 星期几 (0-7, 0和7都表示周日)
│ │ │ └─── 月份 (1-12)
│ │ └───── 日期 (1-31)
│ └─────── 小时 (0-23)
└───────── 分钟 (0-59)
```

#### 6 字段格式（秒级精度）
```
* * * * * *
│ │ │ │ │ │
│ │ │ │ │ └─ 星期几 (0-7)
│ │ │ │ └─── 月份 (1-12)
│ │ │ └───── 日期 (1-31)
│ │ └─────── 小时 (0-23)
│ └───────── 分钟 (0-59)
└─────────── 秒 (0-59)
```

### 1.2 间隔执行

除了 cron 表达式，还支持间隔执行：

```
--every 5m     # 每5分钟
--every 1h     # 每小时
--every 30s    # 每30秒
--every 1d     # 每天一次
```

### 1.3 Cron 特殊字符

| 字符 | 说明 | 示例 |
|------|------|------|
| `*` | 匹配任意值 | `* * * * *` （每分钟） |
| `,` | 列出多个值 | `1,3,5 * * * *` （每小时的第1、3、5分钟） |
| `-` | 指定范围 | `0-5 * * * *` （每小时的0-5分钟） |
| `/` | 指定间隔 | `*/5 * * * *` （每5分钟） |
| `?` | 不指定值（仅用于日期和星期） | `0 0 ? * *` （每天0点0分，不指定日期） |

## 第二部分：快速开始

### 2.1 基本命令

```bash
# 查看所有任务
openclaw cron list

# 查看任务状态
openclaw cron status

# 查看执行历史
openclaw cron runs --id <job-id>
```

### 2.2 创建第一个任务

#### 示例1：简单的 Agent 消息任务
```bash
# 每小时发送一条消息
openclaw cron add \
  --name hourly-reminder \
  --every 1h \
  --message "每小时提醒"
```

#### 示例2：使用 cron 表达式
```bash
# 每天早上9点执行
openclaw cron add \
  --name morning-task \
  --cron "0 9 * * *" \
  --message "早安任务"
```

### 2.3 查看任务详情

```bash
# 列出所有任务
openclaw cron list

# 查看特定任务的执行历史
openclaw cron runs --id <job-id> --limit 10

# 查看任务状态
openclaw cron status
```

## 第三部分：任务管理

### 3.1 添加任务

#### 完整参数示例
```bash
openclaw cron add \
  --name my-scheduled-task \
  --cron "*/5 * * * *" \
  --message "每5分钟执行的任务" \
  --agent main \
  --session isolated \
  --description "这是任务描述" \
  --enabled
```

#### 常用参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--name` | 任务名称 | 必填 |
| `--cron` | Cron 表达式 | 无 |
| `--every` | 间隔执行（与 cron 二选一） | 无 |
| `--message` | Agent 消息内容 | 无 |
| `--agent` | 使用的 Agent | main |
| `--session` | 会话类型 | isolated |
| `--description` | 任务描述 | 无 |
| `--enabled` | 是否启用 | true |
| `--disabled` | 创建时禁用 | false |
| `--exact` | 禁用随机延迟 | false |
| `--timeout` | 超时时间（毫秒） | 30000 |

### 3.2 禁用/启用任务

```bash
# 禁用任务
openclaw cron disable <job-id>

# 启用任务
openclaw cron enable <job-id>
```

### 3.3 编辑任务

```bash
# 编辑任务（使用 JSON 格式）
openclaw cron edit <job-id> \
  --message "更新后的消息" \
  --cron "0 10 * * *"
```

### 3.4 删除任务

```bash
# 删除任务
openclaw cron rm <job-id>
```

### 3.5 手动触发任务

```bash
# 立即执行任务
openclaw cron run <job-id>

# 只在到期时执行
openclaw cron run <job-id> --due
```

## 第四部分：进阶配置

### 4.1 消息投递配置

#### 投递到飞书
```bash
openclaw cron add \
  --name feishu-alert \
  --every 1h \
  --message "每小时检查" \
  --channel feishu \
  --to ou_xxxxxxxxxxxxxxxxxxxx \
  --announce
```

#### 投递到多个渠道
```bash
# Telegram
openclaw cron add \
  --name telegram-alert \
  --every 1h \
  --channel telegram \
  --to @mychannel \
  --message "Telegram 消息"

# WhatsApp
openclaw cron add \
  --name whatsapp-alert \
  --every 1h \
  --channel whatsapp \
  --to +8613800138000 \
  --message "WhatsApp 消息"
```

### 4.2 Agent 任务配置

#### 指定特定 Agent
```bash
# 使用 stock-research Agent
openclaw cron add \
  --name stock-analysis \
  --cron "30 15 * * 1-5" \
  --agent stock-research \
  --message "进行每日股票分析"
```

#### 系统事件任务
```bash
# 触发系统事件
openclaw cron add \
  --name heartbeat \
  --every 1h \
  --system-event "check_system_health" \
  --session main
```

### 4.3 时间和时区配置

```bash
# 指定时区
openclaw cron add \
  --name daily-report \
  --cron "0 18 * * *" \
  --tz "Asia/Shanghai" \
  --message "每日报告"

# 使用 UTC 时间
openclaw cron add \
  --name utc-task \
  --cron "0 0 * * *" \
  --tz "UTC" \
  --message "UTC 任务"
```

### 4.4 一次性任务

```bash
# 在指定时间执行一次
openclaw cron add \
  --name one-time-task \
  --at "2026-03-11T10:00:00" \
  --message "一次性任务" \
  --delete-after-run

# 延迟执行
openclaw cron add \
  --name delayed-task \
  --at "+30m" \
  --message "30分钟后执行" \
  --delete-after-run
```

## 第五部分：实际应用案例

### 5.1 股票监控任务

```bash
# 每分钟监控股票价格
openclaw cron add \
  --name stock-monitor-1m \
  --cron "* * 9-15 * * 1-5" \
  --agent stock-research \
  --message "监控市场动态" \
  --exact

# 每日收盘总结
openclaw cron add \
  --name daily-stock-summary \
  --cron "5 15 * * 1-5" \
  --agent stock-research \
  --message "生成每日收盘总结" \
  --channel feishu \
  --to oc_xxxxxxxxxxxxxxxx \
  --announce
```

### 5.2 系统健康检查

```bash
# 每10分钟检查系统状态
openclaw cron add \
  --name system-health-check \
  --every 10m \
  --message "检查服务器健康状态：CPU、内存、磁盘使用率" \
  --session main

# 每日凌晨3点备份数据
openclaw cron add \
  --name daily-backup \
  --cron "0 3 * * *" \
  --message "执行数据库备份" \
  --timeout 300000
```

### 5.3 定期报告

```bash
# 每周工作总结
openclaw cron add \
  --name weekly-report \
  --cron "0 18 * * 5" \
  --message "生成本周工作总结报告" \
  --channel feishu \
  --to ou_xxxxxxxxxxxxxxxx

# 每月数据统计
openclaw cron add \
  --name monthly-stats \
  --cron "0 0 1 * *" \
  --message "生成上月数据统计" \
  --channel telegram \
  --to @admin-channel
```

### 5.4 提醒和通知

```bash
# 上班提醒
openclaw cron add \
  --name work-reminder \
  --cron "0 9 * * 1-5" \
  --message "提醒：今天是工作日，开始工作啦！" \
  --channel feishu \
  --to ou_xxxxxxxxxxxxxxxx

# 会议提醒
openclaw cron add \
  --name meeting-reminder \
  --cron "0 14 * * 2-4" \
  --message "下午2点有团队会议" \
  --channel feishu \
  --to oc_xxxxxxxxxxxxxxxx
```

### 5.5 定时发送时间

```bash
# 每2分钟发送当前时间（演示用）
openclaw cron add \
  --name time-alert-2m \
  --every 2m \
  --message "发送当前时间，格式：当前时间 [时间]" \
  --channel feishu \
  --to ou_xxxxxxxxxxxxxxxx \
  --announce
```

## 第六部分：执行历史和监控

### 6.1 查看执行历史

```bash
# 查看所有任务的最近执行
openclaw cron runs --limit 20

# 查看特定任务的历史
openclaw cron runs --id <job-id> --limit 50

# 查看 JSON 格式的详细输出
openclaw cron runs --id <job-id> --json
```

### 6.2 理解执行状态

| 状态 | 说明 |
|------|------|
| `idle` | 等待执行 |
| `running` | 正在执行 |
| `ok` | 执行成功 |
| `error` | 执行失败 |
| `timeout` | 执行超时 |

### 6.3 错误诊断

```bash
# 查看失败的任务历史
openclaw cron runs --id <job-id> | grep error

# 查看 Gateway 日志
tail -f /tmp/openclaw/openclaw-$(date +%Y-%m-%d).log | grep cron

# 查看任务配置
cat ~/.openclaw/cron/jobs.json
```

## 第七部分：常见问题

### 7.1 任务不执行

**可能原因：**
1. Gateway 服务未运行
2. 任务被禁用
3. cron 表达式错误
4. 时区配置问题

**解决方案：**
```bash
# 1. 检查 Gateway 状态
openclaw gateway status

# 2. 检查任务是否启用
openclaw cron list

# 3. 验证 cron 表达式
openclaw cron add --name test --cron "*/5 * * * *" --message "测试"

# 4. 检查时区
openclaw cron runs --id <job-id>
```

### 7.2 任务执行失败

**常见错误及解决：**

#### 错误1：Target 未指定
```
Delivering to Feishu requires target <chatId|user:openId|chat:chatId>
```

**解决：**
```bash
# 添加正确的 target 参数
openclaw cron add \
  --name task \
  --every 1h \
  --message "测试" \
  --channel feishu \
  --to ou_xxxxxxxxxxxxxxxx \
  --announce
```

#### 错误2：超时
```
Error: Task timeout after 30000ms
```

**解决：**
```bash
# 增加超时时间
openclaw cron add \
  --name task \
  --every 1h \
  --message "测试" \
  --timeout 120000
```

#### 错误3：权限不足
```
Error: Permission denied
```

**解决：**
```bash
# 检查 OpenClaw 权限
openclaw doctor

# 检查文件权限
ls -la ~/.openclaw/cron/
chmod 644 ~/.openclaw/cron/jobs.json
```

### 7.3 时间不准确

**问题：** 任务执行时间与预期不符

**解决：**
```bash
# 1. 检查系统时区
timedatectl status
date

# 2. 在任务中明确指定时区
openclaw cron add \
  --name task \
  --cron "0 9 * * *" \
  --tz "Asia/Shanghai" \
  --message "测试"

# 3. 使用 --exact 禁用随机延迟
openclaw cron add \
  --name task \
  --cron "0 9 * * *" \
  --exact \
  --message "测试"
```

### 7.4 执行频率过高

**问题：** 任务执行太频繁，占用资源

**解决：**
```bash
# 1. 修改任务间隔
openclaw cron edit <job-id> --cron "*/30 * * * *"

# 2. 使用 --exact 避免重叠执行
openclaw cron edit <job-id> --exact

# 3. 禁用不必要的任务
openclaw cron disable <job-id>

# 4. 删除不需要的任务
openclaw cron rm <job-id>
```

## 第八部分：高级技巧

### 8.1 批量管理任务

```bash
# 列出所有任务 ID
openclaw cron list | awk '{print $1}' | tail -n +2

# 批量禁用所有任务
openclaw cron list | awk '{print $1}' | tail -n +2 | xargs -I {} openclaw cron disable {}

# 删除特定模式的任务
openclaw cron list | grep "test" | awk '{print $1}' | xargs -I {} openclaw cron rm {}
```

### 8.2 备份和恢复任务

```bash
# 备份所有任务
cp ~/.openclaw/cron/jobs.json ~/.openclaw/cron/jobs.json.backup

# 恢复任务
cp ~/.openclaw/cron/jobs.json.backup ~/.openclaw/cron/jobs.json

# 重启 Gateway 使配置生效
openclaw gateway --restart
```

### 8.3 使用脚本管理

创建一个任务管理脚本：

```bash
#!/bin/bash
# cron-manager.sh

case "$1" in
  list)
    openclaw cron list
    ;;
  backup)
    cp ~/.openclaw/cron/jobs.json ~/.openclaw/cron/jobs.json.backup.$(date +%Y%m%d)
    echo "Backup created: jobs.json.backup.$(date +%Y%m%d)"
    ;;
  restore)
    if [ -f "$2" ]; then
      cp "$2" ~/.openclaw/cron/jobs.json
      openclaw gateway --restart
      echo "Restored: $2"
    else
      echo "Usage: $0 restore <backup-file>"
    fi
    ;;
  *)
    echo "Usage: $0 {list|backup|restore <file>}"
    exit 1
    ;;
esac
```

使用：
```bash
chmod +x cron-manager.sh
./cron-manager.sh list
./cron-manager.sh backup
./cron-manager.sh restore jobs.json.backup.20260311
```

### 8.4 监控和告警

```bash
#!/bin/bash
# cron-monitor.sh

# 检查失败的任务
FAILED_TASKS=$(openclaw cron list | grep error | wc -l)

if [ $FAILED_TASKS -gt 0 ]; then
  echo "Warning: $FAILED_TASKS cron tasks failed"
  # 发送告警消息
  openclaw message send \
    --channel feishu \
    --to ou_xxxxxxxxxxxxxxxx \
    --message "⚠️ Cron 任务失败：$FAILED_TASKS 个任务执行失败"
fi
```

添加到系统监控：
```bash
# 每10分钟检查一次
openclaw cron add \
  --name cron-monitor \
  --every 10m \
  --message "检查 cron 任务执行状态" \
  --session main
```

### 8.5 性能优化

#### 1. 避免任务重叠
```bash
# 使用 --exact 确保任务不会重叠
openclaw cron add \
  --name heavy-task \
  --every 30m \
  --exact \
  --message "执行重量级任务"
```

#### 2. 设置合理的超时
```bash
# 根据任务复杂度设置超时
openclaw cron add \
  --name quick-task \
  --every 5m \
  --timeout 30000 \
  --message "快速任务"

openclaw cron add \
  --name heavy-task \
  --every 1h \
  --timeout 300000 \
  --message "重量级任务"
```

#### 3. 使用 stagger 避免峰值
```bash
# 任务在30秒窗口内随机执行
openclaw cron add \
  --name distributed-task \
  --every 1h \
  --stagger 30s \
  --message "分布式任务"
```

## 第九部分：配置文件详解

### 9.1 Cron 配置文件位置

```
~/.openclaw/cron/jobs.json          # 任务配置
~/.openclaw/cron/runs/            # 执行历史
~/.openclaw/cron/state/           # 任务状态
```

### 9.2 任务配置结构

```json
{
  "version": 1,
  "jobs": [
    {
      "id": "unique-job-id",
      "name": "任务名称",
      "enabled": true,
      "createdAtMs": 1773187500000,
      "updatedAtMs": 1773187500000,
      "schedule": {
        "kind": "cron",           // 或 "every"
        "expr": "*/5 * * * *"    // 或 "everyMs": 300000
      },
      "sessionTarget": "isolated",  // 或 "main"
      "wakeMode": "now",          // 或 "next-heartbeat"
      "payload": {
        "kind": "agentTurn",
        "message": "任务消息内容"
      },
      "delivery": {
        "mode": "announce",       // 或 "direct"
        "channel": "feishu",
        "to": "ou_xxxxxxxxxxxx",  // user:openId, chat:chatId
        "accountId": "default"
      },
      "state": {
        "nextRunAtMs": 1773187620000
      }
    }
  ]
}
```

### 9.3 手动编辑配置

```bash
# 1. 备份配置
cp ~/.openclaw/cron/jobs.json ~/.openclaw/cron/jobs.json.backup

# 2. 编辑配置
vim ~/.openclaw/cron/jobs.json

# 3. 验证 JSON 格式
cat ~/.openclaw/cron/jobs.json | jq .

# 4. 重启 Gateway
openclaw gateway --restart
```

## 第十部分：最佳实践

### 10.1 任务命名规范

- 使用描述性的名称
- 包含任务频率或时间
- 使用连字符或下划线分隔单词

```bash
# 好的命名
openclaw cron add --name daily-backup-0300
openclaw cron add --name hourly-health-check
openclaw cron add --name stock-monitor-5m

# 避免的命名
openclaw cron add --name task1
openclaw cron add --name test
openclaw cron add --name cron-job
```

### 10.2 错误处理

```bash
# 1. 设置合理的超时
--timeout 60000

# 2. 使用 --exact 避免重叠
--exact

# 3. 定期检查执行历史
openclaw cron runs --id <job-id> | grep error

# 4. 配置告警通知
--channel feishu --to ou_xxxxxxxxxxxx --announce
```

### 10.3 资源管理

```bash
# 1. 避免过多同时运行的任务
openclaw cron list | grep running | wc -l

# 2. 使用 stagger 分散任务
--stagger 30s

# 3. 在低峰期执行重量级任务
--cron "0 2 * * *"  # 凌晨2点

# 4. 定期清理执行历史
find ~/.openclaw/cron/runs -name "*.jsonl" -mtime +7 -delete
```

### 10.4 安全考虑

```bash
# 1. 定期备份配置
cp ~/.openclaw/cron/jobs.json ~/.openclaw/cron/jobs.json.backup.$(date +%Y%m%d)

# 2. 保护敏感信息
openclaw cron add --message "使用环境变量：$API_KEY"

# 3. 限制任务权限
openclaw config set tools.profile "restricted"

# 4. 监控异常活动
openclaw logs --follow | grep cron
```

### 10.5 文档化

```bash
# 1. 为每个任务添加描述
--description "每日数据库备份，保留7天"

# 2. 记录任务依赖
# 例如：任务B在任务A完成后执行
openclaw cron add --name task-a --cron "0 2 * * *"
openclaw cron add --name task-b --cron "30 2 * * *"

# 3. 维护任务清单
cat > ~/.openclaw/cron/TASKS.md << 'EOF'
# Cron 任务清单

## 生产任务
- daily-backup: 每日备份
- weekly-report: 每周报告

## 开发任务
- test-task: 测试任务

## 监控任务
- health-check: 健康检查
EOF
```

## 总结

通过本指南，你已经掌握了：

✅ **Cron 表达式** 的使用方法  
✅ **任务管理** 的完整流程  
✅ **高级配置** 的各种技巧  
✅ **实际应用** 的典型案例  
✅ **故障排除** 的解决方案  
✅ **最佳实践** 的经验总结  

### 快速参考

| 命令 | 功能 |
|------|------|
| `openclaw cron list` | 列出所有任务 |
| `openclaw cron add` | 添加新任务 |
| `openclaw cron rm <id>` | 删除任务 |
| `openclaw cron disable <id>` | 禁用任务 |
| `openclaw cron enable <id>` | 启用任务 |
| `openclaw cron run <id>` | 手动执行任务 |
| `openclaw cron runs --id <id>` | 查看执行历史 |
| `openclaw cron status` | 查看任务状态 |

### 常用 Cron 表达式

```bash
*/5 * * * *     # 每5分钟
0 * * * *       # 每小时
0 0 * * *       # 每天0点
0 0 * * 0       # 每周日0点
0 0 1 * *       # 每月1号0点
0 9 * * 1-5     # 工作日早上9点
*/10 * * * 1-5   # 工作日每10分钟
0 9-17 * * 1-5  # 工作时间每小时
```

### 相关资源

- [OpenClaw 官方文档 - Cron](https://docs.openclaw.ai/cli/cron)
- [Crontab.guru - Cron 表达式生成器](https://crontab.guru/)
- [OpenClaw GitHub](https://github.com/openclaw/openclaw)
- [社区 Discord](https://discord.gg/clawd)

---

**最后更新：** 2026年3月11日  
**文档版本：** v1.0.0  
**测试环境：** Rocky Linux 10.0 + OpenClaw 2026.3.2
