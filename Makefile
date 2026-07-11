.PHONY: mysql-cli redis-cli

mysql-cli:
	docker exec -it url-mysql mysql -u root -p123456 url


