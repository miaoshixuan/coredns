package dnssec

import (
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/pkg/cache"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
)

func TestCacheSet(t *testing.T) {
	fPriv, rmPriv, _ := test.TempFile(".", privKey)
	fPub, rmPub, _ := test.TempFile(".", pubKey)
	defer rmPriv()
	defer rmPub()

	dnskey, err := ParseKeyFile(fPub, fPriv)
	if err != nil {
		t.Fatalf("failed to parse key: %v\n", err)
	}

	c := cache.New(defaultCap)
	m := testMsg()
	state := request.Request{Req: m, Zone: "miek.nl."}
	k := hash(m.Answer) // calculate *before* we add the sig
	d := New([]string{"miek.nl."}, []*DNSKEY{dnskey}, nil, c)
	d.Sign(state, time.Now().UTC())

	_, ok := d.get(k)
	if !ok {
		t.Errorf("signature was not added to the cache")
	}
}
