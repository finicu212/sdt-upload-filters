cnf ?= config.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

ftp:
	docker run -d \
		--env-file=./config.env \
		-p $(PORT):$(PORT) \
		-p $(PSV_PORTS) \
		-e $(LIST_USER) \
		--name my_ftp \
		delfer/alpine-ftp-server && docker ps

stop:
	docker stop my_ftp; docker rm my_ftp && docker ps