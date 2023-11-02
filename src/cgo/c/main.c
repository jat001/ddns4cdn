#include <stdio.h>

#include "../../../build/ddns4cdn.h"

#ifdef _MSC_VER
#include "../msvc/cgo.c"
#endif

int main(int argc, char *argv[]) {
  char *config = "config.toml";
  if (argc > 1) {
    config = argv[1];
  }

  FILE *fp;
  errno_t err = fopen_s(&fp, config, "r");
  if (err != 0) {
    printf("Error opening file: %s\n", config);
    return 1;
  }

  fseek(fp, 0, SEEK_END);
  long size = ftell(fp);
  fseek(fp, 0, SEEK_SET);

  char *data = malloc(size + 1);
  if (data == NULL) {
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
