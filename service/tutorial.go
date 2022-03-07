package service

import (
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

// Helper

// 获取教程文件夹名称
func GetTutorialFileName(tutorial model.Tutorial) string {
	return strconv.Itoa(int(tutorial.ID)) + "_" + strconv.Itoa(int(tutorial.Version))
}

//  根据信息保存教程
func SaveTutorial(tutorial *model.Tutorial) (err error) {
	err = global.DB.Save(tutorial).Error
	return err
}

// GetTutorialByID 根据教程 ID 查询某个教程
func GetTutorialByID(ID uint64) (tutorial model.Tutorial, notFound bool) {
	err := global.DB.First(&tutorial, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tutorial, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByID: search error")
		return tutorial, true
	} else {
		return tutorial, false
	}
}
