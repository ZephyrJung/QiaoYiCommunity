package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/jwt"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/redis"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/sms"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/wechat"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userDAO     dao.UserDAO
	jwtMgr      *jwt.Manager
	redis       *redis.Client
	smsProvider sms.Provider
	wechatMP    *wechat.MiniProgram
}

func NewAuthService(userDAO dao.UserDAO, jwtMgr *jwt.Manager, redisClient *redis.Client, smsProvider sms.Provider, wechatMP *wechat.MiniProgram) *AuthService {
	return &AuthService{
		userDAO:     userDAO,
		jwtMgr:      jwtMgr,
		redis:       redisClient,
		smsProvider: smsProvider,
		wechatMP:    wechatMP,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginResult struct {
	User  *entity.User `json:"user"`
	Token *TokenPair   `json:"token"`
}

func (s *AuthService) WeChatLogin(ctx context.Context, jsCode string) (*LoginResult, error) {
	session, err := s.wechatMP.JsCode2Session(jsCode)
	if err != nil {
		return nil, err
	}

	var user *entity.User
	if session.UnionID != "" {
		user, _ = s.userDAO.GetByUnionID(ctx, session.UnionID)
	}
	if user == nil && session.OpenID != "" {
		user, _ = s.userDAO.GetByOpenID(ctx, session.OpenID)
	}

	if user == nil {
		newUser := &entity.User{
			ID:       utils.GenerateID(),
			UnionID:  strPtr(session.UnionID),
			OpenID:   strPtr(session.OpenID),
			Nickname: "微信用户",
			Role:     int8(entity.UserRoleResident),
			Status:   int8(entity.UserStatusActive),
		}
		if err := s.userDAO.Create(ctx, newUser); err != nil {
			return nil, fmt.Errorf("create user failed: %w", err)
		}
		user = newUser
	}

	if user.Status == int8(entity.UserStatusBanned) {
		return nil, errors.New("user banned")
	}

	tp, err := s.generateAndStoreTokens(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: user, Token: tp}, nil
}

func (s *AuthService) PhoneSendCode(ctx context.Context, phone string) error {
	if !utils.ValidatePhone(phone) {
		return errors.New("invalid phone number")
	}

	code := sms.GenerateCode()
	key := redis.SMSCodeKey(phone)
	if err := s.redis.Set(ctx, key, code, 5*time.Minute); err != nil {
		return fmt.Errorf("store sms code failed: %w", err)
	}

	return s.smsProvider.SendCode(ctx, phone, code)
}

func (s *AuthService) PhoneLogin(ctx context.Context, phone, code string) (*LoginResult, error) {
	if !utils.ValidatePhone(phone) {
		return nil, errors.New("invalid phone number")
	}

	key := redis.SMSCodeKey(phone)
	stored, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, errors.New("invalid or expired code")
	}
	if stored != code {
		return nil, errors.New("invalid code")
	}

	user, err := s.userDAO.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &entity.User{
				ID:       utils.GenerateID(),
				Phone:    strPtr(phone),
				Nickname: fmt.Sprintf("用户%s", phone[len(phone)-4:]),
				Role:     int8(entity.UserRoleResident),
				Status:   int8(entity.UserStatusActive),
			}
			if err := s.userDAO.Create(ctx, newUser); err != nil {
				return nil, fmt.Errorf("create user failed: %w", err)
			}
			user = newUser
		} else {
			return nil, err
		}
	}

	if user.Status == int8(entity.UserStatusBanned) {
		return nil, errors.New("user banned")
	}

	_ = s.redis.Del(ctx, key)

	tp, err := s.generateAndStoreTokens(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: user, Token: tp}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtMgr.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	key := redis.RefreshTokenKey(claims.UserID, claims.TokenID)
	_, err = s.redis.Get(ctx, key)
	if err != nil {
		return nil, errors.New("refresh token revoked or expired")
	}

	_ = s.redis.Del(ctx, key)

	tp, err := s.generateAndStoreTokens(ctx, claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return tp, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64, tokenID string) error {
	key := redis.RefreshTokenKey(userID, tokenID)
	return s.redis.Del(ctx, key)
}

func (s *AuthService) generateAndStoreTokens(ctx context.Context, userID int64, role int8) (*TokenPair, error) {
	accessToken, refreshToken, tokenID, err := s.jwtMgr.GenerateTokenPair(userID, role)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	key := redis.RefreshTokenKey(userID, tokenID)
	if err := s.redis.Set(ctx, key, tokenID, 30*24*time.Hour); err != nil {
		return nil, fmt.Errorf("store refresh token failed: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    2 * 3600,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
