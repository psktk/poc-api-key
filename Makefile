.PHONY: setup run clean

setup:
	go mod tidy

run:
	go run main.go

clean:
	rm -f products.db
