# DirectTunnelSec
A learning project in which I will improve. This tunnel will be used with IPsec encryption to have a straightforward and direct tunnel. It uses FOB encryption and all addons could be enabled or disabled with command line.
- FOB encryption [ aes 16 24 32 bytes]
- tcp no delay
- buffer size
- ipv4 and ipv6
- good match for ipsec
- later i will work on udp

**usage**
  
 - Server
   
  ipv4 amd64 : ./server_arm64 -listen=800 -local=":5050" -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535

  ipv6 amd64 : ./server_arm64 -listen=800 -local="[::]:5050" -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535

 - Client
   
   ipv4 amd64 : ./client_amd64 -local 5050 -target KharejIPV4:800 -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535
   
   ipv6 amd64 : ./client_amd64 -local 5051 -target [KharejIPV6]:800 -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535

- you can disable encrypt or tcpnodelay by simply removing them.

**Generating AES KEY**

openssl rand -hex 16  or openssl rand -hex 24  or openssl rand -hex 32
