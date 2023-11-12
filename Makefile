bench: cmd/bench/main.go
	go build -o bench cmd/bench/main.go

populator: cmd/populator/main.go
	go build -o populator cmd/populator/main.go 
