package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StudentsService struct {
	repo   repository.Students
	hasher hash.PasswordHasher
}

func NewStudentsService(repo repository.Students, hasher hash.PasswordHasher) *StudentsService {
	return &StudentsService{repo: repo, hasher: hasher}
}

func (s *StudentsService) SignIn(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	student := domain.Student{
		Name:         input.Name,
		Password:     s.hasher.Hash(input.Password),
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		SchoolID:     input.SchoolID,
		Verification: domain.Verification{
			Hash: primitive.NewObjectID(),
		},
	}

	if input.SourceCourseID != "" {
		var err error

		student.SourceCourseID, err = primitive.ObjectIDFromHex(input.SourceCourseID)
		if err != nil {
			return err
		}
	}

	// TODO send emails

	return s.repo.Create(ctx, student)
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	return s.repo.Verify(ctx, hash)
}