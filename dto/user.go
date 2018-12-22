package dto

import "time"

type User struct {
	UserId         int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserName       string `gorm:"UNIQUE_INDEX"`
	Address        string `gorm:"type:varchar(130)"`
	Birthday       string
	Email          string
	Mobile         string
	HashedPassword string
	CreateTime     int64
	Deleted        bool
}

// Load 加载单个用户
func (user *User) Load() error {
	if err := db.Where("user_name=? ", user.UserName).Find(user).Error; err != nil {
		return dbError2CodeError(err)
	}

	return nil
}

// GetList 获取列表
func (user *User) GetList(createTime int64, limit int32, direct string) ([]*User, error) {
	qdb := db.Model(user)
	if direct == "" || createTime == 0 {
		direct = "PREV"
		createTime = time.Now().UnixNano() / int64(time.Millisecond)
	}

	switch direct {
	case "NEXT":
		qdb = qdb.Where("create_time > ?", createTime)
	case "PREV":
		qdb = qdb.Where("create_time < ?", createTime)
	}
	qdb = qdb.Limit(limit)

	users := []*User{}
	if err := qdb.Find(users).Error; err != nil {
		return nil, dbError2CodeError(err)
	}

	return users, nil
}

// Update 获取列表
func (user *User) Update() error {
	qdb := db.Model(user)
	if err := qdb.Where("user_name = ?", user.UserName).Update(map[string]interface{}{
		"birthday": user.Birthday,
		"address":  user.Address,
		"mobile":   user.Mobile,
		"email":    user.Email,
	}).Error; err != nil {
		return dbError2CodeError(err)
	}
	return nil
}

func (user *User) Insert() error {
	if err := db.Create(user).Error; err != nil {
		return dbError2CodeError(err)
	}

	return nil
}
