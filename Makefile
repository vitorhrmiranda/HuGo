dev:
	hugo server -D

release:
	hugo

install:
	./scripts/install_dependencies.sh

push:
	./scripts/push_to_pages.sh
