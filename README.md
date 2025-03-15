# Compare IDs

## Usage

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

- **Merge all test results into a single ata.json file:**

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

