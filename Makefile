build:
	go build -o autoassigner main.go

run-alpha:
	./autoassigner team-alpha

clean:
	rm -f autoassigner