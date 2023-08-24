package models

import (
	"GetHotWord/common"
	"GetHotWord/utils"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type User struct {
	ID       int    `json:"id" `
	UserName string `json:"name"`
	Email    string `json:"email" `
	// Password     string `json:"password,omitempty"`
	Password     string `json:"-"`
	AvatarUrl    string `json:"avatarUrl"`
	Introduce    string `json:"introduce"`
	SaltPassword string `json:"-" gorm:"type:varchar(60);comment:密码hash;<-"`

	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	Devices    []Device  `gorm:"foreignKey:owner"`
	Tags       []Tag     `gorm:"many2many:user_tag_tables"`
}

type APIUser struct {
	ID       int    `json:"id" `
	UserName string `json:"name"`
	Email    string `json:"email" `
	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime *LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

func (user *User) Create() error {
	// err := Db.Create(user)
	// 这个是事务的提交，如果有错误并且返回错误会自动回滚
	err := Db.Transaction(func(tx *gorm.DB) error {
		var tmpUser User
		// 如果已经存在
		if existErr := tx.Where("user_name=?", user.UserName).Find(&tmpUser).Error; existErr == nil {
			return common.NewError(200, 200400, "用户已存在")
		}
		// 转换为加盐密码，如果是前端的话，需要用https的方式传输明文
		saltPw, err := utils.GetHashPassword(user.Password)
		if err != nil {
			return err
		}
		user.SaltPassword = string(saltPw)
		res := tx.Create(&user)
		if res.Error != nil {
			zap.L().Error(res.Error.Error())
			return res.Error
		}
		return nil
	})

	if err != nil {
		zap.L().Error(err.Error())
		return err
	} else {
		return nil
	}
}

func (user *User) Delete() {

}

// 根据用户名获取
func (user *User) GetUserByName(name string) (interface{}, error) {
	ruser := User{}
	// select是只有获取到的字段有值，其余只有默认值
	err := Db.Model(&User{}).Where(&User{UserName: name}).Select([]string{"user_name", "id", "email"}).Find(&ruser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ruser, err

}

// 根据id
func (user *User) GetUserByID() (User, error) {
	var resUser User
	// select是只有获取到的字段有值，其余只有默认值
	err := Db.Where(&User{ID: user.ID}).Preload("Tags").Find(&resUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	return resUser, err
}

// 根据token解析后设置的header 里的id获取
func (user *User) GetUserByTokenID(id int) (interface{}, error) {
	var resUser User
	// select是只有获取到的字段有值，其余只有默认值
	err := Db.Where(&User{ID: id}).Select([]string{"user_name", "id", "email", "avatar_url", "introduce"}).Find(&resUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return resUser, err
}

// 删除用户
func (user *User) DeleteUserById(id int) error {

	// 开启事务
	err := Db.Transaction(func(tx *gorm.DB) error {
		err := Db.Delete(&User{ID: id}).Error
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		return nil
	})
	return err

}
func (user *User) NameExists() (bool, error) {
	return false, nil
}

func (user *User) GetAll() ([]User, error) {
	var resUser []User
	// 先预加载
	err := Db.Find(&resUser).Error

	return resUser, err
}

func (user *User) Update() error {

	upErr := Db.Transaction(func(tx *gorm.DB) error {

		if user.Password != "" {
			//  修改密码
			res, err := utils.GetHashPassword(user.Password)
			if err == nil {
				user.SaltPassword = string(res)
			} else {
				return err
			}
		}
		err := tx.Model(user).Updates(user).Error
		if err != nil {
			return err
		}
		// 删除，替换之前的多对多关系
		tx.Model(&user).Association("Tags").Replace(user.Tags)
		return nil
	})
	// 如果用户插入没错的话
	if upErr == nil {
		// 节点标签+1
		for _, v := range user.Tags {
			v.AddCount()
		}
	}

	return upErr
}

// 检查用户密码是否正确
func (user *User) VertifyNamePasswd(password string, name string) (bool, error, int) {
	var resUser User
	var userID int
	err := Db.Where(&User{UserName: name}).First(&resUser).Error
	fmt.Print(resUser)
	userID = resUser.ID
	if err != nil {
		return false, err, userID
	} else {
		res, judgeErr := utils.JuegeHashPassworCorrect(resUser.SaltPassword, password)
		if judgeErr != nil {
			return res, judgeErr, userID
		} else {
			return res, nil, userID
		}
	}

}
