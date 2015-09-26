RightScale Go Boilerplate
=========================

- Master (public repo):
[![Build Status](https://travis-ci.org/rightscale/go-boilerplate.svg?branch=master)](https://travis-ci.org/rightscale/go-boilerplate)
[![Coverage](https://s3.amazonaws.com/rs-code-coverage/go-boilerplate/cc_badge_master.svg)](https://gocover.io/github.com/rightscale/go-boilerplate)
- Master (private repo):
[![Build Status](https://magnum.travis-ci.com/rightscale/uca.svg?branch=master&token=4Q13wQTY4zqXgU7Edw3B)](https://magnum.travis-ci.com/rightscale/uca)
[![Coverage](https://s3.amazonaws.com/rs-code-coverage/go-boilerplate/cc_badge_master.svg)](https://gocover.io/github.com/rightscale/go-boilerplate)

Downloads:
- https://binaries.rightscale.com/rsbin/go-boilerplate/master/go-boilerplate-linux-amd64.tgz
- https://binaries.rightscale.com/rsbin/go-boilerplate/master/go-boilerplate-darwin-amd64.tgz
- https://binaries.rightscale.com/rsbin/go-boilerplate/v0.1.0/go-boilerplate-linux-amd64.tgz
- https://binaries.rightscale.com/rsbin/go-boilerplate/v0.1.0/go-boilerplate-darwin-amd64.tgz

Getting Started
-----------------
 - Install Go 1.5
 - Ensure your GOPATH is set such that $PWD == $GOPATH/src/github.com/rightscale/go-boilerplate
 - Install dependencies with `make depend`
 - Run tests using `make test`
 - Try it out: `make && ./go-boilerplate`

Go-boilerplate features
-----------------------

The go-boilerplate is a simple web app that has a couple of handlers to index, get, put, delete
key-value pairs in a hash table. The features of the repo are:
 - Simple Makefile and .travis.yml for full lifecycle, from building, testing, code coverage,
   uploads of binaries to S3, badges in README
 - Simple web app with logging, error handling, form parsing, and other middleware

Exercising the boilerplate handlers
-----------------------------------
``` shell
$ curl -i http://localhost:8080/health-check
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Sat, 26 Sep 2015 05:50:13 GMT
Content-Length: 49

go-boilerplate dev - 2015-09-25 22:47:25 - master
$ curl http://localhost:8080/demo/settings
{}
$ curl -XPUT http://localhost:8080/demo/settings/a?value=b
$ curl http://localhost:8080/demo/settings
{"a":"b"}
$
```
And the corresponding log output:
``` shell
$ make && ./go-boilerplate
make: Nothing to be done for `default'.
[2015-09-25 22:49:52] INFO RightScale Go Boilerplate                version="go-boilerplate dev - 2015-09-25 22:47:25 - master" pid=21653
[2015-09-25 22:50:13] DBUG GET /health-check
[2015-09-25 22:50:13] INFO /health-check                            verb=GET id=h/XemzFgiFHC-000001 ip=127.0.0.1:57786 time=79.222µs status=200
[2015-09-25 22:50:19] DBUG GET /demo/settings
[2015-09-25 22:50:19] INFO /demo/settings                           verb=GET id=h/XemzFgiFHC-000002 ip=127.0.0.1:57788 time=202.347µs status=200
[2015-09-25 22:50:23] DBUG PUT /demo/settings/a                     value=b
[2015-09-25 22:50:23] INFO /demo/settings/a                         verb=PUT id=h/XemzFgiFHC-000003 ip=127.0.0.1:57790 time=106.995µs status=0
[2015-09-25 22:50:27] DBUG GET /demo/settings
[2015-09-25 22:50:27] INFO /demo/settings                           verb=GET id=h/XemzFgiFHC-000004 ip=127.0.0.1:57792 time=76.967µs status=200
```
