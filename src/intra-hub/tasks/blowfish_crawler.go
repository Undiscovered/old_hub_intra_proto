package tasks

import (
	"bufio"
	_ "crypto/sha512"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

const (
	blowFishCrawlerTaskName = "blowfishCrawler"
	epitechCDNPath          = "https://cdn.local.epitech.eu/userprofil/profilview/"
	locationFileURL         = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/location"
	blowFishURL             = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/master.passwd.blowfish"
	groupFileURL            = "https://lost-contact.mit.edu/afs/epitech.net/site/etc/group"
)

var (
	studentGroup *models.Group

	specialUsersPath = path.Dir(beego.AppConfigPath) + "/specials"
	specialUsersinfo = make(map[string]*models.Group)
	mapCities        = make(map[string]*models.City)
	mapPromotions    = make(map[string]*models.Promotion)

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

func newUser(blowfish string) (*models.User, string) {
	lineSplitted := strings.Split(blowfish, ":")
	nameSplitted := strings.Split(lineSplitted[7], " ")
	var firstName, lastName string
	if len(nameSplitted) > 0 {
		firstName = nameSplitted[0]
		if len(nameSplitted) > 1 {
			lastName = nameSplitted[1]
		}
	}
	group := studentGroup
	user := &models.User{
		Login:     lineSplitted[0],
		FirstName: strings.Title(firstName),
		LastName:  strings.Title(lastName),
		Password:  lineSplitted[1],
		Picture:   epitechCDNPath + lineSplitted[0] + ".jpg",
		Email:     lineSplitted[0] + "@epitech.eu",
	}
	if g := specialUsersinfo[user.Login]; g != nil {
		group = g
	}
	user.Group = group
	return user, lineSplitted[3]
}

func blowFishCrawler() error {
	beego.Informational("BlowFish run")
	if err := loadUsersFiles(); err != nil {
		return err
	}
	blowfish, location, mapGroup, err := crawlFiles()
	if err != nil {
		return err
	}
	defer blowfish.Close()
	defer location.Close()
	scannerBlowFish := bufio.NewScanner(blowfish)
	scannerLocation := bufio.NewScanner(location)
	orm.Debug = false
	beego.Informational("Inserting users")
	for scannerBlowFish.Scan() {
		scannerLocation.Scan()
		o := orm.NewOrm()
		o.Begin()
		user, groupName := newUser(scannerBlowFish.Text())
		// Set Promotion
		groupName = mapGroup[groupName]
		if promotion := mapPromotions[groupName]; promotion == nil {
			promotion = &models.Promotion{Name: groupName}
			if _, id, err := o.ReadOrCreate(promotion, "Name"); err == nil {
				promotion.Id = int(id)
				user.Promotion = promotion
				mapPromotions[groupName] = promotion
			} else {
				beego.Critical(err)
				o.Rollback()
				return err
			}
		} else {
			user.Promotion = promotion
		}
		// Set City
		cityName := ""
		if len(strings.Split(scannerLocation.Text(), ":")) > 1 {
			cityName = mapLocations[strings.Split(scannerLocation.Text(), ":")[1]]
		}
		if city := mapCities[cityName]; city == nil {
			city = &models.City{Name: cityName}
			if _, id, err := o.ReadOrCreate(city, "Name"); err == nil {
				city.Id = int(id)
				user.City = city
				mapCities[cityName] = city
			} else {
				beego.Critical(err)
				o.Rollback()
				return err
			}
		} else {
			user.City = city
		}
		r, err := o.Raw("INSERT INTO user ("+models.GetUserFields()+") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE password=?", user.Values(), user.Password).Exec()
		if err != nil {
			beego.Critical(err)
			o.Rollback()
			return err
		}
		rowsAffected, err := r.RowsAffected()
		if err != nil {
			beego.Critical(err)
			o.Rollback()
			return err
		}
		if rowsAffected != 0 {
			lastId, err := r.LastInsertId()
			if err != nil {
				beego.Critical(err)
				o.Rollback()
				return err
			}
			user.Id = int(lastId)
		}
		o.Commit()
	}
	beego.Informational("Users inserted")
	orm.Debug = true
	return nil
}

func loadUsersFiles() error {
	specialUsersFile, err := os.Open(specialUsersPath)
	if err != nil {
		return err
	}
	defer specialUsersFile.Close()
	scannerManager := bufio.NewScanner(specialUsersFile)
	group, err := db.GetGroupByNames(models.UserGroupStudent)
	if err != nil {
		return err
	}
	studentGroup = group
	for scannerManager.Scan() {
		lineSplitted := strings.Split(scannerManager.Text(), "=")
		group, err := db.GetGroupByNames(lineSplitted[1])
		if err != nil {
			return err
		}
		specialUsersinfo[lineSplitted[0]] = group
	}
	return nil
}

func init() {
	crawler := toolbox.NewTask(blowFishCrawlerTaskName, "0 0 0/12 * * *", blowFishCrawler)
	toolbox.AddTask(blowFishCrawlerTaskName, crawler)
	toolbox.StartTask()
}
