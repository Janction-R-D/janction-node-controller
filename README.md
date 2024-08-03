# janction-node-controller

## Introduction

The `janction-node-controller` is a component responsible for managing and distributing computational jobs to nodes. It handles job retrieval, node registration, and node health checks through periodic pings.

## Configuration

The parameters and settings for the node controller are specified in the `config.yaml` file. Ensure that the backend URL and other necessary configurations are set correctly.

## Main Features

1. **Job Retrieval**:
   - Fetches jobs from the backend and stores them in a job queue.
   - Distributes jobs to nodes upon request.

2. **Node Registration**:
   - Allows nodes to register with the controller.
   - Stores information about the node's architecture and resource usage capabilities (CPU/GPU).

3. **Node Health Check**:
   - Nodes can send periodic pings to indicate they are still active.
   - Updates the node's last active time in the database.

## API Endpoints

### Job Management

| Method | Endpoint          | Description                       |
| ------ | ----------------- | --------------------------------- |
| GET    | `/api/controller/v1/job` | Retrieves a job for a node     |

### Node Registration

| Method | Endpoint                   | Description                       |
| ------ | -------------------------- | --------------------------------- |
| POST   | `/api/controller/v1/register` | Registers a new node             |

### Node Health Check

| Method | Endpoint                   | Description                       |
| ------ | -------------------------- | --------------------------------- |
| POST   | `/api/controller/v1/ping` | Sends a ping to update node status |

## Usage Instructions

1. **Configure Parameters**:
   Ensure the parameters in the `config.yaml` file are correctly set.

2. **Run the Node Controller**:
   Use the provided setup functions to start the Gin HTTP server and handle incoming requests.

3. **Register Nodes**:
   Nodes should register themselves using the `/register` endpoint with their ID, architecture type, and resource usage capabilities.

4. **Job Retrieval**:
   Nodes can request jobs by calling the `/job` endpoint. The controller will distribute jobs from the job queue.

5. **Node Health Check**:
   Nodes should periodically send pings to the `/ping` endpoint to keep their status updated.

## Contribution

Contributions and suggestions for improvements to the Janction Node Controller are welcome. Please ensure thorough testing and documentation updates before submitting changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
