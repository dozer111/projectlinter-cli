include *.mk

.PHONY: check-helm
check-helm:
	@which helm3 > /dev/null 2>&1 || (echo "You're missing helm3 executable"; exit 1)

.PHONY: check-sops
check-sops:
	@which sops > /dev/null 2>&1 || (echo "You're missing sops executable"; exit 1)

.PHONY: check-golang
check-golang:
	@which go > /dev/null 2>&1 || (echo "You're missing go executable"; exit 1)

.PHONY: check-git-version
check-git-version:
	$(eval GIT_EXPECTED_VERSION := "2.23")
	$(eval GIT_FULL_VERSION := $(shell git version | grep -Eo "([[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+)"))
	@(echo required version ${GIT_EXPECTED_VERSION}; git --version) | sort -Vk3 | tail -1 | grep -v -q required || (printf "\033[0;31mWrong git version\nCurrent: ${GIT_FULL_VERSION}, expected: ^${GIT_EXPECTED_VERSION}\n"; exit 1)

.PHONY: check
check: check-git-version \
	check-golang \
	check-helm \
	check-sops

.PHONY: linters  # run all the linters (only for changed files between local and upstream master)
linters: cs-fix \
	rector

.PHONY: cs-fix
cs-fix: # run php-cs-fixer (only for changed files between local and upstream master)
	@printf '\033[34mrun phpCsFixer:\033[0m\n'
	$(eval CHANGED_FILES = $(shell git diff --name-only origin/master --diff-filter=d | grep -i '\.php$$' | tr '\n' ' '))
	@composer run -- php-cs-fixer --config=.php-cs-fixer.dist.php --diff -v --allow-risky=yes --path-mode=intersection -- $(CHANGED_FILES)


.PHONY: rector
rector: # run rector (only for changed files between local and upstream master)
	@printf "\n\033[34mrun rector:\033[0m\n"
	$(eval CHANGED_FILES = $(shell git diff --name-only origin/master --diff-filter=d | grep -i '\.php$$' | tr '\n' ' '))
	@composer run -- rector -- ${CHANGED_FILES}
