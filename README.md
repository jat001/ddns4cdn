# ddns4cdn [![wakatime](https://wakatime.com/badge/github/jat001/ddns4cdn.svg)](https://wakatime.com/@Jat/projects/jpfnwygket)

自动获取本机 IP 并更新 CDN 源站地址，同时支持 IPV4 和 IPV6。

## 已支持的厂商

- Cloudflare
- 腾讯云 ECDN
- 阿里云 DCDN

## 配置

参考[示例](/config.example.toml)

## 运行

```shell
./ddns4cdn -c config.toml
```

## Build for mobile

```shell
cd src
gomobile bind -o ../ddns4cdn.aar -target=android -androidapi 33 ./worker
```
