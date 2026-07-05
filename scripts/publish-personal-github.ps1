# Publish this project to your PERSONAL GitHub account (IsuruLakshan170).
# Does NOT change global git config — only this repository.
#
# Usage (PowerShell, run from repo root):
#   $env:GH_TOKEN = "ghp_your_personal_access_token"
#   .\scripts\publish-personal-github.ps1
#
# Create a token at: https://github.com/settings/tokens
# Required scope: repo (for private) or public_repo (for public repos)

$ErrorActionPreference = "Stop"

$Owner = "IsuruLakshan170"
$RepoName = "ecommerce-yt"
$Description = "Production-oriented ecommerce REST API - Go, Gin, MongoDB, JWT"
$IsPrivate = $false

if (-not $env:GH_TOKEN) {
    Write-Error @"
GH_TOKEN is not set.

1. Open https://github.com/settings/tokens (logged in as $Owner)
2. Generate a classic token with 'repo' or 'public_repo' scope
3. Run in THIS terminal only (does not affect office git):

   `$env:GH_TOKEN = "ghp_xxxxxxxx"
   .\scripts\publish-personal-github.ps1
"@
}

$headers = @{
    Authorization = "Bearer $env:GH_TOKEN"
    Accept        = "application/vnd.github+json"
    "X-GitHub-Api-Version" = "2022-11-28"
}

$body = @{
    name        = $RepoName
    description = $Description
    private     = $IsPrivate
    auto_init   = $false
} | ConvertTo-Json

Write-Host "Creating repository $Owner/$RepoName ..."
try {
    $repo = Invoke-RestMethod -Method POST -Uri "https://api.github.com/user/repos" -Headers $headers -Body $body -ContentType "application/json"
    Write-Host "Created: $($repo.html_url)"
} catch {
    $status = $_.Exception.Response.StatusCode.value__
    if ($status -eq 422) {
        Write-Host "Repository already exists — continuing with push."
    } else {
        throw
    }
}

$remoteUrl = "https://github.com/$Owner/$RepoName.git"
git remote remove origin 2>$null
git remote add origin $remoteUrl
Write-Host "Remote set to $remoteUrl"

Write-Host "Pushing main branch ..."
$env:GIT_TERMINAL_PROMPT = "0"
git -c credential.helper= push -u "https://$Owner`:$($env:GH_TOKEN)@github.com/$Owner/$RepoName.git" main

Write-Host ""
Write-Host "Done! Repository: https://github.com/$Owner/$RepoName"
Write-Host "Local git identity (this repo only):"
git config --local user.name
git config --local user.email
