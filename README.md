# janction-backend-api

## Configure settings

Create a configuration file (e.g., `config.yaml`) with the necessary settings for your environment.

## API Endpoints

### Authentication

| Method | Endpoint             | Description                         |
| ------ | -------------------- | ----------------------------------- |
| GET    | `/api/v1/auth/nonce` | Get a nonce for the login           |
| POST   | `/api/v1/auth/login` | Login with wallet address and nonce |

### Users

| Method | Endpoint        | Description                |
| ------ | --------------- | -------------------------- |
| GET    | `/api/v1/user/` | Get user by wallet address |

### Nodes

| Method | Endpoint                    | Description                            |
| ------ | --------------------------- | -------------------------------------- |
| POST   | `/api/v1/node/heartbeat`    | Receive a heartbeat signal from a node |
| POST   | `/api/v1/node/job`          | Receive job status submitted by node   |
| GET    | `/api/v1/node/online_count` | Get the count of online nodes          |
| GET    | `/api/v1/node/info`         | Get information for a specific node    |
| GET    | `/api/v1/node/infos`        | Get information for multiple nodes     |
| GET    | `/api/v1/node/logs`         | Get node logs                          |
