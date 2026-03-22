```
{
    "url": "openclaw-config-guide",
    "time": "2026/03/22 12:40",
    "tag": "OpenClaw, 配置"
}
```

# OpenClaw 配置文件完全指南

`openclaw.json` 是 OpenClaw 的核心配置文件，控制着整个系统的行为、模型、技能、Agent 等各个方面。本文按模块详细介绍各个配置字段及其使用场景。

## 一、meta - 元数据

`meta` 字段用于跟踪配置文件的版本和最后修改时间。

```json
"meta": {
  "lastTouchedVersion": "2026.3.13",
  "lastTouchedAt": "2026-03-22T01:25:53.282Z"
}
```

**字段说明：**

- `lastTouchedVersion`：最后修改配置的 OpenClaw 版本
- `lastTouchedAt`：配置文件最后修改的时间戳（ISO 8601 格式）

**使用场景：**

- 用于诊断配置版本兼容性问题
- 跟踪配置文件的更新历史
- 当升级 OpenClaw 后，系统会自动更新这个字段

---

## 二、wizard - 向导配置

`wizard` 字段记录初始化向导的运行信息。

```json
"wizard": {
  "lastRunAt": "2026-03-21T14:00:45.903Z",
  "lastRunVersion": "2026.3.13",
  "lastRunCommand": "onboard",
  "lastRunMode": "local"
}
```

**字段说明：**

- `lastRunAt`：向导最后运行的时间
- `lastRunVersion`：向导运行时的 OpenClaw 版本
- `lastRunCommand`：最后执行的向导命令（如 `onboard`）
- `lastRunMode`：运行模式（local / web）

**使用场景：**

- 用于判断用户是否完成初始化设置
- 当配置问题时，可以查看向导历史
- 新用户安装 OpenClaw 后会自动填充

---

## 三、auth - 认证配置

`auth` 字段定义各个 LLM 提供商的认证方式。

```json
"auth": {
  "profiles": {
    "zai:default": {
      "provider": "zai",
      "mode": "api_key"
    }
  }
}
```

**字段说明：**

- `profiles`：认证配置文件，key 为 `<provider>:<profile>` 格式
  - `provider`：提供商标识（如 zai、openai、anthropic）
  - `mode`：认证模式
    - `api_key`：使用 API Key
    - `oauth`：使用 OAuth 令牌

**使用场景：**

- 配置多个 LLM 提供商的认证信息
- 支持不同的认证方式（API Key、OAuth）
- 为不同的 profile 配置不同的 provider
- 运行时通过 `auth.profiles.zai:default` 指定默认 profile

**示例：**

```json
"auth": {
  "profiles": {
    "zai:default": {
      "provider": "zai",
      "mode": "api_key"
    },
    "openai:gpt4": {
      "provider": "openai",
      "mode": "oauth"
    }
  }
}
```

---

## 四、models - 模型配置

`models` 字段定义 LLM 模型列表和配置。

```json
"models": {
  "mode": "merge",
  "providers": {
    "zai": {
      "baseUrl": "https://open.bigmodel.cn/api/coding/paas/v4",
      "api": "openai-completions",
      "models": [
        {
          "id": "glm-5",
          "name": "GLM-5",
          "reasoning": true,
          "input": ["text"],
          "cost": {
            "input": 0,
            "output": 0,
            "cacheRead": 0,
            "cacheWrite": 0
          },
          "contextWindow": 204800,
          "maxTokens": 131072
        }
      ]
    }
  }
}
```

**字段说明：**

- `mode`：模型配置模式
  - `merge`：合并多个 provider 的模型（默认）
  - `override`：只使用指定 provider 的模型

- `providers`：提供商配置，key 为 provider ID
  - `baseUrl`：API 基础 URL
  - `api`：API 兼容类型（openai-completions 等）
  - `models`：模型列表
    - `id`：模型 ID（用于识别）
    - `name`：模型显示名称
    - `reasoning`：是否支持推理（true/false）
    - `input`：支持的输入类型（text、image、audio）
    - `cost`：成本配置（每百万 token 的价格）
      - `input`：输入成本
      - `output`：输出成本
      - `cacheRead`：读取缓存成本
      - `cacheWrite`：写入缓存成本
    - `contextWindow`：上下文窗口大小（token 数）
    - `maxTokens`：最大输出 token 数

