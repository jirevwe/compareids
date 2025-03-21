# Prelim Results

| ID Type             | 100k Rows (ms)     | 1m Rows (s)        | 10m Rows (s)       |
|---------------------|--------------------|--------------------|--------------------|
| BigSerial           | 276.4461040496826  | 2.9825680255889893 | 26.11849594116211  |
| UUIDv4              | 331.1300277709961  | 4.0545830726623535 | 48.946518898010254 |
| UUIDv4::VARCHAR(36) | 469.35415267944336 | 5.828076124191284  | 160.98254013061523 |
| (BigSerial, UUIDv4) | 356.0628890991211  | 3.38653302192688   | 35.999242067337036 |
| UUIDv7              | 448.045015335083   | 5.243641138076782  | 52.20610499382019  |
| ULID (pgulid)       | 820.4479217529297  | 9.4429190158844    | 95.51932501792908  |
| ULID (pg-ulid)      | 291.0780906677246  | 4.0451600551605225 | 40.463543176651    |
|                     |                    |                    |                    |



## BigSerial
| row count  | table_size | data size | index size | Time to insert (ms) | Insertion Rate | Growth Rate |
|------------|------------|-----------|------------|---------------------|----------------|-------------|
| 100,000    | 6.5 MB     | 4.3 MB    | 2.2 MB     | 292                 | -              | -           |
| 1,000,000  | 64 MB      | 42 MB     | 21 MB      | 3303                | ~11.3x         | ~10x        |
| 10,000,000 | 637 MB     | 422 MB    | 214 MB     | 26372               | ~7.98x         | ~10x        |

## UUIDv4
| row count  | table_size | data size | index size | Time to insert (ms) | Insertion Rate | Growth Rate |
|------------|------------|-----------|------------|---------------------|----------------|-------------|
| 100,000    | 9.4 MB     | 5.09 MB   | 4.3 MB     | 442                 | -              | -           |
| 1,000,000  | 88 MB      | 50 MB     | 38 MB      | 4816                | ~10.89x        | ~9.36x      |
| 10,000,000 | 886 MB     | 498 MB    | 388 MB     | 79622               | ~16.53x        | ~10x        |
