package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "flag"
    "strings"
    "net/http"
    "github.com/elazarl/goproxy"
    "github.com/HessianZ/daisy-proxy/resolver"
)

var (
    listen = flag.String("listen", "localhost:8080", "listen on address")
    ip = flag.String("ip", "", "outgoing address")
    iface = flag.String("if", "", "outgoing interface")
    dnsServers = flag.String("dns", "8.8.8.8 8.8.4.4", "dns servers")
    verbose = flag.Bool("verbose", false, "verbose output")
)

func main() {
    flag.Parse()

   flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "  -h   : show help usage")
    }

    if *ip == "" && *iface != "" {
        *ip = getIpFromInterface(*iface)
    } else {
        log.Fatal("IP address or Interface must be specified")
    }

    var servers []string
    if *dnsServers == "" {
        log.Fatal("DNS servers must be specified")
    }
    servers = strings.Split(*dnsServers, " ")

    r := &resolver.Resolver{Servers: servers, LocalAddr: *ip}

    proxy := goproxy.NewProxyHttpServer()
    proxy.Verbose = *verbose
    proxy.Tr.Dial = func (network, addr string) (c net.Conn, err error) {
        if network == "tcp" {
            var remoteAddr *net.TCPAddr
            var err error

            localAddr, err := net.ResolveTCPAddr(network, *ip + ":0");
            if err != nil {
                return nil, err;
            }
            chunks := strings.Split(addr, ":")
            addrIp, err := r.LookupAddr(chunks[0])

            if err != nil {
                remoteAddr, err = net.ResolveTCPAddr(network, addr);
            } else {
                remoteAddr, err = net.ResolveTCPAddr(network, addrIp + ":" + chunks[1]);
            }

            if err != nil {
                return nil, err;
            }
            return net.DialTCP(network, localAddr, remoteAddr);
        }

        return net.Dial(network, addr);
    }

    go http.ListenAndServe(*listen, proxy)

    log.Printf("DaisyProxy listen on %s outgoing from %s\n", *listen, *ip);

    select {}
}

func getIpFromInterface(name string) (ip string) {
    iface, err := net.InterfaceByName(name)

    if err != nil {
        log.Fatal(err)
    }

    addrs, err := iface.Addrs()

    if err != nil {
        log.Fatal(err)
    }

    return strings.Split(addrs[0].String(), "/")[0]
}
