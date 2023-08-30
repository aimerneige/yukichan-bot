package service

import (
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database"
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database/model"
	"gorm.io/gorm"
)

// DBService 数据库服务
type DBService struct {
	db *gorm.DB
}

// NewDBService 创建数据库服务
func NewDBService() *DBService {
	return &DBService{
		db: database.GetDB(),
	}
}

// CreateELO 创建 ELO
func (s *DBService) CreateELO(uin int64, name string, rate int) error {
	return s.db.Create(&model.ELO{
		Uin:  uin,
		Name: name,
		Rate: rate,
	}).Error
}

// GetELOByUin 获取 ELO
func (s *DBService) GetELOByUin(uin int64) (model.ELO, error) {
	var elo model.ELO
	err := s.db.Where("uin = ?", uin).First(&elo).Error
	return elo, err
}

// GetELORateByUin 获取 ELO 等级分
func (s *DBService) GetELORateByUin(uin int64) (int, error) {
	var elo model.ELO
	err := s.db.Select("rate").Where("uin = ?", uin).First(&elo).Error
	return elo.Rate, err
}

// GetHighestRateList 获取最高的等级分列表
func (s *DBService) GetHighestRateList() ([]model.ELO, error) {
	var eloList []model.ELO
	err := s.db.Order("rate desc").Limit(10).Find(&eloList).Error
	return eloList, err
}

// UpdateELOByUin 更新 ELO 等级分
func (s *DBService) UpdateELOByUin(uin int64, name string, rate int) error {
	return s.db.Model(&model.ELO{}).Where("uin = ?", uin).Update("name", name).Update("rate", rate).Error
}

// CleanELOByUin 清空用户 ELO 等级分
func (s *DBService) CleanELOByUin(uin int64) error {
	return s.db.Model(&model.ELO{}).Where("uin = ?", uin).Update("rate", 100).Error
}

// CreatePGN 创建 PGN
func (s *DBService) CreatePGN(data string, whiteUin int64, blackUin int64, whiteName string, blackName string) error {
	return s.db.Create(&model.PGN{
		Data:      data,
		WhiteUin:  whiteUin,
		BlackUin:  blackUin,
		WhiteName: whiteName,
		BlackName: blackName,
	}).Error
}
