goflags := "-trimpath -buildmode=pie -mod=readonly -modcacherw -buildvcs=false"
target  := "target"
object  := target / "blap"
version := `git log -n 1 --format=%h`
ldflags := env_var_or_default("LDFLAGS", "")

default: build

build:
  mkdir -p "{{target}}"
  go build {{goflags}} -ldflags "{{ldflags}} -X main.version={{version}}" -o "{{object}}" cmd/main.go

check: build
  go test ./...

clean:
  rm -f "{{object}}"
