package service

import (
	"errors"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"gorm.io/gorm"
	"strconv"
)

// Helper

// GetTutorialFileName 获取教程文件夹名称
func GetTutorialFileName(tutorial model.Tutorial) string {
	return strconv.Itoa(int(tutorial.ID)) + "_" + strconv.Itoa(int(tutorial.Version))
}

// SaveTutorial 根据信息保存教程
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

// UpdateTutorial 根据信息更新教程
func UpdateTutorial(tutorial *model.Tutorial, q *model.UpdateTutorialQ) (err error) {
	tutorial.Name, tutorial.Profile, tutorial.Readable, tutorial.Writable, tutorial.Version =
		q.Name, q.Profile, q.Readable, q.Writable, tutorial.Version+1
	err = global.DB.Save(tutorial).Error
	return err
}

// DeleteTutorialByID 根据教程id 删除教程
func DeleteTutorialByID(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(model.Tutorial{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllTutorials 查询所有教程
func GetAllTutorials() (tutorials []model.Tutorial) {
	tutorials = make([]model.Tutorial, 0)
	global.DB.Find(&tutorials)
	return tutorials
}
