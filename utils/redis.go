package utils

import (
	"app/database"
	"fmt"
	"strconv"
)

/*
目前redis中存在以下数据结构：
	like:video_id	->set结构，用于存储喜欢该视频的user_id
	liked:user_id	->set结构，用于存储该user_id喜欢的video_id
	follow:user_id ->set结构，用于存储该user_id主动关注的用户id
	followed:user_id ->set结构，用于关注user_id的用户id（也就是被别人关注）


那么借助redis，可以方便快捷实现以下操作：

	用户点赞视频: sadd like:v%d u%d ; sadd liked:u%d v%d
	用户取消点赞视频：srem like:v%d u%d ; srem liked:u%d v%d
	视频是否被某用户点赞：sismember like:v%d u%d 或者 sismember liked:u%d v%d
	获取视频被点赞的数量：scard like:v%d
	获取用户喜欢的视频数量：scard liked:u%d
	用户喜欢的视频列表：smembers liked:u%d

	u1关注u2：sadd follow:u1 u2 ; sadd followed:u2 u1
	u1取消关注u2：srem follow:u1 u2 ; srem followed:u2 u1
	u1是否关注u2：sismember followed:u2 u1
	用户的关注列表：smembers follow:user_id
	用户的粉丝列表：smembers followed:user_id
	用户的好友列表：sinter follow:user_id followed:user_id （好友的定义是互相关注，所以取交集就行了)

*/

// 获取视频点赞数量
func GetVideoLikeCount(video_id uint) int64 {
	key := fmt.Sprintf("like:%d", video_id)
	res, err := database.RDB.SCard(database.CTX, key).Result()
	if err != nil {
		return 0
	}
	return res
}

// 获取用户喜欢的视频数量
func GetUserLikedVideoCount(user_id uint) int64 {
	key := fmt.Sprintf("liked:%d", user_id)
	res, err := database.RDB.SCard(database.CTX, key).Result()
	if err != nil {
		return 0
	}
	return res
}

// 获取用户喜欢的视频列表，返回的是video_id
func GetUserLikedVideoList(user_id uint) []string {
	key := fmt.Sprintf("liked:%d", user_id)
	res, err := database.RDB.SMembers(database.CTX, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// 视频赞操作
func LikeOrDisLikeVideo(video_id uint, user_id uint, action_type int) bool {
	key1 := fmt.Sprintf("like:%d", video_id)
	key2 := fmt.Sprintf("liked:%d", user_id)
	res := true
	if action_type == 1 {
		_, err := database.RDB.SAdd(database.CTX, key1, strconv.Itoa(int(user_id))).Result()
		if err != nil {
			res = false
		}
		_, err = database.RDB.SAdd(database.CTX, key2, strconv.Itoa(int(video_id))).Result()
		if err != nil {
			res = false
		}
	} else {
		_, err := database.RDB.SRem(database.CTX, key1, strconv.Itoa(int(user_id))).Result()
		if err != nil {
			res = false
		}
		_, err = database.RDB.SRem(database.CTX, key2, strconv.Itoa(int(video_id))).Result()
		if err != nil {
			res = false
		}
	}
	return res
}

// 是否点赞
func VideoIsLiked(video_id, user_id uint) bool {
	key := fmt.Sprintf("like:%d", video_id)
	res, _ := database.RDB.SIsMember(database.CTX, key, strconv.Itoa(int(user_id))).Result()
	return res
}

// 用户关注操作
func FollowOrUnfollowUser(follower, followee uint, action_type int) bool {
	key1 := fmt.Sprintf("follow:%d", follower)
	key2 := fmt.Sprintf("followed:%d", followee)
	res := true
	if action_type == 1 {
		_, err := database.RDB.SAdd(database.CTX, key1, followee).Result()
		if err != nil {
			res = false
		}
		_, err = database.RDB.SAdd(database.CTX, key2, follower).Result()
		if err != nil {
			res = false
		}
	} else {
		_, err := database.RDB.SRem(database.CTX, key1, followee).Result()
		if err != nil {
			res = false
		}
		_, err = database.RDB.SRem(database.CTX, key2, follower).Result()
		if err != nil {
			res = false
		}
	}
	return res
}

// u1是否关注u2
func IsFollowed(u1, u2 uint) bool {
	key := fmt.Sprintf("followed:%d", u2)
	res, _ := database.RDB.SIsMember(database.CTX, key, strconv.Itoa(int(u1))).Result()
	return res
}

// 用户关注列表
func GetFollowList(user_id uint) []string {
	key := fmt.Sprintf("follow:%d", user_id)
	res, err := database.RDB.SMembers(database.CTX, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// 用户粉丝列表
func GetFollowedList(user_id uint) []string {
	key := fmt.Sprintf("followed:%d", user_id)
	res, err := database.RDB.SMembers(database.CTX, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// 用户好友列表
func GetFriendList(user_id uint) []string {
	key1 := fmt.Sprintf("follow:%d", user_id)
	key2 := fmt.Sprintf("followed:%d", user_id)
	res, err := database.RDB.SInter(database.CTX, key1, key2).Result()
	if err != nil {
		return nil
	}
	return res
}

// 用户的关注数量
func GetFollowCount(user_id uint) int64 {
	key := fmt.Sprintf("follow:%d", user_id)
	res, err := database.RDB.SCard(database.CTX, key).Result()
	if err != nil {
		return 0
	}
	return res
}

// 用户的粉丝数量
func GetFollowedCount(user_id uint) int64 {
	key := fmt.Sprintf("followed:%d", user_id)
	res, err := database.RDB.SCard(database.CTX, key).Result()
	if err != nil {
		return 0
	}
	return res
}
