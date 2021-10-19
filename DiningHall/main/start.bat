docker stop dining_hall_container
docker run -d --rm -p 7500:7500 --name dining_hall_container dining_hall_image go run main http://host.docker.internal