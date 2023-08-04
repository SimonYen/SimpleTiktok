package construct

import (
	"app/config"
	"app/database"
	"app/models"
	"app/utils"
	"fmt"
)

func VideoJSON(video_id, operator_id uint) models.VideoJSON {
	//查询视频信息
	v := new(models.Video)
	database.Handler.Where("id = ?", video_id).First(v)
	//卧槽，真的清爽多了
	return models.VideoJSON{
		ID:            v.Id,
		Author:        UserJSON(v.UserID, operator_id),
		PlayURL:       fmt.Sprintf("%s:%d/public/video/%d%s", config.Server.Host, config.Server.Port, v.Id, v.Extension),
		CoverURL:      fmt.Sprintf("%s:%d/public/screenshot/%d.png", config.Server.Host, config.Server.Port, v.Id),
		FavoriteCount: utils.GetVideoLikeCount(v.Id),
		CommentCount:  0,
		IsFavorite:    utils.VideoIsLiked(v.Id, operator_id),
		Title:         v.Title,
	}
}
