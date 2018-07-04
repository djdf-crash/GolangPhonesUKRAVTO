package db

import (
	"config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"path/filepath"
	"time"
)

var DB *gorm.DB

func InitDB() error {

	var err error

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	dataBaseConfig := config.AppConfig.DataBase

	//DB, err = gorm.Open(dataBaseConfig.NameDriver, dataBaseConfig.Path)
	DB, err = gorm.Open(dataBaseConfig.NameDriver, dir+string(os.PathSeparator)+dataBaseConfig.Path)
	if err != nil {
		log.Panic(err.Error())
		return err
	}

	DB.LogMode(config.AppConfig.DataBase.LogMode)

	DB.SingularTable(true)

	DB.AutoMigrate(&Employee{}, &User{}, &Organization{})

	//DB.Model(&Employee{}).AddForeignKey("organization_id", "organizations(id)", "RESTRICT", "RESTRICT")

	return nil
}

func AddUser(user *User) error {

	tx := DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func AddOrganization(org *Organization) error {

	tx := DB.Begin()

	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func AddEmployee(e *Employee) error {

	tx := DB.Begin()

	if err := tx.Create(&e).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func UpdateUser(user *User, lastUpdate time.Time) error {

	user.LastUpdate = lastUpdate

	tx := DB.Begin()

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func UpdateEmployee(e *Employee) error {

	e.LastUpdate = time.Now()

	tx := DB.Begin()

	if err := tx.Save(&e).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func UpdateOrganization(o *Organization) error {

	o.LastUpdatePhones = time.Now()

	tx := DB.Begin()

	if err := tx.Save(&o).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func DeleteEmployeeByOrganizationID(id uint) {

	DB.Model(&Employee{}).Where("organization_id = ?", id).Updates(Employee{LastUpdate: time.Now(), IsDelete: true})
}

func DeleteOrganization(o *Organization) error {

	o.IsDelete = true

	if err := UpdateOrganization(o); err != nil {

		return err
	}

	return nil

}

func DeleteEmployee(e *Employee) error {

	e.IsDelete = true

	if err := UpdateEmployee(e); err != nil {

		return err
	}

	return nil

}

func GetUserByEmailAndDeviceID(email, deviceID string) User {
	var user User
	DB.Where("email = ? AND device_id = ?", email, deviceID).First(&user)

	return user
}

func GetOrganizationByName(name string) Organization {
	var organization Organization

	DB.Where("name LIKE ?", "%"+name+"%").First(&organization)

	return organization
}

func GetOrganizationByNameAndLastUpdate(name string, lastUpdate time.Time) Organization {
	var organization Organization

	DB.Where("name LIKE ? AND last_update_phones >= ?", "%"+name+"%", lastUpdate).Find(&organization)

	return organization
}

func GetOrganizationByID(id uint) Organization {
	var organization Organization

	DB.Where("id = ?", id).First(&organization)

	return organization
}

func GetOrganizationByIDAndLastUpdate(id uint, lastUpdate time.Time) Organization {
	var organization Organization

	DB.Where("id = ? AND last_update_phones >= ?", id, lastUpdate).First(&organization)

	return organization
}

func GetEmployeesByEmail(email string) Employee {
	var emp Employee

	DB.Where("email = ?", email).Find(&emp)

	return emp
}

func GetEmployeesByOrganizationIDLastUpdate(id uint, lastUpdate time.Time) []Employee {
	var e []Employee

	DB.Where("organization_id = ? AND last_update >= ?", id, lastUpdate).Find(&e)

	return e
}

func GetEmployeesByOrganizationID(id uint) []Employee {

	return GetEmployeesByOrganizationIDLastUpdate(id, time.Time{})
}

func GetEmployeeByFullNameANDOrganizationID(id uint, fullName string, post string) Employee {
	var employee Employee

	DB.Where("full_name = ? AND organization_id = ? AND post = ?", fullName, id, post).Find(&employee)

	return employee
}

func GetEmployeesByLastUpdate(lastUpdate time.Time) []Employee {
	var employees []Employee

	DB.Where("last_update >= ?", lastUpdate).Find(&employees)

	return employees
}

func GetAllEmployee() []Employee {
	var employees []Employee

	DB.Find(&employees)

	return employees
}

func GetAllOrganizations() []Organization {
	var orgs []Organization

	DB.Find(&orgs)

	return orgs
}

func GetUserByToken(token string) User {
	var user User
	DB.Where("token = ?", token).Find(&user)

	return user
}
