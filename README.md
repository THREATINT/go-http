# github.com/THREATINT/go-http

## Package http

### ClientIP
Returns the client ip address from HTTP headers. It is aware of `True-Client-IP` and `CF-Connecting-IP` (both from CloudFlare).

### MimeType
Provides a mapping for know filename extension to mime type (e.g. .html -> text/html) based on the built in `[mime.TypeByExtension](https://golang.org/pkg/mime/#TypeByExtension)` and http://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types 
Part of the source code is generated by `GenMimeType.py`

### AcceptLanguage
Use to parse browser `Accept-Language` headers like
`fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5`
Please see e.g. [MDN (Mozilla Developer Network)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language) for general information on Accept-Language.

## License
Release under the MIT License. (see LICENSE)

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/0b2961363d3b4f1eb005a6e936c9534b)](https://www.codacy.com/gh/THREATINT/go-http/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=THREATINT/go-http&amp;utm_campaign=Badge_Grade)