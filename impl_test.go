package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite
	cache *cache
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (cs *CacheSuite) SetupTest() {
	cs.cache = New().(*cache)
}

func (cs *CacheSuite) TestSet() {
	testCases := []struct {
		Desc   string
		Key    string
		Expect interface{}
	}{
		{
			Desc:   "set value when use generally",
			Key:    "use generally",
			Expect: "case1",
		},
		{
			Desc:   "set value when over write",
			Key:    "over write",
			Expect: "case2",
		},
	}

	for _, tc := range testCases {
		cs.cache.Set(tc.Key, tc.Expect)

		res, _ := cs.cache.storage[tc.Key]
		cs.Require().Equal(tc.Expect, res.value, tc.Desc)
	}
}

func (cs *CacheSuite) TestGet() {
	defaultSetup := func(key string, val interface{}) {
		cs.cache.m.Lock()
		defer cs.cache.m.Unlock()

		cs.cache.storage[key] = &data{
			value: val,
			expiredHandle: time.AfterFunc(defaultExpiredTime, func() {
				cs.cache.del(key)
			}),
		}
	}

	testCases := []struct {
		Desc  string
		Key   string
		Err   error
		Data  interface{}
		Setup func(string, interface{})
	}{
		{
			Desc: "string type",
			Key:  "string",
			Data: "string",
		},
		{
			Desc: "int type",
			Key:  "int",
			Data: 1,
		},
		{
			Desc: "bool type",
			Key:  "bool",
			Data: true,
		},
		{
			Desc: "array type",
			Key:  "array",
			Data: []int{1, 2, 3, 4},
		},
		{
			Desc: "object type",
			Key:  "object",
			Data: map[string]string{"test": "test"},
		},
		{
			Desc: "data not found error",
			Key:  "not found",
			Err:  ErrDataNotFound,
			Setup: func(key string, val interface{}) {
			},
		},
		{
			Desc: "data not found error cause by cache expired",
			Key:  "cache expired",
			Err:  ErrDataNotFound,
			Setup: func(key string, val interface{}) {
				defaultExpiredTime = 10 * time.Millisecond
				defaultSetup(key, val)
				time.Sleep(20 * time.Millisecond)
			},
		},
	}

	for _, tc := range testCases {
		switch {
		case tc.Setup != nil:
			tc.Setup(tc.Key, tc.Data)
		default:
			defaultSetup(tc.Key, tc.Data)
		}

		res, err := cs.cache.Get(tc.Key)
		if err != nil {
			cs.Require().Equal(tc.Err, err, tc.Desc)
			return
		}
		cs.Require().Equal(tc.Data, res, tc.Desc)
	}
}
