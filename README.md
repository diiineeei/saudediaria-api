# Saúde Diária API

API em Go (Gin + MongoDB) usada pelo app **Glicemia da Mary**.

## Rotas de glicemia

| Método | Rota | O que faz |
|---|---|---|
| `POST` | `/api/glicemia/:person` | Registra `{ glucose, time, period? }` |
| `GET` | `/api/glicemia/:person/list` | Lista as medições |
| `PUT` | `/api/glicemia/:person/:id` | Edita `{ glucose, time, date?, period }` (`period` vazio remove o momento) |
| `DELETE` | `/api/glicemia/:person/:id` | Exclui a medição |

Momentos aceitos em `period`: `antes_cafe`, `pos_cafe`, `antes_almoco`, `pos_almoco`, `antes_jantar`, `pos_jantar`.

## Deploy automático (GitHub Actions → Cloud Run)

O workflow `.github/workflows/deploy-cloud-run.yml` roda os testes em todo PR e, a cada push no `main`, publica no Cloud Run (serviço `saudediaria`, região `southamerica-east1`). Também dá para rodar manualmente em **Actions → Testar e publicar no Cloud Run → Run workflow**.

### Configuração (uma vez só, pelo navegador)

1. No [Console do Google Cloud](https://console.cloud.google.com/), no projeto da API, abra **IAM e administrador → Contas de serviço → Criar conta de serviço**. Nome sugerido: `github-deploy`.
2. Dê a ela estes papéis:
   - Administrador do Cloud Run (`roles/run.admin`)
   - Usuário da conta de serviço (`roles/iam.serviceAccountUser`)
   - Editor do Cloud Build (`roles/cloudbuild.builds.editor`)
   - Administrador do Artifact Registry (`roles/artifactregistry.admin`)
   - Administrador do Storage (`roles/storage.admin`)
   - Consumidor do Service Usage (`roles/serviceusage.serviceUsageConsumer`)
3. Abra a conta criada → **Chaves → Adicionar chave → Criar nova chave → JSON**. O arquivo é baixado.
4. No GitHub, em **Settings → Secrets and variables → Actions → New repository secret**, crie `GCP_SA_KEY` e cole **todo** o conteúdo do JSON. Depois apague o arquivo do computador.
5. Se o primeiro deploy pedir, ative as APIs **Cloud Build** e **Artifact Registry** no projeto (o erro traz o link).

As variáveis de ambiente do serviço (como `MONGO_CONFIG`) continuam configuradas no Cloud Run; o deploy não mexe nelas.

## Desenvolvimento

```bash
go test ./...
```
