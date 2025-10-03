# proto-go
protohackers



```

echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker build -t proto/smoke-test .      
docker run -p 8080:8080 proto/smoke-test

docker push ghcr.io/mobiljani/smoke-test 
```
