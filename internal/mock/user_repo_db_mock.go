// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"sync"
)

// Ensure, that MockUserRepository does implement users.UserRepository.
// If this is not the case, regenerate this file with moq.
var _ users.UserRepository = &MockUserRepository{}

// MockUserRepository is a mock implementation of users.UserRepository.
//
// 	func TestSomethingThatUsesUserRepository(t *testing.T) {
//
// 		// make and configure a mocked users.UserRepository
// 		mockedUserRepository := &MockUserRepository{
// 			CheckPasswordByUserIDFunc: func(n int, s string) (bool, error) {
// 				panic("mock out the CheckPasswordByUserID method")
// 			},
// 			CreateUserFunc: func(user *models.User) (int, error) {
// 				panic("mock out the CreateUser method")
// 			},
// 			DoesUserExistFunc: func(auth *models.Auth) (int, error) {
// 				panic("mock out the DoesUserExist method")
// 			},
// 			GetAvatarFilenameFunc: func(n int) (string, error) {
// 				panic("mock out the GetAvatarFilename method")
// 			},
// 			GetSettingsFunc: func(n int) (*models.SettingsGet, error) {
// 				panic("mock out the GetSettings method")
// 			},
// 			IsEmailUniqueFunc: func(s string) (bool, error) {
// 				panic("mock out the IsEmailUnique method")
// 			},
// 			IsNicknameUniqueFunc: func(s string) (bool, error) {
// 				panic("mock out the IsNicknameUnique method")
// 			},
// 			UpdateAvatarFunc: func(n int, s string) error {
// 				panic("mock out the UpdateAvatar method")
// 			},
// 			UpdateEmailFunc: func(n int, s string) error {
// 				panic("mock out the UpdateEmail method")
// 			},
// 			UpdateNicknameFunc: func(n int, s string) error {
// 				panic("mock out the UpdateNickname method")
// 			},
// 			UpdatePasswordFunc: func(n int, s string) error {
// 				panic("mock out the UpdatePassword method")
// 			},
// 		}
//
// 		// use mockedUserRepository in code that requires users.UserRepository
// 		// and then make assertions.
//
// 	}
type MockUserRepository struct {
	// CheckPasswordByUserIDFunc mocks the CheckPasswordByUserID method.
	CheckPasswordByUserIDFunc func(n int, s string) (bool, error)

	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(user *models.User) (int, error)

	// DoesUserExistFunc mocks the DoesUserExist method.
	DoesUserExistFunc func(auth *models.Auth) (int, error)

	// GetAvatarFilenameFunc mocks the GetAvatarFilename method.
	GetAvatarFilenameFunc func(n int) (string, error)

	// GetSettingsFunc mocks the GetSettings method.
	GetSettingsFunc func(n int) (*models.SettingsGet, error)

	// IsEmailUniqueFunc mocks the IsEmailUnique method.
	IsEmailUniqueFunc func(s string) (bool, error)

	// IsNicknameUniqueFunc mocks the IsNicknameUnique method.
	IsNicknameUniqueFunc func(s string) (bool, error)

	// UpdateAvatarFunc mocks the UpdateAvatar method.
	UpdateAvatarFunc func(n int, s string) error

	// UpdateEmailFunc mocks the UpdateEmail method.
	UpdateEmailFunc func(n int, s string) error

	// UpdateNicknameFunc mocks the UpdateNickname method.
	UpdateNicknameFunc func(n int, s string) error

	// UpdatePasswordFunc mocks the UpdatePassword method.
	UpdatePasswordFunc func(n int, s string) error

	// calls tracks calls to the methods.
	calls struct {
		// CheckPasswordByUserID holds details about calls to the CheckPasswordByUserID method.
		CheckPasswordByUserID []struct {
			// N is the n argument value.
			N int
			// S is the s argument value.
			S string
		}
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// User is the user argument value.
			User *models.User
		}
		// DoesUserExist holds details about calls to the DoesUserExist method.
		DoesUserExist []struct {
			// Auth is the auth argument value.
			Auth *models.Auth
		}
		// GetAvatarFilename holds details about calls to the GetAvatarFilename method.
		GetAvatarFilename []struct {
			// N is the n argument value.
			N int
		}
		// GetSettings holds details about calls to the GetSettings method.
		GetSettings []struct {
			// N is the n argument value.
			N int
		}
		// IsEmailUnique holds details about calls to the IsEmailUnique method.
		IsEmailUnique []struct {
			// S is the s argument value.
			S string
		}
		// IsNicknameUnique holds details about calls to the IsNicknameUnique method.
		IsNicknameUnique []struct {
			// S is the s argument value.
			S string
		}
		// UpdateAvatar holds details about calls to the UpdateAvatar method.
		UpdateAvatar []struct {
			// N is the n argument value.
			N int
			// S is the s argument value.
			S string
		}
		// UpdateEmail holds details about calls to the UpdateEmail method.
		UpdateEmail []struct {
			// N is the n argument value.
			N int
			// S is the s argument value.
			S string
		}
		// UpdateNickname holds details about calls to the UpdateNickname method.
		UpdateNickname []struct {
			// N is the n argument value.
			N int
			// S is the s argument value.
			S string
		}
		// UpdatePassword holds details about calls to the UpdatePassword method.
		UpdatePassword []struct {
			// N is the n argument value.
			N int
			// S is the s argument value.
			S string
		}
	}
	lockCheckPasswordByUserID sync.RWMutex
	lockCreateUser            sync.RWMutex
	lockDoesUserExist         sync.RWMutex
	lockGetAvatarFilename     sync.RWMutex
	lockGetSettings           sync.RWMutex
	lockIsEmailUnique         sync.RWMutex
	lockIsNicknameUnique      sync.RWMutex
	lockUpdateAvatar          sync.RWMutex
	lockUpdateEmail           sync.RWMutex
	lockUpdateNickname        sync.RWMutex
	lockUpdatePassword        sync.RWMutex
}

