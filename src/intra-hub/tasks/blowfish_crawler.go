package tasks

import (
	"bufio"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"intra-hub/models"
	"io"
	"net/http"
	"strings"
)

const (
	blowFishCrawlerTaskName = "blowfishCrawler"
	epitechCDNPath          = "https://cdn.local.epitech.eu/userprofil/commentview/"
	locationFileURL         = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/location"
	blowFishURL             = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/master.passwd.blowfish"
	groupFileURL            = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/group"
)

var (
	mapLocations = map[string]string{
		"mpl": "Montpellier",
		"prs": "Paris",
		"lyo": "Lyon",
		"stg": "Strasbourg",
		"tls": "Toulouse",
		"lil": "Lille",
		"msl": "Marseille",
		"nce": "Nice",
		"bdx": "Bordeaux",
		"ncy": "Nancy",
		"nts": "Nantes",
		"rns": "Rennes",
	}
)

func crawlFiles() (bodyBlowFish io.ReadCloser, bodyLocation io.ReadCloser, mapGroup map[string]string, err error) {
	res, err := http.Get(blowFishURL)
	if err != nil {
		return
	}
	bodyBlowFish = res.Body
	res, err = http.Get(locationFileURL)
	if err != nil {
		return
	}
	bodyLocation = res.Body
	res, err = http.Get(groupFileURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	mapGroup = make(map[string]string)
	for scanner.Scan() {
		lineSplitted := strings.Split(scanner.Text(), ":")
		mapGroup[lineSplitted[2]] = lineSplitted[0]
	}
	return
}

func newUser(blowfish, location string, mapGroup map[string]string) *models.User {
	locationInfo := strings.Split(location, ":")
	city := ""
	if len(locationInfo) > 1 {
		city = mapLocations[locationInfo[1]]
	}
	lineSplitted := strings.Split(blowfish, ":")
	nameSplitted := strings.Split(lineSplitted[7], " ")
	var firstName, lastName string
	if len(nameSplitted) > 0 {
		firstName = nameSplitted[0]
		if len(nameSplitted) > 1 {
			lastName = nameSplitted[1]
		}
	}
	user := &models.User{
		Login:     lineSplitted[0],
		FirstName: strings.Title(firstName),
		LastName:  strings.Title(lastName),
		Password:  lineSplitted[1],
		Picture:   epitechCDNPath + lineSplitted[0] + ".jpg",
		Promotion: mapGroup[lineSplitted[3]],
		Email:     lineSplitted[0] + "@epitech.eu",
		City:      city,
	}
	return user
}

func blowFishCrawler() error {
	beego.Informational("BlowFish run")
	blowfish, location, mapGroup, err := crawlFiles()
	if err != nil {
		return err
	}
	defer blowfish.Close()
	defer location.Close()
	scannerBlowFish := bufio.NewScanner(blowfish)
	scannerLocation := bufio.NewScanner(location)
	beego.Informational("Inserting users")
	for scannerBlowFish.Scan() {
		scannerLocation.Scan()
		user := newUser(scannerBlowFish.Text(), scannerLocation.Text(), mapGroup)
		_, err := orm.NewOrm().Raw("INSERT INTO user ("+models.GetUserFields()+") VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE password=?", user.Values(), user.Password).Exec()
		if err != nil {
			return err
		}
	}
	beego.Informational("Users inserted")
	return nil
}

func init() {
	crawler := toolbox.NewTask(blowFishCrawlerTaskName, "0 0 0/12 * * *", blowFishCrawler)
	toolbox.AddTask(blowFishCrawlerTaskName, crawler)
	toolbox.StartTask()
}
