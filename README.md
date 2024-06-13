# DirectTunnelSec
A learning project in which I will improve. This tunnel will be used with IPsec encryption to have a straightforward and direct tunnel. It uses FOB encryption and all addons could be enabled or disabled with command line
- FOB encryption [ aes 16 24 32 bytes]
- tcp no delay
- buffer size
- ipv4 and ipv6
- toml config file
- good match for ipsec [https://github.com/Azumi67/6TO4-GRE-IPIP-SIT]
- later i will work on udp

**install**

```
apt update -y
apt install wget -y
apt install unzip -y
wget https://github.com/Azumi67/DirectTunnelSec/releases/download/v1.002/amd64.zip
unzip amd64.zip -d /root/dtunnel
```

**toml file**

- There are server and client configurations examples in explie directory
- look for examples there

**usage**
  
 - Server
   
  ipv4/ipv6 amd64 : ./server_amd64 -config kharej_config.toml

 - Client
   
   ipv4/ipv6 amd64 : ./client_amd64 -config iran_config.toml
   

**Generating AES KEY**

openssl rand -hex 16  or openssl rand -hex 24  or openssl rand -hex 32
