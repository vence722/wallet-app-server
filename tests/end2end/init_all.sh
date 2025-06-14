# shutdown existing docker compose
docker-compose down;
# start docker compose (DB + Redis)
docker-compose up -d;
# wait for DB start
sleep 3;
# create tables
docker exec --user postgres -it end2end-db-1 bash -c "psql -f /wallet-app-server/database/wallet-app.sql";
# insert test data
docker exec --user postgres -it end2end-db-1 bash -c "psql -f /wallet-app-server/tests/end2end/test_data.sql";