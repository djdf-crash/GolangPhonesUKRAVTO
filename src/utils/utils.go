package utils

import (
	"config"
	"crypto/hmac"
	"crypto/sha256"
	"db"
	"encoding/base64"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func CheckerFile() {

	excelFileName := config.AppConfig.RootDirPath + config.AppConfig.SettingsParseFile.PathFile

	fmt.Println(excelFileName)

	if len(excelFileName) == 0 {
		log.Println("No set file path in config!")
		time.Sleep(30 * time.Minute)
		CheckerFile()
		return
	}
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Println(err.Error())
		time.Sleep(24 * time.Hour)
		CheckerFile()
	}

	mapOrg := map[string][]db.Employee{}

	var strOrganization string

	re := regexp.MustCompile(`(\d-)?(\d{3}-){2}(\d{2}-\d{2})`)

	for _, sheet := range xlFile.Sheets {

		strOrganization = ""

		for ind, row := range sheet.Rows {

			if ind <= 2 {
				continue
			}

			employee := &db.Employee{}

			//row.ReadStruct(employee)

			if len(row.Cells) > 0 {

				newOrganization := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnOrganization].String())

				if len(newOrganization) == 0 {
					continue
				}

				strOrganization = newOrganization
			}

			if len(strOrganization) == 0 {
				continue
			}

			if _, ok := mapOrg[strOrganization]; !ok {
				mapOrg[strOrganization] = []db.Employee{}
			}

			fullName := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnFullName].String())
			department := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnDepartment].String())
			section := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnSection].String())
			post := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnPost].String())
			email := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnEmail].String())
			contactInfo := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnContactInfo].String())
			phone := strings.TrimSpace(row.Cells[config.AppConfig.SettingsParseFile.NumberColumnPhone].String())

			if len(fullName) == 0 {
				if len(contactInfo) == 0 {
					continue
				}
				if len(mapOrg[strOrganization]) == 0 {
					continue
				}
				tmpEmployee := mapOrg[strOrganization][len(mapOrg[strOrganization])-1]
				if len(tmpEmployee.RealPhone) != 0 {
					continue
				}

				tmpEmployee.RealPhone = GetRealPhoneSubMatch(contactInfo, re)
				if len(tmpEmployee.RealPhone) == 0 {
					tmpEmployee.RealPhone = GetRealPhoneSubMatch(phone, re)
				}

				if len(tmpEmployee.RealPhone) != 0 {
					mapOrg[strOrganization] = append(mapOrg[strOrganization][:len(mapOrg[strOrganization])-1], tmpEmployee)
				}

				continue
			}

			employee.OrganizationName = strOrganization
			employee.FullName = GetTrimString(fullName, " ")
			employee.Post = GetTrimString(strings.TrimSpace(post), " ")
			employee.Email = GetTrimString(strings.TrimSpace(email), " ")
			employee.ContactInfo = GetTrimString(strings.TrimSpace(contactInfo), " ")
			employee.Phone = GetTrimString(strings.TrimSpace(phone), " ")
			employee.Department = GetTrimString(strings.TrimSpace(department), " ")
			employee.Section = GetTrimString(strings.TrimSpace(section), " ")

			realPhone := GetRealPhoneSubMatch(employee.ContactInfo, re)
			if len(realPhone) == 0 {
				realPhone = GetRealPhoneSubMatch(employee.Phone, re)
			}
			employee.RealPhone = realPhone

			employee.LastUpdate = time.Now()

			mapOrg[strOrganization] = append(mapOrg[strOrganization], *employee)
		}
	}

	for k, v := range mapOrg {

		var newArrayEmployee []db.Employee

		for _, emp := range v {
			if len(emp.RealPhone) == 0 && len(emp.FullName) == 0 && len(emp.Email) == 0 {
				continue
			}
			newArrayEmployee = append(newArrayEmployee, emp)
		}
		if len(newArrayEmployee) != 0 {
			mapOrg[k] = newArrayEmployee
			//fmt.Println(k + " " + strconv.Itoa(len(mapOrg[k])))
		} else {
			delete(mapOrg, k)
		}

	}

	SaveInDB(&mapOrg)

	time.Sleep(30 * time.Minute)

	CheckerFile()
}
func SaveInDB(mapOrg *map[string][]db.Employee) {
	emptyOrganization := db.Organization{}
	emptyEmployee := db.Employee{}

	arrOrg := db.GetAllOrganizations()

	for _, orgDb := range arrOrg {
		find := false
		for k := range *mapOrg {
			if orgDb.Name == k {
				find = true
				if orgDb.IsDelete {
					orgDb.IsDelete = false
					db.UpdateOrganization(&orgDb)
				}
				break
			}
		}
		if !find {
			orgDb.LastUpdatePhones = time.Now()
			db.DeleteEmployeeByOrganizationID(orgDb.ID)
			db.DeleteOrganization(&orgDb)
		}
	}

	for k, v := range *mapOrg {
		org := db.GetOrganizationByName(k)
		if reflect.DeepEqual(emptyOrganization, org) {
			org.Name = k
			org.LastUpdatePhones = time.Now()
			err := db.AddOrganization(&org)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if org.IsDelete {
			org.IsDelete = false
			db.UpdateOrganization(&org)
		}

		var firstUpdateTime time.Time

		for _, emp := range v {

			if len(emp.RealPhone) == 0 && len(emp.FullName) == 0 && len(emp.Email) == 0 {
				continue
			}
			emp.OrganizationID = org.ID
			employeeDb := db.GetEmployeeByFullNameANDOrganizationID(org.ID, emp.FullName, emp.Department, emp.Section, emp.Post)
			if reflect.DeepEqual(emptyEmployee, employeeDb) {
				emp.LastUpdate = time.Now()
				err := db.AddEmployee(&emp)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				var update bool

				if employeeDb.RealPhone != emp.RealPhone {
					update = true
					employeeDb.RealPhone = emp.RealPhone
					if employeeDb.Email != emp.Email {
						employeeDb.Email = emp.Email
					}

					if employeeDb.FullName != emp.FullName {
						employeeDb.FullName = emp.FullName
					}

					if employeeDb.Section != emp.Section {
						employeeDb.Section = emp.Section
					}

					if employeeDb.Department != emp.Department {
						employeeDb.Department = emp.Department
					}

					if employeeDb.Post != emp.Post {
						employeeDb.Post = emp.Post
					}
				}

				if employeeDb.IsDelete {
					update = true
					employeeDb.IsDelete = false
				}

				if update {
					db.UpdateEmployee(&employeeDb)
					firstUpdateTime = time.Now()
				}
			}
		}

		arrEmpDbOrg := db.GetEmployeesByOrganizationID(org.ID)

		for _, empDB := range arrEmpDbOrg {
			find := false
			for _, emp := range v {
				if empDB.RealPhone == emp.RealPhone {
					find = true
					break
				}
			}
			if !find {
				db.UpdateOrganization(&org)
				db.DeleteEmployee(&empDB)
			}
		}

		if !firstUpdateTime.IsZero() {
			db.UpdateOrganization(&org)
		}
	}
	fmt.Println("DONE update!")
}

func GetTrimString(str string, sep string) string {

	var tmpStr string

	arrStr := strings.Split(str, sep)

	for _, v := range arrStr {
		if len(strings.TrimSpace(v)) == 0 {
			continue
		}

		tmpStr = tmpStr + " " + strings.TrimSpace(v)
	}

	return strings.TrimSpace(tmpStr)

}

func GetRealPhoneSubMatch(strPhone string, re *regexp.Regexp) string {

	subMatch := re.FindStringSubmatch(strPhone)

	if len(subMatch) > 2 {
		return subMatch[0]
	}

	return ""

}

func ReadConfigFile(start bool) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	err = config.InitConfig(dir + string(os.PathSeparator) + "config.json")
	if err != nil {
		log.Panic(err.Error())
	}

	if start {
		time.Sleep(30 * time.Minute)
		ReadConfigFile(true)
	}
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
