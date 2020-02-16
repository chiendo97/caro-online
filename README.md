# Caro-online

## Introduction

![demo](docs/demo.gif)

Caro is the tic-tac-to game with 20x20 board.

I wrote this project for learning golang.

## Play online

```bash
export host=caro-game-online.herokuapp.com
go run ./cmd/client/main.go
```

## Installation

* Install [golang](https://golang.org/doc/install), git

```bash
go get https://github.com/chiendo97/caro-online
```

## Play offline

### Running server

Running server on :8080 port

```bass
go run ./cmd/server/main.go
```

### Running client

Run 2 clients on different terminals to play to each others

```bass
go run ./cmd/client/main.go
```

## Future

* Benchmark
* Bot
* API docs
* Player login
* Game history
* Database

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
