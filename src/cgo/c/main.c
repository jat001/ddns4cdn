#include <stdio.h>

#include "../../../build/ddns4cdn.h"

#ifdef _MSC_VER
__pragma(comment(lib, "legacy_stdio_definitions.lib"));

// https://github.com/golang/go/issues/42190#issuecomment-1507839987
void _rt0_amd64_windows_lib();

__pragma(section(".CRT$XCU", read));
__declspec(allocate(".CRT$XCU")) void (*init_lib)() = _rt0_amd64_windows_lib;

__pragma(comment(linker, "/include:init_lib"));
#endif

int main(int argc, char *argv[])
{
    char *config = "config.toml";
    if (argc > 1)
    {
        config = argv[1];
    }

    FILE *fp = fopen(config, "r");
    if (fp == NULL)
    {
        printf("Error opening file %s\n", config);
        return 1;
    }

    fseek(fp, 0, SEEK_END);
    long size = ftell(fp);
    fseek(fp, 0, SEEK_SET);

    char *data = malloc(size + 1);
    if (data == NULL)
    {
        printf("Error allocating memory\n");
        fclose(fp);
        return 1;
    }

    fread(data, 1, size, fp);
    fclose(fp);
    data[size] = '\0';

    Ddns4cdnWorker(data);

    return 0;
}
