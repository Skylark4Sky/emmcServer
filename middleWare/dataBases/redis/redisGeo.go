package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

/**
GEO 地理位置
*/

// GeoOptions 用于GEORADIUS和GEORADIUSBYMEMBER命令的参数
type GeoOptions struct {
	WithCoord bool
	WithDist  bool
	WithHash  bool
	Order     string // ASC从近到远，DESC从远到近
	Count     int
}

// GeoResult 用于GEORADIUS和GEORADIUSBYMEMBER命令的查询结果
type GeoResult struct {
	Name      string
	Longitude float64
	Latitude  float64
	Dist      float64
	Hash      int64
}

// GeoAdd 将给定的空间元素（纬度、经度、名字）添加到指定的键里面，这些数据会以有序集合的形式被储存在键里面，所以删除可以使用`ZREM`。
func (c *Cacher) GeoAdd(key string, longitude, latitude float64, member string) error {
	_, err := redis.Int(c.Do("GEOADD", key, longitude, latitude, member))
	return err
}

// GeoPos 从键里面返回所有给定位置元素的位置（经度和纬度）。
func (c *Cacher) GeoPos(key string, members ...interface{}) ([]*[2]float64, error) {
	args := redis.Args{}
	args = args.Add(key)
	args = args.Add(members...)
	return redis.Positions(c.Do("GEOPOS", args...))
}

// GeoDist 返回两个给定位置之间的距离。
// 如果两个位置之间的其中一个不存在， 那么命令返回空值。
// 指定单位的参数 unit 必须是以下单位的其中一个：
// m 表示单位为米。
// km 表示单位为千米。
// mi 表示单位为英里。
// ft 表示单位为英尺。
// 如果用户没有显式地指定单位参数， 那么 GEODIST 默认使用米作为单位。
func (c *Cacher) GeoDist(key string, member1, member2, unit string) (float64, error) {
	_, err := redis.Float64(c.Do("GEODIST", key, member1, member2, unit))
	return 0, err
}

// GeoRadius 以给定的经纬度为中心， 返回键包含的位置元素当中， 与中心的距离不超过给定最大距离的所有位置元素。
func (c *Cacher) GeoRadius(key string, longitude, latitude, radius float64, unit string, options GeoOptions) ([]*GeoResult, error) {
	args := redis.Args{}
	args = args.Add(key, longitude, latitude, radius, unit)
	if options.WithDist {
		args = args.Add("WITHDIST")
	}
	if options.WithCoord {
		args = args.Add("WITHCOORD")
	}
	if options.WithHash {
		args = args.Add("WITHHASH")
	}
	if options.Order != "" {
		args = args.Add(options.Order)
	}
	if options.Count > 0 {
		args = args.Add("Count", options.Count)
	}

	reply, err := c.Do("GEORADIUS", args...)
	return toGeoResult(reply, err, options)
}

// GeoRadiusByMember 这个命令和 GEORADIUS 命令一样， 都可以找出位于指定范围内的元素， 但是 GEORADIUSBYMEMBER 的中心点是由给定的位置元素决定的， 而不是像 GEORADIUS 那样， 使用输入的经度和纬度来决定中心点。
func (c *Cacher) GeoRadiusByMember(key string, member string, radius float64, unit string, options GeoOptions) ([]*GeoResult, error) {
	args := redis.Args{}
	args = args.Add(key, member, radius, unit)
	if options.WithDist {
		args = args.Add("WITHDIST")
	}
	if options.WithCoord {
		args = args.Add("WITHCOORD")
	}
	if options.WithHash {
		args = args.Add("WITHHASH")
	}
	if options.Order != "" {
		args = args.Add(options.Order)
	}
	if options.Count > 0 {
		args = args.Add("Count", options.Count)
	}

	reply, err := c.Do("GEORADIUSBYMEMBER", args...)
	return toGeoResult(reply, err, options)
}

// GeoHash 返回一个或多个位置元素的 Geohash 表示。
func (c *Cacher) GeoHash(key string, members ...interface{}) ([]string, error) {
	args := redis.Args{}
	args = args.Add(key)
	args = args.Add(members...)
	return redis.Strings(c.Do("GEOHASH", args...))
}

func toGeoResult(reply interface{}, err error, options GeoOptions) ([]*GeoResult, error) {
	values, err := redis.Values(reply, err)
	if err != nil {
		return nil, err
	}
	results := make([]*GeoResult, len(values))
	for i := range values {
		if values[i] == nil {
			continue
		}
		p, ok := values[i].([]interface{})
		if !ok {
			return nil, fmt.Errorf("redisgo: unexpected element type for interface slice, got type %T", values[i])
		}
		geoResult := &GeoResult{}
		pos := 0
		name, err := redis.String(p[pos], nil)
		if err != nil {
			return nil, err
		}
		geoResult.Name = name
		if options.WithDist {
			pos = pos + 1
			dist, err := redis.Float64(p[pos], nil)
			if err != nil {
				return nil, err
			}
			geoResult.Dist = dist
		}
		if options.WithHash {
			pos = pos + 1
			hash, err := redis.Int64(p[pos], nil)
			if err != nil {
				return nil, err
			}
			geoResult.Hash = hash
		}
		if options.WithCoord {
			pos = pos + 1
			pp, ok := p[pos].([]interface{})
			if !ok {
				return nil, fmt.Errorf("redisgo: unexpected element type for interface slice, got type %T", p[i])
			}
			if len(pp) > 0 {
				lat, err := redis.Float64(pp[0], nil)
				if err != nil {
					return nil, err
				}
				lon, err := redis.Float64(pp[1], nil)
				if err != nil {
					return nil, err
				}
				geoResult.Latitude = lat
				geoResult.Longitude = lon
			}
		}
		if err != nil {
			return nil, err
		}
		results[i] = geoResult
	}
	return results, nil
}
