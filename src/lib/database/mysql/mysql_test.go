// запуск теста
// command unix: mysqladmin processlist
// SET GOPATH=C:\Work\projectName
// go test -v lib/database/mysql
// go test -v -bench . lib/database/mysql
package mysql_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"lib"
	"lib/database/mysql"
	"lib/logs"
)

var cntGo = 10
var cntIteration = 1000

type Command struct {
	Id uint64
}

// Пользователи
type TestAccessUsers struct {
	Id              uint64    // Id
	AccessUsers_Id  uint64    // Пользователь
	Login           string    // Логин пользователя
	Password        string    // Пароль пользователя (SHA256)
	Email           string    // Email
	LastName        string    // Фамилия
	Name            string    // Имя
	MiddleName      string    // Отчество
	LastFailed      time.Time // Дата последней не удачной попытки входа
	IsAuthMemory    bool      // Помнить авторизацию в куках
	IsAccess        bool      // Доступ разрешен
	IsCondition     bool      // Условия пользователя
	IsOnline        bool      // Онлайн статус
	IsActivated     bool      // Активированный
	DateOnline      time.Time // Дата последнего посещения
	Date            time.Time // Дата регистрации
	Del             bool      // Запись удалена
	Hash            string    // Контрольная сумма для синхронизации (SHA256)
	CookieActivated string    // Кука активации и идентификации
}

func TestAll(t *testing.T) {

	var cfglogs = new(logs.Cfglogs)
	cfglogs.Debug = true
	cfglogs.DebugDetail = 0
	cfglogs.Level = 6
	cfglogs.Mode = `file`
	cfglogs.Lang = `ru-ru`
	cfglogs.Separator = ` `

	cfglogs.File, _ = os.Getwd()
	cfglogs.File = filepath.Dir(cfglogs.File)
	cfglogs.File = filepath.Dir(cfglogs.File)
	cfglogs.File = filepath.Dir(cfglogs.File) + `/log`
	os.MkdirAll(cfglogs.File, 0777)
	cfglogs.File += `/logs_test.log`

	logs.Init(cfglogs)
	logs.GoStart(nil)

	var err error
	var data []byte
	var messages []string
	var db *mysql.Db
	path, _ := filepath.Abs(`.`)
	path += `/test.sql`

	if data, err = ioutil.ReadFile(path); err != nil {
		t.Error(err.Error())
	}

	var cfg = new(mysql.CfgMysql)
	cfg.Login = `root`
	cfg.Password = `root`
	cfg.Name = `sungora`
	mysql.InitMysql(map[int8]*mysql.CfgMysql{0: cfg})

	if db, err = mysql.NewDb(0); err != nil {
		t.Fatal(err.Error())
		return
	}

	if messages, err = db.QueryByte(data); err != nil {
		t.Error(err.Error())
	}
	t.Log(messages)

	chanenlControlInsert := testInsert(t)
	chanenlControlLoad := testLoad(t)
	chanenlControlLoadSlice := testLoadSlice(t)
	chanenlControlLoadMap := testLoadMap(t)
	chanenlControlQuery := testQuery(t)

	var flag = 5
	for 0 < flag {
		select {
		case <-chanenlControlInsert:
			flag--
		case <-chanenlControlLoad:
			flag--
		case <-chanenlControlLoadSlice:
			flag--
		case <-chanenlControlLoadMap:
			flag--
		case <-chanenlControlQuery:
			flag--
		}
	}

	db.Query(`DROP TABLE TestAccessUsers`)
	db.Free()

	if logs.GoClose() == false {
		t.Fatal(`Ошибка остановки службы логирования`)
	}
}

