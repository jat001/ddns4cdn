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
            path: "Sources"
        ),
    ]
)
