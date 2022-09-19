package repository

import (
	"context"
	"keyvault-demo/domain/entity"
	"log"

	"keyvault-demo/config"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StudentRepo struct {
	db *gorm.DB
}

func NewStudentRepo() StudentRepo {
	db, err := gorm.Open(mysql.Open(config.Configuration().DBURI), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database", err)
	}
	db.AutoMigrate(&entity.Student{})

	return StudentRepo{db}
}

func (r *StudentRepo) FindAll(ctx context.Context) (students []entity.Student, err error) {
	err = r.db.WithContext(ctx).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *StudentRepo) FindByID(ctx context.Context, id uuid.UUID) (student entity.Student, err error) {
	err = r.db.WithContext(ctx).First(&student, id).Error
	return student, err
}

func (r *StudentRepo) Create(ctx context.Context, student entity.Student) (id uuid.UUID, err error) {
	err = r.db.WithContext(ctx).Create(&student).Error
	return student.ID, err
}

func (r *StudentRepo) BulkUpdate(ctx context.Context, students []entity.Student) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, student := range students {
			tx1 := tx.Model(&student).
				Session(&gorm.Session{FullSaveAssociations: true}).
				Updates(&student)
			if tx1.Error != nil {
				return tx1.Error
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
