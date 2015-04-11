Daisy Proxy
===========

A simple proxy based on goproxy.

一个基于goproxy的简单http代理，可以指定出口IP地址。

### Usage
Usage of ./daisy-proxy:
  -dns="8.8.8.8 8.8.4.4": dns servers
  -if="": out interface
  -ip="": out address
  -listen="localhost:8080": listen on this address and port
  -verbose=false: verbose output
  -h   : show help usage
  
### Example
./daisy-proxy -if ppp0
2015/04/11 10:15:40 DaisyProxy listen on localhost:8080 outgoing from 10.0.0.60

### Download

* [DaisyProxy For Mac](http://static.hessian.cn/wp-content/uploads/2014/07/daisy-proxy.gz)
* [DaisyProxy For Windows 64](http://static.hessian.cn/wp-content/uploads/2014/09/daisy-proxy.7z)
