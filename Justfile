PROJECT_NAME := "GzipDate"

alias arc := archive

@_default:
	just _term-wipe
	just --list

# Archive GoReleaser dist
archive:
	#!/bin/sh
	just _term-wipe
	tag="$(git tag --points-at master)"
	app="{{PROJECT_NAME}}"
	arc="${app}_${tag}"

	# echo "app = '${app}'"
	# echo "tag = '${tag}'"
	# echo "arc = '${arc}'"
	if [ -e dist ]; then
		echo "Move dist -> distro/${arc}"
		mv dist "distro/${arc}"

		# echo "cd distro"
		cd distro
		
		printf "pwd = "
		pwd
		
		ls -Alh
	else
		echo "dist directory not found for archiving"
	fi

# Build and install app
build:
	@just _term-wipe
	go build -o gzipdate main.go
	mv gzipdate "${GOBIN}/"
	@# go install main.go


# Build distro
distro:
	#!/bin/sh
	goreleaser
	just archive


# Run code
run +args='':
	@just _term-wipe
	go run main.go {{args}}


_term-wipe:
	#!/bin/sh
	if [[ ${#VISUAL_STUDIO_CODE} -gt 0 ]]; then
		clear
	elif [[ ${KITTY_WINDOW_ID} -gt 0 ]] || [[ ${#TMUX} -gt 0 ]] || [[ "${TERM_PROGRAM}" = 'vscode' ]]; then
		printf '\033c'
	elif [[ "$(uname)" == 'Darwin' ]] || [[ "${TERM_PROGRAM}" = 'Apple_Terminal' ]] || [[ "${TERM_PROGRAM}" = 'iTerm.app' ]]; then
		osascript -e 'tell application "System Events" to keystroke "k" using command down'
	elif [[ -x "$(which tput)" ]]; then
		tput reset
	elif [[ -x "$(which reset)" ]]; then
		reset
	else
		clear
	fi

