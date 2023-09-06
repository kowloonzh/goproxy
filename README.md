# goproxy
a simple http and https forward proxy

# quick start
```
go install github.com/kowloonzh/goproxy@latest

goproxy 

```

# set addr or basic auth
```
auth=$(echo "admin:xxxx" |base64)
goproxy -l :8887 -a $auth
```