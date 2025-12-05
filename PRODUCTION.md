# Production Deployment Guide

## Cookie Configuration on Production

### Overview
On production, cookies are configured with enhanced security:
- **Secure**: `true` (requires HTTPS)
- **SameSite**: `None` (allows cross-origin requests)
- **HttpOnly**: `true` (XSS protection)
- **Domain**: Configurable for cross-subdomain support

---

## Configuration Steps

### 1. Update `config.yaml` for Production

```yaml
env: production
service_name: HeiCRM

jwt:
  secret_key: "your-super-secure-secret-key-change-this"
  alg: HS256

cookie:
  # For cross-subdomain cookies (api.yourdomain.com + crm.yourdomain.com)
  domain: ".yourdomain.com"

  # OR leave empty if frontend and API are on same domain
  domain: ""

database:
  host: your-db-host
  port: 5432
  user: postgres
  password: your-db-password
  name: crm
  sslmode: require  # Enable SSL for production!

redis:
  host: your-redis-host
  port: 6379
  password: your-redis-password

serves:
  auth: "http://auth-service:8080"  # Internal Docker network
```

### 2. Update CORS Origins in API Gateway

Edit `services/apigateway/cmd/gateway/main.go`:

```go
if cfg.Env == "production" || cfg.Env == "prod" {
    allowedOrigins = append(allowedOrigins,
        "https://crm.yourdomain.com",      // Your actual production domain
        "https://app.yourdomain.com",      // Alternative domain if needed
    )
}
```

---

## SSL/TLS Setup (REQUIRED)

### Why HTTPS is Required
- `SameSite=None` cookies **MUST** have `Secure=true`
- `Secure=true` cookies **ONLY** work over HTTPS
- Modern browsers will reject insecure cookies in production

### Option 1: Nginx Reverse Proxy (Recommended)

```nginx
# /etc/nginx/sites-available/crm

server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL certificates (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

    # Proxy to API Gateway
    location / {
        proxy_pass http://localhost:8000;
        proxy_http_version 1.1;

        # Forward headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Cookie support
        proxy_set_header Cookie $http_cookie;
        proxy_pass_header Set-Cookie;
    }
}
```

**Get SSL Certificate:**
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.yourdomain.com
```

### Option 2: Caddy Server (Automatic HTTPS)

```caddyfile
# Caddyfile

api.yourdomain.com {
    reverse_proxy localhost:8000
}
```

Caddy automatically obtains and renews SSL certificates!

### Option 3: CloudFlare (Easiest)

1. Point your domain to CloudFlare
2. Enable "Full (strict)" SSL mode
3. CloudFlare handles SSL termination
4. Your origin can be HTTP (localhost:8000)

---

## Deployment Architecture

### Development (Current)
```
Frontend (http://localhost:3000)
    ↓
API Gateway (http://localhost:8000)
    ↓
Auth Service (http://localhost:8080)

Cookie Settings:
- Secure: false
- SameSite: Lax
- Domain: "" (same-origin)
```

### Production (Recommended)
```
Frontend (https://crm.yourdomain.com)
    ↓
Nginx/Caddy (https://api.yourdomain.com:443)
    ↓
API Gateway (http://localhost:8000)
    ↓
Auth Service (http://localhost:8080)

Cookie Settings:
- Secure: true (HTTPS required)
- SameSite: None (cross-origin)
- Domain: ".yourdomain.com" (cross-subdomain)
```

---

## Domain Configuration Examples

### Scenario 1: Frontend and API on Same Domain
```
Frontend: https://yourdomain.com
API:      https://yourdomain.com/api/v1

config.yaml:
  cookie:
    domain: ""  # Leave empty
```

### Scenario 2: Frontend and API on Different Subdomains (MOST COMMON)
```
Frontend: https://crm.yourdomain.com
API:      https://api.yourdomain.com

config.yaml:
  cookie:
    domain: ".yourdomain.com"  # Note the leading dot!
```

### Scenario 3: Multiple Domains (Advanced)
```
Frontend: https://app.yourdomain.com
API:      https://api.yourdomain.com
Admin:    https://admin.yourdomain.com

config.yaml:
  cookie:
    domain: ".yourdomain.com"  # Works for all subdomains
```

---

## Testing Production Setup Locally

You can test HTTPS locally with self-signed certificates:

### 1. Generate Self-Signed Certificate
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### 2. Update Gin to Use TLS
```go
// services/apigateway/cmd/gateway/main.go
if err := r.RunTLS(":8000", "./cert.pem", "./key.pem"); err != nil {
    log.Fatalf("Failed to start API Gateway: %v", err)
}
```

### 3. Test
```bash
curl -k https://localhost:8000/api/v1/health
```

---

## Checklist Before Going Live

- [ ] `env: production` in `config.yaml`
- [ ] Strong JWT secret key (256+ bits)
- [ ] SSL certificate configured
- [ ] Database SSL enabled (`sslmode: require`)
- [ ] Redis password set
- [ ] CORS origins updated with production domains
- [ ] Cookie domain configured (if using subdomains)
- [ ] Firewall rules configured
- [ ] Environment variables secured (no plaintext secrets)
- [ ] Log rotation enabled
- [ ] Backup strategy in place

---

## Security Notes

### DO:
✅ Use HTTPS in production
✅ Set strong JWT secrets (>32 characters)
✅ Enable database SSL
✅ Use Redis password auth
✅ Set `HttpOnly: true` on cookies (prevents XSS)
✅ Limit CORS origins to your actual domains

### DON'T:
❌ Commit `config.yaml` with real secrets
❌ Use HTTP in production
❌ Use `SameSite=None` without `Secure=true`
❌ Use `AllowOrigins: ["*"]` with credentials
❌ Expose internal service ports (8080) to internet
❌ Use default/weak passwords

---

## Troubleshooting

### Problem: Cookies not visible in browser (Production)
**Causes:**
1. Not using HTTPS → Browser rejects `Secure=true` cookies
2. Wrong domain → Cookie domain doesn't match request domain
3. CORS misconfiguration → `AllowCredentials: true` requires specific origins

**Solution:**
```bash
# Check response headers
curl -v https://api.yourdomain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test"}'

# Should see:
# Set-Cookie: refresh=...; Path=/api/v1/auth; Domain=.yourdomain.com; Secure; HttpOnly; SameSite=None
```

### Problem: CORS errors in production
**Solution:**
- Verify frontend domain is in `allowedOrigins` list
- Check that `AllowCredentials: true` is set
- Ensure `Access-Control-Allow-Origin` is NOT `*` when using credentials

### Problem: Cookie works on login but not on refresh
**Causes:**
1. Cookie path too narrow
2. Cookie expired
3. Redis lost the token

**Solution:**
- Verify path is `/api/v1/auth` (not `/api/v1/auth/refresh`)
- Check Redis: `redis-cli KEYS "refresh_token:*"`
- Verify token TTL matches JWT expiry

---

## Monitoring

Add monitoring to track cookie issues:

```go
// Log cookie settings in production
if isProduction {
    log.Printf("Setting secure cookie: Domain=%s, SameSite=None, Secure=true", cookieDomain)
}
```

---

## Questions?

If you encounter issues:
1. Check browser DevTools → Network → Response Headers
2. Verify `Set-Cookie` header is present
3. Check Application → Cookies in DevTools
4. Review logs in `logs/apigateway.log` and `logs/auth.log`