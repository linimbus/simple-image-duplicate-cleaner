rsrc -manifest exe.manifest -ico static/main.ico
rice embed-go
go build -ldflags="-H windowsgui -w -s" -o simple-image-duplicate-cleaner.exe