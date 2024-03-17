<!--
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-02-28 16:30:44
 * @LastEditTime: 2024-03-17 13:28:22
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /MyIpD/README.md
-->

# MyIpD

Let the server publish its IP addr.

让服务器公布它的IP.

## Current version 当前版本

0.0.2 (TinaSprout)
It is in alpha stage now and the codename(TinaSprout) will be used until the beta stage or stable stage arrived.

该软件目前处于alpha阶段, 版本代号(TinaSprout)将会一直被使用直到beta或稳定阶段到来.

## Why? 为什么?

If you want to serve something like a website, and don't want to expose your server's ip address(to prevent some annoying DDOS attack, etc), you might use a CDN or Cloudflare Argo Tunnel. But if you also set an VPN server like a WireGuard server, you need your server IP. But the IP might change, so you must want to know what is you current ip addr. You can use another domain and sync your IP to the domain's DNS record(s) automatically, but it requires another domain, and its more complex. To use this, you don't need a new domain, there's no more complexities of manage another domain, and it is really easy to set up.

如果你想搭建一个网站或是什么的, 而且你不想暴露你的服务器IP(为了防止如DDOS之类的烦人的攻击), 你可能会用CDN或是Cloudflare Argo Tunnel. 但是如果你同时需要设置一个VPN服务器, 像是一个WireGuard服务器, 你需要你的服务器IP. 但是IP可能会改变, 所以你一定想要知道你现在的IP是什么. 你可以用另一个域名然后将你的IP自动地同步到它的DNS记录, 但是这需要另一个域名, 而且更复杂. 通过使用这个, 你不需要一个新域名, 不会带来需要管理另一个域名而带来的更多复杂性, 而且它十分容易设置.

## How? 怎么用?

1. Download the binary file from Releases or build it yourself.
2. Create a systemd service or something etc.
3. Create a tokens file and put your preferred tokens, one token each line.
4. (Enable and) Start the service, enjoy it.

P.S. You can use "uuidgen" command to generate a UUID as your token, if you are using a GNU/Linux system.

1. 从Releases下载二进制文件或者自己构建它.
2. 创建一个systemd服务或是其他类似的东西.
3. 创建tokens文件并将你喜欢的tokens放进去. 一行一个token.
4. (启用并)启动服务, 享用它.

P.S. 你可以使用"uuidgen"命令来创建一个UUID作为你的token, 如果你在使用GNU/Linux的话.

## Build 构建

To build, just run "go build -ldflags '-s -w' -o myipd". For Windows, you might want to add ".exe" after "myipd".

要构建, 执行"go build -ldflags '-s -w' -o myipd"即可. 如果是Windows操作系统, 你可能想要在"myipd"后面添加".exe".

## Args and default values 参数和默认值

1. -l, --listen, Addr and port to listen on, 0.0.0.0:2170
2. -t, --tokens-file, File which contains tokens, "tokens.conf"
3. -e, --extras-file, Extras file you want to use, "" (empty string)
4. -6, --enable-ipv6, Enable IPv6 support, false
5. --disable-ipv4, Disable IPv4 support, false
6. -q, --quiet, Disable hello and logs, false
7. --no-time, Do NOT print time, false
8. --no-log, Do NOT print logs, false
9. --no-hello, Do NOT print hello info, false
10. --no-tags, Do NOT send tags, false

P.S. No matter -6 is added to your args or not, ALL IPs in your "extras file" will be sent.  
P.S. --no-tags will also let it won't send tags in your "extras file"  
P.S. Related source code was shown in the FYI section below.

1. -l, --listen, 监听的地址和端口, 0.0.0.0:2170
2. -t, --tokens-file, 包含tokens的文件, "tokens.conf"
3. -e, --extras-file, 指定要使用的extras file, "" (空字符串)
4. -6, --enable-ipv6, 启用IPv6支持, false
5. --disable-ipv4, 禁用IPv4支持, false
6. -q, --quiet, 禁用hello和日志, false
7. --no-time, 不要输出和时间, false
8. --no-log, 不要输出日志, false
9. --no-hello, 不要输出hello信息, false
10. --no-tags, 不要发送tags, false

P.S. 无论是否添加了-6到您的参数, 所有的在您"extras file"里的IP将被发送.  
P.S. --no-tags将会同时禁止发送"extras file"中您写的tags.  
P.S. 相关源代码在下方的供参考部分展示.

## Extras file

Extras file is where you can add your own IP addr entries.

Extras file是您可以加入您自己的IP条目的地方

## FYI 供参考

```go
const (
 // Basic info //
 VER          string = "0.0.2"
 VER_CODENAME string = "TinaSprout"
 RELEASE_DATE string = "2024-03-16"
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
```

```go
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
```

```go
func http_handler(w http.ResponseWriter, r *http.Request) {
  //......//
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
  //......//
}
```
