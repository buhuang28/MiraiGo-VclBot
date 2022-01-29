chcp 65001
go build -v -ldflags "-s -w -extldflags '-static' -H windowsgui" -tags tempdll -o vclbot.exe