**使用场景：**

- 定义可用的 LLM 模型及其参数
- 配置模型的成本信息，用于计费和优化
- 控制模型的行为（推理能力、输入类型等）
- 为不同的会话类型选择合适的模型

**示例场景：**

1. **推理模型**：`reasoning: true` 适用于复杂任务
2. **快速模型**：FlashX 版本用于实时对话
3. **多模态模型**：支持 image、audio 输入
4. **成本控制**：通过 cost 字段实现 token 计费

---

## 五、skills - 技能配置

`skills` 字段控制技能的加载、禁用和环境变量。

```json
"skills": {
  "entries": {
    "feishu-bitable": { "enabled": false },
    "web_search": { "enabled": false }
  }
}
```

**字段说明：**

- `entries`：技能配置对象，key 为技能名称
  - `enabled`：技能是否启用（true/false，默认 true）

**技能配置的高级选项：**

```json
"skills": {
  "entries": {
    "skill-name": {
      "enabled": true,
      "env": {
        "API_KEY": "your_api_key_here",
        "DATABASE_URL": "postgresql://localhost:5432/db"
      },
      "apiKey": {
        "source": "env",
        "provider": "default",
        "id": "SKILL_API_KEY"
      }
    }
  }
}
```

**高级字段说明：**

- `env`：为技能注入环境变量（沙盒内不可用）
- `apiKey`：配置技能的 API Key
  - `source`：密钥来源（env、plaintext）
  - `provider`：密钥提供商（default 表示默认提供商）
  - `id`：密钥 ID

**使用场景：**

1. **启用/禁用技能**
   - 禁用不常用的技能，减少干扰
   - 临时禁用有问题的技能
   - 为特定会话启用特定技能

2. **技能环境变量**
   - 为技能提供数据库连接信息
   - 配置第三方 API 密钥
   - 设置技能运行时的配置

3. **技能优先级**
   - `~/.openclaw/skills/`（全局，所有 Agent 共享）
   - `<workspace>/skills/`（Workspace 专属，优先级最高）
   - `skills.load.extraDirs`（额外技能目录，优先级最低）

4. **实际案例**
   ```json
   {
     "skills": {
       "entries": {
         "feishu-calendar": { "enabled": false },    // 禁用日历
         "feishu-create-doc": { "enabled": false }, // 禁用创建文档
         "tavily-search": { "enabled": true }       // 启用 Tavily 搜索
       }
     }
   }
   ```

---

## 六、agents - Agent 配置

`agents` 字段定义 Agent 的默认配置和列表。

```json
"agents": {
  "defaults": {
    "model": {
      "primary": "zai/glm-4.7"
    },
    "models": {
      "zai/glm-5": {
        "alias": "GLM"
      },
      "zai/glm-4.7": {}
    },
    "workspace": "/home/openclaw/.openclaw/workspace",
    "compaction": {
      "mode": "safeguard"
    }
  },
  "list": [
    {
      "id": "main"
    },
    {
      "id": "stock",
      "name": "stock",
      "workspace": "/home/openclaw/.openclaw/workspace-stock",
      "agentDir": "/home/openclaw/.openclaw/agents/stock/agent"
    }
  ]
}
```

**字段说明：**

- `defaults`：所有 Agent 的默认配置
  - `model.primary`：默认使用的模型
  - `models.<model_id>.alias`：模型别名（简短名称）
  - `models.<model_id>.{}`：模型的其他配置（可扩展）
  - `workspace`：默认工作区路径
  - `compaction.mode`：会话压缩模式
    - `off`：不压缩
    - `basic`：基础压缩
    - `aggressive`：激进压缩
    - `safeguard`：安全保护模式（平衡）

- `list`：Agent 列表
  - `id`：Agent ID（唯一标识）
  - `name`：Agent 显示名称
  - `workspace`：Agent 专属工作区路径
  - `agentDir`：Agent 配置目录路径

**使用场景：**

1. **多 Agent 管理**
   - 为不同任务创建专用 Agent（如 stock、main）
   - 每个 Agent 可以有独立的工作区和配置

2. **模型选择**
   - 为不同类型的会话选择合适的模型
   - 通过 alias 简化模型引用

