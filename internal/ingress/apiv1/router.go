package apiv1

import (
	"cmp"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/RobinThrift/belt/internal/auth"
	"github.com/RobinThrift/belt/internal/control"
	"github.com/RobinThrift/belt/internal/domain"
	"github.com/RobinThrift/belt/internal/server/session"
)

type router struct {
	baseURL        string
	memoCtrl       *control.MemoControl
	attachmentCtrl *control.AttachmentControl
	settingsCtrl   *control.SettingsControl
	tokenCtrl      *control.APITokenController
	accountFetcher AccountFetcher
}

type AccountFetcher interface {
	GetAccountForAPIToken(ctx context.Context, value auth.APITokenValue) (*auth.Account, error)
}

func New(baseURL string, mux *http.ServeMux, memoCtrl *control.MemoControl, attachmentCtrl *control.AttachmentControl, settingsCtrl *control.SettingsControl, tokenCtrl *control.APITokenController, accountFetcher AccountFetcher) {
	r := &router{baseURL, memoCtrl, attachmentCtrl, settingsCtrl, tokenCtrl, accountFetcher}

	HandlerWithOptions(NewStrictHandlerWithOptions(r, nil, StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  errorHandlerFunc,
		ResponseErrorHandlerFunc: errorHandlerFunc,
	}), StdHTTPServerOptions{
		BaseRouter:       mux,
		BaseURL:          baseURL + "api/v1",
		ErrorHandlerFunc: errorHandlerFunc,
		Middlewares:      []MiddlewareFunc{recoverer, r.checkAuth},
	})
}

// (GET /memos/{id})
func (r *router) GetMemo(ctx context.Context, req GetMemoRequestObject) (GetMemoResponseObject, error) {
	memo, err := r.memoCtrl.GetMemo(ctx, domain.MemoID(req.Id))
	if err != nil {
		if errors.Is(err, domain.ErrMemoNotFound) {
			return nil, fmt.Errorf("%w: %v", errNotFound, err)
		}

		return nil, err
	}

	return GetMemo200JSONResponse{
		Id:         memo.ID.String(),
		Content:    string(memo.Content),
		IsArchived: memo.IsArchived,
		IsDeleted:  memo.IsDeleted,
		CreatedAt:  memo.CreatedAt,
		CreatedBy:  memo.CreatedBy.String(),
		UpdatedAt:  memo.UpdatedAt,
	}, nil
}

// (GET /memos)
func (r *router) ListMemos(ctx context.Context, req ListMemosRequestObject) (ListMemosResponseObject, error) {
	query := control.ListMemosQuery{
		PageSize:   min(cmp.Or(req.Params.PageSize, 10), 50),
		PageAfter:  req.Params.PageAfter,
		Tag:        req.Params.FilterTag,
		Search:     req.Params.FilterContent,
		IsArchived: req.Params.FilterIsArchived,
		IsDeleted:  req.Params.FilterIsDeleted,
	}

	if req.Params.FilterCreatedAt != nil {
		createdAtOp := "="
		if req.Params.OpCreatedAt != nil {
			createdAtOp = *(*string)(req.Params.OpCreatedAt)
		}

		switch createdAtOp {
		case "=":
			query.CreatedAt = &req.Params.FilterCreatedAt.Time
		case "<=":
			query.MinCreationDate = &req.Params.FilterCreatedAt.Time
		default:
			return nil, fmt.Errorf("%w: unknown operation '%s' for filter[created_at]", errBadRequest, createdAtOp)
		}
	}

	memos, err := r.memoCtrl.ListMemos(ctx, query)
	if err != nil {
		return nil, err
	}

	apiMemos := MemoList{Items: make([]Memo, len(memos.Items))}
	for i, memo := range memos.Items {
		apiMemos.Items[i] = Memo{
			Id:         memo.ID.String(),
			Content:    string(memo.Content),
			IsArchived: memo.IsArchived,
			IsDeleted:  memo.IsDeleted,
			CreatedAt:  memo.CreatedAt,
			CreatedBy:  memo.CreatedBy.String(),
			UpdatedAt:  memo.UpdatedAt,
		}
	}

	if memos.Next != nil {
		apiMemos.Next = memos.Next
	}

	return ListMemos200JSONResponse(apiMemos), nil
}

// (POST /memos)
func (r *router) CreateMemo(ctx context.Context, req CreateMemoRequestObject) (CreateMemoResponseObject, error) {
	cmd := control.CreateMemoCmd{
		Content: []byte(req.Body.Content),
	}

	if req.Body.CreatedAt != nil {
		cmd.CreatedAt = req.Body.CreatedAt
	}

	created, err := r.memoCtrl.CreateMemo(ctx, cmd)
	if err != nil {
		return nil, err
	}

	memo, err := r.memoCtrl.GetMemo(ctx, created)
	if err != nil {
		return nil, err
	}

	return CreateMemo201JSONResponse{
		Id:         memo.ID.String(),
		Content:    string(memo.Content),
		IsArchived: memo.IsArchived,
		IsDeleted:  memo.IsDeleted,
		CreatedAt:  memo.CreatedAt,
		CreatedBy:  memo.CreatedBy.String(),
		UpdatedAt:  memo.UpdatedAt,
	}, nil
}

