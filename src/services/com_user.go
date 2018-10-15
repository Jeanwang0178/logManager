package services

import "logManager/src/models"

func UserServiceUserGetByName(userName string) (*models.User, error) {
	user, err := models.UserGetByName(userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UserServiceUserUpdate(user *models.User, fields ...string) error {
	err := models.UserUpdate(user)
	return err
}

func UserServiceUserGetById(id string) (*models.User, error) {
	user, err := models.UserGetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