3. **工作区隔离**
   - `main` Agent 使用主工作区
   - `stock` Agent 使用独立工作区
   - 避免不同 Agent 的配置冲突

4. **会话压缩**
   - 降低 token 消耗，保留重要上下文
   - 平衡性能和信息保留

5. **实际案例**
   ```json
   {
     "agents": {
       "defaults": {
         "model": {
           "primary": "zai/glm-4.7"  // 默认使用 GLM-4.7
         },
         "workspace": "/home/openclaw/.openclaw/workspace"  // 默认工作区
       },
       "list": [
         {
           "id": "stock",
           "name": "股票助手",
           "workspace": "/home/openclaw/.openclaw/workspace-stock"  // 独立工作区
         }
       ]
     }
   }
   ```

---

## 七、tools - 工具配置

`tools` 字段控制工具的行为和可见性。

```json
"tools": {
  "profile": "full",
  "sessions": {
    "visibility": "all"
  }
}
```

**字段说明：**

- `profile`：工具配置文件
  - `full`：所有工具可用
  - `minimal`：只启用核心工具
  - `sandbox`：沙盒模式限制

- `sessions.visibility`：会话可见性
  - `all`：所有工具对所有会话可见
  - `self`：只对主会话可见
  - `agent`：只对 Agent 会话可见

**使用场景：**

1. **工具权限控制**
   - 限制危险工具的使用（如文件删除、系统命令）
   - 为不同会话类型提供不同级别的工具访问

2. **安全配置**
   - 生产环境使用 `minimal` profile
   - 开发环境使用 `full` profile
   - 敏感操作需要额外确认

3. **会话隔离**
   - 子 Agent 不应该访问所有工具
   - 主会话可以拥有完整工具访问

---

## 八、bindings - 绑定配置

`bindings` 字段定义 Agent 与消息渠道的映射关系。

```json
"bindings": [
  {
    "type": "route",
    "agentId": "stock",
    "match": {
      "channel": "feishu",
      "accountId": "stock"
    }
  }
]
```

**字段说明：**

- `type`：绑定类型
  - `route`：路由绑定，根据条件分配消息
  - `allowlist`：允许列表绑定
  - `blocklist`：阻止列表绑定

- `agentId`：目标 Agent ID

- `match`：匹配条件
  - `channel`：渠道 ID（如 feishu、telegram）
  - `accountId`：账户 ID（区分同一渠道的不同账号）

**使用场景：**

1. **专用 Agent 路由**
   - 将特定渠道的消息路由到专用 Agent
   - 例如：飞书"操盘手机器人"账户消息路由到 stock Agent

2. **多账户管理**
   - 同一渠道支持多个账户
   - 不同账户路由到不同的 Agent

3. **条件匹配**
   - 根据 channel、accountId 组合路由
   - 支持复杂的路由规则

4. **实际案例**
   ```json
   {
     "bindings": [
       {
         "type": "route",
         "agentId": "stock",  // 股票消息发送到 stock Agent
         "match": {
           "channel": "feishu",
           "accountId": "stock"  // 使用"操盘手机器人"账户
         }
       }
     ]
   }
   ```

---

## 九、commands - 命令配置

`commands` 字段控制命令的执行行为。

```json
"commands": {
  "native": "auto",
  "nativeSkills": "auto",
  "restart": true,
  "ownerDisplay": "raw"
}
```

**字段说明：**

- `native`：原生命令处理
  - `auto`：自动识别并处理（默认）
  - `on`：启用原生命令
  - `off`：禁用原生命令

- `nativeSkills`：技能命令处理
  - `auto`：自动加载技能的斜杠命令（默认）
  - `on`：启用技能命令
  - `off`：禁用技能命令

- `restart`：重启行为
  - `true`：命令执行后自动重启 Gateway
  - `false`：不自动重启

- `ownerDisplay`：所有者显示格式
  - `raw`：显示原始命令
  - `masked`：隐藏敏感信息

**使用场景：**

1. **命令启用/禁用**
   - 禁用危险的系统命令
   - 只启用特定的技能命令

2. **调试模式**
   - 关闭自动重启，便于调试
   - 显示原始命令，便于排查问题

