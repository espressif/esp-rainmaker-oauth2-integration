SILENT-$(VERBOSE) := @

# In the short-term, we will support packages as well as executables

# Service Name
SERVICE-NAME := rainmaker_oauth2_integration

# Versioning
MAJOR_VERSION:=1
MINOR_VERSION:=0
REVISION:=0
BUILD_DATE ?= $(shell date +%Y-%m-%dT%H:%M)
BUILD_VERSION ?= $(shell git log --pretty=format:"%h" -1)

help:
	@echo ""
	#command to build package
	@echo "make <PACKAGE_NAME> "
	@echo ""
	#command to build all packages
	@echo "make all "
	@echo ""
	#command to deploy package
	@echo "make <PACKAGE_NAME>-deploy S3-BUCKET=<YOUR_BUCKET_NAME> [STAGE-NAME=<YOUR_STAGE_NAME>]  "
	@echo ""
	#command to build and deploy all packages
	@echo "make deploy S3-BUCKET=<YOUR_BUCKET_NAME> [STAGE-NAME=<YOUR_STAGE_NAME>]  "
	@echo ""
	#command to publish package
	@echo "make <PACKAGE_NAME>-publish S3-BUCKET=<YOUR_BUCKET_NAME> VERSION=<ENTER_VERSION_TO_PUBLISH> "
	@echo ""
	#command to publish all packages
	@echo "make publish S3-BUCKET=<YOUR_BUCKET_NAME> VERSION=<ENTER_VERSION_TO_PUBLISH> "
	@echo ""

STAGE-NAME := dev
VERSION := 1.0.0
REGION := us-east-1

# all packages
ALL_PKGS  :=  espoauth2integration

# all packages
DEPLOY_PKGS :=  espoauth2integration

PUBLISH_PKGS := espoauth2integration

# Add the App version, build number, build date and customer deployment using LDFlags
LDFLAGS += "-s -w -X main.appName=$(SERVICE-NAME) -X main.appVersion=$(MAJOR_VERSION).$(MINOR_VERSION).$(REVISION)-$(BUILD_VERSION) -X main.appBuildDate=$(BUILD_DATE)"

publish:  $(ALL_PKGS)

all: $(ALL_PKGS)

deploy: $(ALL_PKGS)

dep_ensure:
	#dep ensure

define GetValueFromConfig
$(shell node -p "require('./config.json').$(1)")
endef

define pkg_targets

$(1): $$(wildcard src/handlers/$(1)/executables/*)

$(1)-deploy:
	@[ "$(S3-BUCKET)" ] || ( echo ">> Please enter valid s3 bucket name or create a bucket using command: aws s3 mb <YOUR_BUCKET_NAME>"; exit 1)
ifneq ($(filter $(1),$(DEPLOY_PKGS)),)
	$$(SILENT-)sam package --template-file src/handlers/*/$(1).yml --output-template-file $(1)_package.yml --s3-bucket "$(S3-BUCKET)"

	$$(SILENT-)sam deploy --template-file $(1)_package.yml --stack-name $(1) --capabilities CAPABILITY_NAMED_IAM --no-fail-on-empty-changeset --parameter-overrides "StageName=$(STAGE-NAME)" \
	$(2)
endif

deploy: $(1)-deploy

$(1)-publish:
	@[ "$(S3-BUCKET)" ] || ( echo ">> Please enter valid s3 bucket name or create a bucket using command: aws s3 mb <YOUR_BUCKET_NAME>"; exit 1)
ifneq ($(filter $(1),$(PUBLISH_PKGS)),)
	@sam package --template-file src/handlers/*/$(1).yml --output-template-file $(1)_package.yml --s3-bucket "$(S3-BUCKET)" --s3-prefix "$(VERSION)"/$(1)
	$$(SILENT-)sam publish -t $(1)_package.yml --semantic-version "$(VERSION)" --region "$(REGION)"
endif

publish: $(1)-publish

$$(wildcard src/handlers/$(1)/executables/*): dep_ensure
	@echo "[go] $(1) => $$(notdir $$@)"
	$$(SILENT-)GOOS=linux go build -ldflags $(LDFLAGS)  -o bin/handlers/$(1)/$$(notdir $$@) $$@/*.go

endef

$(foreach pkg,$(ALL_PKGS),$(eval $(call pkg_targets,$(pkg))))
