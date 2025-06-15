# Troubleshooting Guide

This comprehensive guide helps you diagnose and resolve common issues with Xelvra P2P Messenger.

## 🔍 Quick Diagnostics

### Run Built-in Diagnostics
```bash
peerchat-cli doctor
```
This command performs comprehensive system checks and provides specific recommendations for any issues found.

### Check System Status
```bash
# View your identity and network info
peerchat-cli id

# Check current connections
peerchat-cli status

# Monitor real-time activity
peerchat-cli listen
```

### View Logs
```bash
# View recent logs
tail -f ~/.xelvra/peerchat.log

# Search for errors
grep -i error ~/.xelvra/peerchat.log

# View all logs
cat ~/.xelvra/peerchat.log
```

## 🌐 Network & Connection Issues

### Issue: No Peers Found During Discovery

**Symptoms:**
- `/discover` command finds no peers
- "No peers discovered" message
- Empty peer list

**Diagnosis:**
```bash
# Check network connectivity
peerchat-cli doctor

# Test discovery manually
peerchat-cli discover --timeout 30

# Monitor discovery process
peerchat-cli listen &
peerchat-cli discover
```

**Solutions:**

1. **Check Firewall Settings**
```bash
# Linux (UFW)
sudo ufw status
sudo ufw allow 42424/udp

# Linux (iptables)
sudo iptables -L | grep 42424
sudo iptables -A INPUT -p udp --dport 42424 -j ACCEPT

# Windows
netsh advfirewall firewall show rule name="Xelvra"
# If no rule exists, create one:
netsh advfirewall firewall add rule name="Xelvra Discovery" dir=in action=allow protocol=UDP localport=42424
```

2. **Network Configuration**
```bash
# Check if you're on the same network as other users
ip route show default  # Linux
route print 0.0.0.0    # Windows
netstat -rn | grep default  # macOS

# Test with different network (mobile hotspot)
```

3. **Port Conflicts**
```bash
# Check if port 42424 is in use
netstat -tulpn | grep 42424  # Linux
netstat -an | findstr 42424  # Windows
lsof -i :42424               # macOS
```

### Issue: Connection Failures

**Symptoms:**
- "Connection failed" errors
- Peers discovered but can't connect
- Timeouts during connection attempts

**Diagnosis:**
```bash
# Test specific peer connection
peerchat-cli connect <peer_id> --verbose

# Check NAT/firewall status
peerchat-cli doctor --network-detailed
```

**Solutions:**

1. **NAT Traversal Issues**
```bash
# Check your public IP
curl ifconfig.me

# Test STUN connectivity
peerchat-cli doctor --stun-test
```

2. **Peer Availability**
```bash
# Verify peer is still online
peerchat-cli discover | grep <peer_id>

# Try connecting to different peer
peerchat-cli discover
peerchat-cli connect <different_peer_id>
```

3. **Version Compatibility**
```bash
# Check your version
peerchat-cli version

# Ensure all users have compatible versions
```

### Issue: Simulation Mode Activated

**Symptoms:**
- "Simulation mode detected" warning
- No real P2P connections possible
- Limited functionality

**Diagnosis:**
```bash
# Check why simulation mode is active
peerchat-cli doctor --simulation-check

# Test network interfaces
ip addr show  # Linux
ipconfig      # Windows
ifconfig      # macOS
```

**Solutions:**

1. **Network Interface Issues**
```bash
# Try different network interface
peerchat-cli start --interface eth0
peerchat-cli start --interface wlan0

# Check available interfaces
peerchat-cli doctor --list-interfaces
```

2. **Connectivity Problems**
```bash
# Test internet connectivity
ping 8.8.8.8
curl -I https://google.com

# Test local network
ping <gateway_ip>
```

## 🔐 Security & Identity Issues

### Issue: Lost Identity/Keys

**Symptoms:**
- "Identity not found" error
- Can't access previous conversations
- Need to reinitialize

**Solutions:**

1. **Restore from Backup**
```bash
# If you have a backup of ~/.xelvra/
cp -r /path/to/backup/.xelvra ~/

# Verify restoration
peerchat-cli id
```

2. **Create New Identity**
```bash
# Backup current (corrupted) identity
mv ~/.xelvra ~/.xelvra.backup

# Create new identity
peerchat-cli init

# Note: Previous message history will be lost
```

### Issue: Key Rotation Problems

**Symptoms:**
- "Key rotation failed" errors
- Messages not decrypting properly
- Authentication failures

**Solutions:**

1. **Manual Key Rotation**
```bash
# Force key rotation
peerchat-cli rotate-keys --force

# Verify new keys
peerchat-cli id --show-keys
```

2. **Reset Encryption State**
```bash
# Clear encryption cache
rm ~/.xelvra/encryption_cache.db

# Restart with fresh encryption state
peerchat-cli start
```

## 📱 Performance Issues

### Issue: High Memory Usage

**Symptoms:**
- System becomes slow
- Out of memory errors
- High RAM consumption

**Diagnosis:**
```bash
# Monitor memory usage
top -p $(pgrep peerchat-cli)  # Linux
tasklist | findstr peerchat  # Windows

# Check memory leaks
peerchat-cli status --memory-detailed
```

**Solutions:**

1. **Optimize Configuration**
```yaml
# Edit ~/.xelvra/config.yaml
performance:
  max_peers: 10           # Reduce from default 50
  connection_timeout: 15s # Reduce from 30s
  max_concurrent_streams: 50  # Reduce from 100

logging:
  level: "warn"          # Reduce logging
  max_size: 5            # Smaller log files
```

