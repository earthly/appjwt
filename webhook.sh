curl -X POST \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  https://api.github.com/repos/pantalasa/lunar/hooks \
  -d '{
    "name": "web",
    "active": true,
    "events": ["push", "pull_request"],
    "config": {
      "url": "https://hub.com/webhook/github",
      "content_type": "json",
      "secret": "secret"
    }
  }'
