all:
	dep ensure
	env GOOS=linux GOARCH=arm GOARM=5 go build
	tar -cvzf parceldrop.tar.gz parceldrop parceldrop.service service.env INSTALL.sh
