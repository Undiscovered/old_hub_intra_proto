package tasks

import (
	"bufio"
	_ "crypto/sha512"
	"io"
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
	locationFileURL         = "./conf/location.txt"
	blowFishURL             = "./conf/master.passwd.blowfish.txt"
	groupFileURL            = "./conf/group.txt"
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

func crawlFiles() (bodyBlowFish io.ReadCloser, bodyLocation io.ReadCloser, bodyGroup io.ReadCloser, err error) {
	bodyBlowFish, err = os.Open(blowFishURL)
	if err != nil {
		return
	}
	bodyLocation, err = os.Open(locationFileURL)
	if err != nil {
		return
	}
	bodyGroup, err = os.Open(groupFileURL)
	if err != nil {
		return
	}
	return
}

func newUser(blowfish string) *models.User {
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
		Promotion: mapPromotions[lineSplitted[3]],
	}
	if g := specialUsersinfo[user.Login]; g != nil {
		group = g
	}
	user.Group = group
	return user
}

func blowFishCrawler() error {
	beego.Informational("BlowFish run")
	if err := loadUsersFiles(); err != nil {
		beego.Error(err)
		return err
	}
	blowfish, location, group, err := crawlFiles()
	if err != nil {
		beego.Error(err)
		return err
	}
	defer blowfish.Close()
	defer location.Close()
	defer group.Close()
	scannerBlowFish := bufio.NewScanner(blowfish)
	scannerLocation := bufio.NewScanner(location)
	scannerGroups := bufio.NewScanner(group)
	beego.Informational("Inserting users")
	o := orm.NewOrm()
	o.Begin()
	for scannerGroups.Scan() {
		lineSplitted := strings.Split(scannerGroups.Text(), ":")
		if len(lineSplitted) > 2 {
			groupName := lineSplitted[0]
			group := &models.Promotion{Name: groupName}
			if _, id, err := o.ReadOrCreate(group, "Name"); err == nil {
				group.Id = int(id)
				mapPromotions[groupName] = group
			} else {
				o.Rollback()
				beego.Error(err)
				return err
			}
		}
	}
	external, err := db.GetPromotionByName("external")
	if err != nil {
		return err
	}
	paris := &models.City{Name: "Paris"}
	if _, id, err := o.ReadOrCreate(paris, "Name"); err == nil {
		paris.Id = int(id)
		mapCities["Paris"] = paris
	}
	for scannerBlowFish.Scan() {
		user := newUser(scannerBlowFish.Text())
		// Set City
		cityName := ""
		if len(strings.Split(scannerLocation.Text(), ":")) > 1 {
			cityName = mapLocations[strings.Split(scannerLocation.Text(), ":")[1]]
		}
		if city := mapCities[cityName]; city == nil && cityName != "" {
			city = &models.City{Name: cityName}
			if _, id, err := o.ReadOrCreate(city, "Name"); err == nil {
				city.Id = int(id)
				user.City = city
				mapCities[cityName] = city
			} else {
				o.Rollback()
				beego.Error(err)
				return err
			}
		} else {
			user.City = city
		}
		if user.City == nil {
			user.City = paris
		}
		if user.Promotion == nil {
			user.Promotion = external
		}
		r, err := o.Raw("INSERT INTO user ("+models.GetUserFields()+") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE password=?", user.Values(), user.Password).Exec()
		if err != nil {
			o.Rollback()
			beego.Error(err)
			return err
		}
		rowsAffected, err := r.RowsAffected()
		if err != nil {
			o.Rollback()
			beego.Error(err)
			return err
		}
		if rowsAffected != 0 {
			lastId, err := r.LastInsertId()
			if err != nil {
				o.Rollback()
				beego.Error(err)
				return err
			}
			user.Id = int(lastId)
		}
	}
	o.Commit()
	beego.Informational("Users inserted")
	return nil
}

func loadUsersFiles() error {
	specialUsersFile, err := os.Open(specialUsersPath)
	if err != nil {
		return err
	}
	defer specialUsersFile.Close()
	scannerManager := bufio.NewScanner(specialUsersFile)
	group, err := db.GetGroupByName(models.UserGroupStudent)
	if err != nil {
		return err
	}
	studentGroup = group
	for scannerManager.Scan() {
		lineSplitted := strings.Split(scannerManager.Text(), "=")
		group, err := db.GetGroupByName(lineSplitted[1])
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
