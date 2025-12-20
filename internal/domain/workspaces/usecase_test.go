package workspaces

import (
	"testing"
)

// mockRepository is a mock implementation of the Repository interface
type mockRepository struct {
	createFunc func(workspace *Workspace) error
	updateFunc func(workspace *Workspace) error
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
	return nil, nil
}

func (m *mockRepository) FindByUserID(userID uint) ([]Workspace, error) {
	return nil, nil
}

func (m *mockRepository) Update(workspace *Workspace) error {
	if m.updateFunc != nil {
		return m.updateFunc(workspace)
	}
	return nil
}

func (m *mockRepository) Delete(id uint) error {
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
			uc := NewUseCase(repo)

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
			uc := NewUseCase(repo)

			err := uc.Update(&Workspace{Privacy: tt.privacy})
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
