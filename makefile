all:
	go build
	./karlmcguire.com
	cd docs/ && python -m SimpleHTTPServer
