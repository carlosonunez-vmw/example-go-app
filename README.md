# Example Go App

This is a really simple Golang application that is destined to run within Tanzu
Application Platform (TAP)

## How To Use It

### Docker

Build the image: `docker build -t app .`

Then create a container from it: `docker run --rm -e APP_USER="Your Name" app`

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
