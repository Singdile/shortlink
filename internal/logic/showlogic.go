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

	// 3. 返回响应
	resp = &types.ShowResponse{
		LongUrl: shorturlMap.Lurl.String,
	}
	return resp, nil
}
