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

> Tips for developers using macOS:
>
> You don't need Xcode (but still need to install command line tools) to build executables and libraries (including Swift).
>
> I created this project without opening Xcode at all, you can do the same.
>
> Xcode sucks. It doesn't provide a human-readable configuration file. When I search on Google for development problems about macOS, most of the results are about Xcode (even on Stack Overflow). Using the mouse to click buttons is not a good way to develop. It is not efficient and not reproducible. I don't want to use a tool that I don't know how it works.

### Executable

```shell
go build -o build/ddns4cdn[.exe] ./src
```

Set `GOOS` and `GOARCH` to build for other platforms.

### Library

#### Static library

```shell
go build -o build/ddns4cdn.a -buildmode=c-archive ./src/cgo/go
```

macOS:

```shell
# C
clang -o build/ddns4cdn_c src/cgo/c/main.c build/ddns4cdn.a -framework CoreFoundation -framework Security -lresolv
# C++
clang++ -o build/ddns4cdn_cpp src/cgo/cpp/main.cc build/ddns4cdn.a -framework CoreFoundation -framework Security -lresolv
# Objective-C
clang -o build/ddns4cdn_objc src/cgo/objc/main.m build/ddns4cdn.a -framework Foundation -framework Security -lresolv
```

Swift:

```shell
cd src/cgo/swift
# the output binary is .build/release/ddns4cdn
swift build -c release -Xlinker ../../../build/ddns4cdn.a -Xlinker -lresolv
```

#### Shared library

```shell
go build -o build/ddns4cdn.so -buildmode=c-shared ./src/cgo/go
```

macOS:

```shell
# C
clang -o build/ddns4cdn_dl_c src/cgo/c/main.c build/ddns4cdn.so
# C++
clang++ -o build/ddns4cdn_dl_cpp src/cgo/cpp/main.cc build/ddns4cdn.so
# Objective-C
clang -o build/ddns4cdn_dl_objc src/cgo/objc/main.m build/ddns4cdn.so -framework Foundation

# set DYLD_LIBRARY_PATH to load shared library
export DYLD_LIBRARY_PATH=$(realpath build)
```

Swift:

```shell
cd src/cgo/swift
# the output binary is .build_dl/release/ddns4cdn
swift build --scratch-path .build_dl -c release -Xlinker ../../../build/ddns4cdn.so

# set DYLD_LIBRARY_PATH to load shared library
export DYLD_LIBRARY_PATH=$(realpath ../../../build)
```

### Go mobile

```shell
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

```shell
cd src/gomobile
# Android
gomobile bind -o ../../build/ddns4cdn.aar -target android -javapkg ddns4cdn -androidapi 33
# iOS and macOS
gomobile bind -o ../../build/Ddns4cdn.xcframework -target ios,macos -prefix Ddns4cdn -iosversion 17
```
