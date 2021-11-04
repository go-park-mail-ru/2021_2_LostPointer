// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"sync"
)

// Ensure, that MockUserUseCase does implement users.UserUseCase.
// If this is not the case, regenerate this file with moq.
var _ users.UserUseCase = &MockUserUseCase{}

// MockUserUseCase is a mock implementation of users.UserUseCase.
//
// 	func TestSomethingThatUsesUserUseCase(t *testing.T) {
//
// 		// make and configure a mocked users.UserUseCase
// 		mockedUserUseCase := &MockUserUseCase{
// 			GetAvatarFilenameFunc: func(n int) (string, *models.CustomError) {
// 				panic("mock out the GetAvatarFilename method")
// 			},
// 			GetSettingsFunc: func(n int) (*models.SettingsGet, *models.CustomError) {
// 				panic("mock out the GetSettings method")
// 			},
// 			LoginFunc: func(auth *models.Auth) (string, *models.CustomError) {
// 				panic("mock out the Login method")
// 			},
// 			LogoutFunc: func(s string) error {
// 				panic("mock out the Logout method")
// 			},
// 			RegisterFunc: func(user *models.User) (string, *models.CustomError) {
// 				panic("mock out the Register method")
// 			},
// 			UpdateSettingsFunc: func(n int, settingsGet *models.SettingsGet, settingsUpload *models.SettingsUpload) *models.CustomError {
// 				panic("mock out the UpdateSettings method")
// 			},
// 		}
//
// 		// use mockedUserUseCase in code that requires users.UserUseCase
// 		// and then make assertions.
//
// 	}
type MockUserUseCase struct {
	// GetAvatarFilenameFunc mocks the GetAvatarFilename method.
	GetAvatarFilenameFunc func(n int) (string, *models.CustomError)

	// GetSettingsFunc mocks the GetSettings method.
	GetSettingsFunc func(n int) (*models.SettingsGet, *models.CustomError)

	// LoginFunc mocks the Login method.
	LoginFunc func(auth *models.Auth) (string, *models.CustomError)

	// LogoutFunc mocks the Logout method.
	LogoutFunc func(s string) error

	// RegisterFunc mocks the Register method.
	RegisterFunc func(user *models.User) (string, *models.CustomError)

	// UpdateSettingsFunc mocks the UpdateSettings method.
	UpdateSettingsFunc func(n int, settingsGet *models.SettingsGet, settingsUpload *models.SettingsUpload) *models.CustomError

	// calls tracks calls to the methods.
	calls struct {
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
		// Login holds details about calls to the Login method.
		Login []struct {
			// Auth is the auth argument value.
			Auth *models.Auth
		}
		// Logout holds details about calls to the Logout method.
		Logout []struct {
			// S is the s argument value.
			S string
		}
		// Register holds details about calls to the Register method.
		Register []struct {
			// User is the user argument value.
			User *models.User
		}
		// UpdateSettings holds details about calls to the UpdateSettings method.
		UpdateSettings []struct {
			// N is the n argument value.
			N int
			// SettingsGet is the settingsGet argument value.
			SettingsGet *models.SettingsGet
			// SettingsUpload is the settingsUpload argument value.
			SettingsUpload *models.SettingsUpload
		}
	}
	lockGetAvatarFilename sync.RWMutex
	lockGetSettings       sync.RWMutex
	lockLogin             sync.RWMutex
	lockLogout            sync.RWMutex
	lockRegister          sync.RWMutex
	lockUpdateSettings    sync.RWMutex
}

