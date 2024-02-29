/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-02-28 16:54:07
 * @LastEditTime: 2024-02-29 15:46:43
 * @LastEditors: FunctionSir
 * @Description: MyIpD. Let the server publish its IP addr(s).
 * @FilePath: /MyIpD/myipd.go
 */

package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

const (
	// Basic info //
	VER          string = "0.0.1"
	VER_CODENAME string = "TinaSprout"
	// Internet IP srcs //
	INTERNET_IP4_SRC = "https://ipv4.icanhazip.com"
	INTERNET_IP6_SRC = "https://ipv6.icanhazip.com"
	// Default values //
	DEFAULT_LISTEN      = "0.0.0.0:2170"
	DEFAULT_TOKENS_FILE = "tokens.conf"
	DEFAULT_ENABLE_IP4  = true
	DEFAULT_ENABLE_IP6  = false
)

var (
	// Setting entries related //
	Listen     string = DEFAULT_LISTEN
	TokensFile string = DEFAULT_TOKENS_FILE
	EnableIp4  bool   = DEFAULT_ENABLE_IP4
	EnableIp6  bool   = DEFAULT_ENABLE_IP6
	// Token related //
	Tokens []string = []string{}
)

func err_handle(err error) bool {
	return err != nil
}

func args_parser() {
	// If nothing provided //
	if len(os.Args) == 1 {
		return
	}
	// Parse the args //
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-l", "--listen":
			Listen = os.Args[i+1]
			i++
		case "-t", "--tokens-file":
			TokensFile = os.Args[i+1]
			i++
		case "-6", "--enable-ipv6":
			EnableIp6 = true
		case "--disable-ipv4":
			EnableIp4 = false
		}
	}
}

func load_tokens() {
	Tokens = []string{} // First, clear old tokens.
	file, err := os.Open(TokensFile)
	err_handle(err)
	tokenScanner := bufio.NewScanner(file)
	for tokenScanner.Scan() {
		Tokens = append(Tokens, tokenScanner.Text())
	}
}

func get_ip() (ip4Addr string, ip6Addr string) {
	// Get IPv4 Addr //
	if EnableIp4 {
		ip4Resp, err := http.Get(INTERNET_IP4_SRC)
		if err_handle(err) {
			ip4Addr = "!ERROR!"
		} else {
			defer ip4Resp.Body.Close()
			tmp, err := io.ReadAll(ip4Resp.Body)
			if err_handle(err) {
				ip4Addr = "!ERROR!"
			} else {
				ip4Addr = strings.Trim(string(tmp), "\n")
			}
		}
	} else {
		ip4Addr = "!DISABLED!"
	}
	// Get IPv6 Addr //
	if EnableIp6 {
		ip6Resp, err := http.Get(INTERNET_IP6_SRC)
		if err_handle(err) {
			ip6Addr = "!ERROR!"
		} else {
			defer ip6Resp.Body.Close()
			tmp, err := io.ReadAll(ip6Resp.Body)
			if err_handle(err) {
				ip6Addr = "!ERROR!"
			} else {
				ip6Addr = strings.Trim(string(tmp), "\n")
			}
		}
	} else {
		ip6Addr = "!DISABLED!"
	}
	return
}

func http_handler(w http.ResponseWriter, r *http.Request) {
	// Get key token and action //
	token := r.URL.Query().Get("token")
	action := r.URL.Query().Get("action")
	// Verify the token //
	if len(token) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized"))
		return
	}
	if !slices.Contains(Tokens, token) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden"))
		return
	}
	// Parse key action //
	switch action {
	case "":
		// Get IP addr(s) //
		lines := []string{}
		ip4Addr, ip6Addr := get_ip()
		if ip4Addr != "!DISABLED!" {
			lines = append(lines, "[internet][ipv4] "+ip4Addr)
		}
		if ip6Addr != "!DISABLED!" {
			lines = append(lines, "[internet][ipv6] "+ip6Addr)
		}
		respStr := ""
		for i := 0; i < len(lines); i++ {
			respStr += lines[i] + "\n"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respStr))
		return
	case "reload-tokens":
		// Reload tokens //
		load_tokens()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
		return
	default:
		// Unknow action //
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 BadRequest"))
		return
	}
}

func main() {
	args_parser() // Parse os.Args.
	load_tokens() // Load tokens from tokens file.
	http.HandleFunc("/", http_handler)
	http.ListenAndServe(Listen, nil)
}
