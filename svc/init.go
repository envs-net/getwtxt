/*
Copyright (c) 2019 Ben Morrison (gbmor)

This file is part of Getwtxt.

Getwtxt is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Getwtxt is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Getwtxt.  If not, see <https://www.gnu.org/licenses/>.
*/

package svc // import "git.sr.ht/~gbmor/getwtxt/svc"

import (
	"html/template"
	"log"
	"os"
	"os/signal"
	"time"

	"git.sr.ht/~gbmor/getwtxt/registry"
	"github.com/spf13/pflag"
)

var (
	// Vers contains the version number set at build time
	Vers         string
	flagVersion  *bool   = pflag.BoolP("version", "v", false, "Display version information, then exit.")
	flagHelp     *bool   = pflag.BoolP("help", "h", false, "Display the quick-help screen.")
	flagMan      *bool   = pflag.BoolP("manual", "m", false, "Display the configuration manual.")
	flagConfFile *string = pflag.StringP("config", "c", "", "The name/path of the configuration file you wish to use.")
	flagAssets   *string = pflag.StringP("assets", "a", "", "The location of the getwtxt assets directory.")
	flagDBPath   *string = pflag.StringP("db", "d", "", "Path to the getwtxt database.")
	flagDBType   *string = pflag.StringP("dbtype", "t", "", "Type of database being used.")
)

// Holds the global configuration
var confObj = &Configuration{}

// Signals to close the log file
var closeLog = make(chan struct{}, 1)

// Used to transmit database pointer
var dbChan = make(chan dbase, 1)

// Used to transmit the wrapped tickers
// corresponding to the in-memory cache
// or the on-disk database.
var dbTickC = make(chan *tick, 1)
var cTickC = make(chan *tick, 1)

// Used to manage the landing page template
var tmpls *template.Template

// Holds the registry data in-memory
var twtxtCache = registry.New(nil)

// List of other registries submitted to this registry
var remoteRegistries = &RemoteRegistries{
	List: make([]string, 0),
}

// In-memory cache of static assets, specifically
// the parsed landing page and the stylesheet.
var staticCache = &staticAssets{}

// Logs an error that should cause a catastrophic
// failure of getwtxt
func errFatal(context string, err error) {
	if err != nil {
		log.Fatalf(context+"%v\n", err.Error())
	}
}

// Logs non-fatal errors.
func errLog(context string, err error) {
	if err != nil {
		log.Printf(context+"%v\n", err.Error())
	}
}

// I'm not using init() because it runs
// even during testing and was causing
// problems.
func initSvc() {
	checkFlags()
	titleScreen()

	initConfig()
	initLogging()
	initDatabase()
	tmpls = initTemplates()
	initPersistence()

	pingAssets()
	watchForInterrupt()
}

// Responds to some command-line flags
func checkFlags() {
	pflag.Parse()
	if *flagVersion {
		titleScreen()
		os.Exit(0)
	}
	if *flagHelp {
		titleScreen()
		helpScreen()
		os.Exit(0)
	}
	if *flagMan {
		titleScreen()
		helpScreen()
		manualScreen()
		os.Exit(0)
	}
}

// Watch for SIGINT aka ^C
// Close the log file then exit
func watchForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sigint := range c {
			log.Printf("Caught %v\n", sigint)

			log.Printf("Pushing to database ...\n")
			pushDB()

			log.Printf("Cleaning up ...\n")
			killTickers()
			killDB()

			confObj.Mu.RLock()
			log.Printf("Closed database connection to %v\n", confObj.DBPath)
			if !confObj.StdoutLogging {
				closeLog <- struct{}{}
			}
			confObj.Mu.RUnlock()

			close(dbChan)
			close(closeLog)

			// Let everything catch up
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}
	}()
}
