# 1Mão - Plataforma de Orquestração de Serviços

**1Mão** é uma plataforma robusta e escalável desenvolvida em Go, destinada à orquestração de microsserviços com foco em desempenho, modularidade e comunicação em tempo real.

## 🧩 Arquitetura Modular

A estrutura modular do projeto segue o padrão de separação de responsabilidades por domínio:

- `cmd/`: ponto de entrada da aplicação.
- `config/`: configurações de banco de dados e cache.
- `delivery/rest/`: handlers e rotas REST.
- `internal/`: lógica de domínio, serviços, repositórios e middleware organizados por contexto (ex: `booking`, `client`, `payment`).
- `pkg/`: bibliotecas reutilizáveis, como autenticação e cache.

## 🔐 Autenticação JWT

A autenticação é feita com tokens JWT, com suporte a middlewares para controle de acesso e segurança.

## 🔄 Comunicação em Tempo Real

Utilizamos WebSockets no módulo de notificações para garantir uma comunicação bidirecional entre clientes e profissionais em tempo real.

## 📦 Integrações

- **Stripe**: Processamento de pagamentos.
- **Redis**: Cache e controle de sessão.
- **Swagger**: Documentação interativa da API.

## 🚀 Executando o Projeto

### Pré-requisitos

- Docker & Docker Compose
- Go 1.21+

### Subindo com Docker

```bash
docker-compose up --build
```

## 🧪 Testes unitários

```bash
go test ./...
```

### Acessando a API

- `http://localhost:8080/api`
- Swagger: `http://localhost/swagger/index.html`

## Teste de chat com WebSocket

Utilize um utilitário para conexões websocket, como o wscat

```bash
wscat -c ws://localhost/ws/chat/<tipo de remetente>/<id do remetente>
```

Ao entrar na interface do wscat, utilize
```bash
{"receiver_id":<id do destinatario>,"receiver_type":"<tipo do destinatário>","content":"mensagem a ser enviada"}
```


## 📁 Documentação

A documentação OpenAPI/Swagger pode ser encontrada em:

- `docs/swagger.yaml`
- `docs/swagger.json`

## 📂 Workflows

CI configurado com GitHub Actions: `.github/workflows/go.yml`

## 🤝 Contribuindo

1. Fork o repositório
2. Crie sua branch (`git checkout -b feature/nome-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para o branch (`git push origin feature/nome-feature`)
5. Abra um Pull Request

---

**1Mão** © 2025 - Plataforma de Orquestração com Go 🚀
