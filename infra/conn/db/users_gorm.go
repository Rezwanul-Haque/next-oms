package db

import (
	"gorm.io/gorm"
	"next-oms/app/domain"
	"next-oms/app/utils/methodsutil"
	"next-oms/app/utils/msgutil"
	"next-oms/infra/conn/db/models"
	"next-oms/infra/errors"
	"strings"
)

func (dc DatabaseClient) SaveUser(user *domain.User) (*domain.User, *errors.RestErr) {
	res := dc.DB.Model(&models.User{}).Create(&user)

	if res.Error != nil {
		dc.lc.Error("error occurred when create user", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (dc DatabaseClient) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := dc.DB.Model(&models.User{}).Where("id = ?", userID).First(&resp)

	if res.RowsAffected == 0 {
		dc.lc.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		dc.lc.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (dc DatabaseClient) UpdateUser(user *domain.User) *errors.RestErr {
	res := dc.DB.Model(&models.User{}).Omit("password", "app_key").Where("id = ?", user.ID).Updates(&user)

	if res.Error != nil {
		dc.lc.Error("error occurred when updating user by user id", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) UpdatePassword(userID uint, companyID uint, updateValues map[string]interface{}) *errors.RestErr {
	res := dc.DB.Model(&models.User{}).Where("id = ? AND company_id = ?", userID, companyID).Updates(&updateValues)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("updating user by user id"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	res := dc.DB.Model(&models.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected == 0 {
		dc.lc.Error("no user found by this email", res.Error)
		return nil, errors.NewError(errors.ErrRecordNotFound)
	}
	if res.Error != nil {
		dc.lc.Error("error occurred when trying to get user by email", res.Error)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (dc DatabaseClient) SetLastLoginAt(user *domain.User) error {
	dbusr := models.User{
		ID: user.ID,
	}

	err := dc.DB.Model(&dbusr).
		Update("last_login_at", user.LastLoginAt).
		Error

	if err != nil {
		dc.lc.Error(err.Error(), err)
		return err
	}

	return nil
}

func (dc DatabaseClient) ResetPassword(userID int, hashedPass []byte) error {
	err := dc.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPass).
		Error

	if err != nil {
		dc.lc.Error("error occur when reset password", err)
		return err
	}

	return nil
}

func (dc DatabaseClient) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	tempUser := &domain.TempVerifyTokenResp{}
	var vtUser domain.VerifyTokenResp

	query := dc.tokenUserFetchQuery()

	res := query.Where("users.id = ?", id).Find(&tempUser)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("get token user"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(tempUser, &vtUser.BaseVerifyTokenResp)
	if err != nil {
		dc.lc.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	vtUser.Permissions = strings.Split(tempUser.Permissions, ",")

	return &vtUser, nil
}

func (dc DatabaseClient) tokenUserFetchQuery() *gorm.DB {
	selections := `
		users.id,
		users.first_name,
		users.last_name,
		users.email,
		users.phone,
		users.profile_pic,
		companies.business_id,
		businesses.name business_name,
		companies.id company_id,
		companies.name company_name,
		(
			CASE
				WHEN 1 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS admin,
		(
			CASE
				WHEN 3 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS super_admin,
		GROUP_CONCAT(DISTINCT permissions.name) AS permissions
	`

	return dc.DB.Table("users").
		Select(selections).
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Joins("LEFT JOIN businesses ON companies.business_id = businesses.id").
		Joins("JOIN roles ON users.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.deleted_at IS NULL").
		Group("users.id")
}
