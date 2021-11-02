#!/bin/bash

moq -out mock/album_repo_db_mock.go -pkg mock album AlbumRepository:MockAlbumRepository
moq -out mock/album_usecase_db_mock.go -pkg mock album AlbumUseCase:MockAlbumUseCase

moq -out mock/artist_repo_db_mock.go -pkg mock artist ArtistRepository:MockArtistRepository
moq -out mock/artist_usecase_db_mock.go -pkg mock artist ArtistUseCase:MockArtistUseCase

moq -out mock/playlist_repo_db_mock.go -pkg mock playlist PlaylistRepository:MockPlayListRepository
moq -out mock/playlist_usecase_db_mock.go -pkg mock playlist PlaylistUseCase:MockPlaylistUseCase

moq -out mock/sessions_repo_db_mock.go -pkg mock sessions SessionRepository:MockSessionRepository

moq -out mock/track_repo_db_mock.go -pkg mock track TrackRepository:MockTrackRepository
moq -out mock/track_usecase_mock.go -pkg mock track TrackUseCase:MockTrackUseCase

moq -out mock/user_repo_db_mock.go -pkg mock users UserRepository:MockUserRepository
moq -out mock/user_usecase_mock.go -pkg mock users UserUseCase:MockUserUseCase

moq -out mock/avatar_repository_mock.go -pkg mock utils/images AvatarRepositoryIFace:MockAvatarRepositoryIFace

moq -out mock/session_checker_mock.go -pkg mock microservices/authorization/delivery SessionCheckerClient:MockSessionCheckerClient