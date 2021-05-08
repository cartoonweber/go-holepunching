build:
	docker build -t batphonghan/holepunching-gov2 .
run:
	docker run -p 8080:8080 batphonghan/holepunching-gov2
deploy:
	gcloud builds submit --tag asia.gcr.io/hazelcoin/holepunching-gov2