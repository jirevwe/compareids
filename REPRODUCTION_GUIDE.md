# Reproduction Guide - Compare IDs

This guide explains how to reproduce the ID comparison results on different machine setups.

## 🎯 Quick Start

1. **Clone and setup**:

   ```bash
   git clone <repository-url>
   cd comapreids
   ./setup.sh
   ```

2. **Choose your scale**:

   ```bash
   # For quick validation (5-10 minutes)
   ./run-small.sh

   # For standard testing (2-4 hours)
   ./run-medium.sh

   # For production-like testing (8-16 hours)
   ./run-large.sh

   # For extreme scale testing (24-72 hours)
   ./run-billion.sh
   ```

3. **View results**:
   ```bash
   # Results are in ./results/data.json
   cat results/data.json | jq '.'
   ```

## 🖥️ Machine Configurations

### Development Machine (4GB RAM, 2 CPU cores)

```bash
# Use small scale testing
./run-small.sh
```

- **Duration**: 5-10 minutes
- **Records**: 1K-10K
- **ID Types**: 4 core types

### Standard Laptop (8GB RAM, 4 CPU cores)

```bash
# Use medium scale testing
./run-medium.sh
```

- **Duration**: 2-4 hours
- **Records**: 1K-1M
- **ID Types**: All 16 types

### High-End Workstation (16GB RAM, 8 CPU cores)

```bash
# Use large scale testing
./run-large.sh
```

- **Duration**: 8-16 hours
- **Records**: 10K-10M
- **ID Types**: All 16 types

### Server/Cloud Instance (32GB+ RAM, 16+ CPU cores)

```bash
# Use billion scale testing
./run-billion.sh
```

- **Duration**: 24-72 hours
- **Records**: 1M-1B
- **ID Types**: 6 core types

## 🔧 Custom Configurations

### Adjusting for Your Hardware

1. **Copy an environment file**:

   ```bash
   cp env.medium my-config.env
   ```

2. **Edit resource limits**:

   ```bash
   nano my-config.env
   ```

   Key settings:

   ```bash
   # Adjust based on your RAM
   MEMORY_LIMIT=8G              # Use 50-75% of your RAM
   POSTGRES_MEMORY_LIMIT=4G     # Use 25-50% of your RAM

   # Adjust based on your CPU cores
   CPU_LIMIT=4                  # Use 50-75% of your cores
   POSTGRES_CPU_LIMIT=2         # Use 25-50% of your cores
   ```

3. **Run with custom config**:
   ```bash
   cp my-config.env .env
   docker-compose up --build
   ```

### Testing Individual ID Types

```bash
# Test only specific ID types
TEST_SCENARIO=medium MAX_ROWS=1000000 docker-compose up --build

# Or modify config files to include only desired ID types
```

## 📊 Expected Results

### Performance Characteristics

| ID Type     | Insert Speed | Storage Size | Index Efficiency |
| ----------- | ------------ | ------------ | ---------------- |
| `bigserial` | Fastest      | Smallest     | Highest          |
| `snowflake` | Fast         | Small        | High             |
| `uuidv7`    | Medium       | Medium       | Medium           |
| `ulid`      | Medium       | Medium       | Medium           |
| `uuidv4`    | Slow         | Large        | Low              |
| `cuid`      | Slow         | Large        | Low              |

### Typical Results (1M records)

- **bigserial**: ~3-5 seconds, ~44MB data, ~22MB index
- **uuidv7**: ~4-6 seconds, ~52MB data, ~32MB index
- **ulid**: ~6-8 seconds, ~68MB data, ~50MB index
- **uuidv4**: ~7-12 seconds, ~52MB data, ~40MB index

## 🚨 Troubleshooting

### Common Issues

1. **Out of Memory**:

   ```bash
   # Reduce memory limits
   MEMORY_LIMIT=4G
   POSTGRES_MEMORY_LIMIT=2G
   ```

2. **Slow Performance**:

   ```bash
   # Check system resources
   docker stats

   # Ensure SSD storage
   # Close other applications
   ```

3. **Test Failures**:

   ```bash
   # Check logs
   docker-compose logs ids

   # Restart with clean state
   docker-compose down -v
   docker-compose up --build
   ```

### Performance Optimization

1. **Use SSD storage** for better I/O performance
2. **Close other applications** to free up resources
3. **Use local Docker** instead of remote Docker
4. **Monitor system resources** during tests
5. **Adjust PostgreSQL settings** for your hardware

## 📈 Comparing Results

### Standard Metrics

When comparing results across different machines:

1. **Normalize by hardware**:

   - CPU cores and speed
   - RAM amount and speed
   - Storage type (SSD vs HDD)

2. **Focus on relative performance**:

   - Which ID types are fastest/slowest
   - Storage efficiency comparisons
   - Index performance differences

3. **Consider use cases**:
   - High-frequency inserts
   - Large dataset storage
   - Query performance
   - Distributed systems

### Example Comparison

```bash
# Run on Machine A
./run-medium.sh

# Run on Machine B
./run-medium.sh

# Compare results
diff results/data.json results-machine-b/data.json
```

## 🔄 Continuous Testing

### Automated Testing

```bash
# Run tests periodically
crontab -e

# Add entry for weekly testing
0 2 * * 0 cd /path/to/comapreids && ./run-medium.sh
```

### CI/CD Integration

```yaml
# Example GitHub Actions workflow
name: ID Performance Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run small scale tests
        run: ./run-small.sh
```

## 📝 Reporting Results

### Sharing Results

1. **Include system specs**:

   - CPU: Model, cores, speed
   - RAM: Amount, type, speed
   - Storage: Type, speed
   - OS: Version, kernel

2. **Include test configuration**:

   - Test scale used
   - Custom modifications
   - Environment variables

3. **Include raw data**:
   - `results/data.json`
   - System logs
   - Performance metrics

### Example Report

```markdown
## Test Results - Machine XYZ

**System Specs:**

- CPU: Intel i7-12700K, 12 cores, 3.6GHz
- RAM: 32GB DDR4-3200
- Storage: NVMe SSD, 1TB
- OS: Ubuntu 22.04 LTS

**Test Configuration:**

- Scale: Large (10K-10M records)
- Duration: 12 hours
- ID Types: All 16 types

**Key Findings:**

- bigserial: Fastest insertion (3.2s for 1M records)
- uuidv7: Best balance of performance and features
- uuidv4: Slowest but most compatible

**Raw Data:** [results/data.json]
```

## 🎉 Success Criteria

Your reproduction is successful if:

1. ✅ Tests complete without errors
2. ✅ Results show expected performance patterns
3. ✅ Relative performance between ID types is consistent
4. ✅ System resources are utilized efficiently
5. ✅ Results are reproducible across runs

## 📞 Support

If you encounter issues:

1. Check this guide first
2. Review the main README.md
3. Check existing issues
4. Create a new issue with:
   - System specifications
   - Test configuration used
   - Error logs
   - Steps to reproduce
