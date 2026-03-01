package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *app) initMongoDB(ctx context.Context) (*mongo.Database, error) {
	// Подключаемся к Mongo (путь берем из конфига)
	mongoURI := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v",
		a.config.DBConfig.User, a.config.DBConfig.Password,
		a.config.DBConfig.Host, a.config.DBConfig.Port, a.config.DBConfig.DB, a.config.DBConfig.User)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// РЕГИСТРИРУЕМ ЗАКРЫТИЕ В CLOZER
	a.closer.AddCloseFunc("mongo", func() error {
		return client.Disconnect(context.Background())
	})

	return client.Database(a.config.DBConfig.DB), nil
}
