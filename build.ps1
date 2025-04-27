Get-Content .env | ForEach-Object {
    $parts = $_ -split '=', 2
    if ($parts.Length -eq 2) {
        $envName = $parts[0].Trim()
        $envValue = $parts[1].Trim()
        [System.Environment]::SetEnvironmentVariable($envName, $envValue)
    }
}
wails build -clean -ldflags="-X auth/utils.embeddedMasterPassword=$env:MASTER_PASSWORD"
