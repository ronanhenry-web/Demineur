# 🧨 GoMineSweeper

A simple terminal-based Minesweeper implementation written in Go.

## 📦 Project Structure

- `minesweeper.go`: core logic for the Minesweeper game.
- `benchmark_test.go`: performance benchmarks.
- `go.mod`: Go module definition.

## 🚀 Getting Started

### Install dependencies

```bash
go mod tidy
```

### Run the game

```bash
go run minesweeper.go
```

### Run benchmarks

```bash
go test -bench=. -benchmem
```

### Build the project

```bash
go build -o gominesweeper
```
