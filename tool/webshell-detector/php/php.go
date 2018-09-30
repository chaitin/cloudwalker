package php

/*
#cgo CFLAGS: -Wall -Wextra -Werror -Wno-unused-parameter -I${SRCDIR}/include/php -I${SRCDIR}/include/php/Zend -I${SRCDIR}/include/php/TSRM -I${SRCDIR}/include/php/main
#cgo LDFLAGS: ${SRCDIR}/lib/libphp7.a -lm

#include <stdint.h>

int init(intptr_t fd_in, intptr_t fd_out);
int execute(void);
*/
import "C"

import (
    "errors"
    "os"
)

var Stdin, Stdout *os.File

func init() {
    var err error
    var stdin, stdout *os.File
    stdin, Stdin, err = os.Pipe()
    if err != nil { panic(err) }
    Stdout, stdout, err = os.Pipe()
    if err != nil { panic(err) }

    if ret := C.init(C.intptr_t(stdin.Fd()), C.intptr_t(stdout.Fd())); ret != 0 {
        panic("cannot initialize PHP runtime")
    }
}

func Start() error {
    if (C.execute() != 0) {
        return errors.New("cannot start PHP runtime")
    }
    return nil
}
