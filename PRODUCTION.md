# 🚀 Production Deployment Guide

## Pre-Deployment Checklist

### ✅ Code Quality
- [x] All tests passing
- [x] Linter checks passing (0 issues)
- [x] Swagger documentation generated
- [x] All handlers documented
- [x] Error handling implemented
- [x] Logging configured

### ✅ Performance Optimizations
- [x] Redis caching with optimized TTLs
- [x] Efficient cache key strategies
- [x] Graceful error handling
- [x] Request timeout configuration
- [x] Connection pooling

### ✅ Security
- [x] Input validation
- [x] CORS configuration
- [x] Error message sanitization
- [x] Trusted proxy configuration
- [x] Request size limits

### ✅ Monitoring & Observability
- [x] Structured logging (Zap)
- [x] Health check endpoints
- [x] Request logging middleware
- [x] Error tracking

## Environment Configuration

### Required Environment Variables

```env
# Application
APP_NAME=jutsu-api
APP_VERSION=1.0.0
APP_ENV=production
APP_DEBUG=false

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
SERVER_IDLE_TIMEOUT=120
SERVER_PREFORK=true

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5

# Logger
LOG_LEVEL=info
LOG_ENCODING=json
LOG_OUTPUT_PATH=/var/log/jutsu-api/app.log
LOG_ERROR_PATH=/var/log/jutsu-api/error.log
```

## Deployment Steps

### 1. Build Production Binary

```bash
make build-prod
```

### 2. Run with Systemd (Linux)

Create `/etc/systemd/system/jutsu-api.service`:

```ini
[Unit]
Description=Jutsu API
After=network.target redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/jutsu-api
ExecStart=/opt/jutsu-api/bin/jutsu-api
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable jutsu-api
sudo systemctl start jutsu-api
```

### 3. Docker Deployment

```bash
docker build -t jutsu-api:latest .
docker run -d \
  --name jutsu-api \
  -p 8080:8080 \
  --env-file .env \
  --restart unless-stopped \
  jutsu-api:latest
```

### 4. Docker Compose

```bash
docker-compose up -d
```

## Performance Tuning

### Cache Strategy

The API uses intelligent caching with the following TTLs:

| Data Type | TTL | Reason |
|-----------|-----|--------|
| Home Info | 15 min | Frequently updated |
| Anime Info | 1 hour | Relatively stable |
| Categories | 30 min | Moderate updates |
| Search/Filter | 5 min | Dynamic results |
| Streaming | 5 min | Links may expire |
| Schedule | 10 min | Daily updates |

### Redis Configuration

For production, configure Redis with:
- Persistence enabled (AOF + RDB)
- Memory limits
- Eviction policy: `allkeys-lru`
- Connection pooling optimized

### Server Configuration

- Enable prefork for multi-core systems
- Set appropriate timeouts
- Configure connection limits
- Enable compression

## Monitoring

### Health Checks

- `/health` - Full health check
- `/ready` - Readiness probe
- `/live` - Liveness probe

### Logging

Logs are structured JSON in production:
- Application logs: `/var/log/jutsu-api/app.log`
- Error logs: `/var/log/jutsu-api/error.log`
- Log rotation: 100MB max, 10 backups, 30 days retention

### Metrics

Monitor:
- Request rate
- Response times
- Error rates
- Cache hit/miss ratios
- Redis connection pool usage

## Scaling

### Horizontal Scaling

1. Deploy multiple instances behind a load balancer
2. Use shared Redis for cache
3. Configure sticky sessions if needed
4. Monitor resource usage

### Vertical Scaling

1. Increase server resources
2. Tune Redis memory limits
3. Adjust connection pool sizes
4. Enable prefork for multi-core

## Backup & Recovery

### Redis Backup

```bash
# Automated backup script
redis-cli --rdb /backup/redis-$(date +%Y%m%d).rdb
```

### Application Backup

- Configuration files
- Environment variables
- Log files
- Swagger documentation

## Security Hardening

1. **Firewall**: Only expose necessary ports
2. **SSL/TLS**: Use reverse proxy (nginx/traefik) for HTTPS
3. **Rate Limiting**: Implement at reverse proxy level
4. **Secrets**: Use secret management (Vault, AWS Secrets Manager)
5. **Updates**: Keep dependencies updated

## Troubleshooting

### Common Issues

1. **Redis Connection Failed**
   - Check Redis is running
   - Verify connection settings
   - Check firewall rules

2. **High Memory Usage**
   - Review cache TTLs
   - Check for memory leaks
   - Adjust Redis eviction policy

3. **Slow Response Times**
   - Check cache hit rates
   - Review database queries
   - Monitor network latency

## Support

For production issues:
1. Check logs: `journalctl -u jutsu-api -f`
2. Review health endpoints
3. Check Redis status
4. Review application metrics

---

**Last Updated**: 2025-01-18

