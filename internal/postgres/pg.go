package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"untitled_folder/internal/cache"
	"untitled_folder/internal/model"
)

//type DB struct {
//	*gorm.DB
//}

func InitDB() (*gorm.DB, error) {
	connStr := "host=localhost user=polaykov.art dbname=wb-orders sslmode=disable"
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		return nil, err
	}
	//db.Migrator().DropTable(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Item{})
	err = db.AutoMigrate(&model.Order{}, &model.Delivery{}, &model.Payment{}, &model.Item{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	return db, nil
}

func LoadCacheFromDB(db *gorm.DB, c *cache.Cache) error {
	var orders []model.Order

	err := db.Preload("Delivery").Preload("Items").Preload("Payment").Find(&orders).Error
	if err != nil {
		return fmt.Errorf("could not load orders from database: %w", err)
	}

	for _, order := range orders {
		c.CacheSet(order)
	}

	return nil
}
