package syncv1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.robinthrift.com/conveyor/internal/auth"
	"go.robinthrift.com/conveyor/internal/control"
	"go.robinthrift.com/conveyor/internal/domain"
	"go.robinthrift.com/conveyor/internal/storage/database/sqlite"
	"go.robinthrift.com/conveyor/internal/storage/filesystem"
	"go.robinthrift.com/conveyor/internal/testhelper"
)

func TestRouter_Attachments(t *testing.T) { //nolint:paralleltest // @TODO: check why these fail when run in parallel
	mux, token := setupSyncV1Router(t)

	content := "attachment content"

	uploadAttachment := func(t *testing.T) {
		t.Helper()

		req := httptest.NewRequest(http.MethodPost, "/api/sync/v1/attachments", bytes.NewReader([]byte(content)))
		req.Header.Add(authHeader, "Bearer "+token)
		req.Header.Add("X-Filepath", "a/b/c/d/test")

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}

	t.Run("Get", func(t *testing.T) { //nolint:paralleltest // @TODO: check why these fail when run in parallel
		uploadAttachment(t)

		req := httptest.NewRequest(http.MethodGet, "/blobs/a/b/c/d/test", nil)
		req.Header.Add(authHeader, "Bearer "+token)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, content, string(body))
	})

	t.Run("Not Found", func(t *testing.T) { //nolint:paralleltest // @TODO: check why these fail when run in parallel
		uploadAttachment(t)

		req := httptest.NewRequest(http.MethodGet, "/blobs/a/b/c/d/not_found", nil)
		req.Header.Add(authHeader, "Bearer "+token)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Invalid Auth", func(t *testing.T) { //nolint:paralleltest // @TODO: check why these fail when run in parallel
		uploadAttachment(t)

		req := httptest.NewRequest(http.MethodGet, "/blobs/a/b/c/d/test", nil)
		req.Header.Add(authHeader, "Bearer INVALID_TOKEN")

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("No Auth", func(t *testing.T) { //nolint:paralleltest // @TODO: check why these fail when run in parallel
		uploadAttachment(t)

		req := httptest.NewRequest(http.MethodGet, "/blobs/a/b/c/d/test", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
}

func setupSyncV1Router(t *testing.T) (http.Handler, string) {
	t.Helper()

	db := testhelper.NewInMemTestSQLite(t)

	accountRepo := sqlite.NewAccountRepo(db)
	syncRepo := sqlite.NewSyncRepo(db)
	authTokenRepo := sqlite.NewAuthTokenRepo(db)

	config := control.AuthConfig{
		Argon2Params:              auth.Argon2Params{KeyLen: 32, Memory: 8192, Threads: 2, Time: 1},
		AuthTokenLength:           32,
		AccessTokenValidDuration:  time.Hour,
		RefreshTokenValidDuration: time.Hour * 2,
	}

	blobDir := t.TempDir()
	blobs := &filesystem.LocalFSBlobStorage{
		BaseDir: blobDir,
		TmpDir:  t.TempDir(),
	}

	accountCtrl := control.NewAccountController(db, accountRepo)
	authCtrl := control.NewAuthController(config, db, accountCtrl, authTokenRepo)
	attachmentCtrl := control.NewAttachmentController(blobs)
	syncCtrl := control.NewSyncController(db, syncRepo, accountCtrl, attachmentCtrl, blobs)

	err := authCtrl.CreateAccount(t.Context(), control.CreateAccountCmd{
		Account: &domain.Account{
			Username: t.Name(),
		},
		PlaintextPasswd: auth.PlaintextPassword(t.Name() + "_init"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = authCtrl.ChangeAccountPassword(t.Context(), control.ChangeAccountPasswordCmd{
		Username:            t.Name(),
		CurrPasswdPlaintext: auth.PlaintextPassword(t.Name() + "_init"),
		NewPasswdPlaintext:  auth.PlaintextPassword(t.Name()),
	})
	if err != nil {
		t.Fatal(err)
	}

	token, err := authCtrl.CreateAuthTokenUsingCredentials(t.Context(), control.CreateAuthTokenUsingCredentialsCmd{
		Username:        t.Name(),
		PlaintextPasswd: auth.PlaintextPassword(t.Name()),
	})
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()

	New(RouterConfig{BasePath: "/"}, mux, syncCtrl, authCtrl, http.Dir(blobDir))

	return mux, token.Plaintext.Export()
}
