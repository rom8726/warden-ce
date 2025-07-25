package throttler

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rom8726/warden/internal/common/throttle"
	"github.com/rom8726/warden/internal/domain"
)

const scriptSrc = `
    local count = redis.call("HINCRBY", KEYS[1], ARGV[1], 1)
    redis.call("ZADD", KEYS[2], count, ARGV[1])
    if redis.call("TTL", KEYS[1]) < 0 then
        redis.call("EXPIRE", KEYS[1], tonumber(ARGV[2]))
    end
    if redis.call("TTL", KEYS[2]) < 0 then
        redis.call("EXPIRE", KEYS[2], tonumber(ARGV[2]))
    end
    return count
`

type Service struct {
	redisClient *redis.Client
	script      *redis.Script
	ttl         time.Duration
	enabled     bool
}

func New(ctx context.Context, redisClient *redis.Client, ttl time.Duration, enabled bool) (*Service, error) {
	srv := &Service{
		redisClient: redisClient,
		script:      redis.NewScript(scriptSrc),
		ttl:         ttl,
		enabled:     enabled,
	}

	if err := srv.script.Load(ctx, redisClient).Err(); err != nil {
		return nil, fmt.Errorf("load LUA script: %w", err)
	}

	return srv, nil
}

func (srv *Service) Provide(ctx context.Context, event *domain.Event) error {
	if !srv.enabled {
		return nil
	}

	keys := []string{throttle.EventsCountMapKey(event.ProjectID), throttle.EventsIndexSetKey(event.ProjectID)}
	fingerPrint := event.FullFingerprint()
	ttl := int(srv.ttl.Seconds())

	return srv.script.Run(ctx, srv.redisClient, keys, fingerPrint, ttl).Err()
}
