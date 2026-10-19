# 服务器部署（Docker + 开机自启）

## 环境

- 目录：`/root/code-mcp`
- 监听：`http://10.217.102.192:8088`
- 默认 GitLab：`http://10.217.102.20`（用户可在请求头 `X-Gitlab-Url` 覆盖）

## 首次 / 更新

```bash
cd /root/code-mcp
./scripts/build.sh          # 需要 /usr/local/go（Go 1.24+）
docker compose up -d --build
```

## systemd

```bash
systemctl enable gitlab-code-reader-mcp.service
systemctl start gitlab-code-reader-mcp.service
systemctl status gitlab-code-reader-mcp.service
```

开机时会先执行 `scripts/build.sh` 再 `docker compose up -d --build`。

## 内网 GitLab 域名（/etc/hosts）

若 GitLab 域名只在宿主机 `/etc/hosts` 里配置（公网 DNS 查不到），必须在 `docker-compose.yml` 的 `extra_hosts` 里为容器写同样映射，例如：

```yaml
extra_hosts:
  - "gitlab.dev.qitegroup.com:10.217.111.3"
```

改完后 `docker compose up -d --force-recreate`。

## 验收

```bash
curl http://10.217.102.192:8088/health
docker exec gitlab-code-reader-mcp wget -qO- --timeout=5 https://gitlab.dev.qitegroup.com 2>&1 | head -1
# 无 Token 应 401
curl -o /dev/null -w "%{http_code}\n" -X POST http://10.217.102.192:8088/mcp
```

客户端配置见 [docs/CLIENT_MCP.md](docs/CLIENT_MCP.md)。



### 运维命令

```shell
systemctl status gitlab-code-reader-mcp.service
docker logs gitlab-code-reader-mcp
cd /root/code-mcp && ./scripts/build.sh && docker compose up -d --build
```

