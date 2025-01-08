# Relatórios dos Testes Automatizados

- Execução dos testes

```bash
❯ go test -timeout 30s  -v ./... -coverprofile=coverage.out
        github.com/vs0uz4/rate-limit/cmd/server         coverage: 0.0% of statements
?       github.com/vs0uz4/rate-limit/internal/contract  [no test files]
?       github.com/vs0uz4/rate-limit/internal/domain/errors     [no test files]
        github.com/vs0uz4/rate-limit/internal/mock              coverage: 0.0% of statements
=== RUN   TestLoadConfigExecPathError
--- PASS: TestLoadConfigExecPathError (0.00s)
=== RUN   TestLoadConfigWithInvalidTokenLimits
--- PASS: TestLoadConfigWithInvalidTokenLimits (0.00s)
=== RUN   TestLoadConfig
--- PASS: TestLoadConfig (0.00s)
=== RUN   TestLoadConfigWithEmptyTokenLimits
--- PASS: TestLoadConfigWithEmptyTokenLimits (0.00s)
=== RUN   TestLoadConfigMissingEnvFile
--- PASS: TestLoadConfigMissingEnvFile (0.00s)
=== RUN   TestLoadConfigMissingRedisHost
--- PASS: TestLoadConfigMissingRedisHost (0.00s)
=== RUN   TestLoadConfigWithDefaultValues
--- PASS: TestLoadConfigWithDefaultValues (0.00s)
PASS
coverage: 100.0% of statements
ok      github.com/vs0uz4/rate-limit/config     0.195s  coverage: 100.0% of statements
=== RUN   TestNewRedisClientSuccess
--- PASS: TestNewRedisClientSuccess (0.00s)
=== RUN   TestNewRedisClientFailure
--- PASS: TestNewRedisClientFailure (0.07s)
PASS
coverage: 40.0% of statements
ok      github.com/vs0uz4/rate-limit/internal/infra/redis       0.424s  coverage: 40.0% of statements
=== RUN   TestRateLimiterTokenConfigured
--- PASS: TestRateLimiterTokenConfigured (0.00s)
=== RUN   TestRateLimiterTokenNotConfigured
--- PASS: TestRateLimiterTokenNotConfigured (0.00s)
=== RUN   TestRateLimiterNoToken
--- PASS: TestRateLimiterNoToken (0.00s)
=== RUN   TestAllowSuccess
--- PASS: TestAllowSuccess (0.00s)
=== RUN   TestAllowLimitExceeded
--- PASS: TestAllowLimitExceeded (0.00s)
=== RUN   TestRateLimiterExpire
=== RUN   TestRateLimiterExpire/Should_set_TTL_when_key_exists_without_expiration
=== RUN   TestRateLimiterExpire/Should_return_error_if_Expire_fails
--- PASS: TestRateLimiterExpire (0.00s)
    --- PASS: TestRateLimiterExpire/Should_set_TTL_when_key_exists_without_expiration (0.00s)
    --- PASS: TestRateLimiterExpire/Should_return_error_if_Expire_fails (0.00s)
=== RUN   TestAllowErrorOnIncr
--- PASS: TestAllowErrorOnIncr (0.00s)
=== RUN   TestAllowErrorOnTTLRateKey
--- PASS: TestAllowErrorOnTTLRateKey (0.00s)
=== RUN   TestAllowErrorOnExpireRateKey
--- PASS: TestAllowErrorOnExpireRateKey (0.00s)
=== RUN   TestAllowErrorOnTTLBlockKey
--- PASS: TestAllowErrorOnTTLBlockKey (0.00s)
=== RUN   TestAllowErrorOnSetNX
Running TestAllowErrorOnSetNX...
Result - Allowed: false, Error: error on setnx
--- PASS: TestAllowErrorOnSetNX (0.00s)
=== RUN   TestAllowBlockKeyTTLNotExpired
--- PASS: TestAllowBlockKeyTTLNotExpired (0.00s)
PASS
coverage: 98.0% of statements
ok      github.com/vs0uz4/rate-limit/internal/rate_limiter      0.520s  coverage: 98.0% of statements
=== RUN   TestPingEndpoint
--- PASS: TestPingEndpoint (0.00s)
=== RUN   TestStartServer
2025/01/08 09:02:53 Starting server on port 9090...
--- PASS: TestStartServer (0.00s)
PASS
coverage: 33.3% of statements
ok      github.com/vs0uz4/rate-limit/internal/webserver 0.881s  coverage: 33.3% of statements
=== RUN   TestExtractIPWithoutPort
--- PASS: TestExtractIPWithoutPort (0.00s)
=== RUN   TestExtractIPIPv6ToIPv4
--- PASS: TestExtractIPIPv6ToIPv4 (0.00s)
=== RUN   TestRateLimiterMiddlewareTokenValid
--- PASS: TestRateLimiterMiddlewareTokenValid (0.00s)
=== RUN   TestRateLimiterMiddlewareTokenInvalid
--- PASS: TestRateLimiterMiddlewareTokenInvalid (0.00s)
=== RUN   TestRateLimiterMiddlewareNoToken
--- PASS: TestRateLimiterMiddlewareNoToken (0.00s)
=== RUN   TestRateLimiterMiddlewareLimitExceeded
--- PASS: TestRateLimiterMiddlewareLimitExceeded (0.00s)
=== RUN   TestRateLimiterMiddlewareError
--- PASS: TestRateLimiterMiddlewareError (0.00s)
PASS
coverage: 54.3% of statements
ok      github.com/vs0uz4/rate-limit/internal/webserver/middleware      0.696s  coverage: 54.3% of statements
```

- Conbertura dos testes

```bash
❯ go tool cover -func=coverage.out

github.com/vs0uz4/rate-limit/cmd/server/main.go:13:                             main                            0.0%
github.com/vs0uz4/rate-limit/config/env.go:23:                                  LoadConfig                      100.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis.go:14:                  NewRedisClient                  100.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis_adapter.go:14:          NewRedisAdapter                 0.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis_adapter.go:18:          Incr                            0.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis_adapter.go:23:          SetNX                           0.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis_adapter.go:28:          TTL                             0.0%
github.com/vs0uz4/rate-limit/internal/infra/redis/redis_adapter.go:33:          Expire                          0.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_persistence_provider.go:17:     NewMockPersistenceProvider      100.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_persistence_provider.go:28:     Incr                            100.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_persistence_provider.go:36:     SetNX                           100.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_persistence_provider.go:43:     TTL                             100.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_persistence_provider.go:53:     Expire                          100.0%
github.com/vs0uz4/rate-limit/internal/mock/mock_rate_limiter.go:12:             Allow                           100.0%
github.com/vs0uz4/rate-limit/internal/rate_limiter/rate_limiter.go:21:          NewRateLimiter                  100.0%
github.com/vs0uz4/rate-limit/internal/rate_limiter/rate_limiter.go:28:          Allow                           100.0%
github.com/vs0uz4/rate-limit/internal/webserver/middleware/rate_limiter.go:11:  extractIP                       100.0%
github.com/vs0uz4/rate-limit/internal/webserver/middleware/rate_limiter.go:24:  RateLimiterMiddleware           100.0%
github.com/vs0uz4/rate-limit/internal/webserver/webserver.go:13:                Start                           100.0%
total:                                                                          (statements)                    84.9%
```

> [!IMPORTANT]
> Como o `redis_adapter.go` está sendo fortemente testado indiretamente através dos testes implementados no `RateLimiter`, não vimos necessidade de implementação de testes dedicados para o mesmo, desta forma evitamos a duplicidade de testes desnecessária, otimizando o desenvolvimento.
