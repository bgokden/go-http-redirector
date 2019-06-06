# go-http-redirector

A basic http service to redirect requests to the defined url.
Requests can be redirected as temporary or permanent based on regex matching.

Environment variables can be used for configuration:

* REDIRECT_URL
* TEMPORARY_REDIRECT_REGEX
* PERMANENT_REDIRECT_REGEX

REDIRECT_URL is mandatory. eg.: https://blog.berkgokden.com

TEMPORARY_REDIRECT_REGEX is a regex string. Default value matches html, php, asp extensions.  

PERMANENT_REDIRECT_REGEX is a regex string. Default value matches jpg, png, bmp, json, js extensions.

Default redirection behaviour is temporary redirection.

Redirections done with 307, 308 codes to keep http method type.

## Build docker image:

```shell
docker build -t berkgokden/redirector .
```

## Run as docker image locally

```shell
docker run -it -p 9090:9090 -p 9091:9091 -e REDIRECT_URL="https://blog.berkgokden.com" berkgokden/redirector
```

Redirection service is running on port 9090

Heath/Readiness service is running on port 9091

These ports are not configurable since it is intended to be run inside docker.

Heath endpoint: /healthy  
Readiness endpoint /ready


## Check redirections with curl:

Permanent redirect:

```shell
$ curl -i localhost:9090/img.jpg
HTTP/1.1 308 Permanent Redirect
Content-Type: text/html; charset=utf-8
Location: https://blog.berkgokden.com/img.jpg
Date: Thu, 06 Jun 2019 23:31:17 GMT
Content-Length: 71

<a href="https://blog.berkgokden.com/img.jpg">Permanent Redirect</a>.
```

Temporary redirect:

```shell
$ curl -i localhost:9090/page.html
HTTP/1.1 307 Temporary Redirect
Content-Type: text/html; charset=utf-8
Location: https://blog.berkgokden.com/page.html
Date: Thu, 06 Jun 2019 23:32:20 GMT
Content-Length: 73

<a href="https://blog.berkgokden.com/page.html">Temporary Redirect</a>.
```
