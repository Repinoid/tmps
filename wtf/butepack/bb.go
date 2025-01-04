package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	var b bytes.Buffer // A Buffer needs no initialization.
	b.Grow(100)
	fmt.Printf(" %v %v\n", b, &b)

	b.WriteString("0123456789")

	var r io.Reader
	r = strings.NewReader("my request")
	//	buf := make([]byte, 18)
	buf, err := io.ReadAll(r)
	if err != nil {
		log.Printf("%v\n%s\n", err, buf)
	}
	sugar.Infoln(
		"buf", string(buf),
		zap.String("buf", string(buf)),
	)
	t:= time.Now()
	
	f:= "Mon, 01 Jan 2006 15:04:05.0000 -07 MST "
	fmt.Println(t.Format(f))
}
