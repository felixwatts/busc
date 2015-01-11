package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/felixwatts/bus"
	"os"
	"strings"
)

var serverAddr string

func init() {
	flag.StringVar(&serverAddr, "server", "localhost:8888", "Server address in the form 'hostname:port'")
}

func main() {

	flag.Parse()

	c, err := bus.Dial(serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Ready. Enter 'help' for help.")
	}

	go printReplies(c)

	bio := bufio.NewReader(os.Stdin)

	for {
		line, _, err := bio.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}

		parts := strings.Split(string(line), " ")

		numParts := len(parts)

		if numParts == 0 {
			fmt.Println("Invalid command. Empty")
			continue
		}

		var cmd = parts[0]

		switch cmd {
		case "quit":
		case "exit":
			return
		case "+":
			if numParts != 2 {
				fmt.Println("Invalid command. Missing key or extra params")
				continue
			}
			key := parts[1]
			_, err := c.Subscribe(key)
			if err != nil {
				fmt.Println(err)
			}

			break
		case "-":
			if numParts != 2 {
				fmt.Println("Invalid command. Missing key or extra params")
				continue
			}
			key := parts[1]
			_, err := c.Unsubscribe(key)
			if err != nil {
				fmt.Println(err)
			}

			break
		case ">":
			if numParts != 3 {
				fmt.Println("Invalid command. Missing key, val or extra params")
				continue
			}
			key := parts[1]
			val := parts[2]
			_, err := c.Publish(key, val)
			if err != nil {
				fmt.Println(err)
			}

		case "|":
			if numParts != 2 {
				fmt.Println("Invalid command. Missing key or extra params")
				continue
			}
			key := parts[1]
			_, err := c.Claim(key)
			if err != nil {
				fmt.Println(err)
			}

			break
		case "help":
			fmt.Print(`Commands:
+ xxx        Subscribe to key 'xxx'
- xxx        Unsubscribe from key 'xxx'
> xxx yyy    Publish value 'yyy' at key 'xxx'
| xxx        Claim the keyspace rooted at 'xxx'
exit         Exit the client
help         Show this help
`)
		default:
			fmt.Println("Invalid command. Try 'help'.")
		}
	}
}

func printReplies(c bus.Client) {
	for {
		m := <-c.Rxc()
		fmt.Print(m)
	}
}
