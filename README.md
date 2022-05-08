# Rabbits

## RabbitMq Docker Setup
### creates a docker network to run rabbitmq
```
docker network create rabbits
```
### Creates host rabbitmq container
```
docker run -d --rm --net rabbits -p 6660:15672 --hostname rabbit-master --name rabbit-master rabbitmq:3
```
### Setups the management Prometheus plugin
```
docker exec -it rabbit-master bash
rabbitmq-plugins enable rabbitmq_management
```
> Prometheus is available on the port 6660  
> default user: guest password: guest

### Building publisher image
> make sure you are inside /publisher folder
```
docker build . -t gkm2806/learning-rabbitmq
```
### Runing image on the rabbitmq network
```
docker run --rm -p 4000:4000 --env-file ./.env --network rabbits gkm2806/learning-rabbitmq
```