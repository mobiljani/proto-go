# proto-go
protohackers



```

echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker build --target runtime -t proto/smoke-test .      
docker run -p 8080:8080 proto/smoke-test

docker push ghcr.io/mobiljani/smoke-test 



https://portal.azure.com/#@ASOS1.onmicrosoft.com/resource/subscriptions/78619d9f-ef51-4241-8268-bcdeced3747c/resourcegroups/proto/providers/Microsoft.ContainerInstance/containerGroups/proto/overview


az container restart --name proto --resource-group proto

```
