# 📝 HashiCorp Vault Notes — Setup, Secrets, and UI

Personal reference notes for setting up Vault on Kubernetes, writing/reading secrets via CLI, and doing the same through the Vault UI.

---

## 1. Install Vault on the cluster

```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

helm install vault hashicorp/vault \
  --set "server.dev.enabled=true" \
  --set "injector.enabled=true" \
  -n vault --create-namespace

kubectl -n vault get pods -w
```

Wait until:
- `vault-0` → `1/1 Running`
- `vault-agent-injector-xxxx` → `1/1 Running`

> `server.dev.enabled=true` runs Vault in **dev mode** — auto-unsealed, in-memory storage, root token pre-generated. Good for learning/demo, **never use in production**.

---

## 2. Configure Kubernetes auth (so pods can log in to Vault)

```bash
kubectl exec -it vault-0 -n vault -- sh
```

Inside the pod (run each as **one line** — multi-line `\` continuations can break in `sh`):

```bash
vault auth enable kubernetes

vault write auth/kubernetes/config kubernetes_host="https://kubernetes.default.svc:443" disable_iss_validation=true
```

Create a policy (defines what a role is allowed to read):

```bash
vault policy write go-crud-app-policy - <<EOF
path "secret/data/go-crud-app" {
  capabilities = ["read"]
}
EOF
```

Create a role (binds a K8s ServiceAccount to that policy):

```bash
vault write auth/kubernetes/role/go-crud-app bound_service_account_names=go-crud-app-sa bound_service_account_namespaces=default policies=go-crud-app-policy ttl=24h
```

---

## 3. Write a secret (CLI)

```bash
vault kv put secret/go-crud-app \
  POSTGRES_USER="myuser" \
  POSTGRES_PASSWORD="mypass" \
  POSTGRES_DB="mydb"
```

Update just one field later (⚠️ `kv put` overwrites the whole secret, so include all keys each time, or use `kv patch` to update selectively):

```bash
vault kv patch secret/go-crud-app POSTGRES_PASSWORD="newpass123"
```

---

## 4. Read/load a secret (CLI)

```bash
vault kv get secret/go-crud-app
```

Get just one field:

```bash
vault kv get -field=POSTGRES_PASSWORD secret/go-crud-app
```

Get as JSON (useful for scripting):

```bash
vault kv get -format=json secret/go-crud-app
```

---

## 5. Verify configuration anytime

```bash
vault auth list                              # confirm kubernetes/ auth method is enabled
vault read auth/kubernetes/config             # confirm host + settings
vault read auth/kubernetes/role/go-crud-app   # confirm role bindings
vault policy read go-crud-app-policy          # confirm policy rule
vault kv get secret/go-crud-app               # confirm secret contents
```

---

## 6. How a pod loads the secret automatically (Vault Agent Injector)

Add these annotations to a pod template — no code changes needed in the app:

```yaml
metadata:
  annotations:
    vault.hashicorp.com/agent-inject: "true"
    vault.hashicorp.com/role: "go-crud-app"
    vault.hashicorp.com/agent-inject-secret-config: "secret/data/go-crud-app"
    vault.hashicorp.com/agent-inject-template-config: |
      {{- with secret "secret/data/go-crud-app" -}}
      export POSTGRES_USER="{{ .Data.data.POSTGRES_USER }}"
      export POSTGRES_PASSWORD="{{ .Data.data.POSTGRES_PASSWORD }}"
      export POSTGRES_DB="{{ .Data.data.POSTGRES_DB }}"
      {{- end }}
```

The sidecar writes the rendered file to `/vault/secrets/config` inside every container in the pod. The app's start command sources it:

```yaml
command: ["/bin/sh", "-c"]
args:
  - . /vault/secrets/config && ./api
```

Check it landed correctly:

```bash
kubectl exec -it deploy/go-crud-app -c app -- cat /vault/secrets/config
```

Check the sidecar's own logs (auth + render status):

```bash
kubectl logs deploy/go-crud-app -c vault-agent
```

**Live-update test** — change the secret in Vault, then re-check the file after a few seconds (no pod restart needed, the agent re-renders automatically):

```bash
vault kv put secret/go-crud-app POSTGRES_USER="myuser" POSTGRES_PASSWORD="rotated" POSTGRES_DB="mydb"
kubectl exec -it deploy/go-crud-app -c app -- cat /vault/secrets/config
```

---

## 7. Using the Vault UI

**Step 1 — port-forward to reach it:**
```bash
kubectl port-forward svc/vault -n vault 8200:8200
```

**Step 2 — get the root token (dev mode):**
```bash
kubectl -n vault logs vault-0 | grep "Root Token"
```
(In dev mode, the root token is printed once in the pod's startup logs.)

**Step 3 — open the UI:**
Go to `http://localhost:8200/ui` in the browser → paste the root token in **Sign in with Token**.

**Step 4 — write a secret via UI:**
- Left sidebar → **Secrets Engines** → click `secret/`
- Click **Create secret +**
- Path: `go-crud-app`
- Add key/value pairs: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- Click **Save**

**Step 5 — read a secret via UI:**
- **Secrets Engines** → `secret/` → click `go-crud-app`
- Values are hidden by default — click the eye icon to reveal, or "Copy" to clipboard
- **Version History** tab shows every past write (KV v2 keeps version history automatically)

**Step 6 — check auth methods / roles via UI:**
- Left sidebar → **Access** → **Auth Methods** → `kubernetes/` → shows config
- **Policies** tab → view `go-crud-app-policy` rule directly
