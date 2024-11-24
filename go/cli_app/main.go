/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/rbieldt/CLI-Todo/go/cli_app/cmd"

type MySqlOptions struct {
	User     string
	Password string
	Database string
	Port     int
}

type GrpcOptions struct {
	Port int
}

type Options struct {
	MySqlOptions *MySqlOptions
	GrpcOptions  *GrpcOptions
}

func main() {
	cmd.Execute()
}
