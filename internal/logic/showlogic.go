// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"database/sql"
	"errors"
	"short/internal/svc"
	"short/internal/types"
	"short/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Show 根据短链接查看对应的长链接
func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 1.参数校验
	shortlinkstr := req.ShortUrl

	// TODO： 优化添加Redis缓存
	res, err := l.svcCtx.RedisClient.GetCtx(l.ctx, shortlinkstr)
	if err != nil { //redis出现错误
		l.Logger.Errorf("redis get err: %v", err)
	}

	if len(res) != 0 { //成功查找缓存，且不为空
		return &types.ShowResponse{LongUrl: res}, nil
	}
	// 2.到数据库中查找
	shorturlMap, err := l.svcCtx.UrlModel.FindOneBySurl(l.ctx, sql.NullString{
		String: shortlinkstr,
		Valid:  true,
	})

	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("404 not found")
		}
		return nil, err
	}

	// 查找成功，记录redis
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, shorturlMap.Surl.String, shorturlMap.Lurl.String)
	if err != nil {
		l.Logger.Errorf("redis set err: %v", err)
	}

	// 3. 返回响应
	resp = &types.ShowResponse{
		LongUrl: shorturlMap.Lurl.String,
	}
	return resp, nil
}
