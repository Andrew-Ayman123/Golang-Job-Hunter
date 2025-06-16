# 1. Load environment variables from .env
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)\s*=\s*(.+)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim('"')
        [System.Environment]::SetEnvironmentVariable($name, $value, "Process")
    }
}

# 2. Run goose migration
goose up

# 3. Generate Swagger docs
swag init -g .\cmd\jobhunter\main.go

# 4. Build the Go application
cd .\cmd\jobhunter\
go build -o ../../main.exe
cd ../../

# 5. Run the app
echo "Running the application..."
./main.exe
