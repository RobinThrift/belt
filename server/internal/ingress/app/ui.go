package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"go.robinthrift.com/belt/internal/version"
)

type pageData struct {
	Title    string
	BaseURL  string
	AssetURL string

	ServerData serverData
}

type serverData struct {
	Error     *uiError    `json:"error,omitempty"`
	BuildInfo uiBuildInfo `json:"buildInfo"`
}

type uiBuildInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	CommitDate string `json:"commitDate"`
	GoVersion  string `json:"goVersion"`
}

type uiError struct {
	Code   int
	Title  string
	Detail string
}

var buildInfo = version.GetBuildInfo()

func render(w http.ResponseWriter, data pageData) error {
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

	encoded, err := json.Marshal(data.ServerData)
	if err != nil {
		return fmt.Errorf("error marshalling ui data to json: %w", err)
	}

	tmpldata.Title = data.Title
	tmpldata.AssetURL = data.AssetURL
	tmpldata.BaseURL = template.HTMLAttr(`content="` + data.BaseURL + `"`)
	tmpldata.UIData = template.HTML(encoded)

	return rootTemplate.Execute(w, tmpldata)
}

func renderErrorPage(w http.ResponseWriter, data pageData) error {
	return errorPageTemplate.Execute(w, data)
}

func (ep uiError) Error() string {
	return fmt.Sprintf("%s: %s", ep.Title, ep.Detail)
}
