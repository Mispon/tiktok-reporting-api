# TikTok Reporting API
Automatically or through a method call in the Rest API downloads marketing data from the TikTok API

## Setup
1. Populate environment variables in docker-compose.yml
    - **ports** - the link of the local port, to the port of the application in the container, eg. `10010:80`, server traffic on port 10010 will be forwarded inside the container to port 80
    - **API_ENDPOINT** - host and port of the application inside the container, eg. `0.0.0.0:1234`, then in ports you will need to specify `10010:1234`
    - **TIKTOK_APP_ID** - TikTok application ID
    - **TIKTOK_APP_SECRET** - TikTok application secret
    - **TIKTOK_APP_TOKEN** - TikTok application actual token, if there is
    - **BQ_PROJECT_ID** - Google bigquery ID
    - **BQ_DATASET_ID** - name of the dataset in google bigquery
    - **BQ_AUC_TABLE_ID** - name of the table for data on auctions in google bigquery
    - **BQ_RES_TABLE_ID** - name of the table for data on reservation in google bigquery
    - **JOB_INTERVAL_HOURS** - how often to run the job, in hours
    - **STATISTIC_DEPTH_DAYS** - with what depth to request data, in days, eg. 7 mean for the last 7 days
2. Fill `credentials.json` config in folder `configs/`
3. If there are `advertising_ids`, you can put it into the file `advert_ids.txt` in folder `configs/`
   each one on the new line. At the start, the application counts the IDs from the file and will collect statistics on them. **
   Auth callback will overwrite the IDs from the file.**

## Run
1. Run `docker compose up -d` in project's root folder
2. Bind the authorization callback `0.0.0.0:80/auth/callback` (change 0.0.0.0:80 to your domain/host:port) in TikTok personal area

## Usage
In the service, the job starts at a specified interval, takes the IDs received either from the file at startup or from the response
in the authorization callback, receives data from the tiktok API and writes it to two tables in bigquery, by auction and reservation.
You can run the job several times a day and for previous dates, more recent statistics will overwrite the old one.
You can get statistics for a specific campaign for a specified period by calling the http api methods: `/report/auction`
and `/report/reservation`,
For example: `http://localhost:80/report/auction?advertiser_id=123456789&start_date=2021-09-29&end_date=2021-10-29`  
Marketing company data will be returned in the reply and will be recorded in bigquery
