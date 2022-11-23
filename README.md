# go-http

## Introduction

### ClientIP
Returns the client ip address from HTTP headers. It is aware of `True-Client-IP` and `CF-Connecting-IP`. (both from CloudFlare)

### MimeType
Provides a mapping for know filename extension to mime type (e.g. .html -> text/html) based on the built in [`mime.TypeByExtension`](https://golang.org/pkg/mime/#TypeByExtension) and [svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types](https://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types)

Part of the source code is generated by `GenMimeType.py`

### AcceptLanguage
Use to parse browser `Accept-Language` headers like
`fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5`
Please see e.g. [MDN (Mozilla Developer Network)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language) for general information on Accept-Language.

## License
Release under the MIT License. (see LICENSE)
