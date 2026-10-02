// hello 是第一個可執行程式。
// 執行方式：go run ./cmd/hello -name 你的名字
package main

import (
	"flag"
	"fmt"

	"github.com/r029httc/llgo/lessons/01-basics"
)

func main() {
	name := flag.String("name", "Gopher", "要打招呼的對象")
	flag.Parse()

	fmt.Println(basics.Greet(*name))
}
