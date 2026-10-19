# Cursor 客户端配置

在**用户电脑**上的 `~/.cursor/mcp.json`（全局）或项目 `.cursor/mcp.json` 中添加：

```json
{
  "mcpServers": {
    "gitlab-code-reader": {
      "url": "http://10.217.102.192:8088/mcp",
      "headers": {
        "X-Gitlab-Token": "glpat-你自己的token",
        "X-Gitlab-Url": "https://gitlab.dev.qitegroup.com"
      }
    }
  }
}
```

- Token 需具备 GitLab `read_api` 权限。
- `X-Gitlab-Url` 可省略，将使用服务器默认 `GITLAB_URL`。
- 也可使用 `Authorization: Bearer glpat-xxx` 代替 `X-Gitlab-Token`。
- 读取不同仓库：在对话中指定 `project_id`（如 `group/project` 或数字 ID）。

**安全提示**：不要将含 Token 的配置提交到 git；内网 HTTP 仅限公司网络使用。
