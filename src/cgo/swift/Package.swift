// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "Ddns4cdn",
    dependencies: [
        .package(
            url: "https://github.com/apple/swift-argument-parser.git",
            from: "1.2.3"
        )
    ],
    targets: [
        .systemLibrary(
            name: "worker",
            path: "Includes"
        ),
        .executableTarget(
            name: "ddns4cdn",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser"),
                .target(name: "worker")
            ],
            path: "Sources"
        )
    ]
)
