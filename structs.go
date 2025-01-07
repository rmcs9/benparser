package bencode

type Benval interface {
    Kind() byte
}

// --------- BENCODED MAPS ---------
type Benmap struct {
    valmap map[string]*Benval  
    raw []byte
}

func (b Benmap) Kind() byte {
    return 0
}

func (b Benmap) Keys() []string {
    keys := make([]string, 0)

    for key := range b.valmap {
        keys = append(keys, key)
    }
    return keys
}

func (b Benmap) Query(key string) (*Benval, bool) {
    val, has := b.valmap[key]
    return val, has
}

// --------- BENCODED LISTS ---------
type Benlist struct {
    vallist []*Benval
    raw []byte
}

func (b Benlist) Kind() byte {
    return 1
}

func (b Benlist) Len() int {
    return len(b.vallist)
}

func (b Benlist) Get(i int) *Benval {
    return b.vallist[i]
}

// --------- BENCODED BYTESTRINGS ---------
type Benstring struct {
    valstring []byte
    raw []byte
}

func (b Benstring) Kind() byte {
    return 2
}

func (b Benstring) Get() []byte {
    return b.valstring
}

// --------- BENCODED INTEGERS ---------
type Benint struct {
    valint int64
    raw []byte
}

func (b Benint) Kind() byte {
    return 3
}

func (b Benint) Get() int64 {
    return b.valint
}

const (
    Map     byte = 0
    List    byte = 1
    String  byte = 2
    Int     byte = 3
)
