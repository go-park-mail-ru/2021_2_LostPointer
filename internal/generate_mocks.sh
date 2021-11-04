#!/bin/bash

go generate moq -out ./mock/album_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/album AlbumRepository:MockAlbumRepository
go generate moq -out ./mock/album_usecase_db_mock.go -pkg mock 2021_2_LostPointer/internal/album AlbumUseCase:MockAlbumUseCase
go generate moq -out ./mock/artist_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/artist ArtistRepository:MockArtistRepository
go generate moq -out ./mock/artist_usecase_db_mock.go -pkg mock 2021_2_LostPointer/internal/artist ArtistUseCase:MockArtistUseCase
go generate moq -out ./mock/playlist_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/playlist PlaylistRepository:MockPlayListRepository
go generate moq -out ./mock/playlist_usecase_db_mock.go -pkg mock 2021_2_LostPointer/internal/playlist PlaylistUseCase:MockPlaylistUseCase
go generate moq -out ./mock/sessions_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/sessions SessionRepository:MockSessionRepository
go generate moq -out ./mock/track_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/track TrackRepository:MockTrackRepository
go generate moq -out ./mock/track_usecase_mock.go -pkg mock 2021_2_LostPointer/internal/track TrackUseCase:MockTrackUseCase
go generate moq -out ./mock/user_repo_db_mock.go -pkg mock 2021_2_LostPointer/internal/users UserRepository:MockUserRepository
go generate moq -out ./mock/user_usecase_mock.go -pkg mock 2021_2_LostPointer/internal/users UserUseCase:MockUserUseCase
go generate moq -out ./mock/avatar_repository_mock.go -pkg mock 2021_2_LostPointer/internal/avatars AvatarRepository:MockAvatarRepository
go generate moq -out ./mock/session_checker_mock.go -pkg mock 2021_2_LostPointer/internal/microservices/authorization/delivery SessionCheckerClient:MockSessionCheckerClient