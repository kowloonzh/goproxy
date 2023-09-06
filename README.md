# goproxy
a simple http and https forward proxy

## quick start
```
go install github.com/kowloonzh/goproxy@latest

goproxy 

```

## set addr or basic auth with param
```
auth=$(echo "admin:xxxx" |base64)

goproxy -l :8887 -a $auth
```

## set  addr or basic auth with ENV
```

auth=$(echo "admin:xxxx" |base64)

GOPROXY_AUTH=${auth} GOPROXY_ADDR=:8887 goproxy
```