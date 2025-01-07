package bencode

import (
	"fmt"
	"strconv"
	"bufio"
	"os"
	"log"
	"strings"
)

var in *bufio.Reader

var typeMap = map[byte]string {
    0 : "benmap", 
    1 : "benlist", 
    2 : "benstring", 
    3 : "benint", 
}

// small diagnostic tool I wrote to debug and observe encoded files. to use, call the launcher function bellow
// with the benval struct from the parser, as well as a string label for the thing that u are parsing

func BencodeDiagnostics(m Benmap, s string) {
	in = bufio.NewReader(os.Stdin)
	for true {
		fmt.Println("PROPERTIES: of " + s + "-----------------------------")
        mkeys := m.Keys()
		for _, key := range mkeys {
			fmt.Printf("KEY: \"%s\"    ", key)
            val, _ := m.Query(key)
			fmt.Printf("TYPE: %s \n", typeMap[(*val).Kind()])
		}
		fmt.Println("-----------------------------------------------------")
		fmt.Println("to exit type exit, otherwise enter a valid key to expand")

		c := readIn()
		if c == "exit" {
			return
		}

		val, valid := m.Query(c)
		for !valid {
			fmt.Println("invalid key, try again or exit")
			c = readIn()
			if c == "exit" {
				return
			}
			val, valid = m.Query(c)
		}
		Launcher(*val, c)
	}
}

func readIn() string {
	c, err := in.ReadString('\n')
	if err != nil {
		log.Fatal("failed to read input")
	}
	return strings.Trim(c, "\n")
}

func Launcher(val Benval, s string) {
	t := val.Kind()

	if t == 0 {
		BencodeDiagnostics(val.(Benmap), s)
	} else if t == 1 {
		listD(val.(Benlist), s)
	} else if t == 2 {
		bytesD(val.(Benstring), s)
	} else if t == 3 {
		intD(val.(Benint), s)
	}
}

func listD(val Benlist, s string) {
	for true {
		fmt.Println("LIST of " + s + "------------------------------------")
		fmt.Println("list size: " + strconv.Itoa(val.Len()))
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("LIST ITEM %d:", i)
			fmt.Printf("    TYPE: %s\n", typeMap[(*val.Get(i)).Kind()])
		}
		fmt.Println("-----------------------------------------------------")
		fmt.Println("input an item to expand, type exit to exit")
		c := readIn()

		if c == "exit" {
			return
		}

		choice, err := strconv.Atoi(c)
		for err != nil || choice > val.Len() {
			fmt.Println("invalid choice, try again")
			c = readIn()
			if c == "exit" {
				return
			}
			choice, err = strconv.Atoi(c)
		}
		id := s + "[" + c + "]"
		Launcher((*val.Get(choice)), id)
	}
}

func intD(val Benint, s string) {
	fmt.Println(s + " VALUE--------------------------------------------")
	fmt.Println(strconv.Itoa(int(val.Get())))
	fmt.Println("------------------------------------------------------")
	fmt.Println("press enter to return...")
	fmt.Scanln()
}

func bytesD(val Benstring, s string) {
	fmt.Println("------------------------------------------------------")
	fmt.Println(s + " BYTES; SIZE: " + strconv.Itoa(len(val.Get())))
    if len(val.Get()) > 150 {
        fmt.Println("size larger than 150 bytes... print to log.txt?")
        c := readIn() 
        if c == "y" {
            os.WriteFile("log.txt", val.Get(), 0644)
        }
    } else {
        fmt.Print(string(val.Get()) + "\n")
    }
	fmt.Println("------------------------------------------------------")
	fmt.Println("press enter to return...")
	fmt.Scanln()
}
