package ui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/RobinThrift/belt/internal/domain"
	"github.com/RobinThrift/belt/internal/version"
)

//go:embed src/html/error.tmpl.html
var errorPageTemplateRaw string

var errorPageTemplate = template.Must(template.New("error.html").Parse(errorPageTemplateRaw))

type PageData struct {
	Title     string
	BaseURL   string
	AssetURL  string
	Icon      string
	CSRFToken string

	ServerData ServerData
}

type Account struct {
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Settings struct {
	Locale   LocaleSettings   `json:"locale,omitempty"`
	Theme    ThemeSettings    `json:"theme,omitempty"`
	Controls ControlsSettings `json:"controls,omitempty"`
}

type LocaleSettings struct {
	Language string `json:"language,omitempty"`
	Region   string `json:"region,omitempty"`
}

type ThemeSettings struct {
	ColourScheme string `json:"colourScheme,omitempty"`
	Mode         string `json:"mode,omitempty"`
	Icon         string `json:"icon,omitempty"`
	ListLayout   string `json:"listLayout,omitempty"`
}

type ControlsSettings struct {
	Vim               bool `json:"vim,omitempty"`
	DoubleClickToEdit bool `json:"doubleClickToEdit,omitempty"`
}

func NewSettings(s *domain.Settings) *Settings {
	return &Settings{
		Locale: LocaleSettings{
			Language: s.Locale.Language,
			Region:   s.Locale.Region,
		},
		Theme: ThemeSettings{
			ColourScheme: s.Theme.ColourScheme,
			Mode:         s.Theme.Mode,
			Icon:         s.Theme.Icon,
			ListLayout:   s.Theme.ListLayout,
		},
		Controls: ControlsSettings{
			Vim:               s.Controls.Vim,
			DoubleClickToEdit: s.Controls.DoubleClickToEdit,
		},
	}
}

type ServerData struct {
	Account    *Account             `json:"account,omitempty"`
	Settings   *Settings            `json:"settings,omitempty"`
	Components ServerDataComponents `json:"components,omitempty"`
	Error      *UIError             `json:"error,omitempty"`
	BuildInfo  BuildInfo            `json:"buildInfo"`
}

type BuildInfo struct {
	Version     string `json:"version"`
	CommitHash  string `json:"commitHash"`
	CommitDate  string `json:"commitDate"`
	GoVersion   string `json:"goVersion"`
	ProjectLink string `json:"projectLink"`
}

type ServerDataComponents struct {
	LoginPage               LoginPage               `json:"LoginPage,omitempty"`
	LoginChangePasswordPage LoginChangePasswordPage `json:"LoginChangePasswordPage,omitempty"`
	SettingsPage            SettingsPage            `json:"SettingsPage,omitempty"`
}

type UIError struct {
	Code   int
	Title  string
	Detail string
}

var buildInfo = version.GetBuildInfo()

func Render(w http.ResponseWriter, data PageData) error {
	var tmpldata struct {
		UIData    template.HTML
		CSRFToken template.HTMLAttr
		AssetURL  string
		Icon      string
		BaseURL   template.HTMLAttr
		Title     string
	}

	data.ServerData.BuildInfo.Version = buildInfo.Version
	data.ServerData.BuildInfo.CommitHash = buildInfo.Hash
	data.ServerData.BuildInfo.CommitDate = buildInfo.Date
	data.ServerData.BuildInfo.GoVersion = buildInfo.GoVersion
	data.ServerData.BuildInfo.ProjectLink = "http://github.com/RobinThrift/belt"

	encoded, err := json.Marshal(data.ServerData)
	if err != nil {
		return fmt.Errorf("error marshalling ui data to json: %w", err)
	}

	tmpldata.Title = data.Title
	tmpldata.AssetURL = data.AssetURL
	tmpldata.Icon = data.Icon
	tmpldata.BaseURL = template.HTMLAttr(`content="` + data.BaseURL + `"`)
	tmpldata.CSRFToken = template.HTMLAttr(`content="` + data.CSRFToken + `"`)
	tmpldata.UIData = template.HTML(encoded)

	return rootTemplate.Execute(w, tmpldata)
}

func RenderErrorPage(w http.ResponseWriter, data PageData) error {
	return errorPageTemplate.Execute(w, data)
}

func (ep UIError) Error() string {
	return fmt.Sprintf("%s: %s", ep.Title, ep.Detail)
}

type LoginPage struct {
	RedirectURL      string            `json:"redirectURL,omitempty"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}

type LoginChangePasswordPage struct {
	RedirectURL      string            `json:"redirectURL,omitempty"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}

type SettingsPage struct {
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}
