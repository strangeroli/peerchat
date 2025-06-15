# Xelvra P2P Messenger - Troubleshooting Guide

This guide helps you diagnose and fix common issues with Xelvra P2P Messenger.

## Quick Diagnostics

Run the built-in diagnostic tool first:
```bash
peerchat-cli doctor
```

This will test:
- Internet connectivity
- DNS resolution
- Firewall and port accessibility
- P2P node startup
- NAT traversal capabilities

## Common Issues

### 1. Firewall Problems ⚠️

**Symptoms:**
- Cannot connect to peers
- "Connection refused" errors
- NAT traversal failures
- Ports showing as "blocked" in diagnostics

**Solutions:**

#### Linux (UFW)
```bash
# Allow P2P ports
sudo ufw allow 4001/tcp
sudo ufw allow 4001/udp
sudo ufw allow 9000/tcp
sudo ufw allow 9000/udp

# Check status
sudo ufw status
```

#### Linux (iptables)
```bash
# Allow incoming P2P connections
sudo iptables -A INPUT -p tcp --dport 4001 -j ACCEPT
sudo iptables -A INPUT -p udp --dport 4001 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 9000 -j ACCEPT
sudo iptables -A INPUT -p udp --dport 9000 -j ACCEPT

# Save rules (Ubuntu/Debian)
sudo iptables-save > /etc/iptables/rules.v4
```

#### Windows Firewall
```powershell
# Run as Administrator
New-NetFirewallRule -DisplayName "Xelvra P2P TCP" -Direction Inbound -Protocol TCP -LocalPort 4001,9000 -Action Allow
New-NetFirewallRule -DisplayName "Xelvra P2P UDP" -Direction Inbound -Protocol UDP -LocalPort 4001,9000 -Action Allow
```

#### macOS
```bash
# Add firewall rules
sudo pfctl -f /etc/pf.conf
# Or disable firewall temporarily for testing
sudo pfctl -d
```

### 2. NAT Traversal Issues

**Symptoms:**
- Can connect to local peers but not external ones
- STUN servers not accessible
- UPnP not available

**Solutions:**

#### Router Configuration
1. **Enable UPnP** in router settings
2. **Port Forwarding** (if UPnP fails):
   - Forward TCP ports 4001, 9000 to your device
   - Forward UDP ports 4001, 9000 to your device

#### Manual Port Forwarding Example
```
External Port: 4001 → Internal IP: 192.168.1.100:4001 (TCP)
External Port: 4001 → Internal IP: 192.168.1.100:4001 (UDP)
External Port: 9000 → Internal IP: 192.168.1.100:9000 (TCP)
External Port: 9000 → Internal IP: 192.168.1.100:9000 (UDP)
```

### 3. Network Connectivity

**Symptoms:**
- "Network unreachable" errors
- DNS resolution failures
- Cannot reach bootstrap nodes

**Solutions:**

#### Check Basic Connectivity
```bash
# Test internet connection
ping 8.8.8.8

# Test DNS resolution
nslookup google.com

# Test specific ports
telnet bootstrap.libp2p.io 443
```

#### DNS Issues
```bash
# Try different DNS servers
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
echo "nameserver 1.1.1.1" | sudo tee -a /etc/resolv.conf
```

### 4. Permission Issues

**Symptoms:**
- "Permission denied" errors
- Cannot create ~/.xelvra directory
- Database access failures

**Solutions:**

```bash
# Fix directory permissions
chmod 700 ~/.xelvra
chmod 600 ~/.xelvra/*

# Recreate config directory
rm -rf ~/.xelvra
peerchat-cli init
```

### 5. Performance Issues

**Symptoms:**
- High CPU usage
- High memory consumption
- Slow message delivery

**Solutions:**

#### Check Resource Usage
```bash
# Monitor process
top -p $(pgrep peerchat)

# Check memory usage
ps aux | grep peerchat

# Check network usage
netstat -i
```

#### Optimize Performance
```bash
# Reduce log level
export XELVRA_LOG_LEVEL=error

# Disable QUIC if causing issues
export XELVRA_DISABLE_QUIC=true

# Restart with optimizations
peerchat-cli stop
peerchat-cli start
```

## Advanced Diagnostics

### Log Analysis

```bash
# View real-time logs
tail -f ~/.xelvra/peerchat.log

# Search for errors
grep -i error ~/.xelvra/peerchat.log

# Search for connection issues
grep -i "connection\|connect\|dial" ~/.xelvra/peerchat.log
```

### Network Testing

```bash
# Test specific peer connection
peerchat-cli connect /ip4/PEER_IP/tcp/4001/p2p/PEER_ID

# Discover local peers
peerchat-cli discover

# Check node status
peerchat-cli status
```

### Database Issues

```bash
# Check database integrity
sqlite3 ~/.xelvra/userdata.db "PRAGMA integrity_check;"

# Backup and recreate database
cp ~/.xelvra/userdata.db ~/.xelvra/userdata.db.backup
rm ~/.xelvra/userdata.db
peerchat-cli init
```

## Environment Variables

Configure these environment variables for troubleshooting:

```bash
# Enable debug logging
export XELVRA_LOG_LEVEL=debug

# Use custom config directory
export XELVRA_CONFIG_DIR=/path/to/config

# Disable QUIC transport
export XELVRA_DISABLE_QUIC=true

# Force IPv4 only
export XELVRA_IPV4_ONLY=true
```

## Getting Help

### Collect Diagnostic Information

Before reporting issues, collect this information:

```bash
# System information
uname -a
peerchat-cli version

# Network diagnostics
peerchat-cli doctor > diagnostics.txt

# Recent logs
tail -100 ~/.xelvra/peerchat.log > recent_logs.txt

# Node status
peerchat-cli status > node_status.txt
```

### Report Issues

1. **GitHub Issues**: https://github.com/Xelvra/peerchat/issues
2. **Include**: Version, OS, logs, steps to reproduce
3. **Attach**: Diagnostic files (remove sensitive information)

## Emergency Recovery

### Complete Reset

```bash
# Stop all processes
peerchat-cli stop
pkill -f peerchat

# Backup important data
cp -r ~/.xelvra ~/.xelvra.backup

# Complete reset
rm -rf ~/.xelvra

# Reinitialize
peerchat-cli init
peerchat-cli start
```

### Restore from Backup

```bash
# Stop node
peerchat-cli stop

# Restore configuration
cp -r ~/.xelvra.backup ~/.xelvra

# Restart
peerchat-cli start
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Network error
- `3`: Configuration error
- `4`: Permission error
- `5`: Peer not found

## Performance Targets

- **Memory**: < 20MB idle, < 50MB active
- **CPU**: < 1% idle, < 5% active
- **Latency**: < 50ms direct, < 200ms relay
- **Energy**: < 20mW idle (mobile)

If your system exceeds these targets, check for:
- Memory leaks in logs
- High CPU processes
- Network congestion
- Inefficient routing
