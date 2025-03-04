# Compare IDs

## Evaluation Criteria

- **How can the ID be generated?**
  - Are there official libraries?
  - Are there plugins? What are the options?
- **Can IDs be generated on the client or Server?**
  - **Write perf:** How do IDs perform under high-concurrency scenarios, where multiple processes/threads are inserting rows simultaneously?
  - **Read perf:** How do different ID types affect read queries, especially on indexed columns with high cardinality?
  - How does this work across replicas with multiple writers?
  - How can we test this?
- **Collision Probability:**
  - What are the theoretical and observed collision risks of each ID type, especially for client-generated IDs?
  - What about server generated IDs? (mention the use of tools like Twitter Snowflake).
  - How do implementations mitigate these risks?
- **Storage:**
  - By what rate does the index size grow per row? (linear, log)
  - By what rate does the data size grow per row? (linear, log)
  - How do inserts affect WAL?
- **Impact on Index Fragmentation**
  - Does the choice of ID lead to index fragmentation, causing bloat or slower operations over time?
- **Human Usability**
  - **Readability:** Are the IDs human-readable? How does this affect use cases where humans interact with or share the IDs?
  - **Copy/Paste Errors:** Are IDs prone to user-input errors (e.g., mistaking similar characters like `O` and `0`)?
- **Compatibility**
  - How do the IDs perform across different database systems (e.g., PostgreSQL, MySQL, SQlite, MongoDB)?
  - Are there constraints for certain databases when using these ID types (e.g., maximum length, indexing efficiency)?

## Usage

This tool has been restructured to generate test data for each ID type one by one, rather than all at once. This allows for more flexibility and better resource management.

### Commands

- **List available ID types:**

  ```
  go run main.go list
  ```

- **Generate test data for a specific ID type:**

  ```
  go run main.go id [id-type] --count [row-count]
  ```

  Example:

  ```
  go run main.go id uuidv4 --count 10000
  ```

- **Merge all test results into a single template_data.json file:**

  ```
  go run main.go merge
  ```

- **Run tests for all ID types with default row counts:**
  ```
  go run main.go all
  ```

### Database Configuration

You can configure the database connection using the following flags:

```
--host string      Database host (default "localhost")
--port int         Database port (default 5432)
--user string      Database user (default "postgres")
--password string  Database password (default "postgres")
--dbname string    Database name (default "postgres")
```

Example:

```
go run main.go id uuidv4 --host mydb.example.com --port 5432 --user myuser --password mypass --dbname mydb
```

### Result dump

ID Type: XID, Count: 1000000, Duration: 5586.00 ms
ID Type: ULID, Count: 1000000, Duration: 6199.00 ms
ID Type: CUID, Count: 1000000, Duration: 6269.00 ms
ID Type: KSUID, Count: 1000000, Duration: 8857.00 ms
ID Type: NanoID, Count: 1000000, Duration: 8755.00 ms
ID Type: UUIDv4, Count: 1000000, Duration: 10150.00 ms
ID Type: UUIDv7, Count: 1000000, Duration: 7269.00 ms
ID Type: TypeID, Count: 1000000, Duration: 6983.00 ms
ID Type: MongoID, Count: 1000000, Duration: 6104.00 ms
ID Type: Snowflake, Count: 1000000, Duration: 5805.00 ms
ID Type: BigSerial, Count: 1000000, Duration: 3347.00 ms
