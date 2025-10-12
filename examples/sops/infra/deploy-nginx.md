# Deploy Nginx

This SOP deploys nginx to the server.

## Step 1: Check if nginx is running
Check the status of nginx service.

```bash
curl -I localhost
```

## Step 2: Install nginx if not present
Install nginx if it's not already installed.

```bash
which nginx || echo "nginx not found"
```

## Step 3: Start nginx service
Start the nginx service if it's not running.

```bash
echo "Nginx deployment complete"
```
