# proto-go
protohackers



``` bash

echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker build --target runtime -t proto/smoke-test .      
docker run -p 8080:8080 proto/smoke-test

docker push ghcr.io/mobiljani/smoke-test 
```

## Smoke Test

https://portal.azure.com/#@ASOS1.onmicrosoft.com/resource/subscriptions/78619d9f-ef51-4241-8268-bcdeced3747c/resourcegroups/proto/providers/Microsoft.ContainerInstance/containerGroups/proto/overview

```
echo "hello hey" | nc 4.250.155.168 8080 
```

``` bash
az container restart --name proto --resource-group proto
```

## Prime time

https://portal.azure.com/#@ASOS1.onmicrosoft.com/resource/subscriptions/78619d9f-ef51-4241-8268-bcdeced3747c/resourceGroups/proto/providers/Microsoft.ContainerInstance/containerGroups/proto-prime/overview

```bash
echo '{"method":"isPrime","number":123}' | nc 20.26.97.157 8080
```
```bash
az container restart --name proto-prime --resource-group proto
```


## Means to an end

https://portal.azure.com/#@ASOS1.onmicrosoft.com/resource/subscriptions/78619d9f-ef51-4241-8268-bcdeced3747c/resourceGroups/proto/providers/Microsoft.ContainerInstance/containerGroups/means-end/overview

```bash
docker build --target runtime -t proto/means-end --build-arg SERVER_DIR=cmd/means-end/main.go .      
docker run -p 8080:8080 proto/means-end
```

```bash
echo '490000303900000065' | nc 85.210.37.5 8080
```

```bash
az container restart --name means-end --resource-group proto
```


Message: [121 0 0 60 0 0 0 100 0]
Message: [111 0 0 240 0 0 0 0 5]
Message: [111 0 0 60 73 0 0 0 144]
Message: [111 0 0 60 72 0 0 0 146]
Message: [111 0 0 60 71 0 0 0 145]
