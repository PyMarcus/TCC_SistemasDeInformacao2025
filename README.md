# 🔬 TCC Sistemas de Informação 2025 - Detecção de Átomos de Confusão com AI

Este projeto é parte de uma pesquisa científica que busca comparar agentes de IA (como o **Gemini AI**) na tarefa de identificar **átomos de confusão** — elementos mínimos de código que induzem a erros de entendimento — no contexto da **Engenharia de Software**.

> 🧩 Baseado nos princípios de **Clean Architecture** e **SOLID**, com tecnologias modernas como **Go**, **PostgreSQL**, e ferramentas de migração e linting.

---

## 🚀 Tecnologias Utilizadas

- **Golang (Go)**
- **PostgreSQL**
- **Golang Migrate** (migrations do banco)
- **Golangci-lint** (linter)
- **Mockgen** (geração de mocks para testes)
- Estruturado com:
  - **Clean Architecture**
  - Princípios **SOLID**
  - Padrão **Repository**
  - Uso de **worker pool** para execução concorrente
  - Logger estruturado com **Zap**

---

## 📦 Estrutura do Projeto
    internal/
    ├── adapters/
    │ ├── config/
    │ ├── db/
    │ ├── http/
    │ ├── repository/
    ├── core/
    │ ├── domain/
    │ ├── ports/
    │ ├── usecase/
    scripts/
    cmd/
    database/
    tests/

## ⚙️ Como Executar o Projeto

### 1. Clone o projeto
```bash
git clone https://github.com/PyMarcus/TCC_SistemasDeInformacao2025.git
cd TCC_SistemasDeInformacao2025
```

### 2. Configure o banco de dados

Edite o arquivo .env com suas credenciais:

    DATABASE_URL=postgres://marcus:marcus123@localhost:5432/marcus_db?sslmode=disable

## 📂 Comandos Úteis (via Makefile)

Criar uma nova migration
```bash
make create-migrations NAME=nome_da_migration
```
Rodar migrações
```bash
make migrate-up
```
Resetar e subir do zero
```bash
make migrate-reset
```

## 🐳 Ferramentas

Instalar migrate
```bash
make migrate-download
```

Instalar mockgen
```bash
make mockgen-download
```

## ⚡ Executar o Sistema
```bash
make run
```

🛠️ Principais Padrões Usados

* Clean Architecture: separação clara entre camadas domain, usecase, adapters.
    
* SOLID: cada módulo com responsabilidade única, injeção de dependências e interfaces.
    
* Worker Pool: executa múltiplas tarefas simultaneamente de forma segura usando goroutines.
    
* Configuração via .env usando loader customizado.
    
* Logger estruturado usando zap.

* Implementa um controle de pausa global com sincronização usando mutex e condições (sync.Cond) para coordenar a execução concorrente das goroutines.

## 🎯 Objetivo da Pesquisa

Este sistema executa a análise de datasets contendo código-fonte java e perguntas relacionadas, criadas a partir de prompt engineering, com objetivo de:

* Identificar pontos de confusão no código.
* Avaliar respostas geradas por modelos AI (ex: Gemini).
* Gerar relatórios para análise científica e estatística.

## 👨‍🔬 Autor

Marcus (PyMarcus)
Graduando em Sistemas de Informação - 2025
Este projeto faz parte do Trabalho de Conclusão de Curso (TCC) em Engenharia de Software.

## 🐘 Banco de Dados

* PostgreSQL

* Gerenciado com golang-migrate

* Dados versionados via migrações SQL

## 🧹 Linting & Testes

* Linting via golangci-lint

* Testes com go test, incluindo:

        Verificação de race conditions

        Cobertura de código

        Mocks com mockgen

## 📊 Logs e Métricas

* Logs estruturados via zap.

* Métricas de tempo de execução exibidas no final de cada pool executado.

## 🏁 Requisitos

* Go 1.20+

* PostgreSQL 13+

* Ferramentas:

        migrate

        golangci-lint

        mockgen
