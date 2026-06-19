binary := "hcptf"
install_dir := env_var_or_default("PREFIX", env_var("HOME") + "/.local") + "/bin"

default:
    @just --list

build:
    go build -o {{binary}} .

install: build
    install -d {{install_dir}}
    install -m 0755 {{binary}} {{install_dir}}/{{binary}}

test:
    go test ./...
