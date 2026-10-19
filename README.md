# gitlab-code-reader-mcp

让 AI 助手只读地查看 GitLab 仓库代码的 MCP 服务。工具只覆盖读代码：列目录、找文件、读文件、搜代码、看符号、diff、blame 和提交历史。

仓库：<https://github.com/CoolLiquid/mcp-read-code>

## 文档

| 文档 | 内容 |
|------|------|
| 本文 | 工具、运行方式、环境变量 |
| [DEPLOY.md](DEPLOY.md) | 这台服务器上的 Docker 与开机自启 |
| [docs/CLIENT_MCP.md](docs/CLIENT_MCP.md) | Cursor 如何连接 HTTP 服务 |

## 工具

| 工具 | 作用 |
|------|------|
| `gl_list_directory` | 按深度列出目录 |
| `gl_find_files` | 用 glob 找文件，例如 `**/*.go` |
| `gl_search_code` | 在仓库里搜代码 |
| `gl_read_file` | 带行号读文件，可指定行范围 |
| `gl_read_multiple` | 一次最多读 10 个文件 |
| `gl_read_symbols` | 小文件返回全文，超过约 300 行返回函数/类签名 |
| `gl_diff` | 查看合并请求或两个 ref 的差异 |
| `gl_blame` | 查看每行的修改者和时间 |
| `gl_commit_history` | 查看提交历史和增删统计 |

大文件默认截断，输出带行号。仓库树、文件内容和项目信息有短期缓存。

## 两种运行方式

由 `MCP_TRANSPORT` 决定，默认 `stdio`。

**stdio**：本机进程，启动时用环境变量里的一个 token。

```bash
go build -o gitlab-code-reader-mcp ./cmd/server/
GITLAB_TOKEN=glpat-xxx GITLAB_URL=https://gitlab.example.com ./gitlab-code-reader-mcp
```

Cursor 配置示例：

```json
{
  "mcpServers": {
    "gitlab-code-reader": {
      "command": "/path/to/gitlab-code-reader-mcp",
      "env": {
        "GITLAB_TOKEN": "glpat-xxx",
        "GITLAB_URL": "https://gitlab.example.com"
      }
    }
  }
}
```

**http**：多人共用一台服务。每个请求自己带 token，服务按「GitLab 地址 + token」复用客户端。

- `GET /health` 返回 `{"status":"ok"}`
- `POST /mcp` 为 MCP 端点
- 必须带 `X-Gitlab-Token`，或 `Authorization: Bearer <token>`，否则返回 401
- `X-Gitlab-Url` 可省略，此时使用服务端的 `GITLAB_URL`

```bash
MCP_TRANSPORT=http LISTEN_ADDR=:8088 GITLAB_URL=https://gitlab.example.com ./gitlab-code-reader-mcp
```

内网部署和 Cursor 的 HTTP 配置见 [DEPLOY.md](DEPLOY.md) 与 [docs/CLIENT_MCP.md](docs/CLIENT_MCP.md)。

## 环境变量

| 变量 | 何时需要 | 默认值 | 说明 |
|------|----------|--------|------|
| `MCP_TRANSPORT` | 否 | `stdio` | `stdio` 或 `http` |
| `GITLAB_TOKEN` | stdio 必填 | — | 个人访问令牌，需要 `read_api` |
| `GITLAB_URL` | 否 | `https://gitlab.com` | 默认 GitLab 地址；HTTP 模式下请求头可覆盖 |
| `LISTEN_ADDR` | http | `:8088` | HTTP 监听地址 |

令牌不要写入仓库。可参考 [.env.example](.env.example)，本地 `.env` 已被 git 忽略。

## 代码结构

```
cmd/server/          入口、HTTP 鉴权、stdio / http 两种启动
internal/gitlab/     GitLab API、缓存、按请求选择客户端
internal/tools/      九个 MCP 工具
internal/symbols/    从源码里提取符号
scripts/build.sh     在本机编译二进制，供 Docker 镜像拷贝
```

## 许可

[MIT](LICENSE)
