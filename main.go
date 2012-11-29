package gored

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Redis struct {
	Connection net.Conn
	Reader     *bufio.Reader
}

func (self *Redis) write(args ...string) {
	cmd := "*" + strconv.Itoa(len(args)) + "\r\n"
	for _, t := range args {
		cmd += fmt.Sprintf("$%d\r\n%v\r\n", len(t), t)
	}
	fmt.Fprint(self.Connection, cmd)
}

func (self *Redis) readline() (string, error) {
	line, err := self.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(line, "\r\n"), nil
}

func (self *Redis) read() (string, error) {
	sym, err := self.Reader.ReadByte()
	switch sym {
	case '-':
		line, err := self.readline()
		if err == nil {
			err = errors.New(line)
		}
		return "", err
	case '+':
		line, err := self.readline()
		if err != nil {
			return "", err
		}
		return line, nil
	case '*':
	default:
		return "", errors.New("Redis protocol error")
	}

	return "NOT IMPL", err
}

func (self *Redis) Ping() (string, error) {
	self.write("PING")
	return self.read()
}

func (self *Redis) Set(key string, value string) (string, error) {
	self.write("SET", key, value)
	return self.read()
}

func (self *Redis) Get(key string) (string, error) {
	self.write("GET", key)
	return self.read()
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
