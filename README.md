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
