package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
	appURL = appName
	appName = strings.ToLower(appName)

	// sanitize appName (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}

	log.Println("Creating new app:", appName)

	// git clone the skeleton app
	color.Green("\t Cloning repository...")

	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/FernandoJVideira/velox-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git folder
	color.Green("\t Removing .git folder...")
	err = os.RemoveAll("./" + appName + "/.git")
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to go .env file
	color.Yellow("\t Creating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", vel.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a makefile
	var source *os.File

	if runtime.GOOS == "windows" {
		source, err = os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))

	} else {
		source, err = os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
	}
	if err != nil {
		exitGracefully(err)
	}

	defer source.Close()

	destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
	if err != nil {
		exitGracefully(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		exitGracefully(err)
	}

	_ = os.Remove(fmt.Sprintf("./%s/Makefile.mac", appName))
	_ = os.Remove(fmt.Sprintf("./%s/Makefile.windows", appName))

	// update the go.mod file
	color.Yellow("\t Creating go.mod file...")
	_ = os.Remove(fmt.Sprintf("./%s/go.mod", appName))

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), fmt.Sprintf("./%s/go.mod", appName))
	if err != nil {
		exitGracefully(err)
	}

	// update the existing .go files with correct name/imports
	color.Yellow("\t Updating source files...")
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project dir
	color.Yellow("\t Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	color.Green("\t Done building" + appName + "!")
}
