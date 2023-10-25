# ddns4cdn [![wakatime](https://wakatime.com/badge/github/jat001/ddns4cdn.svg)](https://wakatime.com/@Jat/projects/jpfnwygket)

自动获取本机 IP 并更新 CDN 源站地址，同时支持 IPV4 和 IPV6。

## Supported Cloud Services

- Cloudflare
- Alibaba Cloud ECDN
- Tencent Cloud DCDN

## Config

See [example](/config.example.toml)

## Run

```shell
./ddns4cdn -c config.toml
```

## Build

### Desktop

```shell
go build -C src -o ../ddns4cdn[.exe]
```

Set `GOOS` and `GOARCH` to build for other platforms.

### Mobile

```shell
cd src
gomobile bind -o ../ddns4cdn.aar -target android -javapkg ddns4cdn -androidapi 33 ./worker
gomobile bind -o ../ddns4cdn.xcframework -target ios -prefix Ddns4cdn -iosversion 17 ./worker
```
