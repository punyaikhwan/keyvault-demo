package service

import (
	"context"
	"keyvault-demo/domain/entity"
	"keyvault-demo/domain/repository"

	"github.com/google/uuid"
)

type StudentService struct {
	repo repository.StudentRepo
}

func NewStudentService() StudentService {
	return StudentService{
		repo: repository.NewStudentRepo(),
	}
}

func (s *StudentService) FindAll(ctx context.Context) (students []entity.Student, err error) {
	students, err = s.repo.FindAll(ctx)
	if err != nil {
		return students, err
	}

	// Decrypt NIK and phone one by one. If fails, fill with error code.
	for i, std := range students {
		// decrypt NIK and phone
		if err = std.DecryptNIK(); err != nil {
			students[i].NIK = "error: " + err.Error()
		}
		if err = std.DecryptPhone(); err != nil {
			students[i].Phone = "error: " + err.Error()
		}
	}

	return students, nil
}

func (s *StudentService) FindByID(ctx context.Context, id uuid.UUID) (student entity.Student, err error) {
	student, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return student, err
	}

	// decrypt NIK and phone
	if err = student.DecryptNIK(); err != nil {
		student.NIK = "error: " + err.Error()
	}
	if err = student.DecryptPhone(); err != nil {
		student.Phone = "error: " + err.Error()
	}

	return student, nil
}

func (s *StudentService) Create(ctx context.Context, student entity.Student) (id uuid.UUID, err error) {
	// encrypt NIK and phone before save
	student.KeyVersion = ""
	if student.ID == uuid.Nil {
		student.GetAndSetID()
	}
	if err = student.EncryptNIK(); err != nil {
		return uuid.Nil, err
	}
	if err = student.EncryptPhone(); err != nil {
		return uuid.Nil, err
	}

	id, err = s.repo.Create(ctx, student)
	return id, err
}

func (s *StudentService) Rotate(ctx context.Context) (numSuccess int, numFails int, err error) {
	stdTemps, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, 0, err
	}

	// Decrypt and reencrypt NIK and phone one by one. If fails, don't rotate
	students := make([]entity.Student, 0)
	for _, std := range stdTemps {
		if err = std.DecryptNIK(); err != nil {
			numFails++
			continue
		}
		if err = std.DecryptPhone(); err != nil {
			numFails++
			continue
		}

		if err = std.EncryptNIK(); err != nil {
			numFails++
			continue
		}
		if err = std.EncryptPhone(); err != nil {
			numFails++
			continue
		}
		students = append(students, std)
	}

	// rotate
	err = s.repo.BulkUpdate(ctx, students)
	if err != nil {
		numFails = len(stdTemps)
	} else {
		numSuccess = len(stdTemps) - numFails
	}
	return numSuccess, numFails, err
}
