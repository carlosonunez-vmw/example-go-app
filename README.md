# Example Go App

This is a really simple Golang application that is destined to run within Tanzu
Application Platform (TAP).

It adds numbers, like this:

```sh
curl -X POST localhost:5000/add?by=10
# {"initial":0,"new":10}
```

If you have cookies enabled, it will increment the last number you added.

```sh
curl --cookie-jar /path/to/cookies \
  -b /path/to/cookies \
  -X POST localhost:5000/add?by=10
# {"initial":30,"new":40}
```

## How To Use It

### Docker

- Build the image: `docker build -t app .`
- Then create a container from it: `docker run --rm -p 5000:5000 app`
- Then add stuff! `curl -X PUT localhost:5000/add?by=10`

### Tanzu Application Platform

> ## ✅ Replace repository references if needed
>
> We'll assume that you're using this repo at
> https://github.com/carlosonunez-vmw/example-go-app. If not, replace
> any references to it with your own repo.

Build the image: `docker build -t app .`

Then push it to the image registry associated with your TAP cluster:

```sh
docker tag app my.registry.local:5000/app &&
docker push my.registry.local:5000/app
```

Finally, use the Tanzu CLI to create a workload from it:

```sh
tanzu apps workload apply example-go-app \
  --git-repo https://github.com/vmware-tanzu/carlosonunez-vmw/example-go-app \
  --git-tag main \
  --type server
```
