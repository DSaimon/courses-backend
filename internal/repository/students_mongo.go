package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type StudentsRepo struct {
	db *mongo.Collection
}

func NewStudentsRepo(db *mongo.Database) *StudentsRepo {
	return &StudentsRepo{
		db: db.Collection(studentsCollection),
	}
}

func (r *StudentsRepo) Create(ctx context.Context, student domain.Student) error {
	_, err := r.db.InsertOne(ctx, student)
	return err
}

func (r *StudentsRepo) GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error) {
	var student domain.Student
	err := r.db.FindOne(ctx, bson.M{"email": email, "password": password, "schoolId": schoolId, "verification.verified": true}).Decode(&student)
	return student, err
}

func (r *StudentsRepo) GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error) {
	var student domain.Student
	err := r.db.FindOne(ctx, bson.M{"session.refreshToken": refreshToken, "schoolId": schoolId,
		"session.expiresAt": bson.M{"$gt": time.Now()}}).Decode(&student)

	return student, err
}

func (r *StudentsRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Student, error) {
	var student domain.Student
	err := r.db.FindOne(ctx, bson.M{"_id": id, "verification.verified": true}).Decode(&student)

	return student, err
}

func (r *StudentsRepo) SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
	return err
}

func (r *StudentsRepo) GiveAccessToCourseAndModule(ctx context.Context, studentId, courseId, moduleId primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$addToSet": bson.M{"availableModules": moduleId, "availableCourses": courseId}})
	return err
}

func (r *StudentsRepo) GiveAccessToCoursesAndModules(ctx context.Context, studentId primitive.ObjectID, courseIds, moduleIds []primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$addToSet": bson.M{"availableModules": bson.M{"$each": moduleIds},
		"availableCourses": bson.M{"$each": courseIds}}})
	return err
}

func (r *StudentsRepo) Verify(ctx context.Context, code string) error {
	codeId, err := primitive.ObjectIDFromHex(code)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOne(ctx,
		bson.M{"verification.code": codeId},
		bson.M{"$set": bson.M{"verification.verified": true}})

	return err
}
