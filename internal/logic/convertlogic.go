// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"database/sql"
	"errors"
	"short/model"
	"short/pkg/base62"

	"short/internal/svc"
	"short/internal/types"
	"short/pkg/connect"
	"short/pkg/murmur3"
	"short/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 长链接转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 参数验证，验证长链接是否有效
	// 1.长链接能否ping通
	ok, err := connect.CheckURL(req.LongUrl)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("URL 无效")
	}

	// 2.验证长链接是否已经有转换后的短链接
	// 2.1 murmur3 hash 值
	murValue := murmur3.Hash(req.LongUrl)

	// 2.2 使用hash值 和 longUrl 查询数据库是否已经存在
	shortUrlMap, err := l.svcCtx.UrlModel.FindByHashAndLurl(l.ctx, req.LongUrl, murValue)

	if err == nil { //说明查找成功，可以直接返回对应的短链接
		return &types.ConvertResponse{ShortUrl: l.svcCtx.Config.Domain + shortUrlMap.Surl.String}, nil
	}

	if err != sqlx.ErrNotFound {
		logx.Errorw("ShortUrlMapModel.FindByHashAndLurl failed", logx.Field("err", err.Error()))
		return nil, err
	}

	// 3.避免循环转链 (输入的不能是一个已有的短链接)
	basePath, err := urltool.GetBaseUrl(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBaseUrl failed", logx.Field("err", err.Error()))
		return nil, err
	}
	_, err = l.svcCtx.UrlModel.FindOneBySurl(context.Background(), sql.NullString{String: basePath, Valid: true})
	if err == nil { //找到了，说明输入的是一个已有的而短链接，决绝转换
		return nil, errors.New("输入的链接已经是短链接，不能循环转换")
	}

	if err != sqlx.ErrNotFound {
		return nil, errors.New("查询短链接 数据库操作失败")
	}
	// 从发号器表获取一个号，生成短链接
	shortlinkString := ""

	// 4. 生成短链
	for {
		// 1. 从表中使用 replace into 语句更新记录，使用自增的id作为短链接的号
		id, err := l.svcCtx.Sequence.Next(context.Background())
		if err != nil {
			logx.Errorw("Sequence.Next failed", logx.Field("err", err.Error()))
			return nil, err
		}
		// 2.将id转换为62进制数作为短链部分
		shortlinkString = base62.Uint2string(id)

		// 3. 对短链字符串进行一些特殊词汇筛选
		if _, ok := l.svcCtx.BlackMap[shortlinkString]; !ok { //没有特殊词汇，则生成的短链符合要求，跳出循环
			break
		}
	}

	// 5. 将长链接和短链接的对应关系写入数据库
	lsurl := &model.ShortUrlMap{
		Lurl: sql.NullString{
			String: req.LongUrl,
			Valid:  true,
		},
		Surl: sql.NullString{
			String: shortlinkString,
			Valid:  true,
		},
		LurlHash: sql.NullInt64{
			Int64: int64(murValue),
			Valid: true,
		},
	}
	_, err = l.svcCtx.UrlModel.Insert(context.Background(), lsurl)
	if err != nil {
		logx.Errorw("UrlModel.Insert failed", logx.Field("err", err.Error()))
		return nil, err
	}

	// 生成的短链接写入布隆过滤器
	if err := l.svcCtx.LinkBloom.Add([]byte(shortlinkString)); err != nil {
		l.Errorf("短链接写入布隆过滤器失败.LinkBloom.Add(%s) failed", shortlinkString)
	}
	l.Infof("短链接写入布隆过滤器成功，%s", shortlinkString)

	// 返回响应, 拼接短域名和短链
	shorturl := l.svcCtx.Config.Domain + shortlinkString
	resp = &types.ConvertResponse{
		ShortUrl: shorturl,
	}
	return resp, nil
}
