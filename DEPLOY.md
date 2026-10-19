# 部署

当前服务跑在本机 Docker 里，由 systemd 负责开机启动。客户端连接方式见 [docs/CLIENT_MCP.md](docs/CLIENT_MCP.md)。

## 当前环境

| 项 | 值 |
|----|----|
| 代码目录 | `/root/code-mcp` |
| 对外地址 | `http://10.217.102.192:8088` |
| 健康检查 | `http://10.217.102.192:8088/health` |
| MCP 路径 | `http://10.217.102.192:8088/mcp` |
| 默认 GitLab | `http://10.217.102.20`（请求头 `X-Gitlab-Url` 可覆盖） |

镜像不在容器内编译。`scripts/build.sh` 用宿主机的 Go（1.24+，`/usr/local/go`）生成 `gitlab-code-reader-mcp`，Dockerfile 再拷贝这个二进制。

## 更新并启动

```bash
cd /root/code-mcp
./scripts/build.sh
docker compose up -d --build
```

## 开机自启

单元文件：`deploy/gitlab-code-reader-mcp.service`。启动前会先执行 `scripts/build.sh`，再 `docker compose up -d --build`。

```bash
systemctl enable gitlab-code-reader-mcp.service
systemctl start gitlab-code-reader-mcp.service
systemctl status gitlab-code-reader-mcp.service
```

## 内网 GitLab 域名

如果 GitLab 域名只写在宿主机 `/etc/hosts` 里，容器 DNS 看不到。在 `docker-compose.yml` 的 `extra_hosts` 中写同样的映射，例如：

```yaml
extra_hosts:
  - "gitlab.dev.qitegroup.com:10.217.111.3"
```

改完后执行：

```bash
docker compose up -d --force-recreate
```

## 检查

```bash
curl http://10.217.102.192:8088/health
curl -o /dev/null -w "%{http_code}\n" -X POST http://10.217.102.192:8088/mcp
docker logs gitlab-code-reader-mcp
```

`/health` 应返回 `{"status":"ok"}`。不带 token 访问 `/mcp` 应返回 `401`。

## 常用命令

```bash
systemctl status gitlab-code-reader-mcp.service
docker logs gitlab-code-reader-mcp
cd /root/code-mcp && ./scripts/build.sh && docker compose up -d --build
```
