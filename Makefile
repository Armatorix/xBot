# XBOT_USERNAME=polski_wojt
XBOT_USERNAME=Armatorix
SSHADDR=ax@hephaestus
SHELL:=/bin/bash
BOT_NAME=xbot-worker-${XBOT_USERNAME}
APP_PATH=~/${BOT_NAME}
WORKER_SYSTEMD_CONFIG_PATH=/etc/systemd/system/${BOT_NAME}.service
WORKER_SYSTEMD_TIMER_PATH=/etc/systemd/system/${BOT_NAME}.timer
.PHONY: run
run:
# run with .env, source .env and run go main
	@echo "Running xBot..."
	@source .env && go run main.go

.PHONY: install-playwright-driver
install-playwright-driver:
	@echo "Installing Playwright driver..."
	@ssh ${SSHADDR} ""
	@ssh ${SSHADDR} "playwright install chromium"

.PHONY: build
build:
	@echo "Building xBot..."
	@mkdir -p bin
	@go build -o ./bin/xbot  .

.PHONY: prepare-user
prepare-user:
	@echo "Preparing user ${XBOT_USERNAME} on ${SSHADDR}..."
	@ssh ${SSHADDR} "mkdir -p ${APP_PATH}"

.PHONY: deploy-systemd-config
deploy-systemd-config:
	@echo "Building xBot systemd config..."
	scp ./worker/systemd.config ${SSHADDR}:${APP_PATH}/systemd.config
	@ssh ${SSHADDR} "sed -i 's|PATH_TO_SCRIPT|${HOME}/${BOT_NAME}/xbot|g' ${HOME}/${BOT_NAME}/systemd.config"
	@ssh ${SSHADDR} "sed -i 's|PATH_TO_ENV|${HOME}/${BOT_NAME}/.env|g' ${HOME}/${BOT_NAME}/systemd.config"
	@ssh ${SSHADDR} "sed -i 's|WORKDIR|${HOME}/${BOT_NAME}|g' ${HOME}/${BOT_NAME}/systemd.config"
	@ssh ${SSHADDR} "sudo mv ${APP_PATH}/systemd.config ${WORKER_SYSTEMD_CONFIG_PATH}"
	@ssh ${SSHADDR} "sudo chown root:root ${WORKER_SYSTEMD_CONFIG_PATH}"
	@ssh ${SSHADDR} "sudo chmod 777 ${WORKER_SYSTEMD_CONFIG_PATH}"

	@echo "Copy systemd timer file..."
	@scp ./worker/systemd.timer ${SSHADDR}:${APP_PATH}/systemd.timer
	@ssh ${SSHADDR} "sudo mv ${APP_PATH}/systemd.timer ${WORKER_SYSTEMD_TIMER_PATH}"
	@ssh ${SSHADDR} "sudo chown root:root ${WORKER_SYSTEMD_TIMER_PATH}"
	@ssh ${SSHADDR} "sudo chmod 644 ${WORKER_SYSTEMD_TIMER_PATH}"

.PHONY: deploy-binary
deploy-binary:
	@echo "Deploying xBot binary to ${SSHADDR}..."
	@ssh ${SSHADDR} "rm -f ${APP_PATH}/xbot || true"
	@scp ./bin/xbot ${SSHADDR}:${APP_PATH}/xbot
	@ssh ${SSHADDR} "chmod +x ${APP_PATH}/xbot"

.PHONY: deploy-restart-service
deploy-restart-service:
	@echo "Using new binary and restarting xBot service on ${SSHADDR}..."
	@ssh ${SSHADDR} "sudo systemctl stop ${BOT_NAME}.timer || true"
	@ssh ${SSHADDR} "sudo systemctl stop ${BOT_NAME}.service || true"
	@ssh ${SSHADDR} "sudo systemctl daemon-reload"
	@ssh ${SSHADDR} "sudo systemctl enable ${BOT_NAME}.service || true"
	@ssh ${SSHADDR} "sudo systemctl enable ${BOT_NAME}.timer"
	@ssh ${SSHADDR} "sudo systemctl start ${BOT_NAME}.timer || true"
	@ssh ${SSHADDR} "sudo systemctl start ${BOT_NAME}.service || true"

.PHONY: deploy-envs
deploy-envs:
	@echo "Deploying .env file to ${SSHADDR}..."
	@scp .env.${XBOT_USERNAME} ${SSHADDR}:${APP_PATH}/.env
	@ssh ${SSHADDR} "chmod 644 ${APP_PATH}/.env"


.PHONY: worker-deploy
worker-deploy:prepare-user build deploy-systemd-config deploy-binary deploy-envs deploy-restart-service

.all:
all:
	$(MAKE) worker-deploy XBOT_USERNAME=polski_wojt
	$(MAKE) worker-deploy XBOT_USERNAME=Armatorix