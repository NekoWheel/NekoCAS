package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	"github.com/unknwon/com"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoCAS/internal/helper"
)

// User 为用户表。
type User struct {
	gorm.Model

	NickName string
	Email    string
	Password string
	Salt     string
	Avatar   string

	IsActive bool
	IsAdmin  bool
}

// EncodePassword 对密码进行加盐处理。
func (u *User) EncodePassword() {
	newPassword := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPassword)
}

// ValidatePassword 检查输入的密码是否正确。
func (u *User) ValidatePassword(password string) bool {
	newUser := &User{Password: password, Salt: u.Salt}
	newUser.EncodePassword()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

// GetActivationCode 返回用户账号激活码，有效期两小时。
func (u *User) GetActivationCode() string {
	code := helper.CreateTimeLimitCode(com.ToStr(u.ID)+u.Email+u.NickName+u.Password+u.Salt, 120, nil)

	// 添加编码后的邮箱信息，用于验证时反查用户信息
	code += hex.EncodeToString([]byte(u.Email))
	return code
}

// VerifyUserActiveCode 检查用户输入的账号激活码是否有效。
func VerifyUserActiveCode(code string) *User {
	if len(code) <= helper.TIME_LIMIT_CODE_LENGTH {
		return nil
	}

	hexStr := code[helper.TIME_LIMIT_CODE_LENGTH:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		if user, err := GetUserByEmail(string(b)); err != nil {
			return nil
		} else {
			prefix := code[:helper.TIME_LIMIT_CODE_LENGTH]
			data := com.ToStr(user.ID) + string(b) + user.NickName + user.Password + user.Salt

			if helper.VerifyTimeLimitCode(data, 120, prefix) {
				return user
			}
		}
	}
	return nil
}

// GetUserSalt 返回随机字符串作为用户的盐。
func GetUserSalt() string {
	return randstr.String(10)
}

// CreateUser 新建一个新的用户。
func CreateUser(u *User) error {
	if err := isUsernameAllowed(u.NickName); err != nil {
		return err
	}

	isExist := IsUserExist(u.NickName)
	if isExist {
		return ErrUserAlreadyExist{arg: u.NickName}
	}

	u.Email = strings.ToLower(u.Email)
	isExist = IsEmailUsed(u.Email)
	if isExist {
		return ErrEmailAlreadyUsed{arg: u.Email}
	}

	u.Avatar = helper.HashEmail(u.Email)
	u.Salt = GetUserSalt()
	u.EncodePassword()

	if err := db.Create(u).Error; err != nil {
		return errors.Wrap(err, "添加用户")
	}
	return nil
}

var ErrBadCredential = errors.New("电子邮箱或密码错误")

// UserAuthenticate 验证用户传入的用户名与密码。
func UserAuthenticate(email string, password string) (*User, error) {
	user := new(User)
	if err := db.Model(&User{}).Where(&User{Email: email}).First(&user).Error; err != nil {
		return nil, ErrBadCredential
	}

	if !user.ValidatePassword(password) {
		return nil, ErrBadCredential
	}
	return user, nil
}

// UpdateUserProfile 修改用户信息。
func UpdateUserProfile(u *User) error {
	if u.Password != "" {
		u.EncodePassword()
	}

	return db.Model(&User{}).Where(&User{
		Model: gorm.Model{ID: u.ID},
	}).Updates(&User{
		NickName: u.NickName,
		Password: u.Password,
		IsActive: u.IsActive,
	}).Error
}

// MustGetUserByID 根据传入的用户 ID 查询对应的用户。
func MustGetUserByID(uid uint) *User {
	var u User
	_ = db.Model(&User{}).Where(&User{
		Model: gorm.Model{
			ID: uid,
		},
	}).First(&u)
	return &u
}

var ErrUserNotFound = errors.New("用户不存在")

// GetUserByEmail 根据给定的电子邮箱地址查询对应的用户。
func GetUserByEmail(email string) (*User, error) {
	var u User
	if err := db.Model(&User{}).Where(&User{
		Email: email,
	}).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "查询用户")
	}

	return &u, nil
}

// GetUserByNickName 根据给定的用户昵称查询对应的用户。
func GetUserByNickName(nickName string) (*User, error) {
	var u User
	if err := db.Model(&User{}).Where(&User{
		NickName: nickName,
	}).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "查询用户")
	}

	return &u, nil
}

// GetUsers 批量获取用户
// options[0] offset
// options[1] limit
func GetUsers(options ...int) []*User {
	var users []*User

	if len(options) == 0 {
		db.Model(&User{}).Find(&users)
	} else {
		offset := 0
		if len(options) > 1 && options[0] > 0 {
			offset = options[0]
		}

		limit := 0
		if len(options) == 2 && options[1] > 0 {
			limit = options[1]
		}
		db.Model(&User{}).Offset(offset).Limit(limit).Find(&users)
	}

	return users
}

// CountUsers 返回用户的总数
func CountUsers() int64 {
	var count int64
	db.Model(&User{}).Count(&count)
	return count
}

func isUsernameAllowed(name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrNameNotAllowed{arg: name}
	}
	return nil
}

// IsUserExist 检查用户昵称是否重复。
func IsUserExist(name string) bool {
	if name == "" {
		return false
	}
	var u User
	err := db.Model(&User{}).Where(&User{NickName: name}).First(&u).Error
	return err == nil
}

// IsEmailUsed 检查邮箱是否重复。
func IsEmailUsed(email string) bool {
	if email == "" {
		return false
	}
	var u User
	err := db.Model(&User{}).Where(&User{Email: email}).First(&u).Error
	return err == nil
}

type ErrUserAlreadyExist struct {
	arg interface{}
}

func IsErrUserAlreadyExist(err error) bool {
	_, ok := err.(ErrEmailAlreadyUsed)
	return ok
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("用户昵称已被使用: %v", err.arg)
}

type ErrEmailAlreadyUsed struct {
	arg interface{}
}

func IsErrEmailAlreadyUsed(err error) bool {
	_, ok := err.(ErrEmailAlreadyUsed)
	return ok
}

func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("电子邮箱已被使用： %v", err.arg)
}

type ErrNameNotAllowed struct {
	arg interface{}
}

func IsErrNameNotAllowed(err error) bool {
	_, ok := err.(ErrNameNotAllowed)
	return ok
}

func (err ErrNameNotAllowed) Error() string {
	return fmt.Sprintf("用户名输入有误： %v", err.arg)
}
