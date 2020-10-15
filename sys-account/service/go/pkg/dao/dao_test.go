package dao_test

import (
	"testing"
	"time"

	"github.com/genjidb/genji"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/getcouragenow/sys/sys-account/service/go/pkg/dao"
	corecfg "github.com/getcouragenow/sys/sys-core/service/go"
	"github.com/getcouragenow/sys/sys-core/service/go/pkg/db"
)

var (
	testDb     *genji.DB
	accdb      *dao.AccountDB
	err        error
	role1ID    = db.UID()
	role2ID    = db.UID()
	account0ID = db.UID()

	defaultDbName          = "getcouragenow.db"
	defaultDbEncryptionKey = "testkey@!" // for test only.
	// TODO: Make config
	defaultDbDir               = "./db"
	defaultDbBackupDir         = "./db/backups"
	defaultDbBackupSchedulSpec = "@every 15s"
	defaultDbRotateSchedulSpec = "@every 1h"

	accs = []dao.Account{
		{
			ID:       account0ID,
			Email:    "2pac@example.com",
			Password: "no_biggie",
			RoleId:   role1ID,
			UserDefinedFields: map[string]interface{}{
				"City": "Compton",
			},
			CreatedAt: time.Now().UTC().Unix(),
			UpdatedAt: time.Now().UTC().Unix(),
			LastLogin: 0,
			Disabled:  false,
		},
		{
			ID:       db.UID(),
			Email:    "bigg@example.com",
			Password: "two_packs",
			RoleId:   role2ID,
			UserDefinedFields: map[string]interface{}{
				"City": "NY",
			},
			CreatedAt: time.Now().UTC().Unix(),
			UpdatedAt: time.Now().UTC().Unix(),
			LastLogin: 0,
			Disabled:  false,
		},
		{
			ID:       db.UID(),
			Email:    "shakur@example.com",
			Password: "no_biggie",
			RoleId:   role1ID,
			UserDefinedFields: map[string]interface{}{
				"City": "Compton LA",
			},
			CreatedAt: time.Now().UTC().Unix(),
			UpdatedAt: time.Now().UTC().Unix(),
			LastLogin: 0,
			Disabled:  false,
		},
	}
)

func init() {
	csc := &corecfg.SysCoreConfig{
		DbConfig: corecfg.DbConfig{
			Name:             defaultDbName,
			EncryptKey:       defaultDbEncryptionKey,
			RotationDuration: 1,
			DbDir:            defaultDbDir,
		},
		CronConfig: corecfg.CronConfig{
			BackupSchedule: defaultDbBackupSchedulSpec,
			RotateSchedule: defaultDbRotateSchedulSpec,
			BackupDir:      defaultDbBackupDir,
		},
	}

	if err := db.InitDatabase(csc); err != nil {
		log.Fatalf("error initializing db: %v", err)
	}

	testDb, _ = db.SharedDatabase()
	log.Println("MakeSchema testing .....")
	accdb, err = dao.NewAccountDB(testDb)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully initialize accountdb :  %v", accdb)
}

func TestAll(t *testing.T) {
	t.Run("Test Account Insert", testAccountInsert)
	t.Run("Test Role Insert", testRolesInsert)
	t.Run("Test Account Query", testQueryAccounts)
	t.Run("Test Role List", testRolesList)
	t.Run("Test Role Get", testRolesGet)
	t.Run("Test Role Update", testRolesUpdate)
	t.Run("Test Account Update", testUpdateAccounts)
}

func testAccountInsert(t *testing.T) {
	t.Log("on inserting accounts")

	for _, acc := range accs {
		err = accdb.InsertAccount(&acc)
		assert.NoError(t, err)
	}
	t.Log("successfully inserted accounts")
}

func testQueryAccounts(t *testing.T) {
	t.Logf("on querying accounts")
	queryParams := []*dao.QueryParams{
		{
			Params: map[string]interface{}{
				"email": "bigg@example.com",
			},
		},
		{
			Params: map[string]interface{}{
				"email": "2pac@example.com",
			},
		},
	}
	var accs []*dao.Account
	for _, qp := range queryParams {
		acc, err := accdb.GetAccount(qp)
		assert.NoError(t, err)
		t.Logf("Account: %v\n", acc)
		accs = append(accs, acc)
	}
	assert.NotEqual(t, accs[0], accs[1])

	for _, qp := range queryParams {
		accs, err = accdb.ListAccount(qp)
		assert.NoError(t, err)
	}
}

func testUpdateAccounts(t *testing.T) {
	accs[0].Email = "makavelli@example.com"
	accs[1].Email = "notorious_big@example.com"

	for _, acc := range accs {
		err = accdb.UpdateAccount(&acc)
		assert.NoError(t, err)
	}
}

func testDeleteAccounts(t *testing.T) {
	assert.NoError(t, accdb.DeleteAccount(accs[0].ID))
}
