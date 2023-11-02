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

  [WinLibs](https://winlibs.com/#download-release)'s build is recommended, tested with MCF and UCRT.

- [Microsoft C++ Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/)  (MSVC, optional)

  `go build` doesn't support MSVC, but you can use MSVC to build exe and dll.

##### macOS / Linux

GCC or Clang

`go build` can use GCC or Clang. Default is GCC, set `CC=clang` to use Clang.

#### Static library

##### Prepare the library

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
# only supports amd64
& "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\Launch-VsDevShell.ps1" -Arch amd64 -HostArch amd64

# C
cl /MD /Fe"build\ddns4cdn_msvc_c.exe" /Fo"build\ddns4cdn_msvc_c.obj" src\cgo\c\main.c /link build\ddns4cdn.a
# C++
cl /EHsc /MD /Fe"build\ddns4cdn_msvc_cpp.exe" /Fo"build\ddns4cdn_msvc_cpp.obj" src\cgo\cpp\main.cc /link build\ddns4cdn.a
```

#### Shared library

##### Prepare the library

```shell
go build -o build/ddns4cdn.so -buildmode=c-shared ./src/cgo/go
```

##### GCC

```shell
# C
gcc -o build/ddns4cdn_dl_c src/cgo/c/main.c build/ddns4cdn.so
# C++
g++ -o build/ddns4cdn_dl_cpp src/cgo/cpp/main.cc build/ddns4cdn.so
```

##### Clang

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

Building a dynamic-link library (.dll) requires a static library (.a) as a link file.

See the previous section for how to build a static library.

```shell

```powershell
# only supports amd64
& "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\Launch-VsDevShell.ps1" -Arch amd64 -HostArch amd64

# DLL
cl /LD /MD /Fe"build\ddns4cdn.dll" /Fo"build\ddns4cdn.obj" src\cgo\msvc\ddns4cdn.c /link /DEF:src\cgo\msvc\ddns4cdn.def build\ddns4cdn.a

# C
cl /DNO_CGO_LIB /Fe"build\ddns4cdn_msvc_dl_c.exe" /Fo"build\ddns4cdn_msvc_dl_c.obj" src\cgo\c\main.c /link build\ddns4cdn.lib
# C++
cl /DNO_CGO_LIB /EHsc /Fe"build\ddns4cdn_msvc_dl_cpp.exe" /Fo"build\ddns4cdn_msvc_dl_cpp.obj" src\cgo\cpp\main.cc /link build\ddns4cdn.lib
```

#### C\#

Compiling C# code requires a dynamic-link library (.dll).

See the previous section for how to build a dll.

```shell
cd src/cgo/csharp
dotnet publish
```

#### Swift

```shell
cd src/cgo/swift
# the output binary is .build/release/ddns4cdn
# and .build/release/ddns4cdn_dl
swift build -c release

# to use ddns4cdn_dl, set DYLD_LIBRARY_PATH to load shared library
# ddns4cdn linked a static library, so you don't need to set it
export DYLD_LIBRARY_PATH=$(realpath ../../../build)
```

### Go mobile

#### Requirements

Android:

Android SDK Command-Line Tools

iOS:

Xcode Command Line Tools

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
