package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

type Game struct {
	Download  string `json:"download"`
	Require []string `json:"require"`
	ExtractTo []struct {
		Src string
		Dst string
	} `json:"extract_to"`
	CopyTo string `json:"copy_to"`
}
var args = os.Args
var gameJson map[string]Game
var emulatorPath string
var shouldSkipOutro = true

func main() {
	if(len(args) > 2) {
		gameJson = readJson(resolveJsonFile())
		emulatorPath = getEmulatorDirectory(args)
		fmt.Print("\033[H\033[2J")
		motd()
		color.New(color.FgDarkGray).Println("> " + args[1], args[2])
		JsonUpdate()
		downloadGame(args[2])
		endingMessage()
		CleanupFile("./fc2roms.zip")
	} else {
		motd()
		color.New(color.FgDarkGray).Println("> not enough args provided, skipping download routine")
		JsonUpdate()
		fmt.Println("erro aqui?")
		endingMessage()
		CleanupFile("./fc2roms.zip")
	}
}

func readJson(file string) map[string]Game {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]Game
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func downloadGame(game string) bool {
	if args[1] == "fbneo" && len(strings.Split(args[2], "_")) > 1 {
		game = strings.Split(args[2], "_")[1]
	}
	selectedGame := gameJson[game]

	if(len(selectedGame.Download) == 0) {
		color.New(color.FgWhite, color.BgRed, color.OpBold).Println("\n ERROR: ROM not found in json ")
		color.New(color.FgRed, color.OpBold).Println(args[2] + " in file " + resolveJsonFile())
		shouldSkipOutro = false
		return false
	}

	split := strings.Split(selectedGame.Download, "/")
	filename := split[len(split)-1]

	if(!GameExists(selectedGame)) {
		if(len(selectedGame.CopyTo) > 0) {
			dir, name := filepath.Split(emulatorPath + "/" + selectedGame.CopyTo)
			cleanPath := path.Clean("./"+strings.Trim(dir, "/"))
			DownloadFile(selectedGame.Download, "./"+cleanPath, name)
		} else {
			DownloadFile(selectedGame.Download, emulatorPath, filename)
		}
		
		if(len(selectedGame.ExtractTo) > 0) {
			for _, item := range selectedGame.ExtractTo {
				unzip(filename, item.Src, item.Dst)
			}

			fmt.Println("cleaning up", filename)
			err := os.Remove(emulatorPath + "/" + filename)
			if err != nil {
				log.Fatal(err)
			}
		}

		
	} else {
		color.New(color.FgGreen).Println("\n" + filename + " was downloaded already")
	}

	if(len(selectedGame.Require) > 0) {
		for _, each := range selectedGame.Require {
			downloadGame(each)
		}
	}
	return true
}

func DownloadFile(gameUrl string, dest string, filename string) {
	req, _ := http.NewRequest("GET", gameUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	filename, err := url.QueryUnescape(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, _ := os.OpenFile(dest + "/" + filename, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	color.Cyanln("\nDownloading " + filename)

	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: "[white]█[reset]",
			BarStart:      "",
			BarEnd:        "",
		}))

	io.Copy(io.MultiWriter(f, bar), resp.Body)
	shouldSkipOutro = false

}

func GameExists(game Game) bool {
	if(len(game.CopyTo) > 0) {
		dir, name := filepath.Split(emulatorPath + "/" + game.CopyTo)
		cleanPath := "./" + path.Clean("./"+strings.Trim(dir, "/"))

		if _, err := os.Stat(cleanPath + "/" + name); !os.IsNotExist(err) {
			return true
		}
	}

	if(len(game.ExtractTo) > 0) {
		for _, item := range game.ExtractTo {
			if _, err := os.Stat(emulatorPath + "/" + item.Dst); os.IsNotExist(err) {
				return false
			}
		}

		return true
	}

	split := strings.Split(game.Download, "/")
	filename := emulatorPath + "/" + split[len(split)-1]

	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

    return true

}

func motd() {
	fmt.Println("")
	color.New(color.FgWhite, color.BgLightMagenta, color.OpBold).Println("                                               ")
	color.New(color.FgWhite, color.BgLightMagenta, color.OpBold).Println("            FRM PLUS by @lofi1048              ")
	color.New(color.FgWhite, color.BgLightMagenta, color.OpBold).Println("                                               ")
	color.New(color.FgDarkGray).Println("-----------------------------------------------")
	color.New(color.FgCyan, color.OpBold).Println("If you enjoy this consider supporting my ko-fi!")
	color.New(color.FgCyan, color.OpUnderscore).Println("https://ko-fi.com/lofi1048")
	color.New(color.FgDarkGray).Println("-----------------------------------------------")
}

