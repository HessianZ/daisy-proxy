package resolver

import (
	"github.com/miekg/dns"
    "fmt"
    "net"
    "log"
)

var (
    GoogleServers = []string{"8.8.8.8", "8.8.4.4"}
    dnsCaches = make(map[string]string)
)

type Resolver struct {
    Servers []string
    LocalAddr string
}


func (r *Resolver) LookupAddr(hostname string) (string, error) {
    if _, ok := dnsCaches[hostname]; ok {
        return dnsCaches[hostname], nil
    }

    msg, err := r.Lookup(dns.TypeA, hostname)
    if err != nil {
        return "", err
    }

    for _, ans := range msg.Answer {
        a, ok := ans.(*dns.A)
        if ok {
            dnsCaches[hostname] = a.A.String()
            return a.A.String(), nil
        }

    }

    return "", err
}

func (r *Resolver) Lookup(qType uint16, name string) (*dns.Msg, error) {
	name = dns.Fqdn(name)
	client := &dns.Client{}
	msg := &dns.Msg{}
	msg.SetQuestion(name, qType)

    var err error
	response := &dns.Msg{}
	for _, server := range r.Servers {
		response, err = r.lookup(msg, client, server + ":53", false)
		if err == nil {
			return response, nil
		}
    }
	return response, fmt.Errorf("Couldn't resolve %s: No server responded", name)
}

func (r *Resolver) lookup(msg *dns.Msg, client *dns.Client, server string, edns bool) (*dns.Msg, error) {
	if edns {
		opt := &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
		}
		opt.SetUDPSize(dns.DefaultMsgSize)
		msg.Extra = append(msg.Extra, opt)
	}

    var response *dns.Msg

    var err error
    if r.LocalAddr == "" {
        response, _, err = client.Exchange(msg, server)
    } else {
        localAddr, err := net.ResolveTCPAddr("tcp", r.LocalAddr + ":0");
        if err != nil {
            return nil, err;
        }
        remoteAddr, err := net.ResolveTCPAddr("tcp", server);
        if err != nil {
            return nil, err;
        }
        c, err := net.DialTCP("tcp", localAddr, remoteAddr);
        if err != nil {
            return nil, err
        }
        co := &dns.Conn{Conn: c}

        err = co.WriteMsg(msg)
        if err != nil {
            return nil, err
        }
        response, err  = co.ReadMsg()
    }

    if response == nil {
        log.Fatal("EEEEEEEEEEEEEE  ", client, server, response, err)
    }

    if err != nil {
        return nil, err
    }


	if msg.Id != response.Id {
		return nil, fmt.Errorf("DNS ID mismatch, request: %d, response: %d", msg.Id, response.Id)
	}

	if response.MsgHdr.Truncated {
		if client.Net == "tcp" {
			return nil, fmt.Errorf("Got truncated message on tcp")
		}

		if edns { // Truncated even though EDNS is used
			client.Net = "tcp"
		}

		return r.lookup(msg, client, server, !edns)
	}

	return response, nil
}



