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
