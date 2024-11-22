package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
	"time"
	"ttv-bot/log"
	"ttv-bot/model"
)

var cacheInvalidation = 5 * time.Minute

type stonFiPoolsCache struct {
	m           []model.Pool
	lastCleared time.Time
	sync.Mutex
}

func newStonFiPoolsCache() *stonFiPoolsCache {
	return &stonFiPoolsCache{
		m:           []model.Pool{},
		lastCleared: time.Time{},
	}
}

func (c *stonFiPoolsCache) get(msg api.Chattable, queryFunc func(msg api.Chattable) ([]model.Pool, error)) []model.Pool {
	c.Lock()
	defer c.Unlock()

	if time.Since(c.lastCleared) >= cacheInvalidation {
		pools, err := queryFunc(msg)
		if err != nil {
			log.Info("query ston.fi pools failed", zap.Error(err))
			return nil
		}
		c.m = pools
		c.lastCleared = time.Now() // 更新缓存清除时间
	}

	return c.m
}

type topHoldersCache struct {
	m           map[string][]model.Jetton
	lastCleared map[string]time.Time
	sync.Mutex
}

func newTopHoldersCache() *topHoldersCache {
	return &topHoldersCache{
		m:           make(map[string][]model.Jetton),
		lastCleared: make(map[string]time.Time),
	}
}

func (c *topHoldersCache) get(msg api.Chattable, addr *model.TonAddr, n int, queryFunc func(addr *model.TonAddr, n int, msg api.Chattable) ([]model.Jetton, error)) []model.Jetton {
	c.Lock()
	defer c.Unlock()
	if v, exists := c.lastCleared[addr.Hex]; !exists || time.Since(v) >= cacheInvalidation {
		jettons, err := queryFunc(addr, n, msg)
		if err != nil {
			log.Info("query topN holders failed", zap.Error(err))
			return nil
		}
		if jettons == nil {
			if _, ok := c.m[addr.Hex]; ok {
				delete(c.m, addr.Hex)
				delete(c.lastCleared, addr.Hex)
			}
			return nil
		}
		c.m[addr.Hex] = jettons
		c.lastCleared[addr.Hex] = time.Now()
	}

	return c.m[addr.Hex]
}

type nftItemsCache struct {
	m           map[string][]model.NFT
	lastCleared map[string]time.Time
	sync.Mutex
}

func newNFTItemsCache() *nftItemsCache {
	return &nftItemsCache{
		m:           make(map[string][]model.NFT),
		lastCleared: make(map[string]time.Time),
	}
}

func (c *nftItemsCache) get(msg api.Chattable, addr *model.TonAddr, limit int, queryFunc func(addr *model.TonAddr, limit int, msg api.Chattable) ([]model.NFT, error)) []model.NFT {
	c.Lock()
	defer c.Unlock()
	if v, exists := c.lastCleared[addr.Hex]; !exists || time.Since(v) >= cacheInvalidation {
		nfts, err := queryFunc(addr, limit, msg)
		if err != nil {
			log.Info("query nft items failed", zap.Error(err))
			return nil
		}
		if nfts == nil {
			if _, ok := c.m[addr.Hex]; ok {
				delete(c.m, addr.Hex)
				delete(c.lastCleared, addr.Hex)
			}
			return nil
		}
		c.m[addr.Hex] = nfts
		c.lastCleared[addr.Hex] = time.Now()
	}

	return c.m[addr.Hex]
}

type nftAssetsCache struct {
	m           map[string][]model.NFT
	lastCleared map[string]time.Time
	sync.Mutex
}

func newNFTAssetsCache() *nftAssetsCache {
	return &nftAssetsCache{
		m:           make(map[string][]model.NFT),
		lastCleared: make(map[string]time.Time),
	}
}

func (c *nftAssetsCache) get(msg api.Chattable, addr *model.TonAddr, limit int, queryFunc func(addr *model.TonAddr, limit int, msg api.Chattable) ([]model.NFT, error)) []model.NFT {
	c.Lock()
	defer c.Unlock()
	if v, exists := c.lastCleared[addr.Hex]; !exists || time.Since(v) >= cacheInvalidation {
		nfts, err := queryFunc(addr, limit, msg)
		if err != nil {
			log.Info("query nft assets failed", zap.Error(err))
			return nil
		}
		if nfts == nil {
			if _, ok := c.m[addr.Hex]; ok {
				delete(c.m, addr.Hex)
				delete(c.lastCleared, addr.Hex)
			}
			return nil
		}
		c.m[addr.Hex] = nfts
		c.lastCleared[addr.Hex] = time.Now()
	}

	return c.m[addr.Hex]
}

type nftCollectionCache struct {
	m           map[string]*model.NFT
	lastCleared map[string]time.Time
	sync.Mutex
}

func newNFTCollectionCache() *nftCollectionCache {
	return &nftCollectionCache{
		m:           make(map[string]*model.NFT),
		lastCleared: make(map[string]time.Time),
	}
}
func (c *nftCollectionCache) get(msg api.Chattable, addr *model.TonAddr, queryFunc func(addr *model.TonAddr, msg api.Chattable) (*model.NFT, error)) *model.NFT {
	c.Lock()
	defer c.Unlock()
	if v, exists := c.lastCleared[addr.Hex]; !exists || time.Since(v) >= cacheInvalidation {
		nft, err := queryFunc(addr, msg)
		if err != nil {
			log.Info("query nft assets failed", zap.Error(err))
			return nil
		}
		if nft == nil {
			if _, ok := c.m[addr.Hex]; ok {
				delete(c.m, addr.Hex)
				delete(c.lastCleared, addr.Hex)
			}
			return nil
		}
		c.m[addr.Hex] = nft
		c.lastCleared[addr.Hex] = time.Now()
	}

	return c.m[addr.Hex]
}
