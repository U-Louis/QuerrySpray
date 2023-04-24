# QuerySpray
 Duplicate requests, get the fastest.  
 QuerySpray is a light app containing a gin server. You send it a request, it sends it multiple times and returns you only the fastest one as soon as it gets it.
 This project was created by non go developers. Help and improvement are welcome !
  

 ## Recommendations
### Weighing the impact / benefit
Please remember that the perpetual growth of the internet traffic is an energy/climate issue. Use this project responsibly.  
This project was develop to enhance an API providing a response in around 150ms, but with a 5000ms every 20 requests. On the tests performed, the best gain was obtained with a multiple of 4. This adds an average 25-30ms by request but in our case, the average response was around 500ms if we included the pikes. So the service is globally more efficient and very much more stable.  
Be sure to perform tests to chose which multiple is most efficient for your project. Postman runners are a great way to do this.
  
### Not getting banned
Be sure to check if there are any limitations on the API you are consuming.
  
## Development purposes
Using a volume to hold the app during development. The compilation is done when launching the docker run.  

### Docker Build
```
docker build -t queryspray:dev .
```
  
### Docker Run
<!-- The -d & tail will keep the container running -->
```
docker run -d -p 8085:8085 -v $(pwd):/app queryspray:dev tail -f /dev/null
```
```
docker run --network host -d -p 8085:8085 -v $(pwd):/app queryspray:dev tail -f /dev/null
```
  

### Refresh (might be a few secs)
```docker exec -it queryspray-dev sh ./build.dev.sh; docker container restart queryspray-dev```
  

### Test with curl (and with throttle-responder here)
```
curl --location --request POST 'http://localhost:8085/spray?multiple=2' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "POST",
    "uri": "http://172.17.0.2:5000/throttle",
    "protocol": "HTTP/1.1",
    "headers": [
        "Content-Type: application/json"
    ],
    "body": "{\n    \"throttle\": 1000,\n    \"id\": \"sample_id\"\n}"
}'
```
  

## Production Purposes
We will use the go app already built in the container. It will reduce the image size to around 15mo.  
Note that the binary `main` have to be compiled differently for different architectures (the default one here is a mac arm64).  
To do so :  
- see build.dev.sh and uncomment/put the architecture you want  
- run the dev image (it will compile the binary `main`)  
- copy the `main` into the dist folder  
- then run the dist image with the according --platform.  


### Docker Build
```
docker build -t queryspray:dist-arm64-v1 -f dockerfile.arm64.dist .
```  
for linux amd64 :  
```
cd dist amd64
docker build --platform linux/amd64 -t queryspray:dist-amd64-v1 -f dockerfile.amd64.dist .
```

### Docker Run
For mac arm64 :
```
docker run -d --name queryspray -p 8085:8085 queryspray:dist-arm64-v1
```
for linux amd64 :  
```
docker run -d --name queryspray -p 8085:8085 queryspray:dist-amd64-v1
```


### Test with curl (and with throttle-responder here)
```
curl --location --request POST 'http://localhost:8085/spray?multiple=2' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "POST",
    "uri": "http://172.17.0.2:5000/throttle",
    "protocol": "HTTP/1.1",
    "headers": [
        "Content-Type: application/json"
    ],
    "body": "{\n    \"throttle\": 1000,\n    \"id\": \"sample_id\"\n}"
}'
```