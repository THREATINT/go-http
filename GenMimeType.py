#!/usr/bin/python3
from collections import OrderedDict

import httplib2
import re

resp, content = httplib2.Http().request("http://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types")

mappings = OrderedDict()

pattern = re.compile("[^\t]+")
for l in str(content, 'utf-8').splitlines():
    if not l.startswith('#'):
        m = pattern.findall(l)
        for i in m[1].split():
            mappings[i] = m[0]

with open('mimeType_gen.go', 'w', encoding='utf-8') as f:
    f.write("package http\n\n")

    f.write("// provides a mapping for know filename extension to mime type based on\n")
    f.write("// http://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types\n\n")
    f.write("// WARNING:    this is generated code, do not edit!\n")
    f.write("//             Please run GenMimeType.py to update this file\n\n\n")

    f.write("var mimeType = map[string]string {\n")
    for extension, mimetype in mappings.items():
        f.write("\t\"")
        f.write(extension)
        f.write("\": \"")
        f.write(mimetype)
        f.write("\",\n")

    f.write("}")
