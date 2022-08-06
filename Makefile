.PHONY: healthplanet-gettoken healthplanet-to-fitbit fitbit-gettoken

fitbit-gettoken:
	go build -o bin/fitbit-gettoken ./cmd/fitbit-gettoken/*.go

healthplanet-gettoken:
	go build -o bin/healthplanet-gettoken ./cmd/healthplanet-gettoken/*.go 

healthplanet-to-fitbit:
	go build -o bin/healthplanet-to-fitbit ./cmd/healthplanet-to-fitbit/*.go 

all: healthplanet-gettoken healthplanet-to-fitbit fitbit-gettoken