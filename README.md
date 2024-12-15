
## Code Structure

-   **`main.go`**: Contains the main logic of the application, including database setup, IP processing, and output.
-   **`setupDatabase`**: Opens a database, creates the table with unique constraint.
-   **`processIPs`**: Opens the file, reads the IPs and attempts to insert them into the database.
-   **`isDuplicateError`**: Verifies if the database returned a duplicate error.
-   **`getMemoryUsage`**: Retrieves information about the memory used by the program.
-   **`main`**: Executes the program, including measuring the time.

Two approaches were considered: **`database insert approach`** vs **`file splitting approach`**
Finally the **`database insert approach`** was implemented
## Approach Comparison

Here's a table comparing the database approach with the file splitting approach:

| Feature  | Database Approach                           | External Sorting (File Splitting)           |
|----------|---------------------------------------------|--------------------------------------------|
| **Memory** |  Moderate, can be controlled                   |  Low, depending on the chunk size                |
| **Disk** |  Proportional to *unique* IPs, can incur overhead due to database operations          |  Proportional to *all* IPs, can be larger than unique IPs if a lot of duplicates are present  |
| **Speed** | Usually Faster if a lot of duplicate IPs |    Slower due to I/O operations, and processing of duplicate IP addresses           |
| **Complexity** | Moderate, uses database functions                   | Higher complexity for implementing file sorting and merging |
| **Accuracy** | 100% guaranteed                             | 100% guaranteed                               |
| **Use Case** |  Suitable for large datasets where you need persistence. | Suitable for very large files when RAM is very limited.      |


## Considerations

-   The program uses SQLite for simplicity and does not require a running database server.
-   You can replace `ip_addresses.txt` with any input file containing IP addresses.
-   The database file will be located in the same directory as the program.
-   The program uses a `DROP TABLE` statement every time it is run, therefore any previous data on the table will be lost. If you need persistence, consider removing the `DROP TABLE` statement.

## Further Improvements

-   Implement batch operations for database inserts for better performance.
-   Add configuration file to configure input and output files.