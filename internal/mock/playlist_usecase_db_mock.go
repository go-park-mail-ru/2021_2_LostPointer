// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/playlist"
	"sync"
)

// Ensure, that MockPlaylistUseCase does implement playlist.PlaylistUseCase.
// If this is not the case, regenerate this file with moq.
var _ playlist.PlaylistUseCase = &MockPlaylistUseCase{}

// MockPlaylistUseCase is a mock implementation of playlist.PlaylistUseCase.
//
// 	func TestSomethingThatUsesPlaylistUseCase(t *testing.T) {
//
// 		// make and configure a mocked playlist.PlaylistUseCase
// 		mockedPlaylistUseCase := &MockPlaylistUseCase{
// 			GetHomeFunc: func(amount int) ([]models.Playlist, *models.CustomError) {
// 				panic("mock out the GetHome method")
// 			},
// 		}
//
// 		// use mockedPlaylistUseCase in code that requires playlist.PlaylistUseCase
// 		// and then make assertions.
//
// 	}
type MockPlaylistUseCase struct {
	// GetHomeFunc mocks the GetHome method.
	GetHomeFunc func(amount int) ([]models.Playlist, *models.CustomError)

	// calls tracks calls to the methods.
	calls struct {
		// GetHome holds details about calls to the GetHome method.
		GetHome []struct {
			// Amount is the amount argument value.
			Amount int
		}
	}
	lockGetHome sync.RWMutex
}

// GetHome calls GetHomeFunc.
func (mock *MockPlaylistUseCase) GetHome(amount int) ([]models.Playlist, *models.CustomError) {
	if mock.GetHomeFunc == nil {
		panic("MockPlaylistUseCase.GetHomeFunc: method is nil but PlaylistUseCase.GetHome was just called")
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
//     len(mockedPlaylistUseCase.GetHomeCalls())
func (mock *MockPlaylistUseCase) GetHomeCalls() []struct {
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
