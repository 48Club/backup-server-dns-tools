package main

import (
	"log"
	"net"

	"github.com/48Club/backup-server-dns-tools/types"
	"github.com/miekg/dns"
)

var (
	srv    *dns.Server
	master *types.RPC
)

func main() {
	srv = &dns.Server{Addr: ":53", Net: "udp"}
	srv.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		if dns.Fqdn(r.Question[0].Name) != dns.Fqdn(config.Server) {
			w.Close()
			return
		}

		m := new(dns.Msg)
		m.SetReply(r)
		for _, question := range r.Question {
			switch question.Qtype {
			case dns.TypeANY:
				fallthrough
			case dns.TypeA:
				ip := config.Master.IP
				if !master.Alive {
					ip = config.Backup
				}
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.ParseIP(ip),
				})
			default:
				m.SetRcode(r, dns.RcodeNameError)
			}
		}

		if err := w.WriteMsg(m); err != nil {
			log.Println("Failed to reply", err)
		}

	})
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
