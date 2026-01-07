package workspaces

import (
	"testing"

	"hrm-app/config"
	"hrm-app/internal/domain/workspacesUsers"
)

// mockRepository is a mock implementation of the Repository interface
type mockRepository struct {
	createFunc              func(workspace *Workspace) error
	updateFunc              func(workspace *Workspace) error
	findByIDFunc            func(id uint) (*Workspace, error)
	deleteFunc              func(id uint) error
	findGuestWorkspacesFunc func(userID uint) ([]Workspace, error)
}

func (m *mockRepository) Create(workspace *Workspace) error {
	if m.createFunc != nil {
		return m.createFunc(workspace)
	}
	return nil
}

func (m *mockRepository) FindAll() ([]Workspace, error) {
	return nil, nil
}

func (m *mockRepository) FindByID(id uint) (*Workspace, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, nil
}

func (m *mockRepository) FindByUserID(userID uint) ([]Workspace, error) {
	return nil, nil
}

func (m *mockRepository) FindByIDs(ids []uint) ([]Workspace, error) {
	return nil, nil
}

func (m *mockRepository) FindByUserAccess(userID uint) ([]Workspace, error) {
	return nil, nil
}

func (m *mockRepository) FindGuestWorkspaces(userID uint) ([]Workspace, error) {
	if m.findGuestWorkspacesFunc != nil {
		return m.findGuestWorkspacesFunc(userID)
	}
	return nil, nil
}

func (m *mockRepository) Update(workspace *Workspace) error {
	if m.updateFunc != nil {
		return m.updateFunc(workspace)
	}
	return nil
}

func (m *mockRepository) Delete(id uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

type mockWorkspacesUsersRepository struct {
	createFunc func(wu *workspacesUsers.WorkspacesUsers) error
}

func (m *mockWorkspacesUsersRepository) Create(wu *workspacesUsers.WorkspacesUsers) error {
	if m.createFunc != nil {
		return m.createFunc(wu)
	}
	return nil
}

func (m *mockWorkspacesUsersRepository) GetByWorkspaceID(workspaceID uint) ([]workspacesUsers.WorkspacesUsers, error) {
	return nil, nil
}

func (m *mockWorkspacesUsersRepository) GetByUserID(userID uint) ([]workspacesUsers.WorkspacesUsers, error) {
	return nil, nil
}

func (m *mockWorkspacesUsersRepository) GetByWorkspaceIDAndUserID(workspaceID, userID uint) (*workspacesUsers.WorkspacesUsers, error) {
	return nil, nil
}

func (m *mockWorkspacesUsersRepository) GetByID(id uint) (workspacesUsers.WorkspacesUsers, error) {
	return workspacesUsers.WorkspacesUsers{}, nil
}

func (m *mockWorkspacesUsersRepository) Delete(id uint) error {
	return nil
}

func (m *mockWorkspacesUsersRepository) Update(wu *workspacesUsers.WorkspacesUsers) error {
	return nil
}

func TestUseCase_Create_PrivacyValidation(t *testing.T) {
	tests := []struct {
		name    string
		privacy string
		wantErr bool
	}{
		{
			name:    "valid public",
			privacy: "public",
			wantErr: false,
		},
		{
			name:    "valid private",
			privacy: "private",
			wantErr: false,
		},
		{
			name:    "valid team",
			privacy: "team",
			wantErr: false,
		},
		{
			name:    "invalid value",
			privacy: "locked",
			wantErr: true,
		},
		{
			name:    "empty value",
			privacy: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				createFunc: func(w *Workspace) error {
					return nil
				},
			}
			wuRepo := &mockWorkspacesUsersRepository{}
			uc := NewUseCase(repo, wuRepo, &config.Config{})

			err := uc.Create(&Workspace{Privacy: tt.privacy})
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Update_PrivacyValidation(t *testing.T) {
	tests := []struct {
		name    string
		privacy string
		wantErr bool
	}{
		{
			name:    "valid public",
			privacy: "public",
			wantErr: false,
		},
		{
			name:    "valid private",
			privacy: "private",
			wantErr: false,
		},
		{
			name:    "valid team",
			privacy: "team",
			wantErr: false,
		},
		{
			name:    "invalid value",
			privacy: "locked",
			wantErr: true,
		},
		{
			name:    "empty value",
			privacy: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				updateFunc: func(w *Workspace) error {
					return nil
				},
			}
			wuRepo := &mockWorkspacesUsersRepository{}
			uc := NewUseCase(repo, wuRepo, &config.Config{})

			err := uc.Update(&Workspace{Privacy: tt.privacy})
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Create_AssociatesUser(t *testing.T) {
	expectedUserID := uint(123)
	expectedWorkspaceID := uint(456)
	var capturedUserID uint
	var capturedWorkspaceID uint

	repo := &mockRepository{
		createFunc: func(w *Workspace) error {
			w.ID = expectedWorkspaceID
			return nil
		},
	}
	wuRepo := &mockWorkspacesUsersRepository{
		createFunc: func(wu *workspacesUsers.WorkspacesUsers) error {
			capturedUserID = wu.UserID
			capturedWorkspaceID = wu.WorkspaceID
			return nil
		},
	}
	uc := NewUseCase(repo, wuRepo, &config.Config{})

	err := uc.Create(&Workspace{Privacy: "public", CreatedBy: expectedUserID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedUserID != expectedUserID {
		t.Errorf("expected UserID %v, got %v", expectedUserID, capturedUserID)
	}
	if capturedWorkspaceID != expectedWorkspaceID {
		t.Errorf("expected WorkspaceID %v, got %v", expectedWorkspaceID, capturedWorkspaceID)
	}
}
func TestUseCase_DeleteByID_Authorization(t *testing.T) {
	creatorID := uint(1)
	otherUserID := uint(2)
	workspaceID := uint(10)

	repo := &mockRepository{
		findByIDFunc: func(id uint) (*Workspace, error) {
			if id == workspaceID {
				return &Workspace{ID: workspaceID, CreatedBy: creatorID}, nil
			}
			return nil, nil
		},
		deleteFunc: func(id uint) error {
			return nil
		},
	}
	wuRepo := &mockWorkspacesUsersRepository{}
	uc := NewUseCase(repo, wuRepo, &config.Config{})

	t.Run("authorized delete", func(t *testing.T) {
		err := uc.DeleteByID(workspaceID, creatorID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("unauthorized delete", func(t *testing.T) {
		err := uc.DeleteByID(workspaceID, otherUserID)
		if err == nil {
			t.Error("expected unauthorized error, got nil")
		}
		if err.Error() != "unauthorized: only the workspace creator can delete this workspace" {
			t.Errorf("expected specific unauthorized error message, got %v", err.Error())
		}
	})

	t.Run("workspace not found", func(t *testing.T) {
		err := uc.DeleteByID(999, creatorID)
		if err == nil {
			t.Error("expected error for non-existent workspace, got nil")
		}
	})
}
