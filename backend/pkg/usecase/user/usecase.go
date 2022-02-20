package user

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/user"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	GetPublicUsers() ([]string, error)
	CreateUser(models.CreaterUserForm) (*models.User, error)
}

func NewUsecase(tp storage.TxProvider, userRepo user.UserRepository) Usecase {
	return defaultUsecase{
		userRepo:   userRepo,
		txProvider: tp,
	}
}

type defaultUsecase struct {
	userRepo   user.UserRepository
	txProvider storage.TxProvider
}

func (u defaultUsecase) CreateUser(form models.CreaterUserForm) (createdUser *models.User, err error) {
	var uow storage.UnitOfWork
	uow, err = u.txProvider.Begin()
	if err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	newUserId := uuid.NewOrderedUUID()
	form.Id = newUserId

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(*form.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPassStr := string(hashedPass)
	form.Password = &hashedPassStr

	err = u.userRepo.WithUnitOfWork(uow).Create(form)
	if err != nil {
		return nil, err
	}

	createdUser, err = u.userRepo.WithUnitOfWork(uow).GetById(newUserId)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u defaultUsecase) GetPublicUsers() ([]string, error) {
	return u.userRepo.GetPublicUsernames()
}
