package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"istomyang.github.com/like-iam/component-base/errors"
	generaloptions "istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (s *datastore) User() store.UserStore {
	return newUser(s)
}

func (s *datastore) Secret() store.SecretStore {
	return newSecret(s)
}

func (s *datastore) Policy() store.PolicyStore {
	return newPolicy(s)
}

func (s *datastore) Run() error {
	return nil
}

func (s *datastore) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

var (
	factory store.Factory
	once    sync.Once
)

func GetMySQLFactoryOr(opts *generaloptions.MySQLOpts) (store.Factory, error) {
	if opts == nil && factory == nil {
		return nil, errors.New("fail to init with no options.")
	}

	var err error
	once.Do(func() {
		var client *gorm.DB
		client, err = newMySqlClient(opts)
		factory = &datastore{db: client}
	})

	if err != nil || factory == nil {
		return nil, fmt.Errorf("reate mysql factory faild: %s", err.Error())
	}

	return factory, nil
}

func newMySqlClient(opts *generaloptions.MySQLOpts) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opts.Pool.MaxIdleConnections)
	sqlDB.SetConnMaxIdleTime(opts.Pool.MaxConnectionLifeTime)
	sqlDB.SetMaxOpenConns(opts.Pool.MaxOpenConnections)

	return db, nil
}
