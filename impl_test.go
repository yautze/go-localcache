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
		Setup  func()
	}{
		{
			Desc:   "set value when use generally",
			Key:    "use generally",
			Expect: "case1",
			Setup: func() {
				cs.cache.Set("use generally", "case1")
			},
		},
		{
			Desc:   "set value when over write",
			Key:    "over write",
			Expect: "case2",
			Setup: func() {
				cs.cache.Set("over write", "case2")
			},
		},
	}

	for _, tc := range testCases {
		if tc.Setup != nil {
			tc.Setup()
		}

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
			Err:  nil,
			Data: "string",
		},
		{
			Desc: "int type",
			Key:  "int",
			Err:  nil,
			Data: 1,
		},
		{
			Desc: "bool type",
			Key:  "bool",
			Err:  nil,
			Data: true,
		},
		{
			Desc: "array type",
			Key:  "array",
			Err:  nil,
			Data: []int{1, 2, 3, 4},
		},
		{
			Desc: "object type",
			Key:  "object",
			Err:  nil,
			Data: map[string]string{"test": "test"},
		},
		{
			Desc: "data not found error",
			Key:  "not found",
			Err:  ErrDataNotFound,
			Data: nil,
			Setup: func(key string, val interface{}) {
			},
		},
		{
			Desc: "data not found error cause by cache expired",
			Key:  "cache expired",
			Err:  ErrDataNotFound,
			Data: nil,
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
		cs.Require().Equal(tc.Data, res, tc.Desc)
		cs.Require().Equal(tc.Err, err, tc.Desc)
	}
}
