# Start Chainlink Node
1/ Start Postgres 
+ postgres -D /usr/local/var/postgres; psql; CREATE USER root; ALTER USER root PASSWORD 'root'
OR
docker run --name postgres-chainlink -v db:/var/lib/postgresql/data -e POSTGRES_PASSWORD=myPostgresPW -d -p 5432:5432 postgres:11.12
docker exec -it postgres-chainlink psql -U postgres -c "CREATE USER chainlink WITH PASSWORD 'myChainlinkPW';"
docker exec -it postgres-chainlink psql -U postgres -c "CREATE DATABASE "chainlink_rinkeby";"
docker exec -it postgres-chainlink psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE "chainlink_rinkeby" TO chainlink;"

2/ Start Ethereum Client Node (in future need to host our own)
+ docker pull ethereum/client-go:latest
+ docker run --name eth -p 8546:8546 -v ~/.geth-rinkeby:/geth -it \
           ethereum/client-go --rinkeby --ws --ipcdisable \
           --ws.addr 0.0.0.0 --ws.origins="*" --datadir /geth
+ docker start -i eth

ETH_CONTAINER_IP=$(docker inspect --format '' $(docker ps -f name=eth -q))
echo "ETH_URL=ws://$ETH_CONTAINER_IP:8546" >> ~/.chainlink-rinkeby/.env

3/ Start Chainlink Node:
cd chainlink && docker run -p 6688:6688 -v chainlink:/chainlink -it --env-file=.env smartcontract/chainlink:0.10.13 local n
OR 
docker run --name chainlink_rinkeby --network host -p 6688:6688 -v chainlink:/chainlink -it --env-file=chainlink/.env smartcontract/chainlink:0.10.8 local n
(pass: HIEU1998!!!1998Hieu)

# Deploy smart contracts
Use online IDE https://remix.ethereum.org/

https://titanwolf.org/Network/Articles/Article?AID=d979ddb0-3f53-4d6a-a82d-9e5729bee3ab#gsc.tab=0
remixd -s /Users/hieuletrung/Documents/repos/side_projects/cloud-morph-host/oracle/contracts --remix-ide http://127.0.0.1:8080
