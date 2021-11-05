// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/album"
	"2021_2_LostPointer/internal/models"
	"sync"
)

// Ensure, that MockAlbumUseCase does implement album.AlbumUseCase.
// If this is not the case, regenerate this file with moq.
var _ album.AlbumUseCase = &MockAlbumUseCase{}

// MockAlbumUseCase is a mock implementation of album.AlbumUseCase.
//
// 	func TestSomethingThatUsesAlbumUseCase(t *testing.T) {
//
// 		// make and configure a mocked album.AlbumUseCase
// 		mockedAlbumUseCase := &MockAlbumUseCase{
// 			GetHomeFunc: func(amount int) ([]models.Album, *models.CustomError) {
// 				panic("mock out the GetHome method")
// 			},
// 		}
//
// 		// use mockedAlbumUseCase in code that requires album.AlbumUseCase
// 		// and then make assertions.
//
// 	}
type MockAlbumUseCase struct {
	// GetHomeFunc mocks the GetHome method.
	GetHomeFunc func(amount int) ([]models.Album, *models.CustomError)

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
func (mock *MockAlbumUseCase) GetHome(amount int) ([]models.Album, *models.CustomError) {
	if mock.GetHomeFunc == nil {
		panic("MockAlbumUseCase.GetHomeFunc: method is nil but AlbumUseCase.GetHome was just called")
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
//     len(mockedAlbumUseCase.GetHomeCalls())
func (mock *MockAlbumUseCase) GetHomeCalls() []struct {
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