func endingMessage() {
	fmt.Print("\n\n")
	color.Style{color.FgLightWhite, color.BgLightGreen, color.OpBold}.Println("                                               ")
	color.Style{color.FgLightWhite, color.BgLightGreen, color.OpBold}.Println("     Done! Press enter to close this screen    ")
	color.Style{color.FgLightWhite, color.BgLightGreen, color.OpBold}.Println("                                               ")

	if(!shouldSkipOutro) {
		fmt.Scanln() 
	}
}

func unzip(zipfile string, src string, dst string) error {
	zipfile = emulatorPath + "/" + zipfile
	dst = emulatorPath + "/" + dst

	r, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != src {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(dst); os.IsNotExist(err) {
			dir, _ := filepath.Split(dst)
			fmt.Println("extracting " + src + " to " + dst + "...")
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			outFile, err := os.Create(dst)
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()
			

			_, err = io.Copy(outFile, rc)
			if err != nil {
				log.Fatal(err)
			}
			
			

			rc.Close()
		}
		break
	}
	return nil
}

func getEmulatorDirectory(args []string) string {

	emulator := args[1]
	gameSplit := strings.Split(args[2],"_")

	fmt.Println()

	if emulator == "fbneo" {
		switch gameSplit[0] {
		case "cv":
			return "./fbneo/ROMs/coleco"
		case "gg":
			return "./fbneo/ROMs/gamegear"
		case "md":
			return "./fbneo/ROMs/megadrive"
		case "msx":
			return "./fbneo/ROMs/msx"
		case "nes":
			return "./fbneo/ROMs/nes"
		case "pce":
			return "./fbneo/ROMs/pce"
		case "sg1k":
			return "./fbneo/ROMs/sg1000"
		case "sms":
			return "./fbneo/ROMs/sms"
		case "tg":
			return "./fbneo/ROMs/tg16"
		default:
			return "./fbneo/ROMs"
		}
	}

	switch args[1] {
		case "flycast":
			return "./flycast/ROMs"
		case "snes9x":
			return "./snes9x/ROMs"
		case "fc1":
			return "./ggpofba/ROMs"
		default: 
			return "./"
	}
}

func resolveJsonFile() string {
	emulator := args[1]
	gameSplit := strings.Split(args[2],"_")

	if (emulator == "fbneo") {
		switch gameSplit[0] {
			case "cv":
				return "./fbneo_cv_roms.json"
			case "gg":
				return "./fbneo_gg_roms.json"
			case "md":
				return "./fbneo_md_roms.json"
			case "msx":
				return "./fbneo_msx_roms.json"
			case "nes":
				return "./fbneo_nes_roms.json"
			case "pce":
				return "./fbneo_pce_roms.json"
			case "sg1k":
				return "./fbneo_sg1k_roms.json"
			case "sms":
				return "./fbneo_sms_roms.json"
			case "tg":
				return "./fbneo_tg_roms.json"
			default:
				return "./fbneo_roms.json"
		}
	}

	return "./" + emulator + "_roms.json"
}

func JsonUpdate() error {
	color.New(color.FgLightYellow).Println("Checking for JSON updates...")
	urlJson := "https://newchallenger.net/fc2/fc2roms.zip"
	urlVersion := "https://newchallenger.net/fc2/jsonversion.txt"
	

	var versionInternal = ""
	var versionExternal = ""

	//reads the version file on server
	rescueStdout := os.Stdout

	resp, err := http.Get(urlVersion)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

	bytesExternal, _ := ioutil.ReadAll(resp.Body)
  	os.Stdout = rescueStdout
    if err != nil {
        return err
    }
	
	versionExternal = string(bytesExternal)

	//reads the version file on local
	bytes, err := ioutil.ReadFile("./jsonversion.txt")
	if err != nil {
        os.Create("./jsonversion.txt")
    }

	versionInternal = string(bytes)

	if(versionInternal != versionExternal) {
		color.New(color.FgGreen).Println("JSON Files update found! ("+ versionExternal + ")")

		f, err := os.Create("./jsonversion.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(versionExternal)
		if err != nil {
			log.Fatal(err)
		}

		//update json files
		DownloadFile(urlJson, "./", "fc2roms.zip")
		unzipAll("./fc2roms.zip", "./")
	}
	return nil
}


func unzipAll(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func CleanupFile(file string) {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		err2 := os.Remove(file)
		if err2 != nil {
			log.Fatal(err)
		}
	}
}