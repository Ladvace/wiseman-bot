runmongo:
	docker rm -f wisemanbotmongo
	docker pull mongo
	
	docker run --name wisemanbotmongo -t -i -p 27017:27017 mongo mongod