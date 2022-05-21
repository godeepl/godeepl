# godeepl
Golang Client to translate md from English to Chinese by https://www.deepl.com/

## dev

```
go mod vendor
go main.go translator --file ./utils/samples/test-1.md \
  -s EN \
  -t ZH \
  -p http://47.243.244.86:8080/translate
```