// (PATCH /memos)
func (r *router) UpdateMemo(ctx context.Context, req UpdateMemoRequestObject) (UpdateMemoResponseObject, error) {
	cmd := control.UpdateMemoCmd{MemoID: domain.MemoID(req.Id)}

	if req.Body.Content != nil {
		cmd.Content = []byte(*req.Body.Content)
	}

	if req.Body.IsArchived != nil {
		cmd.Content = nil
		cmd.IsArchived = req.Body.IsArchived
	}

	if req.Body.IsDeleted != nil {
		cmd.Content = nil
		cmd.IsArchived = nil
		cmd.IsDeleted = req.Body.IsDeleted
	}

	err := r.memoCtrl.UpdateMemo(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrMemoNotFound) {
			return nil, fmt.Errorf("%w: %v", errNotFound, err)
		}
		return nil, err
	}

	return UpdateMemo204Response{}, nil
}

// (DELETE /memos/{id})
func (r *router) DeleteMemo(ctx context.Context, req DeleteMemoRequestObject) (DeleteMemoResponseObject, error) {
	isDeleted := true
	cmd := control.UpdateMemoCmd{MemoID: domain.MemoID(req.Id), IsDeleted: &isDeleted}

	err := r.memoCtrl.UpdateMemo(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrMemoNotFound) {
			return nil, fmt.Errorf("%w: %v", errNotFound, err)
		}
		return nil, err
	}

	return DeleteMemo204Response{}, nil
}

// (GET /tags)
func (r *router) ListTags(ctx context.Context, req ListTagsRequestObject) (ListTagsResponseObject, error) {
	query := control.ListTagsQuery{
		PageSize:  min(cmp.Or(req.Params.PageSize, 10), 50),
		PageAfter: req.Params.PageAfter,
	}

	tags, err := r.memoCtrl.ListTags(ctx, query)
	if err != nil {
		return nil, err
	}

	apiTags := TagList{Items: make([]Tag, len(tags.Items))}
	for i, tag := range tags.Items {
		apiTags.Items[i] = Tag{
			Count: float32(tag.Count),
			Tag:   tag.Tag,
		}
	}

	if tags.Next != nil {
		apiTags.Next = tags.Next
	}

	return ListTags200JSONResponse(apiTags), nil
}

// (GET /attachments)
func (r *router) ListAttachments(ctx context.Context, req ListAttachmentsRequestObject) (ListAttachmentsResponseObject, error) {
	query := control.ListAttachmentsQuery{
		PageSize:  min(cmp.Or(req.Params.PageSize, 10), 50),
		PageAfter: req.Params.PageAfter,
	}

	attachments, err := r.attachmentCtrl.ListAttachments(ctx, query)
	if err != nil {
		return nil, err
	}

	apiAttachments := AttachmentList{Items: make([]Attachment, len(attachments.Items))}
	for i, a := range attachments.Items {
		apiAttachments.Items[i] = Attachment{
			Url:              a.URL(r.baseURL),
			OriginalFilename: a.OriginalFilename,
			ContentType:      a.ContentType,
			Sha256:           fmt.Sprintf("%x", a.Sha256),
			SizeBytes:        a.SizeBytes,
			CreatedBy:        a.CreatedBy.String(),
			CreatedAt:        a.CreatedAt,
		}
	}

	if attachments.Next != nil {
		apiAttachments.Next = attachments.Next
	}

	return ListAttachments200JSONResponse(apiAttachments), nil
}

// (POST /attachments)
func (r *router) CreateAttachment(ctx context.Context, req CreateAttachmentRequestObject) (CreateAttachmentResponseObject, error) {
	content := req.Body
	if req.Params.ContentEncoding != nil && *req.Params.ContentEncoding == "gzip" {
		gr, err := gzip.NewReader(content)
		if err != nil {
			return nil, err
		}
		content = gr
		defer gr.Close()
	}

	id, err := r.attachmentCtrl.CreateAttachment(ctx, control.CreateAttachmentCmd{
		Filename: url.QueryEscape(req.Params.XFilename),
		Content:  content,
	})
	if err != nil {
		return nil, err
	}

	created, err := r.attachmentCtrl.GetAttachment(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAttachmentNotFound) {
			return nil, fmt.Errorf("%w: %v", errNotFound, err)
		}
		return nil, err
	}

	return CreateAttachment201JSONResponse{
		Url:              created.URL(r.baseURL),
		OriginalFilename: created.OriginalFilename,
		ContentType:      created.ContentType,
		Sha256:           fmt.Sprintf("%x", created.Sha256),
		SizeBytes:        created.SizeBytes,
		CreatedAt:        created.CreatedAt,
		CreatedBy:        created.CreatedBy.String(),
	}, nil
}

