#!/bin/bash

# Setup script for Compare IDs project
echo "🚀 Setting up Compare IDs for production-scale testing..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check system resources
echo "📊 Checking system resources..."

# Get system info
TOTAL_MEMORY=$(free -g | awk '/^Mem:/{print $2}')
CPU_CORES=$(nproc)
DISK_SPACE=$(df -h . | awk 'NR==2{print $4}' | sed 's/G//')

echo "System Resources:"
echo "  Memory: ${TOTAL_MEMORY}GB"
echo "  CPU Cores: ${CPU_CORES}"
echo "  Available Disk Space: ${DISK_SPACE}GB"

# Recommend test scale based on resources
if [ "$TOTAL_MEMORY" -ge 32 ] && [ "$CPU_CORES" -ge 16 ]; then
    RECOMMENDED="billion"
    echo "✅ Your system can handle billion-scale testing!"
elif [ "$TOTAL_MEMORY" -ge 16 ] && [ "$CPU_CORES" -ge 8 ]; then
    RECOMMENDED="large"
    echo "✅ Your system can handle large-scale testing!"
elif [ "$TOTAL_MEMORY" -ge 8 ] && [ "$CPU_CORES" -ge 4 ]; then
    RECOMMENDED="medium"
    echo "✅ Your system can handle medium-scale testing!"
else
    RECOMMENDED="small"
    echo "⚠️  Your system is limited. Recommended: small-scale testing"
fi

echo ""
echo "🎯 Recommended test scale: $RECOMMENDED"
echo ""

# Create results directory if it doesn't exist
mkdir -p results

# Make scripts executable
chmod +x run-*.sh

echo "✅ Setup complete!"
echo ""
echo "Quick start commands:"
echo "  Small scale:  ./run-small.sh"
echo "  Medium scale: ./run-medium.sh"
echo "  Large scale:  ./run-large.sh"
echo "  Billion scale: ./run-billion.sh"
echo ""
echo "For your system, we recommend: ./run-$RECOMMENDED.sh"
echo ""
echo "📖 For more information, see README.md"
