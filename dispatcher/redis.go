/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package dispatcher

import redis "github.com/go-redis/redis/v7"

func AddValueToRedisSet(redisClient *redis.Client, setKey, value string) (isSetValueunique bool) {

	if redisClient.SAdd(setKey, value).Val() == 1 {
		isSetValueunique = true
	}

	return
}

func GetValueFromRedisHash(redisClient *redis.Client, hashPath, keyField string) (hashFieldValue string) {

	hashFieldValue = redisClient.HGet(hashPath, keyField).Val()

	return

}

func RedisAddHash(client *redis.Client, id, hashName string, s map[string]interface{}) {
	client.HMSet(hashName+id, s)
}

func RedisUpdateHashField(client *redis.Client, hashID, hashFieldName, hashFieldValue string) {

	client.HSet(hashID, hashFieldName, hashFieldValue)

}

func RedisGetHashFieldValue(client *redis.Client, hashID, hashFieldName string) (hashFieldValue string) {

	hashFieldValue = client.HGet(hashID, hashFieldValue).Val()

	return
}

func SetRedisKeyValue(redisClient *redis.Client, key, value string) {

	err := redisClient.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}

}
