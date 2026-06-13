/*
 * @Author: Vincent Yang
 * @Date: 2024-04-23 00:39:03
 * @LastEditors: Jason Lyu
 * @LastEditTime: 2025-04-08 13:45:00
 * @FilePath: /DeepLX/config.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */

package service

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	IP        string
	Port      int
	Token     string
	DlSession string
	Proxy     string
}

func InitConfig() *Config {
	cfg := &Config{
		IP:   "0.0.0.0",
		Port: 1188,
	}

	// IP flag
	if ip, ok := os.LookupEnv("IP"); ok && ip != "" {
		cfg.IP = ip
	}
	flag.StringVar(&cfg.IP, "ip", cfg.IP, "set up the IP address to bind to")
	flag.StringVar(&cfg.IP, "i", cfg.IP, "set up the IP address to bind to")

	// Port flag
	if port, ok := os.LookupEnv("PORT"); ok && port != "" {
		if _, err := fmt.Sscanf(port, "%d", &cfg.Port); err != nil {
			log.Printf("Warning: failed to parse PORT env var %q: %v. Using default %d", port, err, cfg.Port)
		}
	}
	flag.IntVar(&cfg.Port, "port", cfg.Port, "set up the port to listen on")
	flag.IntVar(&cfg.Port, "p", cfg.Port, "set up the port to listen on")

	// DL Session flag
	dlSession := ""
	if envDlSession, ok := os.LookupEnv("DL_SESSION"); ok {
		dlSession = envDlSession
	}
	flag.StringVar(&cfg.DlSession, "s", dlSession, "set the dl-session for /v1/translate endpoint")

	// Access token flag
	token := ""
	if envToken, ok := os.LookupEnv("TOKEN"); ok {
		token = envToken
	}
	flag.StringVar(&cfg.Token, "token", token, "set the access token for /translate endpoint")

	// HTTP Proxy flag
	proxy := ""
	if envProxy, ok := os.LookupEnv("PROXY"); ok {
		proxy = envProxy
	}
	flag.StringVar(&cfg.Proxy, "proxy", proxy, "set the proxy URL for HTTP requests")

	flag.Parse()

	if cfg.Port < 1 || cfg.Port > 65535 {
		log.Fatalf("Invalid port number: %d. Port must be between 1 and 65535.", cfg.Port)
	}

	return cfg
}
