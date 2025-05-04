package domain

import (
	"github.com/RockkleyPushPost/user-service/domain/dto"
	"github.com/RockkleyPushPost/user-service/entity"
	"github.com/google/uuid"
)

type UserUseCase interface {
	GetByUUID(uuid uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetFriends(userUUID uuid.UUID) ([]entity.User, error)
	DeleteFriend(dto *dto.DeleteFriendDTO) error
}

type FriendshipUseCase interface {
	CreateFriendshipRequest(dto dto.CreateFriendRequestDto) error
	FindRequestByPairUUID(dto dto.FindByPairUUID) (entity.FriendshipRequest, error)
	UpdateFriendshipRequestStatus(dto.UpdateFriendshipRequestDto) error
	DeleteFriendshipRequest(dto.DeleteFriendshipRequestDto) error
	AcceptFriendshipRequest(requestUUID uuid.UUID) error
	DeclineFriendshipRequest(requestUUID uuid.UUID) error
	FindFriendshipRequest(requestDTO *dto.FindFriendshipRequestDTO) ([]*entity.FriendshipRequest, error)
	FindFriendships(friendshipDTO *dto.FindFriendshipDTO) ([]*entity.Friendship, error)
}
