package data

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto"
	"askonce/helpers"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/xiangtao94/golib/flow"
	"strconv"
	"time"
)

/*
登陆态缓存核心逻辑
*/
type SessionCache struct {
	flow.Redis
}

var USER_SESSION_SET = "user_session_set:"

// 生成
func (entity *SessionCache) GenSessionValue() string {
	sessionValue := uuid.NewString()
	return sessionValue
}

// 设置session
func (entity *SessionCache) SetSession(session dto.LoginInfoSession, sessionKey string, expire int64) error {
	if expire == 0 {
		expire = flow.EXPIRE_TIME_1_DAY * 30
	}
	// 调整过期时间，避免白天过期
	expire = normalizeExpire(expire)
	// 设置失效时间，避免redis自动失效失败或ttl设置失败导致一直有效的问题
	if session.LoginTime > 0 {
		session.ExpireTime = session.LoginTime + expire
	}
	m := make(map[string]interface{})
	err := mapstructure.Decode(&session, &m)
	if err != nil {
		entity.LogErrorf("SetSession mapstructure.Decode error, uid:%d, sessionKey:%s, err:%s", session.UserId, sessionKey, err.Error())
		return err
	}
	cacheKey := entity.FormatCacheKey("%s:%s", defines.COOKIE_KEY, sessionKey)
	entity.LogDebugf("redisCacheKey:%s", cacheKey)
	err = helpers.RedisClient.HMSet(cacheKey, m)
	if err != nil {
		entity.LogErrorf("SetSession set session error, uid:%d, sessionKey:%s, err:%s", session.UserId, sessionKey, err.Error())
		return err
	}
	_, err = helpers.RedisClient.Expire(cacheKey, expire)
	if err != nil {
		entity.LogErrorf("SetSession set exipire error, uid:%d, sessionKey:%s, err:%s", session.UserId, sessionKey, err.Error())
	}
	return nil
}

// 设置session
func (entity *SessionCache) UpdateSession(sessionKey string, session dto.LoginInfoSession) error {
	cacheKey := entity.FormatCacheKey("%s:%s", defines.COOKIE_KEY, sessionKey)
	entity.LogDebugf("redisCacheKey:%s", cacheKey)
	exist, err := helpers.RedisClient.Exists(cacheKey)
	if err != nil {
		entity.LogErrorf("UpdateSession sessionKey request redis exists cmd error, sessionKey:%s, err:%s", sessionKey, err.Error())
		return err
	}
	if !exist {
		// 不存在，尝试重建
		if session.ExpireTime > time.Now().Unix() {
			return entity.SetSession(session, sessionKey, session.ExpireTime-time.Now().Unix())
		}
		entity.LogWarnf("UpdateSession zjxussSession not found!")
		return nil
	}
	// 存在，更新
	m := make(map[string]interface{})
	err = mapstructure.Decode(&session, &m)
	if err != nil {
		return err
	}
	err = helpers.RedisClient.HMSet(cacheKey, m)
	if err != nil {
		entity.LogErrorf("SetSession request redis error, uid:%d,sessionKey:%s, err:%s", session.UserId, sessionKey, err.Error())
		return err
	}
	return nil
}

// 获取session
func (entity *SessionCache) GetSession(sessionKey string) (dto.LoginInfoSession, error) {
	cacheKey := entity.FormatCacheKey("%s:%s", defines.COOKIE_KEY, sessionKey)
	m := dto.LoginInfoSession{}
	result, err := redis.Values(helpers.RedisClient.Do("hgetall", cacheKey))
	if err != nil {
		entity.LogErrorf("GetSessionByZJUSS hgetall redis error,err%s", err)
		return m, err
	}

	// 未取到session，视同未登录
	if len(result) == 0 {
		return m, components.ErrorNotLogin
	}
	if err = redis.ScanStruct(result, &m); err != nil {
		entity.LogErrorf("GetSessionByZJUSS scanstruct redis error,err%s", err)
		return m, err
	}

	// 二次校验session有效期确认是否有效，已失效视同未登录
	if m.ExpireTime > 0 && m.ExpireTime < time.Now().Unix() {
		return m, components.ErrorNotLogin
	}
	return m, err
}

// 删除session
func (entity *SessionCache) DelSession(sessionKey string) error {
	cacheKey := entity.FormatCacheKey("%s:%s", defines.COOKIE_KEY, sessionKey)
	// session已经不存在，不需要处理
	exist, err := helpers.RedisClient.Exists(cacheKey)
	if err == nil && !exist {
		return nil
	}
	_, err = helpers.RedisClient.Del(entity.GetCtx(), cacheKey)
	if err != nil {
		entity.LogErrorf("DelSession redis error,err%s", err)
	}
	return nil
}

// 过期时间标准化，避免白天过期，凌晨3点失效
func normalizeExpire(expire int64) int64 {
	curTime := time.Now().Unix()
	expireTime := curTime + expire
	normalizeExpireTime := ((expireTime/86400)+1)*86400 - 5*3600
	return normalizeExpireTime - curTime
}

// 添加到uid_zjxuss有序集合，expire为有效期（秒）
// score设置为失效时间
func (entity *SessionCache) AddUserSessionSet(uid string, sessionKey string, expire int64) {
	cacheKey := entity.FormatCacheKey("%s:%s", USER_SESSION_SET, uid)
	setMap := make(map[string]float64)
	setTime := time.Now().Unix() + expire
	setMap[sessionKey] = float64(setTime)
	_, err := helpers.RedisClient.ZAdd(cacheKey, setMap)
	if err != nil {
		entity.LogErrorf("AddUserSessionSet redis error,err%s", err)
	}

}

// 根据sessionKey删除UserSessionKey中特定值
func (entity *SessionCache) DelUserSessionSet(uid string, sessionKey string) error {
	cacheKey := entity.FormatCacheKey("%s:%s", USER_SESSION_SET, uid)
	_, err := helpers.RedisClient.ZRem(cacheKey, sessionKey)
	if err != nil {
		entity.LogErrorf("DelUserSessionSet redis error,err%s", err)
	}
	return err
}

// 添加到uid_zjxuss有序集合，expire为有效期（秒）
// score设置为失效时间
func (entity *SessionCache) GetUserSessionSet(uid string) ([]string, error) {
	cacheKey := entity.FormatCacheKey("%s:%s", USER_SESSION_SET, uid)
	expireTime := time.Now().Unix()
	expireTimeStr := strconv.FormatInt(expireTime, 10)
	// 获取7天内生成的zjxuss
	res, err := helpers.RedisClient.ZRangeByScore(cacheKey, expireTimeStr, "+inf", false, false, 0, 0)
	if err != nil {
		entity.LogErrorf("GetUserSessionSet request redis error, uid:%d, err:%s", uid, err.Error())
		return nil, err
	}
	sessions := make([]string, 0)
	for _, v := range res {
		sessions = append(sessions, string(v))
	}
	return sessions, nil
}