3. **安全控制**
   - 禁用不需要的命令，减少安全风险
   - 控制 Gateway 的重启行为

---

## 十、session - 会话配置

`session` 字段控制会话行为和权限。

```json
"session": {
  "dmScope": "per-channel-peer"
}
```

**字段说明：**

- `dmScope`：私聊会话作用域
  - `per-channel-peer`：每个渠道的每个用户独立会话
  - `global`：全局会话，所有渠道共享
  - `shared`：共享会话，跨渠道共享上下文

**使用场景：**

1. **会话隔离**
   - 确保不同用户之间的对话独立
   - 防止信息泄露

2. **跨渠道一致性**
   - 需要跨渠道保持对话时使用 `global` 或 `shared`

3. **性能优化**
   - 减少重复的上下文加载
   - 平衡一致性和性能

---

## 十一、channels - 渠道配置

`channels` 字段定义各个消息渠道的配置和权限。

```json
"channels": {
  "feishu": {
    "enabled": true,
    "appId": "cli_a921cf23b4781cb5",
    "domain": "feishu",
    "connectionMode": "websocket",
    "requireMention": true,
    "dmPolicy": "allowlist",
    "allowFrom": [
      "ou_899873fddd054110f5af55ce86df8344"
    ],
    "groupAllowFrom": [
      "ou_899873fddd054110f5af55ce86df8344"
    ],
    "groupPolicy": "allowlist",
    "groups": {
      "*": {
        "enabled": true
      }
    },
    "accounts": {
      "default": {},
      "stock": {
        "appId": "cli_a92658933778dcb1",
        "appSecret": "YOUR_APP_SECRET",
        "botName": "操盘手机器人",
        "dmPolicy": "allowlist",
        "allowFrom": [
          "ou_b12edb2686841407a42abf98a59e27cb"
        ]
      },
      "news": {
        "appId": "cli_a927874d0078dcc2",
        "appSecret": "YOUR_APP_SECRET",
        "botName": "头条资讯",
        "dmPolicy": "allowlist",
        "allowFrom": [
          "ou_b12edb2686841407a42abf98a59e27cb"
        ]
      }
    }
  }
}
```

**字段说明：**

- `enabled`：渠道是否启用

- `appId` / `appSecret`：飞书应用凭证

- `domain`：域名（feishu）

- `connectionMode`：连接模式
  - `websocket`：实时连接
  - `polling`：轮询模式

- `requireMention`：是否需要 @ 才能响应
  - `true`：群聊中需要 @ 机器人
  - `false`：自动响应所有消息

- `dmPolicy` / `groupPolicy`：私聊/群聊策略
  - `allowlist`：只允许列表中的用户/群
  - `blocklist`：阻止列表中的用户/群
  - `all`：允许所有

- `allowFrom` / `groupAllowFrom`：允许的来源用户 ID 列表

- `groups`：群组配置
  - `"*"`：所有群
  - 具体群 ID：单独配置

- `accounts`：账户配置
  - `default`：默认账户配置
  - `<account_id>`：特定账户配置
    - `appId` / `appSecret`：该账户的应用凭证
    - `botName`：机器人名称
    - `dmPolicy` / `allowFrom`：该账户的权限策略

**使用场景：**

1. **多账户支持**
   - 同一渠道配置多个账户（如 default、stock、news）
   - 每个账户有独立的权限和应用凭证

2. **权限控制**
   - `allowlist`：只响应特定用户的消息
   - `blocklist`：阻止特定用户的消息
   - `groupPolicy`：控制群聊中的响应范围

3. **群聊行为**
   - `requireMention: true`：避免刷屏，只在 @ 时响应
   - `groupPolicy: all`：响应所有群聊
   - `groupPolicy: allowlist`：只响应特定群

4. **安全配置**
   - 严格控制敏感操作的权限
   - 为不同的机器人配置不同的访问策略

---

## 十二、gateway - Gateway 配置

`gateway` 字段控制 Gateway 服务的运行参数。

```json
"gateway": {
  "port": 18789,
  "mode": "local",
  "bind": "loopback",
  "auth": {
    "mode": "token",
    "token": "YOUR_GATEWAY_TOKEN"
  },
  "tailscale": {
    "mode": "off",
    "resetOnExit": false
  },
  "nodes": {
    "denyCommands": [
      "camera.snap",
      "camera.clip",
      "screen.record",
      "contacts.add",
      "calendar.add",
      "reminders.add",
      "sms.send"
    ]
  }
}
```

