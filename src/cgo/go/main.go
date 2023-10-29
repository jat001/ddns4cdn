package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/jat001/ddns4cdn/worker"
)

//export Ddns4cdnWorker
func Ddns4cdnWorker(data *C.char) {
	raw := []byte(C.GoString(data))
	C.free(unsafe.Pointer(data))
	worker.Worker(raw)
}

func main() {}
