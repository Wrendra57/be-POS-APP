package template_response

import (
	"github.com/Wrendra57/Pos-app-be/internal/models/domain"
	"github.com/Wrendra57/Pos-app-be/internal/models/webrespones"
)

func ToUserRespone(user domain.User, oauth domain.Oauth) webrespones.UserDetail {
	return webrespones.UserDetail{
		User_id:    user.User_id,
		Email:      oauth.Email,
		Is_enabled: oauth.Is_enabled,
		Username:   oauth.Username,
		Name:       user.Name,
		Gender:     user.Gender,
		Telp:       user.Telp,
		Birthday:   user.Birthday,
		Address:    user.Address,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}
}
