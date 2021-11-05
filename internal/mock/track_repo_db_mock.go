// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"sync"
)

// Ensure, that MockTrackRepository does implement track.TrackRepository.
// If this is not the case, regenerate this file with moq.
var _ track.TrackRepository = &MockTrackRepository{}

// MockTrackRepository is a mock implementation of track.TrackRepository.
//
// 	func TestSomethingThatUsesTrackRepository(t *testing.T) {
//
// 		// make and configure a mocked track.TrackRepository
// 		mockedTrackRepository := &MockTrackRepository{
// 			GetRandomFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
// 				panic("mock out the GetRandom method")
// 			},
// 			IncrementListenCountFunc: func(n int64) error {
// 				panic("mock out the IncrementListenCount method")
// 			},
// 		}
//
// 		// use mockedTrackRepository in code that requires track.TrackRepository
// 		// and then make assertions.
//
// 	}
type MockTrackRepository struct {
	// GetRandomFunc mocks the GetRandom method.
	GetRandomFunc func(amount int, isAuthorized bool) ([]models.Track, error)

	// IncrementListenCountFunc mocks the IncrementListenCount method.
	IncrementListenCountFunc func(n int64) error

	// calls tracks calls to the methods.
	calls struct {
		// GetRandom holds details about calls to the GetRandom method.
		GetRandom []struct {
			// Amount is the amount argument value.
			Amount int
			// IsAuthorized is the isAuthorized argument value.
			IsAuthorized bool
		}
		// IncrementListenCount holds details about calls to the IncrementListenCount method.
		IncrementListenCount []struct {
			// N is the n argument value.
			N int64
		}
	}
	lockGetRandom            sync.RWMutex
	lockIncrementListenCount sync.RWMutex
}

// GetRandom calls GetRandomFunc.
func (mock *MockTrackRepository) GetRandom(amount int, isAuthorized bool) ([]models.Track, error) {
	if mock.GetRandomFunc == nil {
		panic("MockTrackRepository.GetRandomFunc: method is nil but TrackRepository.GetRandom was just called")
	}
	callInfo := struct {
		Amount       int
		IsAuthorized bool
	}{
		Amount:       amount,
		IsAuthorized: isAuthorized,
	}
	mock.lockGetRandom.Lock()
	mock.calls.GetRandom = append(mock.calls.GetRandom, callInfo)
	mock.lockGetRandom.Unlock()
	return mock.GetRandomFunc(amount, isAuthorized)
}

// GetRandomCalls gets all the calls that were made to GetRandom.
// Check the length with:
//     len(mockedTrackRepository.GetRandomCalls())
func (mock *MockTrackRepository) GetRandomCalls() []struct {
	Amount       int
	IsAuthorized bool
} {
	var calls []struct {
		Amount       int
		IsAuthorized bool
	}
	mock.lockGetRandom.RLock()
	calls = mock.calls.GetRandom
	mock.lockGetRandom.RUnlock()
	return calls
}

// IncrementListenCount calls IncrementListenCountFunc.
func (mock *MockTrackRepository) IncrementListenCount(n int64) error {
	if mock.IncrementListenCountFunc == nil {
		panic("MockTrackRepository.IncrementListenCountFunc: method is nil but TrackRepository.IncrementListenCount was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockIncrementListenCount.Lock()
	mock.calls.IncrementListenCount = append(mock.calls.IncrementListenCount, callInfo)
	mock.lockIncrementListenCount.Unlock()
	return mock.IncrementListenCountFunc(n)
}

// IncrementListenCountCalls gets all the calls that were made to IncrementListenCount.
// Check the length with:
//     len(mockedTrackRepository.IncrementListenCountCalls())
func (mock *MockTrackRepository) IncrementListenCountCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockIncrementListenCount.RLock()
	calls = mock.calls.IncrementListenCount
	mock.lockIncrementListenCount.RUnlock()
	return calls
}