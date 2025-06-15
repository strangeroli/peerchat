# Frequently Asked Questions (FAQ)

## 🔐 Security & Privacy

### Q: Is Xelvra really secure?
**A:** Yes! Xelvra uses industry-standard end-to-end encryption with the Signal Protocol, the same technology used by Signal, WhatsApp, and other secure messengers. All messages are encrypted before leaving your device and can only be decrypted by the intended recipient.

### Q: Can anyone intercept my messages?
**A:** No. Even if someone intercepts your encrypted messages, they cannot read them without your private keys. Xelvra also uses forward secrecy, meaning even if keys are compromised in the future, past messages remain secure.

### Q: Does Xelvra collect my data?
**A:** No. Xelvra is completely decentralized - there are no central servers to collect your data. Your messages go directly between devices, and your identity is stored only on your device.

### Q: Can governments or corporations spy on my communications?
**A:** Xelvra is designed to be censorship-resistant and surveillance-resistant. Since there are no central servers, there's no single point that can be compromised or pressured. However, always follow local laws and regulations.

## 🌐 How It Works

### Q: What does "peer-to-peer" mean?
**A:** P2P means your messages travel directly from your device to the recipient's device without going through central servers. It's like having a direct phone line to each person you talk to.

### Q: Do I need internet to use Xelvra?
**A:** For local network communication, you only need to be on the same WiFi network. For internet-wide communication, you need an internet connection. Future versions will support offline mesh networking via Bluetooth.

### Q: How does Xelvra find other users?
**A:** Xelvra uses multiple discovery methods:
- **mDNS**: Finds users on your local network
- **UDP Broadcast**: Discovers nearby users
- **DHT**: Distributed hash table for internet-wide discovery (planned)

### Q: What happens if I'm behind a firewall or NAT?
**A:** Xelvra includes advanced NAT traversal techniques and will automatically use relay servers when direct connections aren't possible, while maintaining end-to-end encryption.

## 💻 Installation & Setup

### Q: What platforms does Xelvra support?
**A:** Currently:
- ✅ **Linux** (all distributions)
- ✅ **macOS** (Intel and Apple Silicon)
- ✅ **Windows** (10/11)
- 📋 **Android/iOS** (planned for GUI version)

### Q: Do I need to be technical to use Xelvra?
**A:** Not at all! While the current CLI version requires basic command-line knowledge, it's designed to be user-friendly. The upcoming GUI version will be even easier to use.

### Q: Can I run multiple instances of Xelvra?
**A:** Yes! Each instance needs its own configuration directory. Use the `--config` flag to specify different directories for each instance.

### Q: Where are my keys and data stored?
**A:** Everything is stored in `~/.xelvra/` on your device:
- `identity.key` - Your private key (keep this secure!)
- `config.yaml` - Configuration settings
- `peerchat.log` - Application logs

## 🔧 Usage & Features

### Q: How do I backup my identity?
**A:** Simply backup the entire `~/.xelvra/` directory to a secure location. This contains your keys, configuration, and identity information. Without this backup, you'll lose access to your identity if your device is lost or damaged.

### Q: Can I change my identity or username?
**A:** Your cryptographic identity (Peer ID and DID) cannot be changed as they're based on your private key. However, you can set a display name that others see. To get a completely new identity, you'd need to run `peerchat-cli init` again (this creates new keys).

### Q: How many people can I chat with at once?
**A:** There's no hard limit, but performance depends on your device and network. The current CLI version handles multiple simultaneous connections efficiently.

### Q: Can I send files?
**A:** Yes! Use the `send-file` command:
```bash
peerchat-cli send-file <peer_id> /path/to/file
```
Files are transferred directly between peers with encryption and integrity verification.

### Q: Is there a group chat feature?
**A:** Group chat is planned for future releases. Currently, you can connect to multiple peers simultaneously, and messages are sent to all connected peers.

## 🛠️ Troubleshooting

### Q: Why can't I find any peers?
**A:** Common solutions:
1. Ensure you're on the same network as other Xelvra users
2. Check firewall settings (allow UDP port 42424)
3. Try a different network (mobile hotspot)
4. Run `peerchat-cli doctor` for diagnostics

### Q: Why do connections keep failing?
**A:** This usually indicates network issues:
1. Check internet connectivity
2. Verify firewall settings
3. Ensure both users have compatible versions
4. Try discovering peers again (addresses may change)

### Q: What does "simulation mode" mean?
**A:** Simulation mode activates when Xelvra can't establish real P2P connections, usually due to network restrictions. Run `peerchat-cli doctor` to diagnose the issue.

### Q: How do I read the logs?
**A:** Logs are stored in `~/.xelvra/peerchat.log`. You can view them with:
```bash
tail -f ~/.xelvra/peerchat.log  # Follow live logs
cat ~/.xelvra/peerchat.log      # View all logs
```

## 🚀 Development & Future

### Q: Is Xelvra open source?
**A:** Yes! Xelvra is licensed under the GNU Affero General Public License v3.0 (AGPLv3). All code is available on [GitHub](https://github.com/Xelvra/peerchat).

### Q: Can I contribute to the project?
**A:** Absolutely! We welcome contributions of all kinds:
- Code contributions
- Documentation improvements
- Bug reports and feature requests
- Community support
- Translations

See our [Contributing Guide](https://github.com/Xelvra/peerchat/blob/main/CONTRIBUTING.md) for details.

### Q: What's the development roadmap?
**A:** Xelvra follows a structured development approach:
- **Epoch 1**: CLI Foundation (current, largely complete)
- **Epoch 2**: API Service (planned)
- **Epoch 3**: GUI Application (planned)
- **Epoch 4**: Advanced Features (future)

### Q: When will the GUI version be available?
**A:** The GUI version is planned for Epoch 3. We're focusing on making the CLI version rock-solid first, as it serves as the foundation for all other components.

### Q: Will there be mobile apps?
**A:** Yes! Mobile apps for Android and iOS are planned as part of the GUI development in Epoch 3.

## 💡 Philosophy & Vision

### Q: Why create another messaging app?
**A:** Existing messaging platforms have fundamental issues:
- Centralized control and single points of failure
- Data collection and privacy violations
- Censorship and content restrictions
- Vendor lock-in and lack of interoperability

Xelvra addresses these issues with true decentralization and user ownership.

### Q: What does "#XelvraFree" mean?
**A:** #XelvraFree represents our commitment to digital freedom - the fundamental right to private communication without surveillance, censorship, or corporate control.

### Q: How is Xelvra funded?
**A:** Xelvra is currently developed as an open-source project. Future funding may come from:
- Community donations
- Hash Token ecosystem (internal virtual credits)
- Transparent crowdfunding
- Never through data collection or advertising

## 🆘 Getting More Help

### Still have questions?
1. **Search this wiki** for more detailed information
2. **Check [GitHub Issues](https://github.com/Xelvra/peerchat/issues)** for known issues
3. **Start a [Discussion](https://github.com/Xelvra/peerchat/discussions)** for community help
4. **Read the [User Manual](User-Manual)** for comprehensive documentation
5. **Join our community** and connect with other users

---

**Don't see your question here?** [Ask in GitHub Discussions](https://github.com/Xelvra/peerchat/discussions) and help us improve this FAQ!
