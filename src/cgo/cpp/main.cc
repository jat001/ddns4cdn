#include <iostream>
#include <fstream>

#include "../../../build/ddns4cdn.h"

#ifdef _MSC_VER
#ifdef __cplusplus
extern "C"
{
#endif

    __pragma(comment(lib, "legacy_stdio_definitions.lib"));

    // https://github.com/golang/go/issues/42190#issuecomment-1507839987
    void _rt0_amd64_windows_lib();

    __pragma(section(".CRT$XCU", read));
    __declspec(allocate(".CRT$XCU")) void (*init_lib)() = _rt0_amd64_windows_lib;

    __pragma(comment(linker, "/include:init_lib"));

#ifdef __cplusplus
}
#endif
#endif

int main(int argc, char *argv[])
{
    std::string config = "config.toml";
    if (argc > 1)
    {
        config = argv[1];
    }

    std::ifstream file(config);
    if (!file.is_open())
    {
        std::cerr << "Error opening file: " << config << std::endl;
        return 1;
    }

    std::string data(
        (std::istreambuf_iterator<char>(file)),
        std::istreambuf_iterator<char>());
    file.close();

    Ddns4cdnWorker(const_cast<char *>(data.c_str()));

    return 0;
}
