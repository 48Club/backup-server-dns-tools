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
		m := new(dns.Msg)
		m.SetReply(r)
		m.RecursionAvailable = true
		for _, question := range r.Question {
			if dns.Fqdn(question.Name) != dns.Fqdn(config.Server) {
				// 拒绝响应预期外的请求
				m.SetRcode(r, dns.RcodeNameError)
				continue
			}
			queryAny := false
			switch question.Qtype {
			case dns.TypeANY:
				queryAny = true
				fallthrough
			case dns.TypeA:
				buildAResp(m, question)
				if !queryAny {
					break
				}
				fallthrough
			case dns.TypeNS:
				buildNSResp(m, question)
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

func buildAResp(m *dns.Msg, q dns.Question) {
	ip := config.Master.IP
	if !master.Alive {
		ip = config.Backup
	}
	m.Answer = append(m.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.ParseIP(ip),
	})
}

func buildNSResp(m *dns.Msg, q dns.Question) {
	for _, ns := range config.RecursiveNS {
		m.Answer = append(m.Answer, &dns.NS{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60},
			Ns:  dns.Fqdn(ns),
		})
	}
}
