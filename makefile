buildback:
	CGO_ENABLED=0 GOOS=linux go build -o ./backend/main ./backend/cmd/server/
buildfront:
	CGO_ENABLED=0 GOOS=linux go build -o ./frontend/main ./frontend/cmd/server/
backdocker:
	docker build -t backend ./backend
frontdocker:
	docker build -t frontend ./frontend
backrun:
	docker run --rm -p8080:8080 -p33061:3306 --name backend backend
frontrun:
	docker run --rm -p8080:8080 -p33061:3306 --name frontend frontend
composebuild:
	docker-compose build
composeup:
	docker-compose up
full: buildback buildfront composebuild composeup
.DEFAULT_GOAL = full