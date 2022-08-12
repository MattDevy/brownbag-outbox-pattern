package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MattDevy/brownbag-outbox-pattern/internal/config"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/pubsub"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/server"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/store"
	"go.uber.org/zap"

	gpub "cloud.google.com/go/pubsub"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	migrate "github.com/rubenv/sql-migrate"

	gsql "database/sql"

	_ "github.com/lib/pq"
)

const (
	migrationDir = "postgresql/migrations"
)

func init() {
	// set up logging
	zcfg := zap.NewDevelopmentConfig()
	zcfg.Encoding = "json"
	zcfg.OutputPaths = []string{"stdout"}
	zcfg.ErrorOutputPaths = []string{"stderr"}
	logger, err := zcfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	defer zap.L().Sync()

	cfg := config.NewConfig()

	// open database connection for outbox subscriber/migrations
	db, err := gsql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		zap.L().Fatal("failed to connect postgres", zap.Error(err))
	}
	defer db.Close()

	// apply migrations
	doMigrate(db, cfg.DBName)

	// initialize outbox subscriber
	sub, err := sql.NewSubscriber(db, sql.SubscriberConfig{
		ConsumerGroup: "default",
		SchemaAdapter: sql.DefaultPostgreSQLSchema{
			GenerateMessagesTableName: func(topic string) string {
				return fmt.Sprintf("outbox_%s", topic)
			},
		},
		OffsetsAdapter: sql.DefaultPostgreSQLOffsetsAdapter{
			GenerateMessagesOffsetsTableName: func(topic string) string {
				return fmt.Sprintf("outbox_offsets_%s", topic)
			},
		},
		InitializeSchema: false,
	}, watermill.NewStdLogger(false, false))
	if err != nil {
		zap.L().Fatal("fail to create watermill subscriber", zap.Error(err))
	}
	defer sub.Close()

	// initialize database connection for service
	s := store.NewStore(cfg.DBConnectionString())
	ps := pubsub.NewEmulator()

	// kick off outbox subscriber
	go outboxToPubSub(context.Background(), sub, ps)

	// run http server
	serv := server.NewServer(s, ps)
	if err := serv.Run(); err != nil {
		zap.L().Fatal("web error", zap.Error(err))
	}
}

func outboxToPubSub(ctx context.Context, sub *sql.Subscriber, ps pubsub.PubSubAPI) {
	msgs, err := sub.Subscribe(ctx, "events")
	if err != nil {
		zap.L().Fatal("failed to get messages", zap.Error(err))
	}

	for msg := range msgs {
		// unmarshal into cloudevent
		e := cloudevents.NewEvent()
		err := json.Unmarshal(msg.Payload, &e)
		if err != nil {
			zap.L().Fatal("bad data", zap.Error(err))
		}

		// only needed to pass instructions to emulator
		ctx := pubsub.ContextFromMetadata(msg.Context(), msg.Metadata)

		// "publish" to pubsub emulator
		res := ps.PublishTopic(ctx, server.OutputTopic, &gpub.Message{
			Data: e.Data(),
		})
		_, err = res.Get(msg.Context())
		if err != nil {
			zap.L().Error("failed to send message to pubsub", zap.Error(err))
		}

		// ack regardless, in production you'd want to do something smarter based on errors
		msg.Ack()
	}
}

func doMigrate(db *gsql.DB, databaseSchema string) {
	migrations := migrate.FileMigrationSource{
		Dir: migrationDir,
	}
	migrate.SetSchema(databaseSchema)
	count, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		zap.L().Panic("failed to apply migrations", zap.Error(err))
	}
	zap.L().Info("migrations applied", zap.Int("count", count))
}
