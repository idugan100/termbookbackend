#! /bin/zsh
#
# linux64
GOOS=linux GOARCH=amd64 go build -o ./builds/lin64/termbook


# linux32 
GOOS=linux GOARCH=386 go build -o ./builds/lin32/termbook



# win64 
GOOS=windows GOARCH=amd64 go build -o ./builds/win64/termbook.exe

# win32 
GOOS=windows GOARCH=386 go build -o ./builds/win32/termbook.exe



# AppI
GOOS=darwin GOARCH=amd64 go build -o ./builds/AppI/termbook

# AppS
GOOS=darwin GOARCH=arm64 go build -o ./builds/AppS/termbook
