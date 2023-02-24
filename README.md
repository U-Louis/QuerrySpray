# QuerySpray
 Duplicate requests, get the fastest

 ## Recommendations
### Weighing the impact / benefit
<!-- TODO -->
### Not getting banned
<!-- TODO -->

## Development purposes
Using a volume to hold the app during development. The compilation is done when launching the docker run.
After a change you may just kill the server head & restart the container. 
### Docker Build
docker build -t queryspray-env-dev .

### Docker Run
<!-- The -d & tail will keep the container running -->
docker run -d --name queryspray-dev -p 8085:8085 -v $(pwd):/app queryspray-env-dev tail -f /dev/null

### Refresh (might be a few secs)
docker exec -it queryspray-dev sh ./build.dev.sh

### Refresh App
./builder.dev.sh

## Production Purposes
Using the app already built in the container

### Docker Build
<!-- TODO -->

### Docker Run
<!-- TODO -->
