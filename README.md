
# googleit

<img src="https://img.shields.io/badge/coverage-83%25-brightgreen.svg?style=flat-square" alt="Code coverage">&nbsp;<a href="https://travis-ci.org/schollz/googleit"><img src="https://img.shields.io/travis/schollz/googleit.svg?style=flat-square" alt="Build Status"></a>&nbsp;<a href="https://godoc.org/github.com/schollz/googleit"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square" alt="Go Doc"></a> 

## Install

```
go get github.com/schollz/googleit
```

## Usage 


```golang
searchTerm := "cats wiki"
numPages := 1
urls, _ := googleit.Search(searchTerm, numPages)
fmt.Println(urls[0])
// https://en.wikipedia.org/wiki/Cat
```

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT
