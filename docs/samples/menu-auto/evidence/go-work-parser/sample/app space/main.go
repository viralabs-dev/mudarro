package main
import("fmt";"example.test/lib";"os";"time")
func main(){fmt.Println(lib.Value());if len(os.Args)>1{for{time.Sleep(time.Second)}}}
