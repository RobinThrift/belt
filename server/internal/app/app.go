package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.robinthrift.com/belt/internal/auth"
	"go.robinthrift.com/belt/internal/control"
	"go.robinthrift.com/belt/internal/ingress/authv1"
	"go.robinthrift.com/belt/internal/ingress/syncv1"
	"go.robinthrift.com/belt/internal/jobs"
	"go.robinthrift.com/belt/internal/server"
	"go.robinthrift.com/belt/internal/storage/database/sqlite"
	"go.robinthrift.com/belt/internal/storage/filesystem"
)

type App struct {
	config Config
	srv    *http.Server
	db     *sqlite.SQLite

	initSetup *initSetup
	jobs      *jobs.System
}

func New(config Config) *App {
	db := &sqlite.SQLite{
		File:         config.Database.Path,
		EnableWAL:    config.Database.EnableWAL,
		Timeout:      config.Database.Timeout,
		DebugEnabled: config.Database.DebugEnabled,
	}

	argon2Params := auth.Argon2Params{
		KeyLen:  config.Argon2.KeyLen,
		Memory:  config.Argon2.Memory,
		Threads: config.Argon2.Threads,
		Time:    config.Argon2.Time,
		Version: config.Argon2.Version,
	}

	accountRepo := sqlite.NewAccountRepo(db)
	syncRepo := sqlite.NewSyncRepo(db)
	authTokenRepo := sqlite.NewAuthTokenRepo(db)
	jobRepo := sqlite.NewJobRepo(db)

	fs := &filesystem.LocalFSBlobStorage{
		BaseDir: config.Blobs.Dir,
		TmpDir:  os.TempDir(),
	}

	authConfig := control.AuthConfig{
		Argon2Params:              argon2Params,
		AuthTokenLength:           32,
		AccessTokenValidDuration:  config.AccessTokenValidDuration,
		RefreshTokenValidDuration: config.RefreshTokenValidDuration,
	}

	accountCtrl := control.NewAccountController(db, accountRepo)
	syncCtrl := control.NewSyncController(db, syncRepo, fs)
	authCtrl := control.NewAuthController(authConfig, db, accountCtrl, authTokenRepo)

	jobSystem := jobs.NewSystem(db, jobRepo, accountCtrl, time.Now, jobFuncs)

	mux := http.NewServeMux()
	srv := server.New(server.Config{Addr: config.Addr}, mux)

	authv1.New(config.BasePath, mux, authCtrl)
	syncv1.New(syncv1.RouterConfig{
		BasePath: config.BasePath,
	}, mux, syncCtrl, authCtrl, http.Dir(config.Blobs.Dir))

	return &App{
		config: config,
		srv:    srv,
		db:     db,
		initSetup: newInitSetup(initSetupConfig{
			InitUsername: config.Init.Username,
			InitPassword: auth.PlaintextPassword(config.Init.Password),
			Argon2params: argon2Params,
		}, db, accountCtrl, authCtrl),
		jobs: jobSystem,
	}
}

func (a *App) Start(ctx context.Context) error {
	err := a.db.Open()
	if err != nil {
		return err
	}

	err = sqlite.RunMigrations(ctx, sqlite.MigrationConfig{LogFormat: a.config.Log.Format, LogLevel: a.config.Log.Level}, a.db.DB)
	if err != nil {
		return err
	}

	err = a.initSetup.exec(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, fmt.Sprintf("starting server on %v", a.config.Addr))
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return a.db.Close()
}

func (a *App) Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping http server")
	return a.srv.Shutdown(ctx)
}