// CheckPasswordByUserID calls CheckPasswordByUserIDFunc.
func (mock *MockUserRepository) CheckPasswordByUserID(n int, s string) (bool, error) {
	if mock.CheckPasswordByUserIDFunc == nil {
		panic("MockUserRepository.CheckPasswordByUserIDFunc: method is nil but UserRepository.CheckPasswordByUserID was just called")
	}
	callInfo := struct {
		N int
		S string
	}{
		N: n,
		S: s,
	}
	mock.lockCheckPasswordByUserID.Lock()
	mock.calls.CheckPasswordByUserID = append(mock.calls.CheckPasswordByUserID, callInfo)
	mock.lockCheckPasswordByUserID.Unlock()
	return mock.CheckPasswordByUserIDFunc(n, s)
}

// CheckPasswordByUserIDCalls gets all the calls that were made to CheckPasswordByUserID.
// Check the length with:
//     len(mockedUserRepository.CheckPasswordByUserIDCalls())
func (mock *MockUserRepository) CheckPasswordByUserIDCalls() []struct {
	N int
	S string
} {
	var calls []struct {
		N int
		S string
	}
	mock.lockCheckPasswordByUserID.RLock()
	calls = mock.calls.CheckPasswordByUserID
	mock.lockCheckPasswordByUserID.RUnlock()
	return calls
}

// CreateUser calls CreateUserFunc.
func (mock *MockUserRepository) CreateUser(user *models.User) (int, error) {
	if mock.CreateUserFunc == nil {
		panic("MockUserRepository.CreateUserFunc: method is nil but UserRepository.CreateUser was just called")
	}
	callInfo := struct {
		User *models.User
	}{
		User: user,
	}
	mock.lockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	mock.lockCreateUser.Unlock()
	return mock.CreateUserFunc(user)
}

// CreateUserCalls gets all the calls that were made to CreateUser.
// Check the length with:
//     len(mockedUserRepository.CreateUserCalls())
func (mock *MockUserRepository) CreateUserCalls() []struct {
	User *models.User
} {
	var calls []struct {
		User *models.User
	}
	mock.lockCreateUser.RLock()
	calls = mock.calls.CreateUser
	mock.lockCreateUser.RUnlock()
	return calls
}

// DoesUserExist calls DoesUserExistFunc.
func (mock *MockUserRepository) DoesUserExist(auth *models.Auth) (int, error) {
	if mock.DoesUserExistFunc == nil {
		panic("MockUserRepository.DoesUserExistFunc: method is nil but UserRepository.DoesUserExist was just called")
	}
	callInfo := struct {
		Auth *models.Auth
	}{
		Auth: auth,
	}
	mock.lockDoesUserExist.Lock()
	mock.calls.DoesUserExist = append(mock.calls.DoesUserExist, callInfo)
	mock.lockDoesUserExist.Unlock()
	return mock.DoesUserExistFunc(auth)
}

