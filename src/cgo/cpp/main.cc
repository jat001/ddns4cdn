#include <iostream>
#include <fstream>

#include "../../../build/ddns4cdn.h"

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
