package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
)

/**
用Redis实时维护商品库存
*/

const (
	INVENTORY_PREFIX = "gift_count_" //所有key设置统一的前缀，方便后续按前缀遍历key
)

// 从Mysql中读出所有奖品的初始库存，存入Redis。如果同时有很多用户来参与抽奖活动，不能并发去Mysql里减库存，mysql扛不住这么高的并发量，Redis可以扛住
func InitGiftInventory() {
	for _, gift := range GetAllGifts() {
		if gift.Count <= 0 {
			slog.Warn("gift count is zero", "id", gift.Id, "name", gift.Name)
			continue //没有库存的商品不参与抽奖
		}
		err := GiftRedis.Set(context.Background(), INVENTORY_PREFIX+strconv.Itoa(gift.Id), gift.Count, 0).Err()
		if err != nil {
			slog.Error("set gift count to redis failed", "gift id", gift.Id, "error", err)
		}
	}
}

// 获取所有奖品剩余的库存量
func GetAllGiftInventory() []*Gift {
	keys, err := GiftRedis.Keys(context.Background(), INVENTORY_PREFIX+"*").Result() //根据前缀获取所有奖品的key
	if err != nil {
		slog.Error("iterate all gift keys failed", "error", err)
		return nil
	}
	gifts := make([]*Gift, 0, len(keys))
	for _, key := range keys { //根据奖品key获得奖品的库存count
		if id, err := strconv.Atoi(key[len(INVENTORY_PREFIX):]); err == nil {
			count, err := GiftRedis.Get(context.Background(), key).Int()
			if err == nil {
				gifts = append(gifts, &Gift{Id: id, Count: count})
			} else {
				slog.Error("gift count is not int", "key", key)
			}
		} else {
			slog.Error("gift id is not int", "key", key)
		}
	}

	return gifts
}

// 获取特定奖品剩余的库存量
func GetGiftInventory(GiftId int) int {
	key := INVENTORY_PREFIX + strconv.Itoa(GiftId)
	count, err := GiftRedis.Get(context.Background(), key).Int()
	if err == nil {
		return count
	} else {
		slog.Error("gift count is not int", "key", key)
		return -1
	}
}

// 奖品对应的库存减1
func ReduceInventory(GiftId int) error {
	key := INVENTORY_PREFIX + strconv.Itoa(GiftId)
	n, err := GiftRedis.Decr(context.Background(), key).Result() //原子操作。返回减1之后的值。如果key不存在则返回-1
	if err != nil {
		slog.Error("decr key failed", "key", key, "error", err)
		return err
	} else {
		if n < 0 {
			msg := fmt.Sprintf("%d已无库存，减1失败", GiftId)
			slog.Error(msg)
			return errors.New(msg)
		}
		return nil
	}
}

// 奖品对应的库存加1
func IncreaseInventory(GiftId int) error {
	key := INVENTORY_PREFIX + strconv.Itoa(GiftId)
	_, err := GiftRedis.Incr(context.Background(), key).Result() //原子加1
	if err != nil {
		slog.Error("incr key failed", "key", key, "error", err)
		return err
	} else {
		return nil
	}
}
