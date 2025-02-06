package benparser


import (
	"os"
	"log"
	"strconv"
)

var fm *fileManager

//top level parse function... takes a filepath to a bencode file
func ParseFile(path string) Benval {
	//read the file 
	file, err := os.ReadFile(path)
	if err != nil {
        log.Fatal("unable to read file:", path, "\n", err)
	}

	//create a new filemanager object
	fm = BuildFM(file)
	//begin parsing
	ret := parseItem()
	fm = nil
	return ret
}

func ParseBytes(bytes []byte) Benval {
    fm = BuildFM(bytes)

    ret := parseItem()
    fm = nil 
    return ret
}

// parses a bencode dictionary
// dictionaries in bencode are formatted as follows:
// dkey:valuekey:value....e 
// where key is always a byteString and value is any of bencodes data types (int, string, list, dict)
func parseDict() Benmap {
    //obtain raw bytes 
    bytes := collect_bytes()

	//ensure that the opening 'd' is present
	if fm.Peek(0) != 'd' {
		log.Fatal("invalid start to bencode object...")
	}
	fm.Absorb(1)
	
    valmap := make(map[string]*Benval)
	//while we are not at the end of the dict..
	for fm.Peek(0) != 'e' {
		//take the key
		key := parseKey()
		//take the value
		value := parseItem()
		//add the pair into the map
		valmap[key] = &value
	}
	//absorb the e
	fm.Absorb(1)
	//return the map
	return Benmap{ valmap, bytes }
}

//a midway point that determines what type the upcoming data is
func parseItem() Benval{
	if fm.Peek(0) == 'd' {
		return parseDict() 
	} else if fm.Peek(0) == 'i' {
		return parseInt()
	} else if fm.Peek(0) == 'l' {
		return parseList()
	} else {
		return parseByteString()
	}
}

// parses bencode ints... ints are formatted as follows:
// i124782385435e
func parseInt() Benint {
    bytes := collect_bytes()

    //absorb the i
    fm.Absorb(1)

	//find the end of the encoded int
	i := fm.Find('e')
	//obtained the sliced bytes
	num := fm.Pop(i)
	//absorb the e
	fm.Absorb(1)
	//convert the bytes to an int and return
    numint, err := strconv.Atoi(string(num))
	if err != nil {
        log.Fatal("unable to parse bencode int: ", string(num))
	}
	return Benint{ int64(numint), bytes }
}

//parses a bencode list, which is encoded as follows:
// lvaluevaluevaluevalue...e where value can be any valid bencode type
func parseList() Benlist {
    bytes := collect_bytes()

	//absorb the l
	fm.Absorb(1)

	//while not at the end of the list, parse new items
	//and add accumulate them in a slice
	var list []*Benval
	for fm.Peek(0) != 'e' {
        item := parseItem()
		list = append(list, &item)
	}
	//absorb the e
	fm.Absorb(1)
	return Benlist{ list, bytes }
}

//parses a bencode bytestring which is formated as follows
// <bytelength>:bytes 
func parseKey() string {
	i := fm.Find(':')
	//obtain the length of the bytestring
	numOfBytes, err := strconv.Atoi(string(fm.Pop(i)))
	if err != nil {
		log.Fatal("unable to parse the length of the key at position ", fm.GetPoint())
	}
	//absorb the colon
	fm.Absorb(1)
	//return the byte string
	return string(fm.Pop(numOfBytes))
}

func parseByteString() Benstring {
	i := fm.Find(':')
	numOfBytes, err := strconv.Atoi(string(fm.Pop(i)))
	if err != nil {
		log.Fatal("unable to parse the length of the bencode bytestring at position ", fm.GetPoint())
	}

	fm.Absorb(1)
    bytes := fm.Pop(numOfBytes)
	return Benstring{ bytes, bytes } 
}
