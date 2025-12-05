package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ReverseProxy creates a reverse proxy handler for the given target URL
func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse target URL
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid proxy target",
			})
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Modify request
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Host = targetURL.Host

			// Get the path parameter from Gin (e.g., /login from /api/v1/auth/login)
			path := c.Param("path")
			if path == "" {
				path = "/"
			}

			req.URL.Path = path
			req.URL.RawQuery = c.Request.URL.RawQuery

			// Forward original headers (including Authorization, Cookie, etc.)
			req.Header = c.Request.Header.Clone()
			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "")
			}
		}

		// ModifyResponse to ensure cookies are properly forwarded to client
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Explicitly copy Set-Cookie headers from backend response
			// This ensures cookies set by the backend service reach the client
			if cookies := resp.Header["Set-Cookie"]; len(cookies) > 0 {
				for _, cookie := range cookies {
					c.Writer.Header().Add("Set-Cookie", cookie)
				}
			}
			return nil
		}

		// Handle errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "service unavailable",
				"details": err.Error(),
			})
		}

		// Serve the proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// ReverseProxyWithPrefix creates a reverse proxy that strips a path prefix
// Example: ReverseProxyWithPrefix("http://localhost:8080", "/api/v1/auth")
// Request: /api/v1/auth/login → Proxy: http://localhost:8080/login
func ReverseProxyWithPrefix(target, stripPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse target URL
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid proxy target",
			})
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Modify request
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Host = targetURL.Host

			// Strip prefix from path
			// /api/v1/auth/login → /login
			req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, stripPrefix)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			req.URL.RawQuery = c.Request.URL.RawQuery

			// Forward original headers
			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "")
			}
		}

		// Handle errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "service unavailable",
				"details": err.Error(),
			})
		}

		// Serve the proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
