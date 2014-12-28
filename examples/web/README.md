WEB Example
===
This example demostrate a practical usage of `delay/web` module.

Usage
---

First of all, start the example:
```
go run examples/web/main.go
```
By default will listen `127.0.0.1:3000`

Now you can

  * Register a message:  
```bash
curl -XPOST http://localhost:3000/emails/someHash --data "Hello"
```
The message will be triggered after two seconds.

* Update a message:  
```bash
curl -XPOST http://localhost:3000/emails/someHash --data "Hello" && \
curl -XPOST http://localhost:3000/emails/someHash --data "Bye"
```  
The message will be triggered after two seconds with payload `Bye` instead (a message with `Hello` payload will be never triggered).

* Cancel a message:  
```bash
curl -XPOST http://localhost:3000/emails/someHash --data "Hello" && \
curl -XDELETE http://localhost:3000/emails/someHash
```
No message will be triggered


