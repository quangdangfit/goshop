package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	mr *miniredis.Miniredis
	r  Redis
}

func (s *RedisTestSuite) SetupTest() {
	mr, err := miniredis.Run()
	s.Require().NoError(err)
	s.mr = mr

	rdb := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	s.r = &redis{cmd: rdb}
}

func (s *RedisTestSuite) TeardownTest() {
	s.mr.Close()
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func TestNew_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	r := New(Config{Address: mr.Addr()})
	assert.NotNil(t, r)
}

func (s *RedisTestSuite) TestIsConnected() {
	tests := []struct {
		name     string
		setup    func() Redis
		expected bool
	}{
		{
			name:     "True",
			setup:    func() Redis { return s.r },
			expected: true,
		},
		{
			name:     "NilCmd",
			setup:    func() Redis { return &redis{cmd: nil} },
			expected: false,
		},
		{
			name: "ServerDown",
			setup: func() Redis {
				s.mr.Close()
				return s.r
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			r := tc.setup()
			assert.Equal(s.T(), tc.expected, r.IsConnected())
		})
	}
}

func (s *RedisTestSuite) TestSet_And_Get() {
	type payload struct {
		Name string `json:"name"`
	}

	err := s.r.Set("key1", payload{Name: "test"})
	assert.NoError(s.T(), err)

	var got payload
	err = s.r.Get("key1", &got)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test", got.Name)
}

func (s *RedisTestSuite) TestGet_NotFound() {
	var v interface{}
	err := s.r.Get("nonexistent", &v)
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestGet_InvalidJSON() {
	s.mr.Set("badkey", "not-json")
	var v map[string]string
	err := s.r.Get("badkey", &v)
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestSetWithExpiration() {
	err := s.r.SetWithExpiration("expkey", "value", 10*time.Second)
	assert.NoError(s.T(), err)

	var got string
	err = s.r.Get("expkey", &got)
	assert.NoError(s.T(), err)
}

func (s *RedisTestSuite) TestRemove() {
	_ = s.r.Set("rmkey", "val")
	err := s.r.Remove("rmkey")
	assert.NoError(s.T(), err)

	var v interface{}
	assert.Error(s.T(), s.r.Get("rmkey", &v))
}

func (s *RedisTestSuite) TestKeys() {
	_ = s.r.Set("prefix:a", "1")
	_ = s.r.Set("prefix:b", "2")

	keys, err := s.r.Keys("prefix:*")
	assert.NoError(s.T(), err)
	assert.Len(s.T(), keys, 2)
}

func (s *RedisTestSuite) TestRemovePattern() {
	tests := []struct {
		name  string
		setup func()
		check func()
	}{
		{
			name: "MatchingKeys",
			setup: func() {
				_ = s.r.Set("pattern:x", "1")
				_ = s.r.Set("pattern:y", "2")
				_ = s.r.Set("other", "3")
			},
			check: func() {
				keys, _ := s.r.Keys("pattern:*")
				assert.Empty(s.T(), keys)
				var v interface{}
				assert.NoError(s.T(), s.r.Get("other", &v))
			},
		},
		{
			name:  "NoMatch",
			setup: func() {},
			check: func() {},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()
			var pattern string
			if tc.name == "NoMatch" {
				pattern = "nomatch:*"
			} else {
				pattern = "pattern:*"
			}
			err := s.r.RemovePattern(pattern)
			assert.NoError(s.T(), err)
			tc.check()
		})
	}
}

func (s *RedisTestSuite) TestRemovePattern_KeysError() {
	s.mr.Close()
	err := s.r.RemovePattern("*")
	assert.Error(s.T(), err)
}

type failOnDelCmdable struct {
	goredis.Cmdable
}

func (f *failOnDelCmdable) Del(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, errors.New("del error"))
}

func (s *RedisTestSuite) TestRemovePattern_RemoveError() {
	mr2, err2 := miniredis.Run()
	s.Require().NoError(err2)
	defer mr2.Close()

	rdb2 := goredis.NewClient(&goredis.Options{Addr: mr2.Addr()})
	baseRedis := &redis{cmd: rdb2}

	err := baseRedis.Set("pattern:key1", "val")
	s.Require().NoError(err)

	r2 := &redis{cmd: &failOnDelCmdable{rdb2}}
	err = r2.RemovePattern("pattern:*")
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestIncr() {
	tests := []struct {
		name     string
		setup    func()
		key      string
		expected int64
		wantErr  bool
	}{
		{
			name:     "FirstTime",
			setup:    func() {},
			key:      "incr_key",
			expected: int64(1),
		},
		{
			name: "SubsequentTimes",
			setup: func() {
				_, _ = s.r.Incr("incr_key2", 60*time.Second)
			},
			key:      "incr_key2",
			expected: int64(2),
		},
		{
			name: "ServerDown",
			setup: func() {
				s.mr.Close()
			},
			key:      "incr_key",
			expected: int64(0),
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()
			count, err := s.r.Incr(tc.key, 60*time.Second)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
			assert.Equal(s.T(), tc.expected, count)
		})
	}
}

// Error path tests — close server to force failures

func (s *RedisTestSuite) TestSet_ServerDown() {
	s.mr.Close()
	err := s.r.Set("k", "v")
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestSetWithExpiration_ServerDown() {
	s.mr.Close()
	err := s.r.SetWithExpiration("k", "v", time.Second)
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestRemove_ServerDown() {
	s.mr.Close()
	err := s.r.Remove("k")
	assert.Error(s.T(), err)
}

func (s *RedisTestSuite) TestKeys_ServerDown() {
	s.mr.Close()
	_, err := s.r.Keys("*")
	assert.Error(s.T(), err)
}
