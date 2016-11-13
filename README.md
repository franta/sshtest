# sshtest
test ssh servers availability

configuration is hardcoded, edit for your needs before using
## usage
```
./ssh google.com localhost example.com
ERROR: localhost dial tcp [::1]:22: getsockopt: connection refused
OK: example.com
ERROR: google.com connection timeout
exit status 1
```