/*
 * USAGE INSTRUCTIONS:
 *
 * make sure libmysqlclient-dev is installed:
 * apt install libmysqlclient-dev
 *
 * Replace "/usr/lib/mysql/plugin" with your MySQL plugins directory (can be found by running "select @@plugin_dir;")
 * go build -buildmode=c-shared -o b64u.so && cp b64u.so /usr/lib/mysql/plugin/
 *
 * Then, on the server:
 * create function`b64u`returns string soname'b64u.so';
 *
 * And use/test like:
 * select`b64u`('eWVldA'); -- outputs yeet
 *
 * Yeet!
 * Brian Leishman
 *
 */

package main

// #include <stdio.h>
// #include <sys/types.h>
// #include <sys/stat.h>
// #include <stdlib.h>
// #include <string.h>
// #include <mysql.h>
// #cgo CFLAGS: -I/usr/include/mysql -fno-omit-frame-pointer -O3
import "C"
import (
	"encoding/base64"
	"strings"
	"unsafe"
)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

var base64URLEncode = strings.NewReplacer("+", "-", "/", "_")

// B64u encodes bytes to URL safe base64
func B64u(b []byte) string {
	return strings.TrimRight(base64URLEncode.Replace(base64.StdEncoding.EncodeToString(b)), "=")
}

//export b64u_deinit
func b64u_deinit(initid *C.UDF_INIT) {}

//export b64u_init
func b64u_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.my_bool {
	if args.arg_count != 1 {
		msg(message, "`b64u` requires 1 parameters: the string to be encoded")
		return 1
	}

	argsTypes := (*[1]uint32)(unsafe.Pointer(args.arg_type))
	argsTypes[0] = C.STRING_RESULT

	initid.maybe_null = 1

	return 0
}

//export b64u
func b64u(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	arguments := 1

	// f, _ := ioutil.TempFile(os.TempDir(), "mysql-b64u")

	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:arguments:arguments]
	// fmt.Fprintf(f, "argsArgs: %#v\n", argsArgs)

	argsLengths := (*[1 << 30]C.int)(unsafe.Pointer(args.lengths))[:1:1]
	// fmt.Fprintf(f, "argsLengths: %#v\n", argsLengths)

	a := make([]string, arguments, arguments)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}
	// fmt.Fprintf(f, "a: %#v\n", a)

	b := B64u([]byte(a[0]))
	// fmt.Fprintf(f, "b: %#v\n", b)

	*length = uint64(len(b))
	return C.CString(string(b))
}

func main() {}
