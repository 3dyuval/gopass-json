CURRENT := $(shell grep 'const version' cmd/gopass-json/main.go | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
MAJOR   := $(word 1, $(subst ., ,$(CURRENT)))
MINOR   := $(word 2, $(subst ., ,$(CURRENT)))
PATCH   := $(word 3, $(subst ., ,$(CURRENT)))

NEXT_PATCH := $(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
NEXT_MINOR := $(MAJOR).$(shell echo $$(($(MINOR)+1))).0
NEXT_MAJOR := $(shell echo $$(($(MAJOR)+1))).0.0

.PHONY: bump patch minor major

bump: patch

patch:
	@$(MAKE) _bump NEXT=$(NEXT_PATCH)

minor:
	@$(MAKE) _bump NEXT=$(NEXT_MINOR)

major:
	@$(MAKE) _bump NEXT=$(NEXT_MAJOR)

_bump:
	@echo "bumping $(CURRENT) → $(NEXT)"
	sed -i 's/const version = "$(CURRENT)"/const version = "$(NEXT)"/' cmd/gopass-json/main.go
	sed -i 's/version-v$(CURRENT)-blue/version-v$(NEXT)-blue/' README.md
	sed -i 's|releases/tag/v$(CURRENT)|releases/tag/v$(NEXT)|' README.md
	git add cmd/gopass-json/main.go README.md
	git commit -m "chore: bump version to $(NEXT)"
	git push origin master
	git tag -a v$(NEXT) -m "v$(NEXT)"
	git push origin v$(NEXT)
