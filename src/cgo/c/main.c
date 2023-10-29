#include <stdio.h>
#include <stdlib.h>

#include "../../../ddns4cdn.h"

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

    Worker(data);

    return 0;
}
