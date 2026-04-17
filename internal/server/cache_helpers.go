package server

import (
	"context"
	"fmt"
	"time"

	"github.com/kiselevos/new_tax/internal/calculate"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

func (s *serverStruct) cacheGet(ctx context.Context, key string, dst proto.Message) (bool, error) {

	if s.redis == nil {
		return false, nil
	}

	b, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			return false, err
		}
		return false, nil
	}

	if err := proto.Unmarshal(b, dst); err != nil {
		return false, err
	}

	return true, nil
}

func (s *serverStruct) cacheSet(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error {
	if s.redis == nil {
		return nil
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, b, ttl).Err()
}

func buildPublicKey(in calculate.CalculateInput) string {
	return fmt.Sprintf(
		"calc:public:v1:gross=%d:tm=%d:nc=%d",
		in.GrossSalary,
		in.TerritorialMultiplier,
		in.NorthernCoefficient,
	)
}

func buildPrivateKey(in calculate.CalculateInput) string {
	start := in.StartDate.UTC().Format("2006-01")
	return fmt.Sprintf(
		"calc:private:v1:gross=%d:tm=%d:nc=%d:start=%s:priv=%t:nonres=%t:bonuses=%v",
		in.GrossSalary,
		in.TerritorialMultiplier,
		in.NorthernCoefficient,
		start,
		in.HasTaxPrivilege,
		in.IsNotResident,
		in.MonthlyBonuses,
	)
}