2. **Restart Periodically**
```bash
# Add to crontab for automatic restart
0 */6 * * * pkill peerchat-cli && sleep 5 && peerchat-cli start --daemon
```

### Issue: High CPU Usage

**Symptoms:**
- System becomes unresponsive
- High CPU utilization
- Fan noise/heat

**Diagnosis:**
```bash
# Monitor CPU usage
htop  # Linux
top   # macOS
# Task Manager on Windows

# Check for busy loops
perf top -p $(pgrep peerchat-cli)  # Linux
```

**Solutions:**

1. **Reduce Discovery Frequency**
```yaml
# Edit ~/.xelvra/config.yaml
network:
  discovery_interval: 60s  # Increase from default 30s
  heartbeat_interval: 30s  # Increase from default 15s
```

2. **Limit Concurrent Operations**
```yaml
performance:
  max_concurrent_streams: 25
  connection_pool_size: 10
```

## 🔧 Configuration Issues

### Issue: Configuration File Problems

**Symptoms:**
- "Invalid configuration" errors
- Settings not taking effect
- Startup failures

**Solutions:**

1. **Validate Configuration**
```bash
# Check configuration syntax
peerchat-cli config --validate

# Show current configuration
peerchat-cli config --show
```

2. **Reset to Defaults**
```bash
# Backup current config
cp ~/.xelvra/config.yaml ~/.xelvra/config.yaml.backup

# Generate default config
peerchat-cli init --reset-config

# Manually edit if needed
nano ~/.xelvra/config.yaml
```

### Issue: Log File Problems

**Symptoms:**
- Logs not being written
- Log files too large
- Permission errors

**Solutions:**

1. **Fix Log Permissions**
```bash
# Check permissions
ls -la ~/.xelvra/peerchat.log

# Fix permissions
chmod 644 ~/.xelvra/peerchat.log
chown $USER:$USER ~/.xelvra/peerchat.log
```

2. **Configure Log Rotation**
```yaml
# Edit ~/.xelvra/config.yaml
logging:
  max_size: 10      # MB
  max_backups: 5    # Keep 5 old files
  max_age: 30       # Days
  compress: true    # Compress old logs
```

## 🐛 Application Crashes

### Issue: Segmentation Faults

**Symptoms:**
- Application crashes with "segfault"
- Core dumps generated
- Sudden termination

**Diagnosis:**
```bash
# Enable core dumps
ulimit -c unlimited

# Run with debugging
peerchat-cli start --debug

# Check system logs
dmesg | grep peerchat  # Linux
# Check Event Viewer on Windows
```

**Solutions:**

1. **Update Dependencies**
```bash
# Rebuild with latest dependencies
cd /path/to/source
go mod tidy
go mod download
go build -o bin/peerchat-cli cmd/peerchat-cli/main.go
```

2. **Check System Resources**
```bash
# Check available memory
free -h  # Linux
vm_stat  # macOS

# Check disk space
df -h
```

### Issue: Panic Errors

**Symptoms:**
- "panic:" messages in logs
- Stack traces
- Unexpected termination

**Solutions:**

1. **Report the Bug**
```bash
# Collect crash information
peerchat-cli version > crash_report.txt
cat ~/.xelvra/peerchat.log >> crash_report.txt
uname -a >> crash_report.txt

# Submit to GitHub Issues with crash_report.txt
```

2. **Temporary Workaround**
```bash
# Run in recovery mode
peerchat-cli start --safe-mode

# Or with minimal features
peerchat-cli start --no-discovery --no-relay
```

## 🔄 Recovery Procedures

### Complete Reset

If all else fails, perform a complete reset:

```bash
# 1. Backup important data
cp -r ~/.xelvra ~/.xelvra.backup.$(date +%Y%m%d)

# 2. Stop all instances
pkill peerchat-cli

# 3. Remove configuration
rm -rf ~/.xelvra

# 4. Reinstall/rebuild
# (Follow installation guide)

# 5. Initialize fresh
peerchat-cli init

# 6. Test basic functionality
peerchat-cli doctor
```

### Selective Reset

Reset specific components:

```bash
# Reset only network configuration
rm ~/.xelvra/network_cache.db

# Reset only logs
rm ~/.xelvra/peerchat.log*

# Reset only peer database
rm ~/.xelvra/peers.db
```

## 📞 Getting Help

### Before Asking for Help

1. **Run diagnostics:**
```bash
peerchat-cli doctor > diagnostics.txt
```

2. **Collect system information:**
```bash
# Linux/macOS
uname -a > system_info.txt
peerchat-cli version >> system_info.txt
cat /etc/os-release >> system_info.txt  # Linux only

# Windows
systeminfo > system_info.txt
peerchat-cli.exe version >> system_info.txt
```

3. **Gather relevant logs:**
```bash
# Last 100 lines of logs
tail -100 ~/.xelvra/peerchat.log > recent_logs.txt
```

### Where to Get Help

1. **[GitHub Issues](https://github.com/Xelvra/peerchat/issues)** - Bug reports and technical issues
2. **[GitHub Discussions](https://github.com/Xelvra/peerchat/discussions)** - General questions and community help
3. **[Wiki FAQ](FAQ)** - Common questions and answers
4. **[User Manual](User-Manual)** - Comprehensive documentation

### Creating Effective Bug Reports

Include this information:
- **System details** (OS, version, architecture)
- **Xelvra version** (`peerchat-cli version`)
- **Steps to reproduce** the issue
- **Expected vs actual behavior**
- **Relevant logs** (sanitized of personal info)
- **Diagnostic output** (`peerchat-cli doctor`)

---

**Still having issues?** Don't hesitate to [ask for help](https://github.com/Xelvra/peerchat/discussions) - the community is here to support you!
