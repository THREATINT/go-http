# go-http

## Introduction
This is a library implemented in Go (Golang) that provides some functionality to deal with HTTP requests and headers.

## Usage
Install in your ${GOPATH} using `go get -u github.com/THREATINT/go-http`.

### AcceptLanguage
Use `ParseAcceptLanguage()` to parse browser `Accept-Language` headers like

`fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5`

Please see e.g. [MDN (Mozilla Developer Network)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language) for general information on Accept-Language.

### ClientIP
`ClientIP()` returns the client ip address from HTTP headers. It works with `X-Forwarded-For` (used by most reverse proxies) and is aware of `True-Client-IP` and `CF-Connecting-IP` (both implemented by CloudFlare).

### MimeType
`MimeTypeByExtension()` provides a mapping for know file extension to mime type (e.g. .html -> text/html) based on the builtin [`mime.TypeByExtension`](https://golang.org/pkg/mime/#TypeByExtension) and [svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types](https://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types).

Part of the source code for this functionality is generated by `GenMimeType.py`. (also provided in this repository)

## Feedback
We would love to hear from you! Please contact us at [help@threatint.com](mailto:help@threatint.com) for feedback and general requests.

Kindly raise an issue in Github if you find a problem in the code.

[Subscribe to our mailing list](https://newsletter.threatint.com/subscription?f=RMr892gtEllhgouxTbPi3QWEepgPIlRgZ763h43B3mbR8gjYGzoBTVvCg88929UbPuxHQwJl09B763PRvyG7633n4hrFx3892A) to learn more about our work.

## License
Release under the MIT License. (see LICENSE)
