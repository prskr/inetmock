.PHONY: book deploy

book:
	@mdbook build

deploy: book
	@rm -rf /tmp/book
	@git worktree remove -f /tmp/book
	@git worktree add /tmp/book gh-pages
	@rm -rf /tmp/book/*
	@cp -rp book/* /tmp/book/
	@cd /tmp/book && \
		git add -A && \
		git commit -m "deployed on $(shell date) by ${USER}" && \
		git push origin gh-pages