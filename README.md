The aim of this project was to gain more practical knowledge on the inner workings of the Domain Name System.<br> I wanted to refresh the knowledge I have of the DNS protocol and understand some  new concepts relating to DNS.
The first step in this journey was reading the [RFC]( https://datatracker.ietf.org/doc/html/rfc7858). It is relatively easy to read and contains diagrams which can help one think about some of the various DNS constructs.
<br> I intended for this to be a full on DNS server similar to this but unfortunately, I had to stop at this point.
At this point, This project simply accepts a response packet that looks like this
```bash
00000000  a2 7e 81 80 00 01 00 01  00 00 00 00 06 67 6f 6f  |.~...........goo|
00000010  67 6c 65 03 63 6f 6d 00  00 01 00 01 c0 0c 00 01  |gle.com.........|
00000020  00 01 00 00 00 61 00 04  d8 3a df ce              |.....a...:..|
0000002c
```
& it parses it to look give the following JSON encoded response.
<img width="520" alt="raycast-untitled" src="https://user-images.githubusercontent.com/46195831/205512178-9d0577d4-2fbe-430d-93d6-806c1b6ea482.png">
<br>
While writting this I spent a lot of time perusing this awesome package. You should definietly take a look if youre doing any DNS related work in Go. 