**字段说明：**

- `port`：Gateway 服务监听端口

- `mode`：运行模式
  - `local`：本地模式（默认）
  - `remote`：远程模式
  - `cluster`：集群模式

- `bind`：绑定地址
  - `loopback`：只监听本地（127.0.0.1）
  - `0.0.0.0`：监听所有接口

- `auth`：Gateway 认证
  - `mode`：认证模式
    - `token`：使用令牌认证
    - `password`：密码认证

- `tailscale`：TailScale 配置
  - `mode`：`off` / `on`
  - `resetOnExit`：退出时是否重置

- `nodes.denyCommands`：节点级禁止的命令列表
  - `camera.snap`：截图
  - `camera.clip`：剪贴板
  - `screen.record`：屏幕录制
  - `contacts.add`：添加联系人
  - `calendar.add`：添加日历
  - `reminders.add`：添加提醒
  - `sms.send`：发送短信

**使用场景：**

1. **端口管理**
   - 指定 Gateway 监听的端口
   - 避免端口冲突
   - 部署在服务器上时使用 `0.0.0.0`

2. **安全控制**
   - 禁用危险命令，防止恶意操作
   - 控制节点的访问权限

3. **认证配置**
   - 配置 Gateway 的访问控制
   - 支持令牌和密码认证

4. **运行模式**
   - `local`：本地开发环境
   - `remote` / `cluster`：生产环境

---

## 十三、plugins - 插件配置

`plugins` 字段管理 OpenClaw 的插件。

```json
"plugins": {
  "allow": [
    "openclaw-lark"
  ],
  "entries": {
    "feishu": {
      "enabled": false
    },
    "openclaw-lark": {
      "enabled": true
    }
  },
  "installs": {
    "openclaw-lark": {
      "source": "npm",
      "spec": "@larksuite/openclaw-lark",
      "installPath": "/home/openclaw/.openclaw/extensions/openclaw-lark",
      "version": "2026.3.17",
      "resolvedName": "@larksuite/openclaw-lark",
      "resolvedVersion": "2026.3.17",
      "resolvedSpec": "@larksuite/openclaw-lark@2026.3.17",
      "integrity": "sha512-...",
      "shasum": "407d4616186cf776c21cf49aa533b80c2ef23530",
      "resolvedAt": "2026-03-22T00:06:34.450Z",
      "installedAt": "2026-03-22T00:08:01.978Z"
    }
  }
}
```

**字段说明：**

- `allow`：允许的插件列表
  - 插件 ID 或包名

- `entries`：插件配置
  - `<plugin_id>`：插件 ID
    - `enabled`：是否启用（true/false）

- `installs`：已安装的插件详情
  - `source`：安装源（npm、git、local）
  - `spec`：包规范（npm 包名或 Git URL）
  - `installPath`：安装路径
  - `version`：安装版本
  - `resolvedName` / `resolvedVersion`：解析的包名和版本
  - `resolvedSpec`：解析的包规范
  - `integrity` / `shasum`：完整性校验值
  - `resolvedAt`：解析时间
  - `installedAt`：安装时间

**使用场景：**

1. **插件启用/禁用**
   - 临时禁用有问题的插件
   - 只启用需要的插件
   - 避免插件之间的冲突

2. **插件管理**
   - 查看已安装的插件及其版本
   - 验证插件的完整性
   - 跟踪插件的安装和解析历史

3. **多插件支持**
   - `allow` 列表定义允许使用的插件
   - 支持多个插件同时运行
   - 插件之间可以共享数据

4. **实际案例**
   ```json
   {
     "plugins": {
       "allow": [
         "openclaw-lark"  // 只允许飞书插件
       ],
       "entries": {
         "openclaw-lark": { "enabled": true }  // 启用飞书插件
       }
     }
   }
   ```

---

## 十四、配置最佳实践

### 1. 技能管理

**推荐做法：**
- 通用技能放在 `~/.openclaw/skills/` 全局目录
- Agent 专属技能放在 `<workspace>/skills/`
- 通过 `skills.entries.<skill>.enabled` 控制启用/禁用