func testInsert(t *testing.T) chan Command {
	var chanenlExit = make(chan Command)
	go func() {

		// инициализация
		var chanenl = make(chan Command)
		var idCnt uint64
		var cntrun = cntGo
		var db, err = mysql.NewDb(0)
		if err != nil {
			t.Error(err)
			return
		}

		// запуск программ
		for i := 0; i < cntrun; i++ {
			go func() {
				for j := 0; j < cntIteration; j++ {
					var user = new(TestAccessUsers)
					user.Login = lib.String.CreatePassword() + fmt.Sprintf(`%d`, i)
					user.Email = lib.String.CreatePassword() + fmt.Sprintf(`%d`, i)
					if _, err := db.Insert(user, `TestAccessUsers`); err != nil {
						t.Error(err.Error())
					} else {
						idCnt++
					}
				}
				chanenl <- Command{}
			}()
		}
		t.Logf("Запущено %d паралельных программ по %d итераций в каждой (проверка Insert)\n", cntrun, cntIteration)

		// завершение
		for cntrun > 0 {
			<-chanenl
			cntrun--
		}
		t.Logf(`Количество добавленных записей: %d`, idCnt)

		db.Free()
		chanenlExit <- Command{}
	}()
	return chanenlExit
}

func testLoad(t *testing.T) chan Command {
	var chanenlExit = make(chan Command)
	go func() {

		// инициализация
		var chanenl = make(chan Command)
		var idCnt uint64
		var cntrun = cntGo
		var db, err = mysql.NewDb(0)
		if err != nil {
			t.Error(err.Error())
			return
		}

		// запуск программ
		for i := 0; i < cntrun; i++ {
			go func() {
				for j := 0; j < cntIteration; j++ {
					var user = new(TestAccessUsers)
					if err := db.Select(user, `SELECT * FROM TestAccessUsers LIMIT 0, 1`); err != nil {
						t.Error(err.Error())
					} else {
						idCnt++
					}
				}
				chanenl <- Command{}
			}()
		}
		t.Logf("Запущено %d паралельных программ по %d итераций в каждой (проверка Select)\n", cntrun, cntIteration)

		// завершение
		for cntrun > 0 {
			<-chanenl
			cntrun--
		}
		t.Logf(`Количество прочитанных записей: %d`, idCnt)

		db.Free()
		chanenlExit <- Command{}

	}()
	return chanenlExit
}

func testLoadSlice(t *testing.T) chan Command {
	var chanenlExit = make(chan Command)
	go func() {

		// инициализация
		var chanenl = make(chan Command)
		var idCnt uint64
		var cntrun = cntGo
		var db, err = mysql.NewDb(0)
		if err != nil {
			t.Error(err.Error())
			return
		}

		// запуск программ
		for i := 0; i < cntrun; i++ {
			go func() {
				for j := 0; j < cntIteration; j++ {
					var user = make([]*TestAccessUsers, 0)
					if err := db.SelectSlice(&user, `SELECT * FROM TestAccessUsers LIMIT 0, 1000`); err != nil {
						t.Error(err.Error())
					} else {
						idCnt += uint64(len(user))
					}
				}
				chanenl <- Command{}
			}()
		}
		t.Logf("Запущено %d паралельных программ по %d итераций в каждой (проверка SelectSlice)\n", cntrun, cntIteration)

		// завершение
		for cntrun > 0 {
			<-chanenl
			cntrun--
		}
		t.Logf(`Количество прочитанных записей: %d`, idCnt)

		db.Free()
		chanenlExit <- Command{}
	}()
	return chanenlExit
}

func testLoadMap(t *testing.T) chan Command {
	var chanenlExit = make(chan Command)
	go func() {

		// инициализация
		var chanenl = make(chan Command)
		var idCnt uint64
		var cntrun = cntGo
		var db, err = mysql.NewDb(0)
		if err != nil {
			t.Error(err.Error())
			return
		}

		// запуск программ
		for i := 0; i < cntrun; i++ {
			go func() {
				for j := 0; j < cntIteration; j++ {
					var user = make(map[uint64]*TestAccessUsers)
					if err := db.SelectMap(&user, `SELECT * FROM TestAccessUsers LIMIT 0, 1000`); err != nil {
						t.Error(err.Error())
					} else {
						idCnt += uint64(len(user))
					}
				}
				chanenl <- Command{}
			}()
		}
		t.Logf("Запущено %d паралельных программ по %d итераций в каждой (проверка SelectMap)\n", cntrun, cntIteration)

		// завершение
		for cntrun > 0 {
			<-chanenl
			cntrun--
		}
		t.Logf(`Количество прочитанных записей: %d`, idCnt)

		db.Free()
		chanenlExit <- Command{}
	}()
	return chanenlExit
}