// DoesUserExistCalls gets all the calls that were made to DoesUserExist.
// Check the length with:
//     len(mockedUserRepository.DoesUserExistCalls())
func (mock *MockUserRepository) DoesUserExistCalls() []struct {
	Auth *models.Auth
} {
	var calls []struct {
		Auth *models.Auth
	}
	mock.lockDoesUserExist.RLock()
	calls = mock.calls.DoesUserExist
	mock.lockDoesUserExist.RUnlock()
	return calls
}

// GetAvatarFilename calls GetAvatarFilenameFunc.
func (mock *MockUserRepository) GetAvatarFilename(n int) (string, error) {
	if mock.GetAvatarFilenameFunc == nil {
		panic("MockUserRepository.GetAvatarFilenameFunc: method is nil but UserRepository.GetAvatarFilename was just called")
	}
	callInfo := struct {
		N int
	}{
		N: n,
	}
	mock.lockGetAvatarFilename.Lock()
	mock.calls.GetAvatarFilename = append(mock.calls.GetAvatarFilename, callInfo)
	mock.lockGetAvatarFilename.Unlock()
	return mock.GetAvatarFilenameFunc(n)
}

// GetAvatarFilenameCalls gets all the calls that were made to GetAvatarFilename.
// Check the length with:
//     len(mockedUserRepository.GetAvatarFilenameCalls())
func (mock *MockUserRepository) GetAvatarFilenameCalls() []struct {
	N int
} {
	var calls []struct {
		N int
	}
	mock.lockGetAvatarFilename.RLock()
	calls = mock.calls.GetAvatarFilename
	mock.lockGetAvatarFilename.RUnlock()
	return calls
}

// GetSettings calls GetSettingsFunc.
func (mock *MockUserRepository) GetSettings(n int) (*models.SettingsGet, error) {
	if mock.GetSettingsFunc == nil {
		panic("MockUserRepository.GetSettingsFunc: method is nil but UserRepository.GetSettings was just called")
	}
	callInfo := struct {
		N int
	}{
		N: n,
	}
	mock.lockGetSettings.Lock()
	mock.calls.GetSettings = append(mock.calls.GetSettings, callInfo)
	mock.lockGetSettings.Unlock()
	return mock.GetSettingsFunc(n)
}

// GetSettingsCalls gets all the calls that were made to GetSettings.
// Check the length with:
//     len(mockedUserRepository.GetSettingsCalls())
func (mock *MockUserRepository) GetSettingsCalls() []struct {
	N int
} {
	var calls []struct {
		N int
	}
	mock.lockGetSettings.RLock()
	calls = mock.calls.GetSettings
	mock.lockGetSettings.RUnlock()
	return calls
}

