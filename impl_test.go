package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite
	cache Cache
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (cs *CacheSuite) SetupTest() {
	cs.cache = New()
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

		res, _ := cs.cache.Get(tc.Key)
		cs.Require().Equal(tc.Expect, res, tc.Desc)
	}
}

func (cs *CacheSuite) TestGet() {
	testCases := []struct {
		Desc string
		Key  string
		Data interface{}
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
	}

	for _, tc := range testCases {
		cs.cache.Set(tc.Key, tc.Data)

		res, _ := cs.cache.Get(tc.Key)
		cs.Require().Equal(tc.Data, res, tc.Desc)
	}
}

func (cs *CacheSuite) TestDataNotFound() {
	key := "Not Found"
	res, _ := cs.cache.Get(key)
	cs.Require().Equal(nil, res, "Not Found")
}

func (cs *CacheSuite) TestCacheExpired() {
	defaultExpiredTime = 10 * time.Millisecond
	key := "test"

	mock_cache := New()
	mock_cache.Set(key, 1)
	time.Sleep(20 * time.Millisecond)

	res, _ := mock_cache.Get(key)
	cs.Require().Equal(nil, res, "cache Expired")
}