**实际案例：**
```json
{
  "skills": {
    "entries": {
      "feishu-calendar": { "enabled": false },  // 禁用不常用功能
      "tavily-search": { "enabled": true }       // 启用搜索功能
    }
  }
}
```

### 2. Agent 配置

**推荐做法：**
- 为不同任务创建专用 Agent（main、stock、news）
- 每个 Agent 使用独立的工作区
- 使用 bindings 路由特定渠道到特定 Agent

**实际案例：**
```json
{
  "agents": {
    "list": [
      {
        "id": "main",
        "workspace": "/home/openclaw/.openclaw/workspace"
      },
      {
        "id": "stock",
        "name": "股票助手",
        "workspace": "/home/openclaw/.openclaw/workspace-stock"
      }
    ]
  },
  "bindings": [
    {
      "type": "route",
      "agentId": "stock",
      "match": {
        "channel": "feishu",
        "accountId": "stock"
      }
    }
  ]
}
```

### 3. 安全配置

**推荐做法：**
- 使用 `tools.profile` 限制危险工具
- 通过 `channels.dmPolicy` 和 `groupPolicy` 控制访问权限
- 在 `gateway.nodes.denyCommands` 中禁用敏感操作

**实际案例：**
```json
{
  "tools": {
    "profile": "minimal"  // 生产环境使用 minimal profile
  },
  "gateway": {
    "nodes": {
      "denyCommands": [
        "camera.snap",  // 禁用截图
        "screen.record",  // 禁用录屏
        "contacts.add"   // 禁用添加联系人
      ]
    }
  }
}
```

### 4. 配置生效

**修改配置后：**
```bash
# 修改配置后，需要重启 Gateway
openclaw gateway restart

# 或使用 systemd
sudo systemctl restart openclaw-gateway
```

**配置验证：**
```bash
# 查看配置是否生效
openclaw gateway status

# 查看日志
journalctl -u openclaw-gateway -f
```

---

## 十五、常见问题

### Q1：修改配置不生效？

**A：** 配置修改后需要重启 Gateway
```bash
openclaw gateway restart
```

### Q2：技能冲突怎么解决？

**A：** 技能按优先级加载：`workspace/skills` → `~/.openclaw/skills/` → bundled
- 同名技能，高优先级覆盖低优先级
- 检查技能的 `name` 字段是否唯一

### Q3：如何添加新模型？

**A：** 在 `models.providers.<provider>.models` 中添加
```json
{
  "models": {
    "providers": {
      "zai": {
        "models": [
          {
            "id": "glm-6",
            "name": "GLM-6",
            "reasoning": true,
            "contextWindow": 262144,
            "maxTokens": 262144
          }
        ]
      }
    }
  }
}
```

### Q4：如何配置多个飞书机器人？

**A：** 在 `channels.feishu.accounts` 中添加
```json
{
  "channels": {
    "feishu": {
      "accounts": {
        "default": { /* 默认机器人配置 */ },
        "stock": { /* 股票机器人配置 */ },
        "news": { /* 资讯机器人配置 */ }
      }
    }
  }
}
```

---

## 十六、总结

`openclaw.json` 是 OpenClaw 的核心配置文件，涵盖了系统的所有方面：

1. **元数据和向导**：`meta`、`wizard`
2. **认证和模型**：`auth`、`models`
3. **技能管理**：`skills`
4. **Agent 配置**：`agents`
5. **工具配置**：`tools`
6. **绑定和命令**：`bindings`、`commands`
7. **会话和渠道**：`session`、`channels`
8. **Gateway 和插件**：`gateway`、`plugins`

**配置原则：**
- 模块化：每个部分独立配置，互不影响
- 安全优先：默认关闭危险功能，需要时手动开启
- 灵活扩展：支持多 Agent、多账户、多插件
- 配置清晰：字段语义明确，易于理解和维护

通过合理配置 `openclaw.json`，可以充分发挥 OpenClaw 的能力，打造符合你需求的 AI 助手系统。

**安全提示：**
- 配置文件中包含敏感信息（API Keys、App Secrets 等）时，请注意保管
- 不要将包含真实敏感信息的配置文件提交到版本控制系统
- 生产环境建议使用更严格的权限配置
