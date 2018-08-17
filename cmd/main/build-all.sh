GOOS=windows GOARCH=amd64 go build -o httpshared-windows-amd64.exe
GOOS=linux   GOARCH=amd64 go build -o httpshared-linux-amd64
GOOS=darwin  GOARCH=amd64 go build -o httpshared-darwin-amd64

zip -r httpshared-windows-amd64.exe.zip  httpshared-windows-amd64.exe
tar cvfj httpshared-linux-amd64.tar.bz2  httpshared-linux-amd64
tar cvfj httpshared-darwin-amd64.tar.bz2 httpshared-darwin-amd64