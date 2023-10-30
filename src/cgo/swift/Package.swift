// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "Ddns4cdn",
    platforms: [
        .macOS(.v14),
    ],
    products: [
        .executable(
            name: "ddns4cdn",
            targets: ["main"]
        ),
        .executable(
            name: "ddns4cdn_dl",
            targets: ["main_dl"]
        ),
    ],
    dependencies: [
        .package(
            url: "https://github.com/apple/swift-argument-parser.git",
            from: "1.2.3"
        ),
    ],
    targets: [
        .systemLibrary(
            name: "worker",
            path: "Includes"
        ),
        .executableTarget(
            name: "main",
            dependencies: [
                .product(
                    name: "ArgumentParser",
                    package: "swift-argument-parser"
                ),
                .target(name: "worker"),
            ],
            path: "Sources",
            linkerSettings: [
                .unsafeFlags([
                    "../../../build/ddns4cdn.a",
                    "-lresolv",
                ]),
            ]
        ),
        .executableTarget(
            name: "main_dl",
            dependencies: [
                .product(
                    name: "ArgumentParser",
                    package: "swift-argument-parser"
                ),
                .target(name: "worker"),
            ],
            // just a symlink to Sources to make SwiftPM happy
            path: "SourcesDl",
            linkerSettings: [
                .unsafeFlags([
                    "../../../build/ddns4cdn.so",
                ]),
            ]
        ),
    ]
)
