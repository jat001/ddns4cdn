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

#### Requirements

##### Windows

- [Mingw-w64](https://www.mingw-w64.org/downloads/) (GCC, required for `go build`)

- Build Tools for Visual Studio (MSVC, optional)

##### macOS

- GCC (`go build` can use GCC or Clang)

- Xcode Command Line Tools (set `CC=clang` to use Clang)

##### Linux

- GCC

#### Static library

```shell
go build -o build/ddns4cdn.a -buildmode=c-archive ./src/cgo/go
```

##### GCC

```shell
# C
gcc -o build/ddns4cdn_c src/cgo/c/main.c build/ddns4cdn.a
# C++
g++ -o build/ddns4cdn_cpp src/cgo/cpp/main.cc build/ddns4cdn.a
```

##### Clang

```shell
# C
clang -o build/ddns4cdn_c src/cgo/c/main.c build/ddns4cdn.a -framework CoreFoundation -framework Security -lresolv
# C++
clang++ -o build/ddns4cdn_cpp src/cgo/cpp/main.cc build/ddns4cdn.a -framework CoreFoundation -framework Security -lresolv
# Objective-C
clang -o build/ddns4cdn_objc src/cgo/objc/main.m build/ddns4cdn.a -framework Foundation -framework Security -lresolv
```

##### MSVC

```powershell
& "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\Launch-VsDevShell.ps1" -Arch amd64 -HostArch amd64
# C
cl /MD /Fe"build\ddns4cdn_msvc_c.exe" /Fo"build\ddns4cdn_msvc_c.exe.obj" src\cgo\c\main.c /link build\ddns4cdn.a
# C++
cl /EHsc /MD /Fe"build\ddns4cdn_msvc_cpp.exe" /Fo"build\ddns4cdn_msvc_cpp.exe.obj" src\cgo\cpp\main.cc /link build\ddns4cdn.a
```

#### Shared library

```shell
go build -o build/ddns4cdn.so -buildmode=c-shared ./src/cgo/go
```

##### GCC (Windows & Linux)

```shell
# C
gcc -o build/ddns4cdn_dl_c src/cgo/c/main.c build/ddns4cdn.so
# C++
g++ -o build/ddns4cdn_dl_cpp src/cgo/cpp/main.cc build/ddns4cdn.so
```

##### Clang (macOS)

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

##### MSVC

```powershell
& "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\Launch-VsDevShell.ps1" -Arch amd64 -HostArch amd64
cl /LD /MD /Fe"build\ddns4cdn.dll" /Fo"build\ddns4cdn.dll.obj" src\cgo\msvc\ddns4cdn.c /link /DEF:src\cgo\msvc\ddns4cdn.def build\ddns4cdn.a

# C
cl /Fe"build\ddns4cdn_msvc_dl_c.exe" /Fo"build\ddns4cdn_msvc_dl_c.exe.obj" src\cgo\c\main.c /link build\ddns4cdn.lib
# C++
cl /EHsc /Fe"build\ddns4cdn_msvc_dl_cpp.exe" /Fo"build\ddns4cdn_msvc_dl_cpp.exe.obj" src\cgo\cpp\main.cc /link build\ddns4cdn.lib
```

#### Swift

```shell
cd src/cgo/swift
# the output binary is .build_dl/release/ddns4cdn
# and .build_dl/release/ddns4cdn_dl
swift build -c release

# to use ddns4cdn_dl, set DYLD_LIBRARY_PATH to load shared library
# ddns4cdn linked a static library, so you don't need to set it
export DYLD_LIBRARY_PATH=$(realpath ../../../build)
```

### Go mobile

#### Requirements

Android:

iOS:

#### Install gomobile

```shell
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

#### Build

```shell
cd src/gomobile
# Android
gomobile bind -o ../../build/ddns4cdn.aar -target android -javapkg ddns4cdn -androidapi 33
# iOS and macOS
gomobile bind -o ../../build/Ddns4cdn.xcframework -target ios,macos -prefix Ddns4cdn -iosversion 17
```
