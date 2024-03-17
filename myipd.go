/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-02-28 16:54:07
 * @LastEditTime: 2024-03-17 13:30:56
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
	"time"
)

type ExtrasEntry struct {
	TagStr string
	IpAddr string
}

const (
	// Basic info //
	VER          string = "0.0.2"
	VER_CODENAME string = "TinaSprout"
	RELEASE_DATE string = "2024-03-17"
	// Internet IP srcs //
	INTERNET_IP4_SRC string = "https://ipv4.icanhazip.com"
	INTERNET_IP6_SRC string = "https://ipv6.icanhazip.com"
	// Default values //
	DEFAULT_LISTEN      string = "0.0.0.0:2170"
	DEFAULT_TOKENS_FILE string = "tokens.conf"
	DEFAULT_EXTRAS_FILE string = ""
	DEFAULT_ENABLE_IP4  bool   = true
	DEFAULT_ENABLE_IP6  bool   = false
	DEFAULT_NO_TAGS     bool   = false
	DEFAULT_PRINT_HELLO bool   = true
	DEFAULT_PRINT_LOG   bool   = true
	DEFAULT_PRINT_TIME  bool   = true
)

var (
	// Setting entries related //
	Listen     string = DEFAULT_LISTEN
	TokensFile string = DEFAULT_TOKENS_FILE
	ExtrasFile string = DEFAULT_EXTRAS_FILE
	EnableIp4  bool   = DEFAULT_ENABLE_IP4
	EnableIp6  bool   = DEFAULT_ENABLE_IP6
	NoTags     bool   = DEFAULT_NO_TAGS
	PrintHello bool   = DEFAULT_PRINT_HELLO
	PrintLog   bool   = DEFAULT_PRINT_LOG
	PrintTime  bool   = DEFAULT_PRINT_TIME
	// Token related //
	Tokens []string      = []string{}
	Extras []ExtrasEntry = []ExtrasEntry{}
)

func print_log(logStr string) {
	if PrintLog {
		if PrintTime {
			shim := " "
			logStr = time.Now().String() + shim + logStr
		}
		println(logStr)
	}
}

func err_handle(err error) bool {
	if err != nil {
		print_log("[E] " + err.Error())
	}
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
		case "-e", "--extras-file":
			ExtrasFile = os.Args[i+1]
			i++
		case "-6", "--enable-ipv6":
			EnableIp6 = true
		case "--disable-ipv4":
			EnableIp4 = false
		case "-q", "--quiet":
			PrintHello = false
			PrintLog = false
		case "--no-time":
			PrintTime = false
		case "--no-log":
			PrintLog = false
		case "--no-hello":
			PrintHello = false
		case "--no-tags":
			NoTags = true
		}
	}
}

func load_tokens() {
	Tokens = []string{} // First, clear old tokens.
	file, err := os.Open(TokensFile)
	err_handle(err)
	tokenScanner := bufio.NewScanner(file)
	for tokenScanner.Scan() {
		token := tokenScanner.Text()
		if token[0] != '#' {
			Tokens = append(Tokens, token)
		}
	}
}

func load_extras() {
	if ExtrasFile == "" {
		return
	}
	Extras = []ExtrasEntry{} // First, clear old extras.
	file, err := os.Open(ExtrasFile)
	err_handle(err)
	extraScanner := bufio.NewScanner(file)
	for extraScanner.Scan() {
		splited := strings.Split(extraScanner.Text(), " ")
		if splited[0][0] != '#' {
			switch len(splited) {
			case 2:
				Extras = append(Extras, ExtrasEntry{TagStr: splited[0], IpAddr: splited[1]})
			case 1:
				Extras = append(Extras, ExtrasEntry{TagStr: "", IpAddr: splited[0]})
			default:
				print_log("[W] An entry in extras file will not be load: format error")
			}
		}
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

func get_req_ip(r *http.Request) string {
	xForwardFor := r.Header.Get("X-Forwarded-For")
	if xForwardFor != "" {
		return strings.Split(strings.Split(xForwardFor, ", ")[0], ":")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func http_handler(w http.ResponseWriter, r *http.Request) {
	// Get token and action //
	token := r.URL.Query().Get("token")
	action := r.URL.Query().Get("action")
	tags := r.URL.Query().Get("tags")
	// Secondary tags switch //
	tagSwitch := true
	// Verify the token //
	if len(token) == 0 {
		print_log("[W] " + get_req_ip(r) + " triggered a 401 Unauthorized error")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized"))
		return
	}
	if !slices.Contains(Tokens, token) {
		print_log("[W] " + get_req_ip(r) + " triggered a 403 Forbidden error")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden"))
		return
	}
	// Parse tags key //
	switch tags {
	case "", "true", "yes", "on", "1", "t":
		tagSwitch = true
	case "false", "no", "off", "0", "f":
		tagSwitch = false
	}
	// Parse action key //
	switch action {
	case "", "get-ip", "get-ip-addr":
		print_log("[I] " + get_req_ip(r) + " requested get-ip-addr")
		// Get IP addr(s) //
		lines := []string{}
		ip4Addr, ip6Addr := get_ip()
		// Internet IPv4 //
		if ip4Addr != "!DISABLED!" {
			tagStr := "[internet][ipv4]"
			shim := " "
			if NoTags || !tagSwitch {
				tagStr = ""
				shim = ""
			}
			lines = append(lines, tagStr+shim+ip4Addr)
		}
		// Internet IPv6 //
		if ip6Addr != "!DISABLED!" {
			tagStr := "[internet][ipv6]"
			shim := " "
			if NoTags || !tagSwitch {
				tagStr = ""
				shim = ""
			}
			lines = append(lines, tagStr+shim+ip6Addr)
		}
		// Extras //
		for i := 0; i < len(Extras); i++ {
			tagStr := Extras[i].TagStr
			extIpAddr := Extras[i].IpAddr
			shim := " "
			if NoTags || !tagSwitch {
				tagStr = ""
				shim = ""
			}
			lines = append(lines, tagStr+shim+extIpAddr)
		}
		// Make the response string //
		respStr := ""
		for i := 0; i < len(lines); i++ {
			respStr += lines[i] + "\n"
		}
		// Write header and response //
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respStr))
		return
	case "reload", "reload-all":
		// Reload all //
		print_log("[I] " + get_req_ip(r) + " requested reload-all")
		load_tokens()
		load_extras()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
		return
	case "reload-tokens":
		// Reload tokens //
		print_log("[I] " + get_req_ip(r) + " requested reload-tokens")
		load_tokens()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
		return
	case "reload-extras":
		// Reload extras //
		print_log("[I] " + get_req_ip(r) + " requested reload-extras")
		load_extras()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
		return
	default:
		// Unknow action //
		print_log("[W] " + get_req_ip(r) + " requested unknow action")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 BadRequest"))
		return
	}
}

func print_hello() {
	if PrintHello {
		println("MyIpD Version " + VER + " (" + VER_CODENAME + "), a libre software under AGPLv3")
		println("This software comes with absolutely NO warranties, use it at your own risk")
		println("Code and document was written by FunctionSir with peace and love, released at " + RELEASE_DATE)
	}
}

func main() {
	print_hello() // Print hello info.
	args_parser() // Parse os.Args.
	load_tokens() // Load tokens from tokens file.
	load_extras() // Load extras entries from extras file.
	// HTTP Service //
	print_log("[I] The HTTP server will listen on " + Listen)
	http.HandleFunc("/", http_handler)
	err := http.ListenAndServe(Listen, nil)
	err_handle(err)
}
