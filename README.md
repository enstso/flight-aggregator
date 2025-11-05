# Flight Aggregator

An HTTP API written in Go that **aggregates flights coming from two Node services** (`j-server1` and `j-server2`), normalizes the data, and exposes endpoints to **search** and **sort** flights by **price**, **departure date**, or **total travel time**.
CORS is enabled for all endpoints.

---

## Table of contents

* [High-level architecture](#high-level-architecture)
* [Project layout](#project-layout)
* [How it works](#how-it-works)
* [Getting started (Docker Compose)](#getting-started-docker-compose)
* [Environment](#environment)
* [API](#api)
* [Testing](#testing)
* [Development tips](#development-tips)
* [Troubleshooting & gotchas](#troubleshooting--gotchas)

---

## High-level architecture

```
+-------------------+          +---------------------+
|  j-server1 (Node) |  /flights|  j-server2 (Node)   |
|  :4001            |<-------->|  :4002              |
+-------------------+          +---------------------+
         \                            /
          \                          /
           \                        /
            v                      v
        +----------------------------------+
        |   Go API :8080 (Flight Aggreg.)  |
        |  internal/db      -> fetch JSON  |
        |  internal/repo    -> normalize   |
        |  internal/service -> sorting     |
        +----------------------------------+
                         |
                         v
                   HTTP clients
```

---

## Project layout

```
server/
  internal/
    config/          # viper-based env loader (SERVER1_URL, SERVER2_URL)
    db/              # HTTP client helpers (GetJSON)
    domain/          # core models + repository interface
    handler/         # HTTP handlers (/flights, /health, …)
    health/          # health check response types + handler
    repo/            # repos reading j-server1 & j-server2 payloads + Multi aggregator
    service/         # sorting (price, travel time, departure date)
    test/            # unit tests (testify mocks)
  main.go            # routes, CORS, server bootstrap
  Dockerfile
j-server1/           # Node stub service #1 (exposes /flights)
j-server2/           # Node stub service #2 (exposes /flight_to_book)
docker-compose.yaml  # all 3 services
.env                 # ports and service names
Makefile             # make test
```

---

## How it works

* **Config** (`internal/config`): uses **Viper** to read `../.env` (because the Go service builds from `./server`) and sets `SERVER1_URL` and `SERVER2_URL`.
* **Fetch** (`internal/db`): `GetJSON(ctx, url)` fetches upstream JSON with a 5s timeout.
* **Repositories** (`internal/repo`):

    * `RepoFlights` parses `j-server1`’s `/flights` list.
    * `RepoFlightToBook` parses `j-server2`’s `/flight_to_book` list.
    * Both map their different payloads into the common **domain** model `Flight`.
    * `Multi` composes any number of repositories and queries them uniformly.
* **Service layer** (`internal/service`): implements:

    * `SortByPrice`
    * `SortByTimeTravel` (first departure → last arrival)
    * `SortByDepartureDate`
* **Handlers** (`internal/handler`): HTTP endpoints that:

    * build a `repo.Multi` by fetching from both sources,
    * run queries/sorts,
    * return **JSON** or appropriate errors.
* **Health** (`/health`): pings both upstream services; returns `200` only if both are up.

---

## Getting started (Docker Compose)

> You said you run `docker compose` **from the repo root** (outside the Go project). That’s correct.

```bash
# 1) Build all images
docker compose build

# 2) Start in background
docker compose up -d

# 3) View logs
docker compose logs -f server
```

**.env example (root):**

```dotenv
JSERVER1_PORT=4001
JSERVER1_NAME=j-server1
JSERVER2_PORT=4002
JSERVER2_NAME=j-server2
SERVER_PORT=3001
SERVER1_URL=http://localhost:4001/
SERVER2_URL=http://localhost:4002/
```

> The Go server currently listens on **:8080** (see `main.go`).


---

## Environment

* `SERVER1_URL` → base URL used to fetch `flights` (e.g., `http://localhost:4001/`)
* `SERVER2_URL` → base URL used to fetch `flight_to_book` (e.g., `http://localhost:4002/`)

Other variables in `.env` configure the Node services and Compose port mappings.

---

## API

Base URL (hosted by the Go API):

* **Local (default code):** `http://localhost:8080`

### Health

**GET** `/health`

* Checks both upstream servers.
* **200**: both OK → `{ "Status": 200, "Message": "Health Ok" }`
* **200 with 503 payload**: if one is down → `{ "Status": 503, "Message": "Health Not Ok" }`

### List all flights

**GET** `/flights`

* Aggregates data from both Node services.
* **200** `[]Flight`

`Flight` (normalized) schema:

```json
{
  "id": "string",
  "status": "string",
  "passengerName": "string",
  "segments": [
    {
      "flightNumber": "string",
      "from": "IATA",
      "to": "IATA",
      "depart": "RFC3339 timestamp",
      "arrive": "RFC3339 timestamp"
    }
  ],
  "total": { "amount": 123.45, "currency": "USD" },
  "source": "flights | flight_to_book"
}
```

### Find by ID

**GET** `/flights/id/{id}`

* **200** `Flight`
* **404** if not found

### Find by flight number

**GET** `/flights/number/{flightNumber}`

* **200** `Flight`
* **404** if not found

### Find by passenger name

**GET** `/flights/passengerName/{name}`

* **200** `[]Flight`
* **404** if none

### Find by destination

**GET** `/flights/destination`
**Body:** JSON (**yes, GET + body** by design)

```json
{
  "departure": "CDG",
  "arrival": "HND"
}
```

* **200** `[]Flight`
* **400** invalid JSON
* **404** if none
* **502** if an upstream service fails

**cURL**

```bash
curl -X GET "http://localhost:8080/flights/destination" \
  -H "Content-Type: application/json" \
  -d '{"departure":"CDG","arrival":"HND"}'
```

### Find by exact price

**GET** `/flights/price/{amount}`

* **200** `[]Flight`
* **404** if none

### Sorted list

**GET** `/flights/sorted?type=price|time|duration|departure`

* `price` → by `total.amount` ascending
* `time`/`duration` → by total travel time ascending
* `departure`/`depart`/`departure_date` → by earliest segment departure

**cURL examples**

```bash
curl "http://localhost:8080/flights/sorted?type=price"
curl "http://localhost:8080/flights/sorted?type=time"
curl "http://localhost:8080/flights/sorted?type=departure"
```

### Common error codes

* **400** – bad input (e.g., invalid JSON on `/flights/destination`)
* **404** – not found (`ErrFlightNotFound` / `ErrFlightsNotFound`)
* **405** – method not allowed (only GET is supported)
* **500** – internal encode/processing error
* **502** – upstream fetch failed

---

## Testing

All tests live under `server/internal/test/` and use **Testify** (mocks + assertions).

* **Run all tests** (as requested):

  ```bash
  make test
  ```

  > `Makefile` target runs: `go test -v ./internal/test/...`

* Areas covered:

    1. **Service layer**

        * `SortByPrice` sorts ascending, handles errors & empty lists.
        * `SortByTimeTravel` sorts by total travel time; also covers multi-segment connections.
        * `SortByDepartureDate` sorts by first segment departure, including flights with no segments.
        * `TotalTravelTime` unit cases (single/multi/no segments).
    2. **Repository aggregator (`repo.Multi`)**

        * `List` aggregates across multiple repos, handles repo errors and empty repos.
        * `FindByID`, `FindByNumber`, `FindByPassenger`, `FindByDestination`, `FindByPrice`

            * success from first/second repo
            * proper propagation of `ErrFlightNotFound` / `ErrFlightsNotFound`
            * proper aggregation across repos
    3. **Mocks**

        * `MockFlightsRepository` implements the `domain.FlightsRepository` interface using `testify/mock`.

**Install gotestsum (optional):**

```bash
go install gotest.tools/gotestsum@latest
```

---

## Development tips

* **Hot reload**: `air` is included; edit `air.conf` as needed.
* **CORS**: enabled for `GET, POST, PUT, DELETE, OPTIONS` with `*` origin.
* **Code style**: Handlers are thin; repository/service layers handle data parsing and logic.

---

## Troubleshooting & gotchas

* **.env path:** `config.Load()` reads `../.env` (because the Docker build context for the Go app is `./server`). Keep `.env` at the repo root (where `docker-compose.yaml` lives).

* **GET with body:** `/flights/destination` expects a JSON body even though it’s a GET. Some HTTP clients strip bodies on GET—use `curl` or Postman as shown above.

* **Upstream endpoints used by the Go API:**

    * `SERVER1_URL + "flights"`
    * `SERVER2_URL + "flight_to_book"`

---

## License

MIT (or your preferred license).

---

## Credits

* **Go**, **Viper**, **Testify**
* Sample Node services: `j-server1`, `j-server2`
