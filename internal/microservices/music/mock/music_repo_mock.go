// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"2021_2_LostPointer/internal/microservices/music"
	"2021_2_LostPointer/internal/microservices/music/proto"
	"sync"
)

// Ensure, that MockMusicStorage does implement music.MusicStorage.
// If this is not the case, regenerate this file with moq.
var _ music.MusicStorage = &MockMusicStorage{}

// MockMusicStorage is a mock implementation of music.MusicStorage.
//
// 	func TestSomethingThatUsesMusicStorage(t *testing.T) {
//
// 		// make and configure a mocked music.MusicStorage
// 		mockedMusicStorage := &MockMusicStorage{
// 			AlbumDataFunc: func(n int64) (*proto.AlbumPageResponse, error) {
// 				panic("mock out the AlbumData method")
// 			},
// 			AlbumTracksFunc: func(n int64, b bool) ([]*proto.AlbumTrack, error) {
// 				panic("mock out the AlbumTracks method")
// 			},
// 			ArtistAlbumsFunc: func(n1 int64, n2 int64) ([]*proto.Album, error) {
// 				panic("mock out the ArtistAlbums method")
// 			},
// 			ArtistInfoFunc: func(n int64) (*proto.Artist, error) {
// 				panic("mock out the ArtistInfo method")
// 			},
// 			ArtistTracksFunc: func(n1 int64, b bool, n2 int64) ([]*proto.Track, error) {
// 				panic("mock out the ArtistTracks method")
// 			},
// 			DoesPlaylistExistFunc: func(n int64) (bool, error) {
// 				panic("mock out the DoesPlaylistExist method")
// 			},
// 			FindAlbumsFunc: func(s string) ([]*proto.Album, error) {
// 				panic("mock out the FindAlbums method")
// 			},
// 			FindArtistsFunc: func(s string) ([]*proto.Artist, error) {
// 				panic("mock out the FindArtists method")
// 			},
// 			FindTracksByFullWordFunc: func(s string, b bool) ([]*proto.Track, error) {
// 				panic("mock out the FindTracksByFullWord method")
// 			},
// 			FindTracksByPartialFunc: func(s string, b bool) ([]*proto.Track, error) {
// 				panic("mock out the FindTracksByPartial method")
// 			},
// 			IncrementListenCountFunc: func(n int64) error {
// 				panic("mock out the IncrementListenCount method")
// 			},
// 			IsPlaylistOwnerFunc: func(n1 int64, n2 int64) (bool, error) {
// 				panic("mock out the IsPlaylistOwner method")
// 			},
// 			PlaylistInfoFunc: func(n int64) (*proto.PlaylistData, error) {
// 				panic("mock out the PlaylistInfo method")
// 			},
// 			PlaylistTracksFunc: func(n int64) ([]*proto.Track, error) {
// 				panic("mock out the PlaylistTracks method")
// 			},
// 			RandomAlbumsFunc: func(n int64) (*proto.Albums, error) {
// 				panic("mock out the RandomAlbums method")
// 			},
// 			RandomArtistsFunc: func(n int64) (*proto.Artists, error) {
// 				panic("mock out the RandomArtists method")
// 			},
// 			RandomTracksFunc: func(n int64, b bool) (*proto.Tracks, error) {
// 				panic("mock out the RandomTracks method")
// 			},
// 			UserPlaylistsFunc: func(n int64) ([]*proto.PlaylistData, error) {
// 				panic("mock out the UserPlaylists method")
// 			},
// 		}
//
// 		// use mockedMusicStorage in code that requires music.MusicStorage
// 		// and then make assertions.
//
// 	}
type MockMusicStorage struct {
	// AlbumDataFunc mocks the AlbumData method.
	AlbumDataFunc func(n int64) (*proto.AlbumPageResponse, error)

	// AlbumTracksFunc mocks the AlbumTracks method.
	AlbumTracksFunc func(n int64, b bool) ([]*proto.AlbumTrack, error)

	// ArtistAlbumsFunc mocks the ArtistAlbums method.
	ArtistAlbumsFunc func(n1 int64, n2 int64) ([]*proto.Album, error)

	// ArtistInfoFunc mocks the ArtistInfo method.
	ArtistInfoFunc func(n int64) (*proto.Artist, error)

	// ArtistTracksFunc mocks the ArtistTracks method.
	ArtistTracksFunc func(n1 int64, b bool, n2 int64) ([]*proto.Track, error)

	// DoesPlaylistExistFunc mocks the DoesPlaylistExist method.
	DoesPlaylistExistFunc func(n int64) (bool, error)

	// FindAlbumsFunc mocks the FindAlbums method.
	FindAlbumsFunc func(s string) ([]*proto.Album, error)

	// FindArtistsFunc mocks the FindArtists method.
	FindArtistsFunc func(s string) ([]*proto.Artist, error)

	// FindTracksByFullWordFunc mocks the FindTracksByFullWord method.
	FindTracksByFullWordFunc func(s string, b bool) ([]*proto.Track, error)

	// FindTracksByPartialFunc mocks the FindTracksByPartial method.
	FindTracksByPartialFunc func(s string, b bool) ([]*proto.Track, error)

	// IncrementListenCountFunc mocks the IncrementListenCount method.
	IncrementListenCountFunc func(n int64) error

	// IsPlaylistOwnerFunc mocks the IsPlaylistOwner method.
	IsPlaylistOwnerFunc func(n1 int64, n2 int64) (bool, error)

	// PlaylistInfoFunc mocks the PlaylistInfo method.
	PlaylistInfoFunc func(n int64) (*proto.PlaylistData, error)

	// PlaylistTracksFunc mocks the PlaylistTracks method.
	PlaylistTracksFunc func(n int64) ([]*proto.Track, error)

	// RandomAlbumsFunc mocks the RandomAlbums method.
	RandomAlbumsFunc func(n int64) (*proto.Albums, error)

	// RandomArtistsFunc mocks the RandomArtists method.
	RandomArtistsFunc func(n int64) (*proto.Artists, error)

	// RandomTracksFunc mocks the RandomTracks method.
	RandomTracksFunc func(n int64, b bool) (*proto.Tracks, error)

	// UserPlaylistsFunc mocks the UserPlaylists method.
	UserPlaylistsFunc func(n int64) ([]*proto.PlaylistData, error)

	// calls tracks calls to the methods.
	calls struct {
		// AlbumData holds details about calls to the AlbumData method.
		AlbumData []struct {
			// N is the n argument value.
			N int64
		}
		// AlbumTracks holds details about calls to the AlbumTracks method.
		AlbumTracks []struct {
			// N is the n argument value.
			N int64
			// B is the b argument value.
			B bool
		}
		// ArtistAlbums holds details about calls to the ArtistAlbums method.
		ArtistAlbums []struct {
			// N1 is the n1 argument value.
			N1 int64
			// N2 is the n2 argument value.
			N2 int64
		}
		// ArtistInfo holds details about calls to the ArtistInfo method.
		ArtistInfo []struct {
			// N is the n argument value.
			N int64
		}
		// ArtistTracks holds details about calls to the ArtistTracks method.
		ArtistTracks []struct {
			// N1 is the n1 argument value.
			N1 int64
			// B is the b argument value.
			B bool
			// N2 is the n2 argument value.
			N2 int64
		}
		// DoesPlaylistExist holds details about calls to the DoesPlaylistExist method.
		DoesPlaylistExist []struct {
			// N is the n argument value.
			N int64
		}
		// FindAlbums holds details about calls to the FindAlbums method.
		FindAlbums []struct {
			// S is the s argument value.
			S string
		}
		// FindArtists holds details about calls to the FindArtists method.
		FindArtists []struct {
			// S is the s argument value.
			S string
		}
		// FindTracksByFullWord holds details about calls to the FindTracksByFullWord method.
		FindTracksByFullWord []struct {
			// S is the s argument value.
			S string
			// B is the b argument value.
			B bool
		}
		// FindTracksByPartial holds details about calls to the FindTracksByPartial method.
		FindTracksByPartial []struct {
			// S is the s argument value.
			S string
			// B is the b argument value.
			B bool
		}
		// IncrementListenCount holds details about calls to the IncrementListenCount method.
		IncrementListenCount []struct {
			// N is the n argument value.
			N int64
		}
		// IsPlaylistOwner holds details about calls to the IsPlaylistOwner method.
		IsPlaylistOwner []struct {
			// N1 is the n1 argument value.
			N1 int64
			// N2 is the n2 argument value.
			N2 int64
		}
		// PlaylistInfo holds details about calls to the PlaylistInfo method.
		PlaylistInfo []struct {
			// N is the n argument value.
			N int64
		}
		// PlaylistTracks holds details about calls to the PlaylistTracks method.
		PlaylistTracks []struct {
			// N is the n argument value.
			N int64
		}
		// RandomAlbums holds details about calls to the RandomAlbums method.
		RandomAlbums []struct {
			// N is the n argument value.
			N int64
		}
		// RandomArtists holds details about calls to the RandomArtists method.
		RandomArtists []struct {
			// N is the n argument value.
			N int64
		}
		// RandomTracks holds details about calls to the RandomTracks method.
		RandomTracks []struct {
			// N is the n argument value.
			N int64
			// B is the b argument value.
			B bool
		}
		// UserPlaylists holds details about calls to the UserPlaylists method.
		UserPlaylists []struct {
			// N is the n argument value.
			N int64
		}
	}
	lockAlbumData            sync.RWMutex
	lockAlbumTracks          sync.RWMutex
	lockArtistAlbums         sync.RWMutex
	lockArtistInfo           sync.RWMutex
	lockArtistTracks         sync.RWMutex
	lockDoesPlaylistExist    sync.RWMutex
	lockFindAlbums           sync.RWMutex
	lockFindArtists          sync.RWMutex
	lockFindTracksByFullWord sync.RWMutex
	lockFindTracksByPartial  sync.RWMutex
	lockIncrementListenCount sync.RWMutex
	lockIsPlaylistOwner      sync.RWMutex
	lockPlaylistInfo         sync.RWMutex
	lockPlaylistTracks       sync.RWMutex
	lockRandomAlbums         sync.RWMutex
	lockRandomArtists        sync.RWMutex
	lockRandomTracks         sync.RWMutex
	lockUserPlaylists        sync.RWMutex
}