// GetAvatarFilename calls GetAvatarFilenameFunc.
func (mock *MockUserUseCase) GetAvatarFilename(n int) (string, *models.CustomError) {
	if mock.GetAvatarFilenameFunc == nil {
		panic("MockUserUseCase.GetAvatarFilenameFunc: method is nil but UserUseCase.GetAvatarFilename was just called")
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
//     len(mockedUserUseCase.GetAvatarFilenameCalls())
func (mock *MockUserUseCase) GetAvatarFilenameCalls() []struct {
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
func (mock *MockUserUseCase) GetSettings(n int) (*models.SettingsGet, *models.CustomError) {
	if mock.GetSettingsFunc == nil {
		panic("MockUserUseCase.GetSettingsFunc: method is nil but UserUseCase.GetSettings was just called")
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
//     len(mockedUserUseCase.GetSettingsCalls())
func (mock *MockUserUseCase) GetSettingsCalls() []struct {
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

// Login calls LoginFunc.
func (mock *MockUserUseCase) Login(auth *models.Auth) (string, *models.CustomError) {
	if mock.LoginFunc == nil {
		panic("MockUserUseCase.LoginFunc: method is nil but UserUseCase.Login was just called")
	}
	callInfo := struct {
		Auth *models.Auth
	}{
		Auth: auth,
	}
	mock.lockLogin.Lock()
	mock.calls.Login = append(mock.calls.Login, callInfo)
	mock.lockLogin.Unlock()
	return mock.LoginFunc(auth)
}

// LoginCalls gets all the calls that were made to Login.
// Check the length with:
//     len(mockedUserUseCase.LoginCalls())
func (mock *MockUserUseCase) LoginCalls() []struct {
	Auth *models.Auth
} {
	var calls []struct {
		Auth *models.Auth
	}
	mock.lockLogin.RLock()
	calls = mock.calls.Login
	mock.lockLogin.RUnlock()
	return calls
}

// Logout calls LogoutFunc.
func (mock *MockUserUseCase) Logout(s string) error {
	if mock.LogoutFunc == nil {
		panic("MockUserUseCase.LogoutFunc: method is nil but UserUseCase.Logout was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockLogout.Lock()
	mock.calls.Logout = append(mock.calls.Logout, callInfo)
	mock.lockLogout.Unlock()
	return mock.LogoutFunc(s)
}

// LogoutCalls gets all the calls that were made to Logout.
// Check the length with:
//     len(mockedUserUseCase.LogoutCalls())
func (mock *MockUserUseCase) LogoutCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockLogout.RLock()
	calls = mock.calls.Logout
	mock.lockLogout.RUnlock()
	return calls
}

// Register calls RegisterFunc.
func (mock *MockUserUseCase) Register(user *models.User) (string, *models.CustomError) {
	if mock.RegisterFunc == nil {
		panic("MockUserUseCase.RegisterFunc: method is nil but UserUseCase.Register was just called")
	}
	callInfo := struct {
		User *models.User
	}{
		User: user,
	}
	mock.lockRegister.Lock()
	mock.calls.Register = append(mock.calls.Register, callInfo)
	mock.lockRegister.Unlock()
	return mock.RegisterFunc(user)
}

// RegisterCalls gets all the calls that were made to Register.
// Check the length with:
//     len(mockedUserUseCase.RegisterCalls())
func (mock *MockUserUseCase) RegisterCalls() []struct {
	User *models.User
} {
	var calls []struct {
		User *models.User
	}
	mock.lockRegister.RLock()
	calls = mock.calls.Register
	mock.lockRegister.RUnlock()
	return calls
}

// UpdateSettings calls UpdateSettingsFunc.
func (mock *MockUserUseCase) UpdateSettings(n int, settingsGet *models.SettingsGet, settingsUpload *models.SettingsUpload) *models.CustomError {
	if mock.UpdateSettingsFunc == nil {
		panic("MockUserUseCase.UpdateSettingsFunc: method is nil but UserUseCase.UpdateSettings was just called")
	}
	callInfo := struct {
		N              int
		SettingsGet    *models.SettingsGet
		SettingsUpload *models.SettingsUpload
	}{
		N:              n,
		SettingsGet:    settingsGet,
		SettingsUpload: settingsUpload,
	}
	mock.lockUpdateSettings.Lock()
	mock.calls.UpdateSettings = append(mock.calls.UpdateSettings, callInfo)
	mock.lockUpdateSettings.Unlock()
	return mock.UpdateSettingsFunc(n, settingsGet, settingsUpload)
}

// UpdateSettingsCalls gets all the calls that were made to UpdateSettings.
// Check the length with:
//     len(mockedUserUseCase.UpdateSettingsCalls())
func (mock *MockUserUseCase) UpdateSettingsCalls() []struct {
	N              int
	SettingsGet    *models.SettingsGet
	SettingsUpload *models.SettingsUpload
} {
	var calls []struct {
		N              int
		SettingsGet    *models.SettingsGet
		SettingsUpload *models.SettingsUpload
	}
	mock.lockUpdateSettings.RLock()
	calls = mock.calls.UpdateSettings
	mock.lockUpdateSettings.RUnlock()
	return calls
}