package boardsUsers

import (
	"errors"
	"hrm-app/config"
	"hrm-app/internal/pkg/utils"
)

// BoardRepository defines minimal interface needed to verify board ownership
type BoardRepository interface {
	FindByID(id uint) (*BoardInfo, error)
}

// BoardInfo contains minimal board information needed for authorization
type BoardInfo struct {
	ID          uint
	CreatedBy   uint
	WorkspaceID uint
}

// WorkspaceRepository defines minimal interface needed to verify workspace passcode
type WorkspaceRepository interface {
	FindByID(id uint) (*WorkspaceInfo, error)
}

// WorkspaceInfo contains minimal workspace information needed for passcode verification
type WorkspaceInfo struct {
	ID       uint
	PassCode string
}

type UseCase interface {
	Create(boardUsers *BoardsUsers, requestingUserID uint) error
	GetByBoardID(boardID uint) ([]BoardsUsers, error)
	GetByUserID(userID uint) ([]BoardsUsers, error)
	GetByID(id uint) (BoardsUsers, error)
	Delete(id uint) error
	Update(boardUsers *BoardsUsers) error
	HasAccess(boardID, userID uint) (bool, error)
	Join(userID uint, token string) error
	GenerateJoinToken(boardID, userID uint) (string, error)
}

type usecase struct {
	repo          Repository
	boardRepo     BoardRepository
	workspaceRepo WorkspaceRepository
	cfg           *config.Config
}

func NewUseCase(repo Repository, boardRepo BoardRepository, workspaceRepo WorkspaceRepository, cfg *config.Config) UseCase {
	return &usecase{
		repo:          repo,
		boardRepo:     boardRepo,
		workspaceRepo: workspaceRepo,
		cfg:           cfg,
	}
}

func (u *usecase) Create(boardUsers *BoardsUsers, requestingUserID uint) error {
	board, err := u.boardRepo.FindByID(boardUsers.BoardID)

	if err != nil {
		return errors.New("board not found")
	}

	// Check if the requesting user is the creator of the board
	if board.CreatedBy != requestingUserID {
		return errors.New("unauthorized: only board creator can add users")
	}

	// Check if user is already assigned to the board
	existingUser, _ := u.repo.GetByBoardIDAndUserID(boardUsers.BoardID, boardUsers.UserID)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already assigned to this board")
	}

	return u.repo.Create(boardUsers)
}

func (u *usecase) GetByBoardID(boardID uint) ([]BoardsUsers, error) {
	return u.repo.GetByBoardID(boardID)
}

func (u *usecase) GetByUserID(userID uint) ([]BoardsUsers, error) {
	boardsUsers, err := u.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(boardsUsers) == 0 {
		return nil, errors.New("board users not found")
	}

	return boardsUsers, nil
}

func (u *usecase) GetByID(id uint) (BoardsUsers, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) Update(boardUsers *BoardsUsers) error {
	return u.repo.Update(boardUsers)
}

func (u *usecase) HasAccess(boardID, userID uint) (bool, error) {
	return u.repo.HasAccess(boardID, userID)
}

func (u *usecase) GenerateJoinToken(boardID, userID uint) (string, error) {
	board, err := u.boardRepo.FindByID(boardID)
	if err != nil {
		return "", errors.New("board not found")
	}

	// Only board creator can generate join token
	if board.CreatedBy != userID {
		return "", errors.New("unauthorized: only board creator can generate join token")
	}

	workspace, err := u.workspaceRepo.FindByID(board.WorkspaceID)
	if err != nil {
		return "", errors.New("workspace not found")
	}

	return utils.GenerateJoinToken(u.cfg, boardID, "board", workspace.PassCode)
}

func (u *usecase) Join(userID uint, token string) error {
	claims, err := utils.ValidateJoinToken(u.cfg, token)
	if err != nil {
		return err
	}

	if claims.EntityType != "board" {
		return errors.New("invalid token type")
	}

	boardID := claims.EntityID
	passCode := claims.PassCode

	// Fetch the board to get workspace_id
	board, err := u.boardRepo.FindByID(boardID)
	if err != nil {
		return errors.New("board not found")
	}

	// Fetch the parent workspace to verify passcode
	workspace, err := u.workspaceRepo.FindByID(board.WorkspaceID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Check if the passcode matches the workspace passcode
	if workspace.PassCode != passCode {
		return errors.New("invalid passcode")
	}

	// Check if user is already assigned to the board
	existingUser, _ := u.repo.GetByBoardIDAndUserID(boardID, userID)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already assigned to this board")
	}

	boardUsers := &BoardsUsers{
		BoardID: boardID,
		UserID:  userID,
	}

	return u.repo.Create(boardUsers)
}
