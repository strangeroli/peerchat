#!/bin/bash

# Xelvra P2P Messenger - System Service Installation Script
# This script installs Xelvra as a systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="peerchat"
SERVICE_USER="xelvra"
SERVICE_GROUP="xelvra"
INSTALL_DIR="/opt/xelvra"
DATA_DIR="/var/lib/xelvra"
LOG_DIR="/var/log/xelvra"

echo -e "${BLUE}🚀 Xelvra P2P Messenger - Service Installation${NC}"
echo "=============================================="

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   echo "Usage: sudo ./scripts/install-service.sh"
   exit 1
fi

# Check if systemd is available
if ! command -v systemctl &> /dev/null; then
    echo -e "${RED}❌ systemd is not available on this system${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Installation Configuration:${NC}"
echo "  Service name: $SERVICE_NAME"
echo "  User/Group: $SERVICE_USER:$SERVICE_GROUP"
echo "  Install directory: $INSTALL_DIR"
echo "  Data directory: $DATA_DIR"
echo "  Log directory: $LOG_DIR"
echo

# Confirm installation
read -p "Do you want to proceed with the installation? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⚠️ Installation cancelled${NC}"
    exit 0
fi

echo -e "${BLUE}📦 Creating system user and directories...${NC}"

# Create system user and group
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd --system --shell /bin/false --home-dir "$DATA_DIR" --create-home "$SERVICE_USER"
    echo -e "${GREEN}✅ Created system user: $SERVICE_USER${NC}"
else
    echo -e "${YELLOW}⚠️ User $SERVICE_USER already exists${NC}"
fi

# Create directories
mkdir -p "$INSTALL_DIR"/{bin,configs}
mkdir -p "$DATA_DIR"
mkdir -p "$LOG_DIR"

echo -e "${BLUE}📁 Setting up directory permissions...${NC}"

# Set ownership and permissions
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$DATA_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$LOG_DIR"

chmod 755 "$INSTALL_DIR"
chmod 700 "$DATA_DIR"
chmod 755 "$LOG_DIR"

echo -e "${BLUE}📋 Installing binary and configuration...${NC}"

# Copy binary
if [[ -f "bin/peerchat-cli" ]]; then
    cp "bin/peerchat-cli" "$INSTALL_DIR/bin/"
    chmod 755 "$INSTALL_DIR/bin/peerchat-cli"
    echo -e "${GREEN}✅ Installed binary: $INSTALL_DIR/bin/peerchat-cli${NC}"
else
    echo -e "${RED}❌ Binary not found: bin/peerchat-cli${NC}"
    echo "Please build the project first: go build -o bin/peerchat-cli ./cmd/peerchat-cli"
    exit 1
fi

# Copy systemd service file
if [[ -f "configs/systemd/peerchat.service" ]]; then
    cp "configs/systemd/peerchat.service" "/etc/systemd/system/"
    echo -e "${GREEN}✅ Installed systemd service file${NC}"
else
    echo -e "${RED}❌ Service file not found: configs/systemd/peerchat.service${NC}"
    exit 1
fi

echo -e "${BLUE}🔧 Configuring systemd service...${NC}"

# Reload systemd
systemctl daemon-reload

# Enable service
systemctl enable "$SERVICE_NAME"
echo -e "${GREEN}✅ Service enabled for auto-start${NC}"

echo -e "${BLUE}🔑 Initializing Xelvra identity...${NC}"

# Initialize identity as service user
sudo -u "$SERVICE_USER" XELVRA_CONFIG_DIR="$DATA_DIR" "$INSTALL_DIR/bin/peerchat-cli" init
echo -e "${GREEN}✅ Identity initialized${NC}"

echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
echo
echo -e "${BLUE}📖 Service Management Commands:${NC}"
echo "  Start service:    sudo systemctl start $SERVICE_NAME"
echo "  Stop service:     sudo systemctl stop $SERVICE_NAME"
echo "  Restart service:  sudo systemctl restart $SERVICE_NAME"
echo "  Check status:     sudo systemctl status $SERVICE_NAME"
echo "  View logs:        sudo journalctl -u $SERVICE_NAME -f"
echo
echo -e "${BLUE}📁 Important Directories:${NC}"
echo "  Binary:           $INSTALL_DIR/bin/peerchat-cli"
echo "  Configuration:    $DATA_DIR/config.yaml"
echo "  Identity:         $DATA_DIR/identity.key"
echo "  Logs:             $LOG_DIR/peerchat.log"
echo
echo -e "${YELLOW}⚠️ Next Steps:${NC}"
echo "1. Review configuration: $DATA_DIR/config.yaml"
echo "2. Start the service: sudo systemctl start $SERVICE_NAME"
echo "3. Check status: sudo systemctl status $SERVICE_NAME"
echo
echo -e "${GREEN}✅ Xelvra P2P Messenger service is ready!${NC}"
