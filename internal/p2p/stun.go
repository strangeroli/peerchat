package p2p

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pion/stun"
	"github.com/sirupsen/logrus"
)

// STUNClient handles STUN server communication for NAT discovery
type STUNClient struct {
	servers []string
	logger  *logrus.Logger
}

// NewSTUNClient creates a new STUN client
func NewSTUNClient(logger *logrus.Logger) *STUNClient {
	return &STUNClient{
		servers: []string{
			"stun.l.google.com:19302",
			"stun1.l.google.com:19302",
			"stun2.l.google.com:19302",
			"stun.cloudflare.com:3478",
			"stun.nextcloud.com:443",
		},
		logger: logger,
	}
}

// DiscoverNAT discovers NAT type and public IP address
func (s *STUNClient) DiscoverNAT(ctx context.Context, localPort int) (*NATInfo, error) {
	s.logger.Info("Starting NAT discovery via STUN...")

	natInfo := &NATInfo{
		Type:        "unknown",
		STUNServers: s.servers,
		LocalPort:   localPort,
		UsingRelay:  false,
	}

	// Get local IP
	localIP, err := s.getLocalIP()
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get local IP")
		natInfo.LocalIP = "127.0.0.1"
	} else {
		natInfo.LocalIP = localIP
	}

	// Try each STUN server
	for _, server := range s.servers {
		s.logger.WithField("server", server).Debug("Trying STUN server")
		
		publicIP, publicPort, natType, err := s.querySTUNServer(ctx, server, localPort)
		if err != nil {
			s.logger.WithError(err).WithField("server", server).Debug("STUN query failed")
			continue
		}

		natInfo.PublicIP = publicIP
		natInfo.PublicPort = publicPort
		natInfo.Type = natType

		s.logger.WithFields(logrus.Fields{
			"public_ip":   publicIP,
			"public_port": publicPort,
			"nat_type":    natType,
			"server":      server,
		}).Info("NAT discovery successful")

		return natInfo, nil
	}

	s.logger.Warn("All STUN servers failed, assuming symmetric NAT")
	natInfo.Type = "symmetric"
	return natInfo, fmt.Errorf("all STUN servers failed")
}

// querySTUNServer queries a single STUN server
func (s *STUNClient) querySTUNServer(ctx context.Context, server string, localPort int) (string, int, string, error) {
	// Create UDP connection
	conn, err := net.DialTimeout("udp", server, 5*time.Second)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to connect to STUN server: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log error but don't fail the function - STUN query can continue
			_ = err // Explicitly ignore error
		}
	}()

	// Set deadline
	deadline, ok := ctx.Deadline()
	if ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return "", 0, "", fmt.Errorf("failed to set deadline: %w", err)
		}
	} else {
		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return "", 0, "", fmt.Errorf("failed to set deadline: %w", err)
		}
	}

	// Create STUN binding request
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	// Send request
	_, err = conn.Write(message.Raw)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to send STUN request: %w", err)
	}

	// Read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to read STUN response: %w", err)
	}

	// Parse response
	response := &stun.Message{Raw: buf[:n]}
	if err := response.Decode(); err != nil {
		return "", 0, "", fmt.Errorf("failed to decode STUN response: %w", err)
	}

	// Extract mapped address
	var mappedAddr stun.XORMappedAddress
	if err := mappedAddr.GetFrom(response); err != nil {
		// Try regular mapped address as fallback
		var regularAddr stun.MappedAddress
		if err := regularAddr.GetFrom(response); err != nil {
			return "", 0, "", fmt.Errorf("no mapped address in STUN response")
		}
		mappedAddr.IP = regularAddr.IP
		mappedAddr.Port = regularAddr.Port
	}

	publicIP := mappedAddr.IP.String()
	publicPort := mappedAddr.Port

	// Determine NAT type (simplified)
	localIP, _ := s.getLocalIP()
	natType := s.determineNATType(localIP, localPort, publicIP, publicPort)

	return publicIP, publicPort, natType, nil
}

// getLocalIP gets the local IP address
func (s *STUNClient) getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log error but don't fail the function - local IP detection can continue
			_ = err // Explicitly ignore error
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// determineNATType determines NAT type based on addresses
func (s *STUNClient) determineNATType(localIP string, localPort int, publicIP string, publicPort int) string {
	// If public IP equals local IP, no NAT
	if publicIP == localIP {
		return "none"
	}

	// If ports are the same, likely full cone NAT
	if publicPort == localPort {
		return "full_cone"
	}

	// Check if it's a private IP range
	if s.isPrivateIP(localIP) {
		// Different ports suggest port-restricted or symmetric NAT
		// This is a simplified detection - real implementation would need multiple tests
		return "port_restricted"
	}

	return "symmetric"
}

// isPrivateIP checks if an IP is in private range
func (s *STUNClient) isPrivateIP(ip string) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
