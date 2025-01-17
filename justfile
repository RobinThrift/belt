local_bin  := absolute_path("./.bin")
version    := env_var_or_default("VERSION", "dev")

go_ldflags := env_var_or_default("GO_LDFLAGS", "") + " -X 'github.com/RobinThrift/belt/internal.Version=" + version + "'"
go_buildflags := env_var_or_default("GO_BUILD_FLAGS", "-trimpath")
go_test_reporter := env("GO_TEST_REPORTER", "pkgname-and-test-fails")
go_lint_reporter := env("GO_LINT_REPORTER", "colored-line-number")

import ".scripts/database.justfile"
import ".scripts/api.justfile"
import ".scripts/oci.justfile"

_default:
    @just --list

# run the binary locally and restart on file changes
run: (_install-tool "watchexec")
    @mkdir -p test/manual
    @BELT_LOG_LEVEL="debug" BELT_LOG_FORMAT="console" \
        BELT_ADDR="localhost:8081" \
        BELT_SECURE_COOKIES="false" \
        BELT_DATABASE_PATH="./test/manual/belt.db" \
        BELT_DATABASE_DEBUG_ENABLED="true" \
        BELT_ATTACHMENTS_DIR="./test/manual/attachments" \
        BELT_INIT_USERNAME="user" \
        BELT_INIT_PASSWORD="password" \
        {{ local_bin }}/watchexec -r -e go -- go run -trimpath -tags dev ./bin/belt

run-prod:
    cd ui && just install build
    just build
    @BELT_LOG_LEVEL="debug" BELT_LOG_FORMAT="console" \
        BELT_ADDR="localhost:8081" \
        BELT_SECURE_COOKIES="false" \
        BELT_DATABASE_PATH="./test/manual/belt.db" \
        BELT_DATABASE_DEBUG_ENABLED="true" \
        BELT_ATTACHMENTS_DIR="./test/manual/attachments" \
        BELT_INIT_USERNAME="user" \
        BELT_INIT_PASSWORD="password" \
        ./build/belt

build:
    go build -ldflags="{{go_ldflags}}" {{ go_buildflags }} -o build/belt ./bin/belt

export_report := env_var_or_default("EXPORT_REPORT", "false")
test *flags="-v -failfast -timeout 15m ./...": (_install-tool "gotestsum")
    @mkdir -p ui/build && touch ./ui/build/index.js
    {{ local_bin }}/gotestsum \
        {{ if export_report == "true" { "--junitfile 'tests.junit.xml' --junitfile-project-name 'RobinThrift/belt' --junitfile-hide-empty-pkg" } else { "" } }} \
        --format {{ go_test_reporter }} \
        --format-hide-empty-pkg \
        -- {{ go_buildflags }} {{ flags }}

test-watch *flags="-v -failfast -timeout 15m ./...": (_install-tool "gotestsum")
    @mkdir -p ui/build && touch ./ui/build/index.js
    {{ local_bin }}/gotestsum --format {{ go_test_reporter }} --format-hide-empty-pkg --watch -- {{ go_buildflags }} {{ flags }}

# lint using staticcheck and golangci-lint
lint: (_install-tool "staticcheck") (_install-tool "golangci-lint") (_install-tool "sqlc") (_install-tool "vacuum")
    @mkdir -p ui/build && touch ./ui/build/index.js
    {{ local_bin }}/staticcheck ./...
    {{ local_bin }}/golangci-lint run --out-format={{ go_lint_reporter }} ./...
    {{ local_bin }}/sqlc -f internal/storage/database/sqlite/sqlc.yaml vet
    {{ local_bin }}/vacuum lint --ruleset .scripts/vacuum.yaml --details --fail-severity warn --no-banner --all-results ./api/apiv1/apiv1.openapi3.yaml

fmt:
    @go fmt ./...

generate: _gen-sqlc _gen-api-v1-server


clean:
    rm -rf {{ local_bin }} build ui/build ui/node_modules test/manual
    go clean -cache

# generate a release with the given tag
release tag:
    just changelog {{tag}}
    git add CHANGELOG
    git commit -m "Releasing version {{tag}}"
    git tag {{tag}}
    git push
    git push origin {{tag}}

# generate a changelog using https://github.com/orhun/git-cliff
changelog tag: (_install-tool "git-cliff")
    git-cliff --config CHANGELOG/cliff.toml -o CHANGELOG/CHANGELOG-{{tag}}.md --unreleased --tag {{ tag }} 
    echo "- [CHANGELOG-{{tag}}.md](./CHANGELOG-{{tag}}.md)" >> CHANGELOG/README.md


_install-tool tool:
    @cd ./.scripts/toolfetcher && go run . -to {{ local_bin }} -versionfile ../TOOL_VERSIONS {{ tool }}