// AlbumData calls AlbumDataFunc.
func (mock *MockMusicStorage) AlbumData(n int64) (*proto.AlbumPageResponse, error) {
	if mock.AlbumDataFunc == nil {
		panic("MockMusicStorage.AlbumDataFunc: method is nil but MusicStorage.AlbumData was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockAlbumData.Lock()
	mock.calls.AlbumData = append(mock.calls.AlbumData, callInfo)
	mock.lockAlbumData.Unlock()
	return mock.AlbumDataFunc(n)
}

// AlbumDataCalls gets all the calls that were made to AlbumData.
// Check the length with:
//     len(mockedMusicStorage.AlbumDataCalls())
func (mock *MockMusicStorage) AlbumDataCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockAlbumData.RLock()
	calls = mock.calls.AlbumData
	mock.lockAlbumData.RUnlock()
	return calls
}

// AlbumTracks calls AlbumTracksFunc.
func (mock *MockMusicStorage) AlbumTracks(n int64, b bool) ([]*proto.AlbumTrack, error) {
	if mock.AlbumTracksFunc == nil {
		panic("MockMusicStorage.AlbumTracksFunc: method is nil but MusicStorage.AlbumTracks was just called")
	}
	callInfo := struct {
		N int64
		B bool
	}{
		N: n,
		B: b,
	}
	mock.lockAlbumTracks.Lock()
	mock.calls.AlbumTracks = append(mock.calls.AlbumTracks, callInfo)
	mock.lockAlbumTracks.Unlock()
	return mock.AlbumTracksFunc(n, b)
}

// AlbumTracksCalls gets all the calls that were made to AlbumTracks.
// Check the length with:
//     len(mockedMusicStorage.AlbumTracksCalls())
func (mock *MockMusicStorage) AlbumTracksCalls() []struct {
	N int64
	B bool
} {
	var calls []struct {
		N int64
		B bool
	}
	mock.lockAlbumTracks.RLock()
	calls = mock.calls.AlbumTracks
	mock.lockAlbumTracks.RUnlock()
	return calls
}

// ArtistAlbums calls ArtistAlbumsFunc.
func (mock *MockMusicStorage) ArtistAlbums(n1 int64, n2 int64) ([]*proto.Album, error) {
	if mock.ArtistAlbumsFunc == nil {
		panic("MockMusicStorage.ArtistAlbumsFunc: method is nil but MusicStorage.ArtistAlbums was just called")
	}
	callInfo := struct {
		N1 int64
		N2 int64
	}{
		N1: n1,
		N2: n2,
	}
	mock.lockArtistAlbums.Lock()
	mock.calls.ArtistAlbums = append(mock.calls.ArtistAlbums, callInfo)
	mock.lockArtistAlbums.Unlock()
	return mock.ArtistAlbumsFunc(n1, n2)
}

// ArtistAlbumsCalls gets all the calls that were made to ArtistAlbums.
// Check the length with:
//     len(mockedMusicStorage.ArtistAlbumsCalls())
func (mock *MockMusicStorage) ArtistAlbumsCalls() []struct {
	N1 int64
	N2 int64
} {
	var calls []struct {
		N1 int64
		N2 int64
	}
	mock.lockArtistAlbums.RLock()
	calls = mock.calls.ArtistAlbums
	mock.lockArtistAlbums.RUnlock()
	return calls
}

// ArtistInfo calls ArtistInfoFunc.
func (mock *MockMusicStorage) ArtistInfo(n int64) (*proto.Artist, error) {
	if mock.ArtistInfoFunc == nil {
		panic("MockMusicStorage.ArtistInfoFunc: method is nil but MusicStorage.ArtistInfo was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockArtistInfo.Lock()
	mock.calls.ArtistInfo = append(mock.calls.ArtistInfo, callInfo)
	mock.lockArtistInfo.Unlock()
	return mock.ArtistInfoFunc(n)
}

// ArtistInfoCalls gets all the calls that were made to ArtistInfo.
// Check the length with:
//     len(mockedMusicStorage.ArtistInfoCalls())
func (mock *MockMusicStorage) ArtistInfoCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockArtistInfo.RLock()
	calls = mock.calls.ArtistInfo
	mock.lockArtistInfo.RUnlock()
	return calls
}

// ArtistTracks calls ArtistTracksFunc.
func (mock *MockMusicStorage) ArtistTracks(n1 int64, b bool, n2 int64) ([]*proto.Track, error) {
	if mock.ArtistTracksFunc == nil {
		panic("MockMusicStorage.ArtistTracksFunc: method is nil but MusicStorage.ArtistTracks was just called")
	}
	callInfo := struct {
		N1 int64
		B  bool
		N2 int64
	}{
		N1: n1,
		B:  b,
		N2: n2,
	}
	mock.lockArtistTracks.Lock()
	mock.calls.ArtistTracks = append(mock.calls.ArtistTracks, callInfo)
	mock.lockArtistTracks.Unlock()
	return mock.ArtistTracksFunc(n1, b, n2)
}

// ArtistTracksCalls gets all the calls that were made to ArtistTracks.
// Check the length with:
//     len(mockedMusicStorage.ArtistTracksCalls())
func (mock *MockMusicStorage) ArtistTracksCalls() []struct {
	N1 int64
	B  bool
	N2 int64
} {
	var calls []struct {
		N1 int64
		B  bool
		N2 int64
	}
	mock.lockArtistTracks.RLock()
	calls = mock.calls.ArtistTracks
	mock.lockArtistTracks.RUnlock()
	return calls
}

// DoesPlaylistExist calls DoesPlaylistExistFunc.
func (mock *MockMusicStorage) DoesPlaylistExist(n int64) (bool, error) {
	if mock.DoesPlaylistExistFunc == nil {
		panic("MockMusicStorage.DoesPlaylistExistFunc: method is nil but MusicStorage.DoesPlaylistExist was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockDoesPlaylistExist.Lock()
	mock.calls.DoesPlaylistExist = append(mock.calls.DoesPlaylistExist, callInfo)
	mock.lockDoesPlaylistExist.Unlock()
	return mock.DoesPlaylistExistFunc(n)
}

// DoesPlaylistExistCalls gets all the calls that were made to DoesPlaylistExist.
// Check the length with:
//     len(mockedMusicStorage.DoesPlaylistExistCalls())
func (mock *MockMusicStorage) DoesPlaylistExistCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockDoesPlaylistExist.RLock()
	calls = mock.calls.DoesPlaylistExist
	mock.lockDoesPlaylistExist.RUnlock()
	return calls
}

// FindAlbums calls FindAlbumsFunc.
func (mock *MockMusicStorage) FindAlbums(s string) ([]*proto.Album, error) {
	if mock.FindAlbumsFunc == nil {
		panic("MockMusicStorage.FindAlbumsFunc: method is nil but MusicStorage.FindAlbums was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockFindAlbums.Lock()
	mock.calls.FindAlbums = append(mock.calls.FindAlbums, callInfo)
	mock.lockFindAlbums.Unlock()
	return mock.FindAlbumsFunc(s)
}

// FindAlbumsCalls gets all the calls that were made to FindAlbums.
// Check the length with:
//     len(mockedMusicStorage.FindAlbumsCalls())
func (mock *MockMusicStorage) FindAlbumsCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockFindAlbums.RLock()
	calls = mock.calls.FindAlbums
	mock.lockFindAlbums.RUnlock()
	return calls
}

// FindArtists calls FindArtistsFunc.
func (mock *MockMusicStorage) FindArtists(s string) ([]*proto.Artist, error) {
	if mock.FindArtistsFunc == nil {
		panic("MockMusicStorage.FindArtistsFunc: method is nil but MusicStorage.FindArtists was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockFindArtists.Lock()
	mock.calls.FindArtists = append(mock.calls.FindArtists, callInfo)
	mock.lockFindArtists.Unlock()
	return mock.FindArtistsFunc(s)
}

// FindArtistsCalls gets all the calls that were made to FindArtists.
// Check the length with:
//     len(mockedMusicStorage.FindArtistsCalls())
func (mock *MockMusicStorage) FindArtistsCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockFindArtists.RLock()
	calls = mock.calls.FindArtists
	mock.lockFindArtists.RUnlock()
	return calls
}

// FindTracksByFullWord calls FindTracksByFullWordFunc.
func (mock *MockMusicStorage) FindTracksByFullWord(s string, b bool) ([]*proto.Track, error) {
	if mock.FindTracksByFullWordFunc == nil {
		panic("MockMusicStorage.FindTracksByFullWordFunc: method is nil but MusicStorage.FindTracksByFullWord was just called")
	}
	callInfo := struct {
		S string
		B bool
	}{
		S: s,
		B: b,
	}
	mock.lockFindTracksByFullWord.Lock()
	mock.calls.FindTracksByFullWord = append(mock.calls.FindTracksByFullWord, callInfo)
	mock.lockFindTracksByFullWord.Unlock()
	return mock.FindTracksByFullWordFunc(s, b)
}

// FindTracksByFullWordCalls gets all the calls that were made to FindTracksByFullWord.
// Check the length with:
//     len(mockedMusicStorage.FindTracksByFullWordCalls())
func (mock *MockMusicStorage) FindTracksByFullWordCalls() []struct {
	S string
	B bool
} {
	var calls []struct {
		S string
		B bool
	}
	mock.lockFindTracksByFullWord.RLock()
	calls = mock.calls.FindTracksByFullWord
	mock.lockFindTracksByFullWord.RUnlock()
	return calls
}

// FindTracksByPartial calls FindTracksByPartialFunc.
func (mock *MockMusicStorage) FindTracksByPartial(s string, b bool) ([]*proto.Track, error) {
	if mock.FindTracksByPartialFunc == nil {
		panic("MockMusicStorage.FindTracksByPartialFunc: method is nil but MusicStorage.FindTracksByPartial was just called")
	}
	callInfo := struct {
		S string
		B bool
	}{
		S: s,
		B: b,
	}
	mock.lockFindTracksByPartial.Lock()
	mock.calls.FindTracksByPartial = append(mock.calls.FindTracksByPartial, callInfo)
	mock.lockFindTracksByPartial.Unlock()
	return mock.FindTracksByPartialFunc(s, b)
}

// FindTracksByPartialCalls gets all the calls that were made to FindTracksByPartial.
// Check the length with:
//     len(mockedMusicStorage.FindTracksByPartialCalls())
func (mock *MockMusicStorage) FindTracksByPartialCalls() []struct {
	S string
	B bool
} {
	var calls []struct {
		S string
		B bool
	}
	mock.lockFindTracksByPartial.RLock()
	calls = mock.calls.FindTracksByPartial
	mock.lockFindTracksByPartial.RUnlock()
	return calls
}

// IncrementListenCount calls IncrementListenCountFunc.
func (mock *MockMusicStorage) IncrementListenCount(n int64) error {
	if mock.IncrementListenCountFunc == nil {
		panic("MockMusicStorage.IncrementListenCountFunc: method is nil but MusicStorage.IncrementListenCount was just called")
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
//     len(mockedMusicStorage.IncrementListenCountCalls())
func (mock *MockMusicStorage) IncrementListenCountCalls() []struct {
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

// IsPlaylistOwner calls IsPlaylistOwnerFunc.
func (mock *MockMusicStorage) IsPlaylistOwner(n1 int64, n2 int64) (bool, error) {
	if mock.IsPlaylistOwnerFunc == nil {
		panic("MockMusicStorage.IsPlaylistOwnerFunc: method is nil but MusicStorage.IsPlaylistOwner was just called")
	}
	callInfo := struct {
		N1 int64
		N2 int64
	}{
		N1: n1,
		N2: n2,
	}
	mock.lockIsPlaylistOwner.Lock()
	mock.calls.IsPlaylistOwner = append(mock.calls.IsPlaylistOwner, callInfo)
	mock.lockIsPlaylistOwner.Unlock()
	return mock.IsPlaylistOwnerFunc(n1, n2)
}

// IsPlaylistOwnerCalls gets all the calls that were made to IsPlaylistOwner.
// Check the length with:
//     len(mockedMusicStorage.IsPlaylistOwnerCalls())
func (mock *MockMusicStorage) IsPlaylistOwnerCalls() []struct {
	N1 int64
	N2 int64
} {
	var calls []struct {
		N1 int64
		N2 int64
	}
	mock.lockIsPlaylistOwner.RLock()
	calls = mock.calls.IsPlaylistOwner
	mock.lockIsPlaylistOwner.RUnlock()
	return calls
}

// PlaylistInfo calls PlaylistInfoFunc.
func (mock *MockMusicStorage) PlaylistInfo(n int64) (*proto.PlaylistData, error) {
	if mock.PlaylistInfoFunc == nil {
		panic("MockMusicStorage.PlaylistInfoFunc: method is nil but MusicStorage.PlaylistInfo was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockPlaylistInfo.Lock()
	mock.calls.PlaylistInfo = append(mock.calls.PlaylistInfo, callInfo)
	mock.lockPlaylistInfo.Unlock()
	return mock.PlaylistInfoFunc(n)
}

// PlaylistInfoCalls gets all the calls that were made to PlaylistInfo.
// Check the length with:
//     len(mockedMusicStorage.PlaylistInfoCalls())
func (mock *MockMusicStorage) PlaylistInfoCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockPlaylistInfo.RLock()
	calls = mock.calls.PlaylistInfo
	mock.lockPlaylistInfo.RUnlock()
	return calls
}

// PlaylistTracks calls PlaylistTracksFunc.
func (mock *MockMusicStorage) PlaylistTracks(n int64) ([]*proto.Track, error) {
	if mock.PlaylistTracksFunc == nil {
		panic("MockMusicStorage.PlaylistTracksFunc: method is nil but MusicStorage.PlaylistTracks was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockPlaylistTracks.Lock()
	mock.calls.PlaylistTracks = append(mock.calls.PlaylistTracks, callInfo)
	mock.lockPlaylistTracks.Unlock()
	return mock.PlaylistTracksFunc(n)
}

// PlaylistTracksCalls gets all the calls that were made to PlaylistTracks.
// Check the length with:
//     len(mockedMusicStorage.PlaylistTracksCalls())
func (mock *MockMusicStorage) PlaylistTracksCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockPlaylistTracks.RLock()
	calls = mock.calls.PlaylistTracks
	mock.lockPlaylistTracks.RUnlock()
	return calls
}

// RandomAlbums calls RandomAlbumsFunc.
func (mock *MockMusicStorage) RandomAlbums(n int64) (*proto.Albums, error) {
	if mock.RandomAlbumsFunc == nil {
		panic("MockMusicStorage.RandomAlbumsFunc: method is nil but MusicStorage.RandomAlbums was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockRandomAlbums.Lock()
	mock.calls.RandomAlbums = append(mock.calls.RandomAlbums, callInfo)
	mock.lockRandomAlbums.Unlock()
	return mock.RandomAlbumsFunc(n)
}

// RandomAlbumsCalls gets all the calls that were made to RandomAlbums.
// Check the length with:
//     len(mockedMusicStorage.RandomAlbumsCalls())
func (mock *MockMusicStorage) RandomAlbumsCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockRandomAlbums.RLock()
	calls = mock.calls.RandomAlbums
	mock.lockRandomAlbums.RUnlock()
	return calls
}

// RandomArtists calls RandomArtistsFunc.
func (mock *MockMusicStorage) RandomArtists(n int64) (*proto.Artists, error) {
	if mock.RandomArtistsFunc == nil {
		panic("MockMusicStorage.RandomArtistsFunc: method is nil but MusicStorage.RandomArtists was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockRandomArtists.Lock()
	mock.calls.RandomArtists = append(mock.calls.RandomArtists, callInfo)
	mock.lockRandomArtists.Unlock()
	return mock.RandomArtistsFunc(n)
}

// RandomArtistsCalls gets all the calls that were made to RandomArtists.
// Check the length with:
//     len(mockedMusicStorage.RandomArtistsCalls())
func (mock *MockMusicStorage) RandomArtistsCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockRandomArtists.RLock()
	calls = mock.calls.RandomArtists
	mock.lockRandomArtists.RUnlock()
	return calls
}

// RandomTracks calls RandomTracksFunc.
func (mock *MockMusicStorage) RandomTracks(n int64, b bool) (*proto.Tracks, error) {
	if mock.RandomTracksFunc == nil {
		panic("MockMusicStorage.RandomTracksFunc: method is nil but MusicStorage.RandomTracks was just called")
	}
	callInfo := struct {
		N int64
		B bool
	}{
		N: n,
		B: b,
	}
	mock.lockRandomTracks.Lock()
	mock.calls.RandomTracks = append(mock.calls.RandomTracks, callInfo)
	mock.lockRandomTracks.Unlock()
	return mock.RandomTracksFunc(n, b)
}

// RandomTracksCalls gets all the calls that were made to RandomTracks.
// Check the length with:
//     len(mockedMusicStorage.RandomTracksCalls())
func (mock *MockMusicStorage) RandomTracksCalls() []struct {
	N int64
	B bool
} {
	var calls []struct {
		N int64
		B bool
	}
	mock.lockRandomTracks.RLock()
	calls = mock.calls.RandomTracks
	mock.lockRandomTracks.RUnlock()
	return calls
}

// UserPlaylists calls UserPlaylistsFunc.
func (mock *MockMusicStorage) UserPlaylists(n int64) ([]*proto.PlaylistData, error) {
	if mock.UserPlaylistsFunc == nil {
		panic("MockMusicStorage.UserPlaylistsFunc: method is nil but MusicStorage.UserPlaylists was just called")
	}
	callInfo := struct {
		N int64
	}{
		N: n,
	}
	mock.lockUserPlaylists.Lock()
	mock.calls.UserPlaylists = append(mock.calls.UserPlaylists, callInfo)
	mock.lockUserPlaylists.Unlock()
	return mock.UserPlaylistsFunc(n)
}

// UserPlaylistsCalls gets all the calls that were made to UserPlaylists.
// Check the length with:
//     len(mockedMusicStorage.UserPlaylistsCalls())
func (mock *MockMusicStorage) UserPlaylistsCalls() []struct {
	N int64
} {
	var calls []struct {
		N int64
	}
	mock.lockUserPlaylists.RLock()
	calls = mock.calls.UserPlaylists
	mock.lockUserPlaylists.RUnlock()
	return calls
}