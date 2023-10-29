import ArgumentParser
import Foundation
import Worker

@main
struct Ddns4cdn: ParsableCommand {
    @Option(name: [.short, .customLong("config")], help: "config file path")
    var config = "config.toml"

    mutating func run() throws {
        let data = try String(contentsOfFile: config)
        Ddns4cdnWorker(strdup(data))
    }
}
