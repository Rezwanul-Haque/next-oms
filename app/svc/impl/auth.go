package impl

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"next-oms/app/domain"
	"next-oms/app/repository"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/app/utils/consts"
	"next-oms/app/utils/methodsutil"
	"next-oms/app/utils/msgutil"
	"next-oms/infra/config"
	"next-oms/infra/conn/cache"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
	"strconv"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	ctx   context.Context
	lc    logger.LogClient
	urepo repository.IUsers
	tSvc  svc.IToken
}

func NewAuthService(ctx context.Context, lc logger.LogClient, urepo repository.IUsers, tokenSvc svc.IToken) svc.IAuth {
	return &auth{
		ctx:   ctx,
		lc:    lc,
		urepo: urepo,
		tSvc:  tokenSvc,
	}
}

func (as *auth) Login(req *serializers.LoginReq) (*serializers.LoginResp, error) {
	var user *domain.User
	var err error

	if user, err = as.urepo.GetUserByEmail(req.Email); err != nil {
		return nil, errors.ErrInvalidEmail
	}

	loginPass := []byte(req.Password)
	hashedPass := []byte(*user.Password)

	if err = bcrypt.CompareHashAndPassword(hashedPass, loginPass); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrInvalidPassword
	}

	var token *serializers.JwtToken

	if token, err = as.tSvc.CreateToken(user.ID); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	if err = as.tSvc.StoreTokenUuid(user.ID, token); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrStoreTokenUuid
	}

	if err = as.urepo.SetLastLoginAt(user); err != nil {
		as.lc.Error("error occur when trying to set last login", err)
		return nil, errors.ErrUpdateLastLogin
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(user.ID, false); err != nil {
		return nil, err
	}

	res := &serializers.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    token.AccessExpiry,
		User:         userResp,
	}
	return res, nil
}

func (as *auth) Logout(user *serializers.LoggedInUser) error {
	return as.tSvc.DeleteTokenUuid(
		config.Cache().Redis.AccessUuidPrefix+user.AccessUuid,
		config.Cache().Redis.RefreshUuidPrefix+user.RefreshUuid,
	)
}

func (as *auth) RefreshToken(refreshToken string) (*serializers.LoginResp, error) {
	oldToken, err := as.parseToken(refreshToken, consts.RefreshTokenType)
	if err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	if !as.userBelongsToTokenUuid(int(oldToken.UserID), oldToken.RefreshUuid, consts.RefreshTokenType) {
		return nil, errors.ErrInvalidRefreshToken
	}

	var newToken *serializers.JwtToken

	if newToken, err = as.tSvc.CreateToken(oldToken.UserID); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	if err = as.tSvc.DeleteTokenUuid(
		config.Cache().Redis.AccessUuidPrefix+oldToken.AccessUuid,
		config.Cache().Redis.RefreshUuidPrefix+oldToken.RefreshUuid,
	); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrDeleteOldTokenUuid
	}

	if err = as.tSvc.StoreTokenUuid(newToken.UserID, newToken); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrStoreTokenUuid
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(newToken.UserID, false); err != nil {
		return nil, err
	}

	res := &serializers.LoginResp{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		User:         userResp,
	}

	return res, nil
}

func (as *auth) VerifyToken(accessToken string) (*serializers.VerifyTokenResp, error) {
	token, err := as.parseToken(accessToken, consts.AccessTokenType)
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}

	if !as.userBelongsToTokenUuid(int(token.UserID), token.AccessUuid, consts.AccessTokenType) {
		return nil, errors.ErrInvalidAccessToken
	}

	var resp *serializers.VerifyTokenResp

	if resp, err = as.getTokenResponse(token); err != nil {
		return nil, err
	}

	return resp, nil
}

func (as *auth) getUserInfoWithParam(userID uint, checkInCache bool) (*serializers.UserWithParamsResp, error) {
	userResp := &serializers.UserResp{}
	userWithParams := serializers.UserWithParamsResp{}

	userCacheKey := config.Cache().Redis.UserPrefix + strconv.Itoa(int(userID))
	var err error

	if checkInCache {
		if err = cache.Client().GetStruct(as.ctx, userCacheKey, &userResp); err == nil {
			as.lc.Info("User served from cache")
			return nil, nil
		}

		as.lc.Error(err.Error(), err)
	}

	user, getErr := as.urepo.GetUserByID(userID)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodsutil.StructToStruct(user, &userWithParams)
	if err != nil {
		as.lc.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	if err := cache.Client().Set(as.ctx, userCacheKey, userWithParams, 0); err != nil {
		as.lc.Error("setting user data on redis key", err)
	}

	return &userWithParams, nil
}

func (as *auth) parseToken(token, tokenType string) (*serializers.JwtToken, error) {
	claims, err := as.parseTokenClaim(token, tokenType)
	if err != nil {
		return nil, err
	}

	tokenDetails := &serializers.JwtToken{}

	if err := methodsutil.MapToStruct(claims, &tokenDetails); err != nil {
		as.lc.Error(err.Error(), err)
		return nil, err
	}

	if tokenDetails.UserID == 0 || tokenDetails.AccessUuid == "" || tokenDetails.RefreshUuid == "" {
		as.lc.Info(fmt.Sprintf("%v", claims))
		return nil, errors.ErrInvalidRefreshToken
	}

	return tokenDetails, nil
}

func (as *auth) parseTokenClaim(token, tokenType string) (jwt.MapClaims, error) {
	secret := config.Jwt().AccessTokenSecret

	if tokenType == consts.RefreshTokenType {
		secret = config.Jwt().RefreshTokenSecret
	}

	parsedToken, err := methodsutil.ParseJwtToken(token, secret)
	if err != nil {
		as.lc.Error(err.Error(), err)
		return nil, errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
		return nil, errors.ErrInvalidAccessToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	return claims, nil
}

func (as *auth) getTokenResponse(token *serializers.JwtToken) (*serializers.VerifyTokenResp, error) {
	var resp *serializers.VerifyTokenResp
	var err error
	tokenCacheKey := config.Cache().Redis.TokenPrefix + strconv.Itoa(int(token.UserID))

	if err = cache.Client().GetStruct(as.ctx, tokenCacheKey, &resp); err == nil {
		as.lc.Info("Token user served from cache")
		return resp, nil
	}

	as.lc.Error(err.Error(), err)

	user, getErr := as.urepo.GetTokenUser(token.UserID)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodsutil.StructToStruct(user, &resp)
	if err != nil {
		as.lc.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user to verify token response"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	if err := cache.Client().Set(as.ctx, tokenCacheKey, resp, 0); err != nil {
		as.lc.Error("setting user data on redis key", err)
	}

	return resp, err
}

func (as *auth) userBelongsToTokenUuid(userID int, uuid, uuidType string) bool {
	prefix := config.Cache().Redis.AccessUuidPrefix

	if uuidType == consts.RefreshTokenType {
		prefix = config.Cache().Redis.RefreshUuidPrefix
	}

	redisKey := prefix + uuid

	redisUserId, err := cache.Client().GetInt(as.ctx, redisKey)
	if err != nil {
		switch err {
		case redis.Nil:
			as.lc.Error(redisKey, errors.NewError(" not found in redis"))
		default:
			as.lc.Error(err.Error(), err)
		}
		return false
	}

	if userID != redisUserId {
		return false
	}

	return true
}
