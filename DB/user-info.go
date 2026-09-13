package DB

// InsertUserInfo 初始化用户的info
func InsertUserInfo(userID uint) error {
	var userInfo = UserInfo{
		UserID: userID,
	}
	err := MySQL.Table("user_infos").Create(&userInfo).Error
	return err
}

// UpdateUserInfo 更新用户的Info
func UpdateUserInfo(userID uint, postform UserInfo) error {
	if err := MySQL.Table("user_infos").Where("user_id=?", userID).Updates(&postform).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserName 更新用户的昵称
func UpdateUserName(userID uint, New string) error {
	if err := MySQL.Table("users").Where("id=?", userID).Update("nickname", New).Error; err != nil {
		return err
	}
	return nil
}
