package store

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StoreInterface interface {
	CreateItem(ctx context.Context, item *Item) error
	CreateItemTx(ctx context.Context, item *Item, fn func(context.Context) error) error
	CreateItemWithHooks(ctx context.Context, item *Item) error
}

type Store struct {
	db *gorm.DB
}

func NewStore(dsn string) *Store {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Panic("failed to open db", zap.Error(err))
	}
	return &Store{db: db}
}

func (s Store) CreateItem(ctx context.Context, item *Item) error {
	return s.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Create(item).Error
}

func (s Store) CreateItemTx(ctx context.Context, item *Item, fn func(context.Context) error) error {
	// if either db insert or push to queue failed, db will rollback
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Session(&gorm.Session{SkipHooks: true}).Create(item).Error; err != nil {
			return err
		}

		return fn(ctx)
	})
}

func (s Store) CreateItemWithHooks(ctx context.Context, item *Item) error {
	return s.db.Session(&gorm.Session{SkipHooks: false}).WithContext(ctx).Model(item).Create(item).Error
}
