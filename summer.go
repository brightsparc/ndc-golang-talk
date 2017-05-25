package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Summer interface
type Summer interface {
	Sum() (int, error)
}

// Config class contains logger
type summerConfig struct {
	Logger *log.Logger
}

// Summer type ctonains reader and config
type summer struct {
	Reader io.Reader
	Config summerConfig
}

// NewSummer initializes the summer with optional logger config
func NewSummer(reader io.Reader, config summerConfig) Summer {
	if config.Logger == nil {
		config.Logger = log.New(ioutil.Discard, "", 0) // zero as default
	}
	return &summer{
		Reader: reader,
		Config: config,
	}
}

// Sum implmenetation enumerates the file
func (s *summer) Sum() (int, error) {
	total := 0

	scanner := bufio.NewScanner(s.Reader)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return total, err
		}
		total += num
		s.Config.Logger.Println(num)
	}

	return total, nil
}

func main() {
	// Defer will run before function exists
	defer func(t time.Time) {
		fmt.Println("done", time.Since(t))
	}(time.Now())

	reader := strings.NewReader("1\n2\n3\n4\nA")
	summer := NewSummer(reader, summerConfig{
		Logger: log.New(os.Stdout, "", log.Lmicroseconds),
	})
	fmt.Println(summer.Sum())
}
