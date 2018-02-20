.PHONY: all
all: build package terraform

dist:
	mkdir -p dist

.PHONY: build
build: dist
	cd dist && GOOS=linux go build ..

.PHONY: package
package: dist
	cd dist && zip strava-commute.zip strava-commute config.json

.PHONY: clean
clean:
	rm -rf dist

.PHONY: default_region
default_region:
	$(eval export AWS_DEFAULT_REGION ?= eu-west-1)

.PHONY: terraform
terraform: default_region
	$(if ${STRAVA_API_TOKEN},,$(error must set STRAVA_API_TOKEN))
	cd terraform && terraform apply \
		-var "strava_api_token=${STRAVA_API_TOKEN}"

.PHONY: register
register: default_region
	$(if ${STRAVA_CLIENT_ID},,$(error must set STRAVA_CLIENT_ID))
	$(if ${STRAVA_CLIENT_SECRET},,$(error must set STRAVA_CLIENT_SECRET))
	$(eval export URL=$$(shell \
		cd terraform && AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} terraform output url \
	))
	curl -XPOST \
		-F "client_id=${STRAVA_CLIENT_ID}" \
		-F "client_secret=${STRAVA_CLIENT_SECRET}" \
		-F 'verify_token=STRAVA' \
		-F "callback_url=${URL}" \
		https://api.strava.com/api/v3/push_subscriptions

.PHONY: init
init: default_region
	$(if ${BUCKET_SUFFIX},,$(error must set BUCKET_SUFFIX))
	$(eval BUCKET_NAME=strava-commute-${BUCKET_SUFFIX})
	aws s3api create-bucket \
		--acl private \
		--bucket "$(BUCKET_NAME)" \
		--create-bucket-configuration "LocationConstraint=${AWS_DEFAULT_REGION}"
	cd terraform && terraform init \
		--backend-config "bucket=$(BUCKET_NAME)"
