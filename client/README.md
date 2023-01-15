# Jarvis Client

## Todo
- [ ] Write as Dockerfile service that runs actions on cron or input (run on proxmox host)
- [ ] Send heartbeats from client -> server
    - IP address
    - Service expects a heartbeat every X seconds
- [ ] Syslog parsing from hosts (pfsense, proxmox)
- [ ] Configuration file shared with `jarvis`
- [ ] Telegram Listener 
   - [ ] Ask bot to do common debug/admin tasks:
       - [ ] Renew DHCP Leases
       - [ ] Restart Talescale/Wireguard/VPN
       - [ ] WoL
