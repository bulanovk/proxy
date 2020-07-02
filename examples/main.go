// Copyright (c) 2012-today José Nieto, https://xiam.io
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"github.com/bulanovk/proxy/pkg/proxy"
	"log"
	"os"
	"sync"
)

const Version = "2.0.0"

const (
	defaultAddress = `:8080`
	defaultPort    = uint(1080)
	defaultTLSPort = uint(10443)
)

var (
	flagHelp        = flag.Bool("h", false, "Shows usage options.")
	flagAddress     = flag.String("addr", defaultAddress, "Bind address.")
	flagTLSCertFile = flag.String("ca-cert", "", "Path to root CA certificate.")
	flagTLSKeyFile  = flag.String("ca-key", "", "Path to root CA key.")
	flagDNS         = flag.String("dns", "", "Custom DNS server that bypasses the OS settings")
)

func main() {

	flag.Parse()

	if *flagHelp {
		fmt.Printf("Usage: hyperfox [options]\n\n")
		flag.PrintDefaults()
		return
	}

	fmt.Printf("Hyperfox v%s, by José Nieto\n\n", Version)

	// Opening database.
	var err error
	if err != nil {
		log.Fatal("Failed to setup database: ", err)
	}

	os.Setenv(proxy.EnvTLSCert, *flagTLSCertFile)
	os.Setenv(proxy.EnvTLSKey, *flagTLSKeyFile)

	// Creating proxy.
	p := proxy.NewProxy()

	if *flagDNS != "" {
		if err := p.SetCustomDNS(*flagDNS); err != nil {
			log.Fatalf("unable to set custom DNS server: %v", err)
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.StartSplit(*flagAddress); err != nil {
			log.Fatalf("Failed to bind to %s (TLS): %v", *flagAddress, err)
		}
	}()

	wg.Wait()
}
