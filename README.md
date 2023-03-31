# QuerySpray
 Duplicate requests, get the fastest.    
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
After a change you may just kill the server head & restart the container. 

### Docker Build
```docker build -t queryspray-env-dev .```
  
### Docker Run
<!-- The -d & tail will keep the container running -->
```docker run -d --name queryspray-dev -p 8085:8085 -v $(pwd):/app queryspray-env-dev tail -f /dev/null```
```docker run --network host -d --name queryspray-dev -p 8085:8085 -v $(pwd):/app queryspray-env-dev tail -f /dev/null```
  

### Refresh (might be a few secs)
```docker exec -it queryspray-dev sh ./build.dev.sh; docker container restart queryspray-dev```
  

### Test with curl and throttle-responder (magic uri, might change)
```curl --location --request POST 'http://localhost:8085/spray?multiple=2' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "POST",
    "uri": "http://172.17.0.2:5000/throttle",
    "protocol": "HTTP/1.1",
    "headers": [
        "Content-Type: application/json"
    ],
    "body": "{\n    \"throttle\": 1000,\n    \"id\": \"sample_id\"\n}"
}'```
  

## Production Purposes
Using the app already built in the container


### Docker Build
```docker build -t queryspray-dist -f dockerfile.dist .```
Or for amd64 :
```docker build --platform linux/amd64 -t queryspray-dist-amd64 -f dockerfile.dist .```

### Docker Run
```docker run -d --name queryspray-dist -p 8086:8085 -v $(pwd):/app queryspray-dist tail -f /dev/null```
  

### Test with curl and throttle-responder (magic uri, might change)
```curl --location --request POST 'http://localhost:8086/spray?multiple=2' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "POST",
    "uri": "http://172.17.0.2:5000/throttle",
    "protocol": "HTTP/1.1",
    "headers": [
        "Content-Type: application/json"
    ],
    "body": "{\n    \"throttle\": 1000,\n    \"id\": \"sample_id\"\n}"
}'```