func testQuery(t *testing.T) chan Command {
	var chanenlExit = make(chan Command)
	go func() {

		// инициализация
		var chanenl = make(chan Command)
		var idCnt uint64
		var cntrun = cntGo
		var db, err = mysql.NewDb(0)
		if err != nil {
			t.Error(err.Error())
			return
		}

		// запуск программ
		for i := 0; i < cntrun; i++ {
			go func() {
				for j := 0; j < cntIteration; j++ {
					if err := db.Query(`UPDATE TestAccessUsers SET Hash = ? WHERE Id = 1`, lib.String.CreatePasswordHash(lib.String.CreatePassword())); err != nil {
						t.Error(err.Error())
					} else {
						idCnt++
					}
				}
				chanenl <- Command{}
			}()
		}
		t.Logf("Запущено %d паралельных программ по %d итераций в каждой (проверка Query)\n", cntrun, cntIteration)

		// завершение
		for cntrun > 0 {
			<-chanenl
			cntrun--
		}
		t.Logf(`Количество обновленных записей: %d`, idCnt)

		db.Free()
		chanenlExit <- Command{}
	}()
	return chanenlExit
}

//CREATE TABLE `mdsrecords` (
//  `id` bigint(20) NOT NULL AUTO_INCREMENT,
//  `tuzik` bigint(20) unsigned DEFAULT NULL,
//  `IdEditor` varchar(100) DEFAULT NULL,
//  `DateCreate` datetime DEFAULT NULL,
//  `DateEdit` datetime DEFAULT NULL,
//  `NameCreator` varchar(100) DEFAULT NULL,
//  `NameEditor` varchar(100) DEFAULT NULL,
//  `AuthorName` varchar(100) DEFAULT NULL,
//  `RecordName` varchar(100) DEFAULT NULL,
//  `Date` datetime DEFAULT NULL,
//  `Stantion` varchar(100) DEFAULT NULL,
//  `Black` tinyint(1) NOT NULL,
//  `Verdikt` varchar(100) DEFAULT NULL,
//  `Following` bigint(20) unsigned DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8

//type MdsRecords struct {
//	Id          uint64 `db:"id"`
//	IdCreator   uint64 `db:"tuzik"`
//	IdEditor    string
//	DateCreate  time.Time
//	DateEdit    time.Time
//	NameCreator string
//	NameEditor  string
//	AuthorName  string
//	RecordName  string
//	Date        time.Time
//	Stantion    string
//	Black       bool
//	Verdikt     string
//	Following   uint64
//}

//func TestLoad(t *testing.T) {
//	var err error
//	//var data []byte
//	//var messages []string
//	var db *mysql.Db

//	var cfg = new(mysql.CfgMysql)
//	cfg.Login = `root`
//	cfg.Password = `root`
//	cfg.Name = `zegota`
//	mysql.InitMysql(map[int8]*mysql.CfgMysql{0: cfg})

//	if db, err = mysql.NewDb(0); err != nil {
//		t.Fatal(err.Error())
//		return
//	}
//	defer db.Free()

//	var data []*MdsRecords
//	if err = db.LoadArray(&data, "SELECT * FROM `mdsrecords`"); err != nil {
//		t.Error(err.Error())
//		return
//	}
//	logs.Dumper(data)

//	var obj = new(MdsRecords)
//	if err = db.Load(obj, "SELECT * FROM `mdsrecords` WHERE id = 78"); err != nil {
//		t.Error(err.Error())
//		return
//	}
//	logs.Dumper(obj)

//	t.Log("done ok")
//	logs.Dumper()
//}
