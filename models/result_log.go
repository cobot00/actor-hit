package models

import (
  "log"
  "math"

  "github.com/jinzhu/gorm"
)

type ResultLog struct {
  IP  string
  Hit uint
}

func InsertResultLog(db *gorm.DB, ip string, hit uint) {
  resultLog := ResultLog{
    IP:  ip,
    Hit: hit,
  }
  db.Table("result_log").Create(&resultLog)
}

func HitAverage(db *gorm.DB, ip string) int {
  var resultLogs []ResultLog
  db.Table("result_log").
    Where("ip = ?", ip).Order("created_at DESC").Limit(5).Find(&resultLogs)

  var hitTotal uint
  for _, resultLog := range resultLogs {
    hitTotal += resultLog.Hit
  }

  total := float64(len(resultLogs) * 10)
  log.Println("[HitAverage] total: %f, hitTotal: %d", total, hitTotal)

  hitRate := float64(hitTotal) / total * 100
  return int(math.Ceil(hitRate))
}
