// Package auth determines and asserts client permissions to access and modify
// resources.
package auth

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bakape/meguca/config"
)

var (
	// IsReverseProxied specifies, if the server is deployed behind a reverse
	// proxy.
	IsReverseProxied bool

	// ReverseProxyIP specifies the IP of a non-localost reverse proxy. Used for
	// filtering in XFF IP determination.
	ReverseProxyIP string
)

// User contains ID, password hash and board-related data of a registered user
// account
type User struct {
	ID       string    `gorethink:"id"`
	Password []byte    `gorethink:"password"`
	Sessions []Session `gorethink:"sessions"`
}

// Session contains the token and expiry time of a single authenticated login
// session
type Session struct {
	Token   string    `gorethink:"token"`
	Expires time.Time `gorethink:"expires"`
}

// Ident is used to verify a client's access and write permissions. Contains its
// IP and logged in user data, if any.
type Ident struct {
	UserID string
	IP     string
}

// LookUpIdent determine access rights of an IP
func LookUpIdent(req *http.Request) Ident {
	ident := Ident{
		IP: GetIP(req),
	}
	return ident
}

// IsBoard confirms the string is a valid board
func IsBoard(board string) bool {
	if board == "all" {
		return true
	}
	return IsNonMetaBoard(board)
}

// IsNonMetaBoard returns wheather a valid board is a classic board and not
// some other path that emulates a board
func IsNonMetaBoard(board string) bool {
	for _, b := range config.Get().Boards {
		if board == b {
			return true
		}
	}
	return false
}

// GetIP extracts the IP of a request, honouring reverse proxies, if set
func GetIP(req *http.Request) string {
	if IsReverseProxied {
		for _, h := range [...]string{"X-Forwarded-For", "X-Real-Ip"} {
			addresses := strings.Split(req.Header.Get(h), ",")

			// March from right to left until we get a public address.
			// That will be the address right before our reverse proxy.
			for i := len(addresses) - 1; i >= 0; i-- {
				// Header can contain padding spaces
				ip := strings.TrimSpace(addresses[i])

				// Filter the reverse proxy IPs
				switch {
				case ip == ReverseProxyIP:
				case !net.ParseIP(ip).IsGlobalUnicast():
				default:
					return ip
				}
			}
		}
	}
	return req.RemoteAddr
}
