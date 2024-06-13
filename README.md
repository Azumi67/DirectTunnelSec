# DirectTunnelSec
A learning project in which I will improve. This tunnel will be used with IPsec encryption to have a straightforward and direct tunnel. It uses FOB encryption and all addons could be enabled or disabled with command line
- FOB encryption [ aes 16 24 32 bytes]
- tcp no delay
- buffer size
- ipv4 and ipv6
- toml config file - having some problems with this. i will add it later
- good match for ipsec [https://github.com/Azumi67/6TO4-GRE-IPIP-SIT]
- later i will work on udp

**install**

```
apt update -y
apt install wget -y
apt install unzip -y
wget https://github.com/Azumi67/DirectTunnelSec/releases/download/v1.001/server_amd64
wget https://github.com/Azumi67/DirectTunnelSec/releases/download/v1.001/client_amd64
chmod +x server_amd64
chmod +x client_amd64

```

**toml file**

- it will be added later

**usage**
  
- Server
   
  ipv4 amd64 : ./server_arm64 -listen=800 -local=":5050" -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535
  
  ipv6 amd64 : ./server_arm64 -listen=800 -local="[::]:5050" -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535
  
 - Client
   
   ipv4 amd64 : ./client_amd64 -local 5050 -target KharejIPV4:800 -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535
   
   ipv6 amd64 : ./client_amd64 -local 5051 -target [KharejIPV6]:800 -noDelay=true -encrypt -key ATMZE1uD7dmgNDnERJLSFw== -buffer 65535

- no ecryption : just remove -encrypt & - key

**Generating AES KEY**
  
openssl rand -hex 16  or openssl rand -hex 24  or openssl rand -hex 32

- too many open files [ulimit]

```
ulimit -u  #displaying the maximum user process
ulimit -f  #displaying the maximum file size a user can have
ulimit -n
ulimit -n <new input> for example : ulimit -n 65536
nano /etc/security/limits.conf
*         hard    nofile      <your input>
*         soft    nofile      <your input>

```
