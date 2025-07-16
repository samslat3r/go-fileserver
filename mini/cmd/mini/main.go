package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"mini/internal/config"
	"mini/internal/server"
	"mini/internal/tlsutil"
)

func main() {
	var (
		cfgPath string
		host	string
		port	int
		verbose	bool
		selfSigned	bool 
	)

	flag.StringVar(&cfgPath, "config", "", "Path to YAML config file")
	flag.StringVar(&host, "host", "", "Override host bind address")
	flag.IntVar(&port, "port", 0, "Override listen port")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&selfSigned, "selfsigned", false, "Force regenerate of self signed cert/key")
	flag.Parse()

	over := make(map[string]any)
	if host != "" { over["host"] = host }
	if port != 0 { over["port"] = port }
	if verbose { over["verbose"] = true }
	if selfSigned { over["tls.self_signed"] = true }

	cfg, err := config.Load(cfgPath, over)
	if err != nil { log.Fatal(err) }

	// generate certificate and key if they should exist but don't
	if cfg.TLS.Enabled && cfg.TLS.SelfSigned {
		if err := tlsutil.EnsureSelfSigned(cfg.TLS.Cert, cfg.TLS.Key, cfg.Host); err != nil {
			log.Fatal(err)
		}
	}

	h := server.NewHandlers(cfg)
	r := server.NewRouter(h)

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	srv := &http.Server{
		Addr:	addr,
		Handler:	r,
	}

	fmt.Printf("Serving on %s (TLS: %v)\n", addr, cfg.TLS.Enabled)
	if cfg.TLS.Enabled {
		log.Fatal(srv.ListenAndServeTLS(cfg.TLS.Cert, cfg.TLS.Key))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}