# benparser 

benparser is a quick bencode format parser that I wrote to assist with a couple of other projects I am currently working on. 

Bencode's specification can be found on its wikipedia page, here: https://en.wikipedia.org/wiki/Bencode

## Use 

To use, go get the package 

``` go get github.com/rmcs9/benparser ```

Call either ``` benparser.ParseFile() ``` or ``` benparser.ParseBytes()``` 

Additionally, there is an included diagnostic function to view and navigate the structured contents of an encoded file. this function can be found in ```diag.go``` and can be accessed via ```benparser.Launcher()```
