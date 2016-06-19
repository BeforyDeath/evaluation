package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {

	err := Config.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	if Config.LogFile {
		logFile, err := os.OpenFile("core.log", os.O_CREATE|os.O_RDWR, 0666) //os.O_APPEND |
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.SetFormatter(&log.JSONFormatter{})
	}

	//return

	err = GoogleDrive.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	//UpdateXLSX()
	go RefreshEvery(1*time.Minute, UpdateXLSX)

	log.Info("Server started...")
	//log.Info(runtime.NumGoroutine())
	http.HandleFunc("/api/secret", GetDataSecret)
	log.Fatal(http.ListenAndServe(Config.Listen, nil))
}

func GetDataSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SecretColumns)
	return
}

func RefreshEvery(d time.Duration, f func() error) {
	for _ = range time.Tick(d) {
		//log.Infof("goroutine: %v ", d.String())
		err := f()
		if err != nil {
			log.Info(err)
			return
		}
	}
}

func UpdateXLSX() error {
	log.Info("Get modifiedTime")

	meta, err := GoogleDrive.Service.Files.Get(Config.FileId).Fields("modifiedTime").Do()
	if err != nil {
		return err
	}

	if meta.ModifiedTime != GoogleDrive.modifiedTime {
		err = GetXLSXFile()
		if err != nil {
			return err
		}

		log.Infof("Save new modifTime: %v", meta.ModifiedTime)
		GoogleDrive.modifiedTime = meta.ModifiedTime

		err = UpdateJSON()
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func GetXLSXFile() error {
	log.Info("Get files")

	resp, err := GoogleDrive.Service.Files.Export(Config.FileId, Config.MimeType).Download()
	if err != nil {
		return err
	}

	out, err := os.Create("upload/table.xlsx")
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	log.Info("Save new file")

	return nil
}

func UpdateJSON() error {

	err := XlsxTable.ParseXLSX("upload/table.xlsx")
	if err != nil {
		return err
	}

	err = XlsxTable.Save("public/data/table.json")
	if err != nil {
		return err
	}

	return nil
}
