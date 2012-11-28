package gored

import (
	"strconv"
	"log"
	"net"
	"fmt"
	"bufio"
)

type Redis struct{
	Connection net.Conn
}

func (self *Redis) Run(args ...string) {
	//build command
	cmd := "*" + strconv.Itoa(len(args)) + "\r\n"
	for _, t := range args {
		cmd += fmt.Sprintf("$%d\r\n%v\r\n", len(t), t)
	}
	fmt.Fprint(self.Connection, cmd)
	r := bufio.NewReader(self.Connection)
	log.Println(r.ReadString('\n'))
}

func (self *Redis) Close() {
	self.Connection.Close()
}

//TODO: ability to pass a host and port
func New() (*Redis, error) {
	c, e := net.Dial("tcp", "localhost:6379")
	if e != nil {
		return nil, e
	}
	r := &Redis{
		Connection: c,
	}

	return r, nil
}

func Hello(){
	log.Println("Hello redis ;)")
}
