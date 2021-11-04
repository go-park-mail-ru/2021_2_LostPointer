// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/models"
	"sync"
)

// Ensure, that MockArtistUseCase does implement artist.ArtistUseCase.
// If this is not the case, regenerate this file with moq.
var _ artist.ArtistUseCase = &MockArtistUseCase{}

// MockArtistUseCase is a mock implementation of artist.ArtistUseCase.
//
// 	func TestSomethingThatUsesArtistUseCase(t *testing.T) {
//
// 		// make and configure a mocked artist.ArtistUseCase
// 		mockedArtistUseCase := &MockArtistUseCase{
// 			GetHomeFunc: func(amount int) ([]models.Artist, *models.CustomError) {
// 				panic("mock out the GetHome method")
// 			},
// 			GetProfileFunc: func(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
// 				panic("mock out the GetProfile method")
// 			},
// 		}
//
// 		// use mockedArtistUseCase in code that requires artist.ArtistUseCase
// 		// and then make assertions.
//
// 	}
type MockArtistUseCase struct {
	// GetHomeFunc mocks the GetHome method.
	GetHomeFunc func(amount int) ([]models.Artist, *models.CustomError)

	// GetProfileFunc mocks the GetProfile method.
	GetProfileFunc func(id int, isAuthorized bool) (*models.Artist, *models.CustomError)

	// calls tracks calls to the methods.
	calls struct {
		// GetHome holds details about calls to the GetHome method.
		GetHome []struct {
			// Amount is the amount argument value.
			Amount int
		}
		// GetProfile holds details about calls to the GetProfile method.
		GetProfile []struct {
			// ID is the id argument value.
			ID int
			// IsAuthorized is the isAuthorized argument value.
			IsAuthorized bool
		}
	}
	lockGetHome    sync.RWMutex
	lockGetProfile sync.RWMutex
}

// GetHome calls GetHomeFunc.
func (mock *MockArtistUseCase) GetHome(amount int) ([]models.Artist, *models.CustomError) {
	if mock.GetHomeFunc == nil {
		panic("MockArtistUseCase.GetHomeFunc: method is nil but ArtistUseCase.GetHome was just called")
	}
	callInfo := struct {
		Amount int
	}{
		Amount: amount,
	}
	mock.lockGetHome.Lock()
	mock.calls.GetHome = append(mock.calls.GetHome, callInfo)
	mock.lockGetHome.Unlock()
	return mock.GetHomeFunc(amount)
}

// GetHomeCalls gets all the calls that were made to GetHome.
// Check the length with:
//     len(mockedArtistUseCase.GetHomeCalls())
func (mock *MockArtistUseCase) GetHomeCalls() []struct {
	Amount int
} {
	var calls []struct {
		Amount int
	}
	mock.lockGetHome.RLock()
	calls = mock.calls.GetHome
	mock.lockGetHome.RUnlock()
	return calls
}

// GetProfile calls GetProfileFunc.
func (mock *MockArtistUseCase) GetProfile(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
	if mock.GetProfileFunc == nil {
		panic("MockArtistUseCase.GetProfileFunc: method is nil but ArtistUseCase.GetProfile was just called")
	}
	callInfo := struct {
		ID           int
		IsAuthorized bool
	}{
		ID:           id,
		IsAuthorized: isAuthorized,
	}
	mock.lockGetProfile.Lock()
	mock.calls.GetProfile = append(mock.calls.GetProfile, callInfo)
	mock.lockGetProfile.Unlock()
	return mock.GetProfileFunc(id, isAuthorized)
}

// GetProfileCalls gets all the calls that were made to GetProfile.
// Check the length with:
//     len(mockedArtistUseCase.GetProfileCalls())
func (mock *MockArtistUseCase) GetProfileCalls() []struct {
	ID           int
	IsAuthorized bool
} {
	var calls []struct {
		ID           int
		IsAuthorized bool
	}
	mock.lockGetProfile.RLock()
	calls = mock.calls.GetProfile
	mock.lockGetProfile.RUnlock()
	return calls
}