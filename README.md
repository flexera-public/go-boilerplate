RightScale Go Boilerplate
=========================




Exercising the boilerplate handlers:
``` shell
$ curl -i http://localhost:8080/health-check
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Sat, 26 Sep 2015 05:50:13 GMT
Content-Length: 49

go-boilerplate dev - 2015-09-25 22:47:25 - master
$ curl -i http://localhost:8080/demo/settings
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 26 Sep 2015 05:50:19 GMT
Content-Length: 2

{}
$ curl -i -XPUT http://localhost:8080/demo/settings/a?value=b
HTTP/1.1 200 OK
Date: Sat, 26 Sep 2015 05:50:23 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8

$ curl -i http://localhost:8080/demo/settings
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 26 Sep 2015 05:50:27 GMT
Content-Length: 9

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
