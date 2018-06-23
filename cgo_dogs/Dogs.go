package main

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Foundation
//#cgo LDFLAGS: -framework Realm
//#include "Dogs.h"
import "C"

func main() {
	C.CountDogs()
}
