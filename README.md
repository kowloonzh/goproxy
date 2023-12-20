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

## set HTTP_PROXY and NO_PRPXY
```
export no_proxy="localhost,127.0.0.1,.localdomain.com,*.localdomain.com,10.0.0.1/8"
export http_proxy="http://proxy.cn:80"
export https_proxy=$http_proxy

goproxy
```