
APP_YAML_DIR := $(HOME)/go/appengine

all:

devserver:
	python2.7 `which dev_appserver.py` $(APP_YAML_DIR)/app.yaml

deploy:
	gcloud app deploy
