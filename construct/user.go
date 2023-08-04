package construct

import (
	"app/config"
	"app/database"
	"app/models"
	"app/utils"
	"fmt"
	"strconv"
)

func UserJSON(user_id, operator_id uint) models.UserJSON {
	//查询用户
	u := new(models.User)
	database.Handler.Where("id = ?", user_id).First(u)
	//查询发布的视频
	var videos []models.Video
	database.Handler.Where("user_id = ?", user_id).Find(&videos)
	//计算获赞数量
	like_sum := 0
	for _, v := range videos {
		like_sum += int(utils.GetVideoLikeCount(v.Id))
	}
	return models.UserJSON{
		Avatar:          fmt.Sprintf("%s:%d/public/avatar/%d.png", config.Server.Host, config.Server.Port, u.Id),
		BackgroundImage: fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
		FavoriteCount:   utils.GetUserLikedVideoCount(user_id),
		FollowCount:     utils.GetFollowCount(u.Id),
		FollowerCount:   utils.GetFollowedCount(u.Id),
		ID:              u.Id,
		IsFollow:        utils.IsFollowed(operator_id, u.Id),
		Name:            u.Username,
		Signature:       "Hello World!",
		TotalFavorited:  strconv.Itoa(like_sum),
		WorkCount:       int64(len(videos)),
	}
}
