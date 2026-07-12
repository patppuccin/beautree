#!/usr/bin/env bash

set -e

GITIGNORE=$(cat << 'EOF'
ignored/
*.secret
mode_nodules/
EOF
)

CMD="$1"
shift

ROOT=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --dir)
            ROOT="$2"
            shift 2
            ;;
        *)
            echo "Unknown argument: $1"
            echo "Usage: scaffold.sh [create|cleanup] --dir PATH"
            exit 1
            ;;
    esac
done

if [[ -z "$ROOT" ]]; then
    echo "Error: --dir is required"
    echo "Usage: scaffold.sh [create|cleanup] --dir PATH"
    exit 1
fi

create() {
    mkdir -p "$ROOT"

    echo "$GITIGNORE" > "$ROOT/.gitignore"

    touch "$ROOT/.hidden-config"
    touch "$ROOT/.env"
    touch "$ROOT/README.md"
    touch "$ROOT/definitely-not-production.toml"

    mkdir -p "$ROOT/ignored"
    touch "$ROOT/ignored/my-passwords-totally-safe-here.txt"
    touch "$ROOT/my-aws-keys.secret"

    mkdir -p "$ROOT/src/core"
    mkdir -p "$ROOT/src/api"
    touch "$ROOT/src/main.go"
    touch "$ROOT/src/core/the-algorithm.go"
    touch "$ROOT/src/api/move-fast-break-things.go"

    mkdir -p "$ROOT/mode_nodules/left-pad/src/utils"
    touch "$ROOT/mode_nodules/left-pad/index.js"
    touch "$ROOT/mode_nodules/left-pad/src/left-pad.js"
    touch "$ROOT/mode_nodules/left-pad/src/utils/pad-left-utils.js"

    mkdir -p "$ROOT/assets"
    touch "$ROOT/assets/vim-vs-emacs-100hr-documentary.mp4"
    touch "$ROOT/assets/logo-final-v3-ACTUALLY-FINAL.png"

    mkdir -p "$ROOT/docs"
    touch "$ROOT/docs/guide.md"
    touch "$ROOT/docs/todo-rewrite-in-rust.md"

    ln -sf "./src" "$ROOT/link-to-src"
    ln -sf "./nonexistent" "$ROOT/broken-link"

    echo "Demo created at $ROOT"
}

cleanup() {
    rm -rf "$ROOT"
    echo "Demo cleaned up at $ROOT"
}

case "$CMD" in
    create)  create  ;;
    cleanup) cleanup ;;
    *)
        echo "Usage: scaffold.sh [create|cleanup] --dir PATH"
        exit 1
        ;;
esac