// (DELETE /attachments/{filename})
func (r *router) DeleteAttachment(ctx context.Context, req DeleteAttachmentRequestObject) (DeleteAttachmentResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// (GET /settings)
func (r *router) GetSettings(ctx context.Context, req GetSettingsRequestObject) (GetSettingsResponseObject, error) {
	settings, err := r.settingsCtrl.Get(ctx)
	if err != nil {
		return nil, err
	}

	return GetSettings200JSONResponse{
		LocaleLanguage:            settings.Locale.Language,
		LocaleRegion:              settings.Locale.Region,
		ThemeColourScheme:         settings.Theme.ColourScheme,
		ThemeMode:                 settings.Theme.Mode,
		ThemeIcon:                 settings.Theme.Icon,
		ThemeListLayout:           settings.Theme.ListLayout,
		ControlsVim:               settings.Controls.Vim,
		ControlsDoubleClickToEdit: settings.Controls.DoubleClickToEdit,
	}, nil
}

// (PATCH /settings)
func (r *router) UpdateSettings(ctx context.Context, req UpdateSettingsRequestObject) (UpdateSettingsResponseObject, error) {
	cmd := control.SetSettingsCmd{
		LocaleLanguage:            req.Body.LocaleLanguage,
		LocaleRegion:              req.Body.LocaleRegion,
		ThemeColourScheme:         req.Body.ThemeColourScheme,
		ThemeMode:                 req.Body.ThemeMode,
		ThemeIcon:                 req.Body.ThemeIcon,
		ThemeListLayout:           req.Body.ThemeListLayout,
		ControlsVim:               req.Body.ControlsVim,
		ControlsDoubleClickToEdit: req.Body.ControlsDoubleClickToEdit,
	}

	err := r.settingsCtrl.Set(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return UpdateSettings204Response{}, nil
}

// List API Tokens paginated
// (GET /apitokens)
func (r *router) ListAPITokens(ctx context.Context, req ListAPITokensRequestObject) (ListAPITokensResponseObject, error) {
	var pageAfter *int64
	if req.Params.PageAfter != nil {
		p, err := strconv.ParseInt(*req.Params.PageAfter, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid pageAfter", errBadRequest)
		}
		pageAfter = &p
	}

	query := control.ListAPITokenQuery{
		PageSize:  req.Params.PageSize,
		PageAfter: pageAfter,
	}

	tokens, err := r.tokenCtrl.ListAPITokens(ctx, query)
	if err != nil {
		return nil, err
	}

	apiTokens := APITokenList{Items: make([]APIToken, len(tokens.Items))}
	for i, token := range tokens.Items {
		apiTokens.Items[i] = APIToken{
			Name:      token.Name,
			CreatedAt: token.CreatedAt,
			ExpiresAt: token.ExpiresAt,
		}
	}

	if tokens.Next != nil {
		next := fmt.Sprint(*tokens.Next)
		apiTokens.Next = &next
	}

	return ListAPITokens200JSONResponse(apiTokens), nil
}

// Create a new API Token
// (POST /apitokens)
func (r *router) CreateAPIToken(ctx context.Context, req CreateAPITokenRequestObject) (CreateAPITokenResponseObject, error) {
	value, err := r.tokenCtrl.CreateAPIToken(ctx, &auth.APIToken{
		Name:      req.Body.Name,
		ExpiresAt: req.Body.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return CreateAPIToken201JSONResponse{Token: value}, nil
}

// Delete API Token
// (DELETE /apitokens/{name})
func (r *router) DeleteAPIToken(ctx context.Context, req DeleteAPITokenRequestObject) (DeleteAPITokenResponseObject, error) {
	err := r.tokenCtrl.DeleteAPIToken(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return DeleteAPIToken204Response{}, nil
}

func (router *router) checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var account *auth.Account
		var err error
		if token, ok := apiTokenFromHeader(r.Header); ok {
			account, err = router.accountFetcher.GetAccountForAPIToken(r.Context(), token)
		} else {
			account, _ = session.Get[*auth.Account](r.Context(), "account")
		}

		if account == nil || errors.Is(err, auth.ErrUnauthorized) {
			errorHandlerFunc(w, r, auth.ErrUnauthorized)
			return
		}

		if err != nil {
			errorHandlerFunc(w, r, &Error{
				Code:  http.StatusInternalServerError,
				Title: http.StatusText(http.StatusInternalServerError),
				Type:  "belt/api/v1/InternalServerError",
			})
			return
		}

		ctx := auth.CtxWithAccount(r.Context(), account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const authHeader = "Authorization"

func apiTokenFromHeader(header http.Header) (auth.APITokenValue, bool) {
	token := header.Get(authHeader)
	if token == "" || token == "Bearer" {
		return nil, false
	}

	return auth.APITokenValue(strings.TrimPrefix(token, "Bearer ")), true
}
