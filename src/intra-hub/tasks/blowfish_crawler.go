package tasks

import (
	"bufio"
	_ "crypto/sha512"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"intra-hub/db"
	"intra-hub/models"
	"io"
	"net/http"
	"os"
	"path"
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
	studentGroup *models.Group

	specialUsersPath = path.Dir(beego.AppConfigPath) + "/specials"
	specialUsersinfo = make(map[string][]*models.Group)
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
	groups := make([]*models.Group, 1)
	groups[0] = studentGroup
	user := &models.User{
		Login:     lineSplitted[0],
		FirstName: strings.Title(firstName),
		LastName:  strings.Title(lastName),
		Password:  lineSplitted[1],
		Picture:   epitechCDNPath + lineSplitted[0] + ".jpg",
		Email:     lineSplitted[0] + "@epitech.eu",
	}
	if groupsToAdd, ok := specialUsersinfo[user.Login]; ok {
		groups = append(groups, groupsToAdd...)
	}
	user.Groups = groups
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
			id, err := o.Insert(promotion)
			if err != nil {
				o.Rollback()
                return err
			}
			promotion.Id = int(id)
            user.Promotion = promotion
			mapPromotions[groupName] = promotion
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
			id, err := o.Insert(city)
			if err != nil {
				o.Rollback()
				return err
			}
			city.Id = int(id)
            user.City = city
			mapCities[cityName] = city
		} else {
			user.City = city
		}
		r, err := o.Raw("INSERT INTO user ("+models.GetUserFields()+") VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE password=?", user.Values(), user.Password).Exec()
		if err != nil {
			o.Rollback()
			return err
		}
		rowsAffected, err := r.RowsAffected()
		if err != nil {
			o.Rollback()
			return err
		}
		if rowsAffected != 0 {
			lastId, err := r.LastInsertId()
			if err != nil {
				o.Rollback()
				return err
			}
			user.Id = int(lastId)
			m2m := o.QueryM2M(user, "Groups")
			if _, err := m2m.Add(user.Groups); err != nil {
				o.Rollback()
				return err
			}
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
	groups, err := db.GetGroupsByNames(models.UserGroupStudent)
	if err != nil {
		return err
	}
	studentGroup = groups[0]
	for scannerManager.Scan() {
		lineSplitted := strings.Split(scannerManager.Text(), "=")
		groups, err := db.GetGroupsByNames(strings.Split(lineSplitted[1], ",")...)
		if err != nil {
			return err
		}
		specialUsersinfo[lineSplitted[0]] = groups
	}
	return nil
}

func init() {
	crawler := toolbox.NewTask(blowFishCrawlerTaskName, "0 0 0/12 * * *", blowFishCrawler)
	toolbox.AddTask(blowFishCrawlerTaskName, crawler)
	toolbox.StartTask()
}
