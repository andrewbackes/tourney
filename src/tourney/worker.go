/*******************************************************************************

 Project: Tourney

 Module: worker
 Created: 12/3/2014
 Author(s): Andrew Backes
 Description:

 TODO:
 	- "\" vs "/" in file paths will be a problem on windows.
 	- send logs back to server

*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

// Wrapper for functions to be played with net/rpc :
type Worker struct {
	// To be used by the Work Manager:
	Address net.Addr
	RPC     *rpc.Client

	// To be used by the Worker:
	serverAddr string
	localPath  string
}

func (W *Worker) DownloadEngine(EngineFileName string, rpcResponse *string) error {

	// Make the file locally:
	LocalEngineFilePath := W.localPath + "/" + EngineFileName
	LocalFile, err := os.Create(LocalEngineFilePath)
	defer LocalFile.Close()
	if err != nil {
		fmt.Println("Error creating local file.")
		return err
	}
	err = LocalFile.Chmod(0755)
	if err != nil {
		fmt.Println("Error modifying file permissions.")
		return err
	}

	// Get it from the server:
	httpFile, err := http.Get("http://" + strings.Split(W.serverAddr, ":")[0] + ":9001" + "/" + EngineFileName)
	defer httpFile.Body.Close()
	if err != nil {
		return err
	}

	// Save if locally:
	size, err := io.Copy(LocalFile, httpFile.Body)
	if err != nil {
		return err
	}
	fmt.Println(EngineFileName, "downloaded.", size, "bytes.")

	return nil
}

func (W *Worker) LocalizeEngines(WorkingGame *Game, rpcResponse *string) error {
	// TODO: should check other locations also.

	fmt.Println("Localizing game engines...")

	// check if folder exists:
	if err := os.MkdirAll(W.localPath, os.ModePerm); err != nil { //!os.IsExist(err) {
		fmt.Println("Could not make directory:", W.localPath, " - ", err)
		return err
	}

	for color := 0; color <= 1; color++ {
		// figure out the engine name and paths:
		parsed := strings.SplitAfter(WorkingGame.Player[color].Path, "/")
		EngineFileName := parsed[len(parsed)-1]
		LocalPath := W.localPath + "/" + EngineFileName

		// Update file locations in the Game object:
		WorkingGame.Player[color].Path = "./" + LocalPath
		EngineValidated := false
		for attempts := 0; attempts < 3; attempts++ {
			// Verify file existence and integrity:
			if err := WorkingGame.Player[color].ValidateEngineFile(); err != nil {
				if _, err2 := os.Stat(WorkingGame.Player[color].Path); err2 == nil {
					//file exists but is corrupt. delete it.
					fmt.Println("Engine file corrupt.")
					os.Remove(WorkingGame.Player[color].Path)
				}
				// Download engine file:
				fmt.Println("Downloading", EngineFileName, "from the server... ")
				if err3 := W.DownloadEngine(EngineFileName, new(string)); err3 != nil {
					fmt.Println("Failed. - ", err3)
					//return err3
				}
			} else {
				// Engine file is verified.
				EngineValidated = true
				break
			}
		}
		if !EngineValidated {
			return errors.New("Engine file's integrety could not be validated.")
		}
		/*
			// check if the file exists:
			fmt.Println("Looking for", LocalPath, "...")
			if _, test := os.Stat(LocalPath); os.IsNotExist(test) {
				// file doesnt exist, so get it from the server:
				fmt.Println("Downloading ", EngineFileName, " from the server... ")
				if err := W.DownloadEngine(EngineFileName, new(string)); err != nil {
					fmt.Println("Failed. - ", err)
					return err
				}
				// md5 check:
				if err := WorkingGame.Player[color].ValidateEngineFile(); err != nil {
					fmt.Println("MD5:", err)
				} else {
					fmt.Println("MD5 Ok.")
				}
			} else if test == nil {
				// file does exist
				// TODO: should md5 check it!
			}
		*/

		fmt.Println("File Integrity Verified. Using:", "./"+LocalPath)

	}

	return nil
}

func (W *Worker) PlayGame(G Game, CompletedGame *Game) error {
	fmt.Println("Recieved Round", G.Round)

	// Copy the game so that we keep what we got from the server intact.
	WorkingGame := G

	// Identify what engine files need to be used.
	if err := W.LocalizeEngines(&WorkingGame, new(string)); err == nil {

		// engines are localized, so play:
		if err2 := PlayGame(&WorkingGame); err2 != nil {
			fmt.Println(err2)
		}

	} else {
		fmt.Print(err)
		WorkingGame.ResultDetail = fmt.Sprintln(err)
	}

	// Return game:
	*CompletedGame = WorkingGame
	fmt.Println("Game results sent back to the server.")

	// Communicate with the console:
	//WorkingGame.PrintHUD()

	return nil
}

func ConnectAndWait(address string) {
	// First connect to the host:
	fmt.Print("\nConnecting to ", address, " ... ")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Success.")
	fmt.Println("Waiting on server...")

	// Establish RPC serving:
	ThisWorker := &Worker{
		serverAddr: address,
		localPath:  "worker",
	}
	rpc.Register(ThisWorker)
	rpc.ServeConn(conn)

	// Server disconnected.
	fmt.Println("Connection closed.")
}

func WorkForDirtyBit() {
	// TEMPORARY !!!!!!!!!!
	// TODO: REMOVE THIS !!

	// get the current IP from dirty-bit.com:
	fmt.Println("Resolving IP address from www.dirty-bit.com...")
	res, err := http.Get("http://www.dirty-bit.com/tourney/ip.txt")
	if err != nil {
		fmt.Println(err)
	}
	ip, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	ConnectAndWait(string(ip) + ":9000")

}
