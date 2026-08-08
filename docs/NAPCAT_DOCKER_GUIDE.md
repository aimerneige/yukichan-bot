# NapCat-Docker 启动指南

本文介绍如何使用 [NapCat-Docker](https://github.com/NapNeko/NapCat-Docker) 登录一个 QQ 账号，并通过 OneBot V11 协议连接 `yukichan-bot`。

## 前置准备

请先准备好：

- 已安装 Docker 与 Docker Compose。
- 一个用于机器人的 QQ 账号。
- 当前项目根目录下已有 `config/blacklist.yaml` 与 `config/imgsave.yaml`。

## 配置 OneBot V11 连接

NapCat 在本示例中作为 OneBot V11 WebSocket 服务端，`yukichan-bot` 作为客户端主动连接 NapCat。

项目根目录已提供 `config/config.napcat.yaml`，compose 会将它挂载为容器内的 `/app/config/config.yaml`。该文件将 `ws_server` 指向 compose 网络里的 NapCat 服务名：

```yaml
common:
  ws_server: "ws://napcat:3001"
  ws_token: ""
  nickname:
    - "ゆき酱"
  command_prefix: "/"
  super_users:
    - 1227427929
```

如果你要改用自己的 `config/config.yaml`，请将 `common.ws_server` 设置为 `ws://napcat:3001`。

## 使用 compose 文件

项目根目录的 `docker-compose.yaml` 会同时启动 NapCat 与 `yukichan-bot`，两个容器位于同一个 Docker 私有网络中。宿主机只暴露 NapCat WebUI 端口 `6099`，用于登录和配置 NapCat。

```yaml
services:
  napcat:
    image: mlikiowa/napcat-docker:latest
    container_name: napcat
    restart: unless-stopped
    mac_address: 02:42:ac:11:00:02
    environment:
      - NAPCAT_UID=${NAPCAT_UID:-1000}
      - NAPCAT_GID=${NAPCAT_GID:-1000}
      - MODE=ws
    ports:
      - "6099:6099"
    volumes:
      - ./napcat/config:/app/napcat/config
      - ./napcat/qq:/app/.config/QQ
    networks:
      - yukichan-net

  yukichan:
    build: .
    container_name: yukichan
    restart: unless-stopped
    depends_on:
      - napcat
    volumes:
      - ./config:/app/config
      - ./config/config.napcat.yaml:/app/config/config.yaml:ro
      - ./data:/app/data
      - ./temp:/app/temp
    networks:
      - yukichan-net

networks:
  yukichan-net:
    driver: bridge
```

说明：

- `MODE=ws` 会让 NapCat-Docker 使用 WebSocket 服务端模板，默认监听容器内 `3001` 端口。
- `ws://napcat:3001` 中的 `napcat` 是 compose 服务名，只在同一个 Docker 网络内解析。
- `config/config.napcat.yaml` 是 compose 专用的 OneBot V11 配置。
- `./napcat/qq` 会持久化 QQ 登录数据，避免每次重建容器都重新登录。
- `./napcat/config` 会持久化 NapCat 配置，便于后续在 WebUI 中调整。

## 启动容器

在项目根目录执行：

```bash
NAPCAT_UID=$(id -u) NAPCAT_GID=$(id -g) docker compose up -d
```

查看 NapCat 日志：

```bash
docker logs -f napcat
```

NapCat-Docker 的 WebUI token 会出现在容器日志中；如果没有手动设置，按当前 NapCat-Docker 文档说明以日志为准。

## 登录 QQ 账号

1. 打开 `http://<宿主机 IP>:6099/webui`。
2. 使用 `docker logs napcat` 中显示的 WebUI token 登录。
3. 在 WebUI 中按照页面提示登录机器人 QQ 账号。
4. 登录完成后，确认网络配置中已启用 WebSocket 服务端，监听地址为 `0.0.0.0`，端口为 `3001`。
5. 重启或等待 `yukichan-bot` 重连。

连接成功后，可以在 QQ 私聊或群聊中发送项目支持的指令，例如 `老婆`、`下棋` 等，确认链路已打通。

## 常见问题

| 现象 | 处理方式 |
| --- | --- |
| `yukichan-bot` 无法连接 NapCat | 检查 compose 是否挂载了 `config/config.napcat.yaml`，并确认配置里的地址是 `ws://napcat:3001`。 |
| 打不开 WebUI | 确认 `napcat` 服务正在运行，并检查宿主机 `6099` 端口是否被防火墙或其他进程占用。 |
| 重建容器后需要重新登录 | 确认 compose 中已挂载 `./napcat/qq:/app/.config/QQ`，且目录权限与 `NAPCAT_UID` / `NAPCAT_GID` 匹配。 |
| QQ 群消息没有响应 | 确认机器人账号已在群内，并检查消息是否符合插件触发规则。 |
