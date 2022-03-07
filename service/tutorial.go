package service

import (
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
)

//  根据信息保存教程
func SaveTutorial(tutorial *model.Tutorial) (err error) {
	err = global.DB.Save(tutorial).Error
	return err
}
