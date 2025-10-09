# proto-go
protohackers



```

echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker build --target runtime -t proto/smoke-test .      
docker run -p 8080:8080 proto/smoke-test

docker push ghcr.io/mobiljani/smoke-test 


https://us-east-1.console.aws.amazon.com/ecs/v2/clusters/blushing-butterfly-b0m7c1/services/proto-smoke-test-service-zx278tbo/tasks/aba7ff263e0546c997c4394407be5d5a/configuration?selectedContainer=proto-smoke-test

```
