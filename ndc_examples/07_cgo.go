package main

/*
// Everything in comments above the import "C" is C code and will be compiles with the GCC.
#include "stdio.h"
int addInC(int a, int b) {
  return a + b;
}
*/
import "C"
import "fmt"

func main() {
	a, b := 3, 5
	c := C.addInC(C.int(a), C.int(b))
	fmt.Printf("Add in C: %d + %d = %d", a, b, int(c))
}