// IsEmailUnique calls IsEmailUniqueFunc.
func (mock *MockUserRepository) IsEmailUnique(s string) (bool, error) {
	if mock.IsEmailUniqueFunc == nil {
		panic("MockUserRepository.IsEmailUniqueFunc: method is nil but UserRepository.IsEmailUnique was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockIsEmailUnique.Lock()
	mock.calls.IsEmailUnique = append(mock.calls.IsEmailUnique, callInfo)
	mock.lockIsEmailUnique.Unlock()
	return mock.IsEmailUniqueFunc(s)
}

// IsEmailUniqueCalls gets all the calls that were made to IsEmailUnique.
// Check the length with:
//     len(mockedUserRepository.IsEmailUniqueCalls())
func (mock *MockUserRepository) IsEmailUniqueCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockIsEmailUnique.RLock()
	calls = mock.calls.IsEmailUnique
	mock.lockIsEmailUnique.RUnlock()
	return calls
}

// IsNicknameUnique calls IsNicknameUniqueFunc.
func (mock *MockUserRepository) IsNicknameUnique(s string) (bool, error) {
	if mock.IsNicknameUniqueFunc == nil {
		panic("MockUserRepository.IsNicknameUniqueFunc: method is nil but UserRepository.IsNicknameUnique was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockIsNicknameUnique.Lock()
	mock.calls.IsNicknameUnique = append(mock.calls.IsNicknameUnique, callInfo)
	mock.lockIsNicknameUnique.Unlock()
	return mock.IsNicknameUniqueFunc(s)
}

// IsNicknameUniqueCalls gets all the calls that were made to IsNicknameUnique.
// Check the length with:
//     len(mockedUserRepository.IsNicknameUniqueCalls())
func (mock *MockUserRepository) IsNicknameUniqueCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockIsNicknameUnique.RLock()
	calls = mock.calls.IsNicknameUnique
	mock.lockIsNicknameUnique.RUnlock()
	return calls
}

// UpdateAvatar calls UpdateAvatarFunc.
func (mock *MockUserRepository) UpdateAvatar(n int, s string) error {
	if mock.UpdateAvatarFunc == nil {
		panic("MockUserRepository.UpdateAvatarFunc: method is nil but UserRepository.UpdateAvatar was just called")
	}
	callInfo := struct {
		N int
		S string
	}{
		N: n,
		S: s,
	}
	mock.lockUpdateAvatar.Lock()
	mock.calls.UpdateAvatar = append(mock.calls.UpdateAvatar, callInfo)
	mock.lockUpdateAvatar.Unlock()
	return mock.UpdateAvatarFunc(n, s)
}

// UpdateAvatarCalls gets all the calls that were made to UpdateAvatar.
// Check the length with:
//     len(mockedUserRepository.UpdateAvatarCalls())
func (mock *MockUserRepository) UpdateAvatarCalls() []struct {
	N int
	S string
} {
	var calls []struct {
		N int
		S string
	}
	mock.lockUpdateAvatar.RLock()
	calls = mock.calls.UpdateAvatar
	mock.lockUpdateAvatar.RUnlock()
	return calls
}

// UpdateEmail calls UpdateEmailFunc.
func (mock *MockUserRepository) UpdateEmail(n int, s string) error {
	if mock.UpdateEmailFunc == nil {
		panic("MockUserRepository.UpdateEmailFunc: method is nil but UserRepository.UpdateEmail was just called")
	}
	callInfo := struct {
		N int
		S string
	}{
		N: n,
		S: s,
	}
	mock.lockUpdateEmail.Lock()
	mock.calls.UpdateEmail = append(mock.calls.UpdateEmail, callInfo)
	mock.lockUpdateEmail.Unlock()
	return mock.UpdateEmailFunc(n, s)
}

// UpdateEmailCalls gets all the calls that were made to UpdateEmail.
// Check the length with:
//     len(mockedUserRepository.UpdateEmailCalls())
func (mock *MockUserRepository) UpdateEmailCalls() []struct {
	N int
	S string
} {
	var calls []struct {
		N int
		S string
	}
	mock.lockUpdateEmail.RLock()
	calls = mock.calls.UpdateEmail
	mock.lockUpdateEmail.RUnlock()
	return calls
}

// UpdateNickname calls UpdateNicknameFunc.
func (mock *MockUserRepository) UpdateNickname(n int, s string) error {
	if mock.UpdateNicknameFunc == nil {
		panic("MockUserRepository.UpdateNicknameFunc: method is nil but UserRepository.UpdateNickname was just called")
	}
	callInfo := struct {
		N int
		S string
	}{
		N: n,
		S: s,
	}
	mock.lockUpdateNickname.Lock()
	mock.calls.UpdateNickname = append(mock.calls.UpdateNickname, callInfo)
	mock.lockUpdateNickname.Unlock()
	return mock.UpdateNicknameFunc(n, s)
}

// UpdateNicknameCalls gets all the calls that were made to UpdateNickname.
// Check the length with:
//     len(mockedUserRepository.UpdateNicknameCalls())
func (mock *MockUserRepository) UpdateNicknameCalls() []struct {
	N int
	S string
} {
	var calls []struct {
		N int
		S string
	}
	mock.lockUpdateNickname.RLock()
	calls = mock.calls.UpdateNickname
	mock.lockUpdateNickname.RUnlock()
	return calls
}

// UpdatePassword calls UpdatePasswordFunc.
func (mock *MockUserRepository) UpdatePassword(n int, s string) error {
	if mock.UpdatePasswordFunc == nil {
		panic("MockUserRepository.UpdatePasswordFunc: method is nil but UserRepository.UpdatePassword was just called")
	}
	callInfo := struct {
		N int
		S string
	}{
		N: n,
		S: s,
	}
	mock.lockUpdatePassword.Lock()
	mock.calls.UpdatePassword = append(mock.calls.UpdatePassword, callInfo)
	mock.lockUpdatePassword.Unlock()
	return mock.UpdatePasswordFunc(n, s)
}

// UpdatePasswordCalls gets all the calls that were made to UpdatePassword.
// Check the length with:
//     len(mockedUserRepository.UpdatePasswordCalls())
func (mock *MockUserRepository) UpdatePasswordCalls() []struct {
	N int
	S string
} {
	var calls []struct {
		N int
		S string
	}
	mock.lockUpdatePassword.RLock()
	calls = mock.calls.UpdatePassword
	mock.lockUpdatePassword.RUnlock()
	return calls
}