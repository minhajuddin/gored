package gored

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Redis struct {
	Connection net.Conn
	Reader     *bufio.Reader
}

func (self *Redis) run(args ...string) {
	//build command
	cmd := "*" + strconv.Itoa(len(args)) + "\r\n"
	for _, t := range args {
		cmd += fmt.Sprintf("$%d\r\n%v\r\n", len(t), t)
	}
	fmt.Fprint(self.Connection, cmd)
}

func (self *Redis) Ping() (string, error) {
	self.run("PING")
	op, err := self.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(op, "+\r\n"), nil
}

func (self *Redis) Set(key string, value string) error {
	self.run("SET", key, value)
	line, err := self.Reader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Println(line)
	return nil
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
		Reader:     bufio.NewReader(c),
	}

	return r, nil
}

func Hello() {
	log.Println("Hello redis ;)")
}
