# 连接 HTTP 服务

服务部署见 [DEPLOY.md](../DEPLOY.md)。下面是在用户电脑上配置 Cursor。

在 `~/.cursor/mcp.json`（全局）或项目 `.cursor/mcp.json` 中加入：

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

- Token 需要 GitLab `read_api` 权限。每人使用自己的 token。
- 不写 `X-Gitlab-Url` 时，使用服务器上的 `GITLAB_URL`。
- 也可以用 `Authorization: Bearer glpat-xxx` 代替 `X-Gitlab-Token`。
- 读哪个仓库由对话里的 `project_id` 决定，可以是 `group/project` 或数字 ID。
- 不要把带 token 的配置提交到 git。这个地址是内网 HTTP，只在公司网络里使用。

本机直接跑 stdio 的配置写在 [README](../README.md) 的「两种运行方式」里。
