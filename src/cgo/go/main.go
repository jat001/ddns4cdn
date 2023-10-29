package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/jat001/ddns4cdn/worker"
)

//export Worker
func Worker(data *C.char) {
	defer C.free(unsafe.Pointer(data))
	worker.Worker([]byte(C.GoString(data)))
}

func main() {}
