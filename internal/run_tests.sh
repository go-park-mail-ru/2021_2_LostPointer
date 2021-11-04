#!/bin/bash

go test -coverprofile cover.out 2021_2_LostPointer/internal/{album,artist,csrf,microservices/authorization/usecase,playlist,sessions,track,users,avatars}/... 2021_2_LostPointer/pkg/...
go tool cover -func cover.out | grep total