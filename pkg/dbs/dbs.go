package dbs

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const DatabaseTimeout = 5 * time.Second

//go:generate mockery
type Database interface {
	GetDB() *gorm.DB
	AutoMigrate(models ...any) error
	WithTransaction(function func() error) error
	Create(ctx context.Context, doc any) error
	CreateInBatches(ctx context.Context, docs any, batchSize int) error
	Update(ctx context.Context, doc any) error
	Delete(ctx context.Context, value any, opts ...FindOption) error
	FindById(ctx context.Context, id string, result any) error
	FindOne(ctx context.Context, result any, opts ...FindOption) error
	Find(ctx context.Context, result any, opts ...FindOption) error
	Count(ctx context.Context, model any, total *int64, opts ...FindOption) error
}

type Query struct {
	Query string
	Args  []any
}

func NewQuery(query string, args ...any) Query {
	return Query{
		Query: query,
		Args:  args,
	}
}

type database struct {
	db *gorm.DB
}

func NewDatabase(uri string) (Database, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// Set up connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(200)

	return &database{
		db: db,
	}, nil
}

func (d *database) AutoMigrate(models ...any) error {
	return d.db.AutoMigrate(models...)
}

func (d *database) WithTransaction(function func() error) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	originalDB := d.db
	d.db = tx

	if err := function(); err != nil {
		d.db = originalDB
		tx.Rollback()
		return err
	}

	d.db = originalDB
	return tx.Commit().Error
}

func (d *database) Create(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.WithContext(ctx).Create(doc).Error
}

func (d *database) CreateInBatches(ctx context.Context, docs any, batchSize int) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.WithContext(ctx).CreateInBatches(docs, batchSize).Error
}

func (d *database) Update(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.WithContext(ctx).Save(doc).Error
}

func (d *database) Delete(ctx context.Context, value any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	return query.WithContext(ctx).Delete(value).Error
}

func (d *database) FindById(ctx context.Context, id string, result any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	if err := d.db.WithContext(ctx).Where("id = ? ", id).First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *database) FindOne(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *database) Find(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).Find(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *database) Count(ctx context.Context, model any, total *int64, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).Model(model).Count(total).Error; err != nil {
		return err
	}

	return nil
}

func (d *database) GetDB() *gorm.DB {
	return d.db
}

func (d *database) applyOptions(opts ...FindOption) *gorm.DB {
	query := d.db

	opt := getOption(opts...)

	if len(opt.preloads) != 0 {
		for _, preload := range opt.preloads {
			query = query.Preload(preload)
		}
	}

	if opt.query != nil {
		for _, q := range opt.query {
			query = query.Where(q.Query, q.Args...)
		}
	}

	if opt.order != "" {
		query = query.Order(opt.order)
	}

	if opt.offset != 0 {
		query = query.Offset(opt.offset)
	}

	if opt.limit != 0 {
		query = query.Limit(opt.limit)
	}

	return query
}
