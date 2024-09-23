package main

import _ "go.uber.org/mock/gomock"

//go:generate mockgen -package mocks_core -destination mocks/repositories.go github.com/thomas-marquis/wordpress-simple-backup/internal/core Repository
